/**
 * 分析模块业务常量定义
 */

/**
 * 表名常量
 */
export const TableNames = {
  EVENT: 'event_data',
  USER: 'user_data',
  TAG: 'user_tag_data',
  USER_GROUP: 'user_group_data',
  USER_PROPERTY: 'ups',
} as const

/**
 * 指标计算公式
 */
export enum FormulaTypes {
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
  LastValue = 16,
  FirstValue = 17,
  BitmapUnion = 18,
  CountDistinctDailyUserId = 19,
}

/**
 * 过滤操作符
 */
export enum OperatorTypes {
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
  DynamicDates = 19,
  NDayRegiste = 20,
}

/**
 * 可视化模式
 */
export enum VizModes {
  Line = 'line',
  Bar = 'bar',
  Pie = 'pie',
  Table = 'table',
}

/**
 * 自定义公式 Token 类型
 */
export enum FormulaTokenTypes {
  Metric = 'metric',
  Operator = 'operator',
  Number = 'number',
}

/**
 * 自定义公式默认值
 */
export const FormulaDefaults = {
  METRIC_PREFIX: 'm', // 默认指标别名前缀，如 m0, m1
  UNKNOWN_FORMULA: '未命名公式',
  TOTAL_TIMES_FIELD: '__TOTAL_TIMES__',
  TOTAL_USERS_FIELD: '__TOTAL_USERS__',
  DEFAULT_FORMAT: 'int', // 默认格式化为整数
} as const
