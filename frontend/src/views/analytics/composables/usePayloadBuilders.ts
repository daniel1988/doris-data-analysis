import type { EventAnalysisReq, FunnelAnalysisReq, GlobalFilterGroups, Group, Order, RetentionAnalysisReq, TimeGrain, UserPropertyAnalysisReq } from '@/types/doris/analysis'
import { EVENT_TABLE, Operator as OP } from '@/types/doris/common'
import type { EventMetric } from '@/types/doris/metric'
import { deepNormalizeMetric } from '../utils/payloadNormalizer'
import { useCustomFormula } from './useCustomFormula'

export function usePayloadBuilders() {
    const { buildCustomMetric } = useCustomFormula()

    /**
     * 解析过滤值，处理逗号分隔的字符串为数组
     */
    const parseFilterValue = (f: any) => {
        if ([OP.IsNull, OP.IsNotNull].includes(f.operator)) return undefined
        const raw = f.value
        if (typeof raw !== 'string') return raw
        if ([OP.Between, OP.In, OP.NotIn].includes(f.operator)) {
            return raw.split(',').map((s: string) => s.trim()).filter(Boolean)
        }
        return raw
    }

    const buildGlobalFilterGroups = (filter_groups: GlobalFilterGroups, dynamicDateFilter?: { operator: number; value: [string, string] } | null) => {
        const finalFilterGroups: GlobalFilterGroups = {
            ...filter_groups,
            global_filters: {
                ...filter_groups.global_filters,
                filters: (filter_groups.global_filters.filters || []).map(f => ({
                    ...f,
                    value: parseFilterValue(f)
                }))
            }
        }

        // 处理对比时间
        if (Array.isArray(filter_groups.comparison_query_dates) && filter_groups.comparison_query_dates.length === 2) {
            finalFilterGroups.comparison_query_dates = filter_groups.comparison_query_dates
        }

        // 处理动态时间
        if (dynamicDateFilter) {
            finalFilterGroups.global_filters.filters.push({
                column: { table: EVENT_TABLE, field: 'e_event_time', alias: '' },
                operator: dynamicDateFilter.operator,
                value: dynamicDateFilter.value
            })
            // 动态时间下清空静态日期
            finalFilterGroups.query_dates = []
        }

        return finalFilterGroups
    }

    const buildPayload = (
        project_alias: string,
        event_metrics: EventMetric[],
        filter_groups: GlobalFilterGroups,
        time_grain: TimeGrain,
        groups: Group[],
        orders: Order[],
        dynamicDateFilter?: { operator: number; value: [string, string] } | null,
        page_size: number = 1000,
        page_num: number = 1
    ): EventAnalysisReq => {
        // 转换自定义指标并规范化
        const finalMetrics = event_metrics.map(metric => {
            const metricCopy = JSON.parse(JSON.stringify(metric))

            // 如果是前端 UI 维护的自定义指标 Token 序列，先构建出 API 结构
            if (metricCopy.isCustom && metricCopy.custom_web_metric) {
                const webMetric = {
                    tokens: metricCopy.custom_web_metric.eventMetrics,
                    format: metricCopy.custom_web_metric.format || ''
                }
                metricCopy.custom_metric = buildCustomMetric(webMetric as any)
            }

            // 统一规范化处理（包括自定义指标内部的 event_metrics 和普通指标）
            return deepNormalizeMetric(metricCopy)
        })

        const finalFilterGroups = buildGlobalFilterGroups(filter_groups, dynamicDateFilter)

        return {
            project_alias,
            event_metrics: finalMetrics,
            filter_groups: finalFilterGroups,
            time_grain,
            groups,
            orders,
            page_size,
            page_num
        }
    }

    const buildRetentionPayload = (
        project_alias: string,
        init_event_metric: EventMetric,
        end_event_metric: EventMetric,
        filter_groups: GlobalFilterGroups,
        time_grain: TimeGrain,
        groups: Group[],
        day_n_array: number[],
        dynamicDateFilter?: { operator: number; value: [string, string] } | null
    ): RetentionAnalysisReq => {
        const finalFilterGroups = buildGlobalFilterGroups(filter_groups, dynamicDateFilter)

        // 规范化留存分析指标
        const initMetricCopy = deepNormalizeMetric(JSON.parse(JSON.stringify(init_event_metric)))
        const endMetricCopy = deepNormalizeMetric(JSON.parse(JSON.stringify(end_event_metric)))

        return {
            project_alias,
            init_event_metric: initMetricCopy,
            end_event_metric: endMetricCopy,
            global_filter_groups: finalFilterGroups,
            time_grain,
            groups,
            day_n_array
        }
    }

    const buildFunnelPayload = (
        project_alias: string,
        event_metrics: EventMetric[],
        filter_groups: GlobalFilterGroups,
        time_grain: TimeGrain,
        groups: Group[],
        dynamicDateFilter?: { operator: number; value: [string, string] } | null
    ): FunnelAnalysisReq => {
        const finalFilterGroups = buildGlobalFilterGroups(filter_groups, dynamicDateFilter)

        // 规范化漏斗分析步骤指标
        const finalSteps = event_metrics.map(metric => {
            return deepNormalizeMetric(JSON.parse(JSON.stringify(metric)))
        })

        // 转化窗口期：天 -> 秒
        const windowInSeconds = (time_grain.window_num || 1) * 24 * 3600

        return {
            project_alias,
            steps: finalSteps,
            global_filter_groups: finalFilterGroups,
            time_grain,
            groups,
            window: windowInSeconds,
            page_size: 1000,
            page_num: 1
        }
    }

    const buildPropertyPayload = (
        project_alias: string,
        metric: any,
        filter_groups: GlobalFilterGroups,
        groups: Group[],
        user_groups: any[],
        group_type: number,
        dynamicDateFilter?: { operator: number; value: [string, string] } | null
    ): UserPropertyAnalysisReq => {
        // 规范化指标
        const normalizedMetric = deepNormalizeMetric(JSON.parse(JSON.stringify(metric)))

        // 转换全局过滤（处理逗号分隔值等）
        const finalFilterGroups = buildGlobalFilterGroups(filter_groups, dynamicDateFilter)

        // 转换人群分群结构 (从前端包装的 GlobalFilterGroups 提取内部的 FilterGroup)
        const processedUserGroups = user_groups.map((ug: any) => ({
            alias: ug.alias,
            filter_group: buildGlobalFilterGroups(ug.filter_group).global_filters
        }))

        return {
            project_alias,
            metric: normalizedMetric.metric,
            filter_groups: finalFilterGroups,
            groups,
            user_groups: processedUserGroups,
            group_type,
            page_size: 1000,
            page_num: 1
        }
    }

    return {
        buildPayload,
        buildRetentionPayload,
        buildFunnelPayload,
        buildPropertyPayload
    }
}

