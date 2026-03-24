/**
 * Doris 相关表名常量
 */
export const EVENT_TABLE = 'event_data'
export const USER_TABLE = 'user_data'
export const TAG_TABLE = 'user_tag_data'
export const USER_GROUP_TABLE = 'user_group_data'
export const USER_PROPERTY_TABLE = 'ups'

export const EVENT_TABLE_ALIAS = 'events'
export const USER_TABLE_ALIAS = 'users'

/**
 * 过滤操作符枚举
 */
export enum Operator {
    Empty = 0,
    EqualTo = 1,
    NotEqualTo = 2,
    LessThan = 3,
    LessOrEqual = 4,
    GreaterThan = 5,
    GreaterOrEqual = 6,
    IsNull = 7,
    IsNotNull = 8,
    Between = 9,
    In = 10,
    NotIn = 11,
    Like = 12,
    DateDiff = 13,
    NotLike = 14,
    BelongTo = 15,
    NotBelongTo = 16,
    StartWith = 17,
    EndWith = 18,
    DynamicDates = 19, // 动态日期
    NDayRegiste = 20, // N 天前注册
}

/**
 * 指标计算公式枚举
 */
export enum Formula {
    Empty = 0,
    Sum = 1,
    Count = 2,
    CountDistinct = 3,
    CountDistinctUserId = 4,
    Any = 5,
    Max = 6,
    Min = 7,
    Avg = 8,
    IsNull = 9,
    IsNotNull = 10,
    IsTrue = 11,
    IsFalse = 12,
    First = 13,
    Last = 14,
    CountDays = 15,
    LastValue = 16, // 指定列最大值： 即按照事件时间倒序排序第一条
    FirstValue = 17, // 指定列最小值： 即按照事件时间正序排序第一条
    BitmapUnion = 18,
    CountDistinctDailyUserId = 19, // 按天去重计数
}

/**
 * 格式化类型
 */
export const FormatDefault = 'raw'
export const FormatInt = 'int'
export const FormatDecimal = 'decimal'
export const FormatPercent = 'percent'

/**
 * 数据类型枚举
 */
export enum DataType {
    String = 'string',
    Number = 'number',
    DateTime = 'datetime',
    Boolean = 'boolean'
}

/**
 * 定义操作符的展示文本
 */
export const OperatorLabels: Record<Operator, string> = {
    [Operator.Empty]: '为空',
    [Operator.EqualTo]: '等于',
    [Operator.NotEqualTo]: '不等于',
    [Operator.LessThan]: '小于',
    [Operator.LessOrEqual]: '小于等于',
    [Operator.GreaterThan]: '大于',
    [Operator.GreaterOrEqual]: '大于等于',
    [Operator.IsNull]: '为空 (NULL)',
    [Operator.IsNotNull]: '不为空 (NOT NULL)',
    [Operator.Between]: '区间 (Between)',
    [Operator.In]: '包含 (IN)',
    [Operator.NotIn]: '不包含 (NOT IN)',
    [Operator.Like]: '模糊匹配',
    [Operator.DateDiff]: '日期差',
    [Operator.NotLike]: '不模糊匹配',
    [Operator.BelongTo]: '属于',
    [Operator.NotBelongTo]: '不属于',
    [Operator.StartWith]: '以...开头',
    [Operator.EndWith]: '以...结尾',
    [Operator.DynamicDates]: '动态日期',
    [Operator.NDayRegiste]: 'N天前注册',
}

/**
 * 定义不同 DataType 支持的操作符列表
 */
export const DataTypeOperators: Record<DataType, Operator[]> = {
    [DataType.String]: [
        Operator.EqualTo, Operator.NotEqualTo,
        Operator.In, Operator.NotIn,
        Operator.Like, Operator.NotLike,
        Operator.StartWith, Operator.EndWith,
        Operator.IsNull, Operator.IsNotNull
    ],
    [DataType.Number]: [
        Operator.EqualTo, Operator.NotEqualTo,
        Operator.GreaterThan, Operator.GreaterOrEqual,
        Operator.LessThan, Operator.LessOrEqual,
        Operator.Between, Operator.In, Operator.NotIn,
        Operator.IsNull, Operator.IsNotNull
    ],
    [DataType.DateTime]: [
        Operator.Between, Operator.EqualTo, Operator.NotEqualTo,
        Operator.GreaterThan, Operator.GreaterOrEqual,
        Operator.LessThan, Operator.LessOrEqual,
        Operator.IsNull, Operator.IsNotNull
    ],
    [DataType.Boolean]: [
        Operator.EqualTo, Operator.NotEqualTo
    ]
}

/**
 * 基础列定义
 */
export interface Column {
    table: string
    field: string
    alias: string
}

/**
 * 排序定义
 */
export interface Order {
    field: string
    desc: boolean
}
