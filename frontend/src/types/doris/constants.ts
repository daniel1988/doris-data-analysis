import { Formula, Operator, EVENT_TABLE } from '../doris/common'

export const MetricFormulas = Formula
export const FilterOperators = Operator

export const MetricFormulaLabels: Record<number, string> = {
    [Formula.Empty]: '空（不聚合）',
    [Formula.Sum]: '求和',
    [Formula.Count]: '总次数',
    [Formula.CountDistinct]: '去重数',
    [Formula.CountDistinctUserId]: '用户数',
    [Formula.Any]: '任意值',
    [Formula.Max]: '最大值',
    [Formula.Min]: '最小值',
    [Formula.Avg]: '平均值',
    [Formula.IsNull]: '为空',
    [Formula.IsNotNull]: '不为空',
    [Formula.IsTrue]: '为真',
    [Formula.IsFalse]: '为假',
    [Formula.First]: '首次',
    [Formula.Last]: '末次',
    [Formula.CountDays]: '按天计数',
    [Formula.LastValue]: '最后值',
    [Formula.FirstValue]: '最先值',
    [Formula.BitmapUnion]: '位图并集',
    [Formula.CountDistinctDailyUserId]: '按天去重数'
}

export const OperatorLabels: Record<number, string> = {
    [Operator.Empty]: '空',
    [Operator.EqualTo]: '等于',
    [Operator.NotEqualTo]: '不等于',
    [Operator.LessThan]: '小于',
    [Operator.LessOrEqual]: '小于等于',
    [Operator.GreaterThan]: '大于',
    [Operator.GreaterOrEqual]: '大于等于',
    [Operator.IsNull]: '为空',
    [Operator.IsNotNull]: '不为空',
    [Operator.Between]: '区间',
    [Operator.In]: '包含',
    [Operator.NotIn]: '不包含',
    [Operator.Like]: '匹配',
    [Operator.DateDiff]: '相对事件发生时刻',
    [Operator.NotLike]: '不匹配',
    [Operator.BelongTo]: '属于',
    [Operator.NotBelongTo]: '不属于',
    [Operator.StartWith]: '以...开头',
    [Operator.EndWith]: '以...结尾',
    [Operator.DynamicDates]: '动态日期',
    [Operator.NDayRegiste]: '与当天间隔的天数'
}

export const TableNames = {
    EVENT_TABLE: EVENT_TABLE
}

export const FORMULA_OPTIONS = Object.entries(MetricFormulaLabels).map(([value, label]) => ({
    value: Number(value),
    label
}))

export const OPERATOR_OPTIONS = Object.entries(OperatorLabels).map(([value, label]) => ({
    value: Number(value),
    label
}))
