package doris

import (
	"fmt"
	"strconv"
	"strings"
)

type Filter struct {
	Column   Column      `json:"column"`
	Value    interface{} `json:"value"`
	Operator int         `json:"operator"`
}

func BuildFilter(filter Filter) string {
	switch filter.Operator {
	case OperGreaterOrEqual:
		return fmt.Sprintf("%s >= '%v'", filter.Column.Field, filter.Value)

	case OperGreaterThan:
		return fmt.Sprintf("%s > '%v'", filter.Column.Field, filter.Value)

	case OperLessOrEqual:
		return fmt.Sprintf("%s <= '%v'", filter.Column.Field, filter.Value)
	case OperLessThan:
		return fmt.Sprintf("%s < '%v'", filter.Column.Field, filter.Value)
	case OperEqualTo:
		return fmt.Sprintf("%s = '%v'", filter.Column.Field, filter.Value)
	case OperIn:
		return fmt.Sprintf("%s IN (%v)", filter.Column.Field, formatFilterValue(filter.Value))
	case OperNotEqualTo, OperNotIn:
		return fmt.Sprintf("%s NOT IN (%v)", filter.Column.Field, formatFilterValue(filter.Value))
	case OperIsNotNull:
		return fmt.Sprintf("(%s !='' OR %s IS NOT NULL)", filter.Column.Field, filter.Column.Field)
	case OperIsNull:
		return fmt.Sprintf("(%s ='' OR %s IS NULL)", filter.Column.Field, filter.Column.Field)
	case OperBetween:
		if arrValue, ok := filter.Value.([]interface{}); ok {
			return fmt.Sprintf("%s >= %v AND %s < %v", filter.Column.Field, formatFilterValue(arrValue[0]), filter.Column.Field, formatFilterValue(arrValue[1]))
		}
	case OperLike:
		value := fmt.Sprintf("%v", filter.Value)
		if !strings.Contains(value, "%") {
			value = "%" + value + "%"
		}
		return fmt.Sprintf("%s LIKE '%v'", filter.Column.Field, value)
	case OperStartWith:
		value := fmt.Sprintf("%v", filter.Value)
		value = value + "%"
		return fmt.Sprintf("%s LIKE '%v'", filter.Column.Field, value)
	case OperEndWith:
		value := fmt.Sprintf("%v", filter.Value)
		value = "%" + value
		return fmt.Sprintf("%s LIKE '%v'", filter.Column.Field, value)
	case OperNotLike:
		value := fmt.Sprintf("%v", filter.Value)
		if !strings.Contains(value, "%") {
			value = "%" + value + "%"
		}
		return fmt.Sprintf("%s NOT LIKE '%v'", filter.Column.Field, value)
	case OperDateDiff:
		if value, ok := filter.Value.(map[string]interface{}); ok {
			opt := fmt.Sprintf("%v", value["type"])
			if opt == "1" { //当天
				return fmt.Sprintf("date(events.e_event_time) = date(%s)", filter.Column.Field)
			} else if opt == "2" { // 间隔区间天数
				arrayValue := value["values"].([]interface{})
				timeField := "e_event_time"
				return fmt.Sprintf("(DATEDIFF(%s, %s) BETWEEN %v AND %v)", filter.Column.Field, timeField, arrayValue[0], arrayValue[1])
			}
		}
	case OperDynamicDates:
		if value, ok := filter.Value.([]string); ok {
			oper := fmt.Sprintf("%v", value[0])
			switch oper {
			case "0": // 过去 n 天
				days, _ := strconv.Atoi(fmt.Sprintf("%v", value[1]))
				return fmt.Sprintf("%s > DATE_SUB(CURRENT_DATE(), INTERVAL %v DAY) AND %s < CURRENT_DATE()",
					filter.Column.Field, days, filter.Column.Field)

			case "1": // 最近 n 天
				days, _ := strconv.Atoi(fmt.Sprintf("%v", value[1]))
				return fmt.Sprintf("%s > DATE_SUB(CURRENT_DATE(), INTERVAL %v DAY)", filter.Column.Field, days)
			case "2": // 自某日至今
				return fmt.Sprintf("%s > '%v')", filter.Column.Field, value[1])
			}
		}
	case OperNDayRegiste:
		return fmt.Sprintf("DATEDIFF(CURDATE(), %s) = %v", filter.Column.Field, filter.Value)
	}
	return "true"
}

// 格式化字符串查询值
func formatFilterValue(value interface{}) string {
	if strings.Contains(fmt.Sprintf("%v", value), "convert_tz") {
		return fmt.Sprintf("%v", value)
	}

	switch value := value.(type) {
	case float64, float32, int64, int32, int16, int8:
		return fmt.Sprintf("%v", value)
	case []interface{}:
		arrayValues := make([]string, 0)
		for _, v := range value {
			arrayValues = append(arrayValues, fmt.Sprintf("\"%v\"", v))
		}
		return strings.Join(arrayValues, ",")
	case []string:
		arrayValues := make([]string, 0)
		for _, v := range value {
			arrayValues = append(arrayValues, fmt.Sprintf("\"%v\"", v))
		}
		return strings.Join(arrayValues, ",")
	}

	strValue := fmt.Sprintf("%v", value)
	if strings.Contains(strValue, ",") {
		arrayValues := strings.Split(strValue, ",")
		for k, v := range arrayValues {
			arrayValues[k] = fmt.Sprintf("\"%v\"", v)
		}
		return strings.Join(arrayValues, ",")
	}

	return fmt.Sprintf("\"%v\"", value)
}
