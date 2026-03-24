package doris

import (
	"fmt"
	"strings"
)

const (
	FormatDefault string = "raw"
	FormatInt     string = "int"
	FormatDecimal string = "decimal"
	FormatPercent string = "percent"
)

type Metric struct {
	Column  Column `json:"column"`
	Formula int    `json:"formula"`
	Format  string `json:"format"`
}

func BuildMetric(metric Metric) string {
	selectFormula := ""
	selectField := metric.Column.Field
	if strings.Contains(metric.Column.Field, "CASE") {
		selectField = metric.Column.Field
	}
	switch metric.Formula {
	case Formula_Count:
		selectFormula = fmt.Sprintf(`COUNT(%s)`, selectField)
	case Formula_Count_Distinct, Formula_Count_Distinct_UserId:
		selectFormula = fmt.Sprintf("COUNT(DISTINCT %s)", selectField)
	case Formula_Sum:
		selectField = fmt.Sprintf("CAST(%s AS DECIMAL(20,6))", selectField)
		selectFormula = fmt.Sprintf("SUM(%s)", selectField)

	case Formula_Max:
		selectField = fmt.Sprintf("CAST(%s AS DOUBLE)", selectField)
		selectFormula = fmt.Sprintf("MAX(%s)", selectField)
	case Formula_Min:
		selectField = fmt.Sprintf("CAST(%s AS DOUBLE)", selectField)
		selectFormula = fmt.Sprintf("MIN(%s)", selectField)
	case Formula_Avg:
		selectField = fmt.Sprintf("CAST(%s AS DOUBLE)", selectField)
		selectFormula = fmt.Sprintf("AVG(%s)", selectField)
	case Formula_Any:
		selectFormula = fmt.Sprintf("ANY_VALUE(%s)", selectField)
	case Formula_Count_Days:
		selectFormula = fmt.Sprintf("COUNT(DISTINCT DATE(%s))", selectField)
	case Formula_Bitmap_Union:
		selectFormula = fmt.Sprintf("BITMAP_UNION(BITMAP_HASH(%s))", selectField)
		return selectFormula
	case Formula_Count_Distinct_Daily_UserId:
		selectFormula = fmt.Sprintf("COUNT(DISTINCT CONCAT(DATE(e_event_time), '#', %s))", selectField)
	default:
		selectFormula = fmt.Sprintf(`COUNT(%s)`, selectField)
	}
	// 指标格式化
	switch metric.Format {
	case FormatInt:
		selectFormula = fmt.Sprintf("ROUND(IFNULL(%s, 0))", selectFormula)
	case FormatDecimal:
		selectFormula = fmt.Sprintf("ROUND(IFNULL(%s, 0), 2)", selectFormula)
	case FormatPercent:
		selectFormula = fmt.Sprintf("CONCAT(ROUND((%s)*100, 4), '%s')", selectFormula, "%")
	default:
		// selectFormula = fmt.Sprintf("IFNULL(%s, 0)", selectFormula)
	}
	return selectFormula
}

func RoundFormat(selectFormula string, format string) string {
	// 指标格式化
	switch format {
	case "int":
		selectFormula = fmt.Sprintf("ROUND(IFNULL(%s, 0))", selectFormula)
	case "decimal":
		selectFormula = fmt.Sprintf("ROUND(IFNULL(%s, 0), 2)", selectFormula)
	case "percent":
		selectFormula = fmt.Sprintf("CONCAT(ROUND((%s)*100, 4), '%s')", selectFormula, "%")
	default:
		selectFormula = fmt.Sprintf("IFNULL(%s, 0)", selectFormula)
	}
	return selectFormula
}
