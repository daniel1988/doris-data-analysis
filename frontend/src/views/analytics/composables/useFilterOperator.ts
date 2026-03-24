import { DataType, DataTypeOperators, Operator, OperatorLabels } from '@/types/doris/common'
import { computed, unref, type Ref } from 'vue'

export function useFilterOperator(dataTypeRef: Ref<DataType | string | undefined>) {
  // 计算当前数据类型可用的操作符选项
  const availableOperators = computed(() => {
    const type = unref(dataTypeRef) as DataType || DataType.String
    const ops = DataTypeOperators[type] || DataTypeOperators[DataType.String]
    
    return ops.map(op => ({ 
      value: op, 
      label: OperatorLabels[op] || `Unknown(${op})` 
    }))
  })

  // 校验并纠正当前操作符
  const validateOperator = (currentOperator: Operator | number): Operator => {
    const allowed = availableOperators.value.map(o => o.value)
    
    // 如果当前操作符在合法列表中，直接返回
    if (allowed.includes(currentOperator as Operator)) {
      return currentOperator as Operator
    }
    
    // 否则，返回该类型支持的第一个（或默认）操作符
    // 注意：如果是 DateTime，默认是 Between，其他一般是 EqualTo
    const type = unref(dataTypeRef) as DataType
    if (type === DataType.DateTime && allowed.includes(Operator.Between)) {
      return Operator.Between
    }
    
    return allowed.length > 0 ? allowed[0] : Operator.EqualTo
  }

  return {
    availableOperators,
    validateOperator
  }
}
