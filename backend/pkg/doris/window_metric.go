package doris

import (
	"fmt"
	"strings"
)

// sum(field) over (partition by date(e_event_time),e_openid order by e_event_time)

const (
	Window_Func_First_Value = iota + 1
	Window_Func_Last_Value
	Window_Func_Count
	Window_Func_Sum
	Window_Func_Avg
)

type WindowMetric struct {
	Column     Column    `json:"column"`
	WindowFunc int       `json:"window_func"`
	TimeGrain  TimeGrain `json:"time_grain"`
	Groups     []Group   `json:"groups"`
}

func BuildWindowMetricFormula(wm *WindowMetric) string {
	windowSqls := make([]string, 0)
	switch wm.WindowFunc {
	case Window_Func_First_Value:
		windowSqls = append(windowSqls, fmt.Sprintf("FIRST_VALUE(%s)", wm.Column.Field))
	case Window_Func_Last_Value:
		windowSqls = append(windowSqls, fmt.Sprintf("LAST_VALUE(%s)", wm.Column.Field))
	case Window_Func_Count:
		windowSqls = append(windowSqls, fmt.Sprintf("COUNT(%s)", wm.Column.Field))
	case Window_Func_Sum:
		windowSqls = append(windowSqls, fmt.Sprintf("SUM(%s)", wm.Column.Field))
	case Window_Func_Avg:
		windowSqls = append(windowSqls, fmt.Sprintf("AVG(%s)", wm.Column.Field))
	}

	partitionSqls := make([]string, 0)
	if len(wm.Groups) > 0 {
		for _, group := range wm.Groups {
			partitionSqls = append(partitionSqls, group.Column.Field)
		}
	}
	tgFormula := BuildTimeGrainFormula(wm.TimeGrain)
	partitionSqls = append(partitionSqls, tgFormula)

	orderSqls := make([]string, 0)
	orderSqls = append(orderSqls, tgFormula)
	windowSqls = append(windowSqls, fmt.Sprintf("OVER(PARTITION BY %s ORDER BY %s)",
		strings.Join(partitionSqls, ","), strings.Join(orderSqls, ",")))

	return strings.Join(windowSqls, " ")
}
