package doris

import (
	"fmt"
	"strings"
	"time"
)

const (
	Tg_Interval_Empty = iota + 1
	Tg_Interval_Day
	Tg_Interval_Week
	Tg_Interval_Month
	Tg_Interval_Quarter
	Tg_Interval_Year
	Tg_Interval_Hour
	Tg_Interval_Minute
)

type TimeGrain struct {
	Column    Column `json:"column"`
	Interval  int    `json:"interval"`
	WindowNum int    `json:"window_num"`
}

func BuildTimeGrainV2(tg TimeGrain) (string, string) {
	if tg.Interval == Tg_Interval_Empty {
		return fmt.Sprintf("'%s' as 日期", tg.Column.Alias), "日期"
	}
	tgField := tg.Column.Field
	if tg.Column.Table == EVENT_TABLE || tg.Column.Table == "" {
		tgField = fmt.Sprintf("%s.%s", EVENT_TABLE_ALIAS, tg.Column.Field)
	} else if tg.Column.Table == USER_TABLE {
		tgField = fmt.Sprintf("%s.%s", USER_TABLE_ALIAS, tg.Column.Field)
	}

	formula := fmt.Sprintf("DATE(%s) AS %s", tgField, tg.Column.Alias)
	switch tg.Interval {
	case Tg_Interval_Minute:
		dateFormat := "%Y-%m-%d %H:%i"
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s') AS %s", tgField, dateFormat, tg.Column.Alias)
	case Tg_Interval_Hour:
		dateFormat := "%Y-%m-%d %H:00"
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s') AS %s", tgField, dateFormat, tg.Column.Alias)
	case Tg_Interval_Day:
		dateFormat := "%Y-%m-%d"
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s') AS %s", tgField, dateFormat, tg.Column.Alias)
	case Tg_Interval_Week:
		dateFormat := "%Y-%m-%d"
		formula = fmt.Sprintf("DATE_FORMAT(TO_MONDAY(%s), '%s') AS %s", tgField, dateFormat, tg.Column.Alias)
	case Tg_Interval_Month:
		dateFormat := "%Y-%m"
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s') AS %s", tgField, dateFormat, tg.Column.Alias)
	case Tg_Interval_Quarter:
		formula = fmt.Sprintf("CONCAT(YEAR(%s),' ', QUARTER(%s),'季') AS %s", tgField, tgField, tg.Column.Alias)
	case Tg_Interval_Year:
		formula = fmt.Sprintf("YEAR(%s) AS %s", tgField, tg.Column.Alias)
	}
	return formula, tg.Column.Alias
}

func BuildTimeGrainFormula(tg TimeGrain) string {
	if tg.Interval == Tg_Interval_Empty {
		return fmt.Sprintf("'%s'", tg.Column.Alias)
	}
	tgField := tg.Column.Field

	if tg.Column.Table == EVENT_TABLE || tg.Column.Table == "" {
		tgField = fmt.Sprintf("%s.%s", EVENT_TABLE_ALIAS, tg.Column.Field)
	} else if tg.Column.Table == USER_TABLE {
		tgField = fmt.Sprintf("%s.%s", USER_TABLE_ALIAS, tg.Column.Field)
	}
	if strings.Contains(tgField, "convert_tz") {
		tgField = tg.Column.Field
	}

	formula := fmt.Sprintf("DATE(%s)", tgField)
	switch tg.Interval {
	case Tg_Interval_Minute:
		dateFormat := "%Y-%m-%d %H:%i"
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s')", tgField, dateFormat)
	case Tg_Interval_Hour:
		dateFormat := "%Y-%m-%d %H:00"
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s')", tgField, dateFormat)
	case Tg_Interval_Day:
		dateFormat := "%Y-%m-%d"
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s')", tgField, dateFormat)
	case Tg_Interval_Week:
		dateFormat := "%Y-%m-%d"
		formula = fmt.Sprintf("DATE_FORMAT(TO_MONDAY(%s), '%s')", tgField, dateFormat)
	case Tg_Interval_Month:
		dateFormat := "%Y-%m"
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s')", tgField, dateFormat)
	case Tg_Interval_Quarter:
		formula = fmt.Sprintf("CONCAT(YEAR(%s),' ', QUARTER(%s),'季') ", tgField, tgField)
	case Tg_Interval_Year:
		formula = fmt.Sprintf("YEAR(%s)", tgField)
	}
	return formula
}

func BuildTimeGrain(tg TimeGrain) (string, string) {
	if tg.Interval == Tg_Interval_Empty {
		return fmt.Sprintf("'%s'", tg.Column.Alias), "日期"
	}
	formula := fmt.Sprintf("DATE(%s)", tg.Column.Field)
	switch tg.Interval {
	case Tg_Interval_Minute:
		dateFormat := "%Y-%m-%d %H:%i"
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s')", tg.Column.Field, dateFormat)
	case Tg_Interval_Hour:
		dateFormat := "%Y-%m-%d %H:00"
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s')", tg.Column.Field, dateFormat)
	case Tg_Interval_Day:
		caseConvert := fmt.Sprintf(`CASE date_format(%s, '%%a')
        WHEN 'Mon' THEN '一'
        WHEN 'Tue' THEN '二'
        WHEN 'Wed' THEN '三'
        WHEN 'Thu' THEN '四'
        WHEN 'Fri' THEN '五'
        WHEN 'Sat' THEN '六'
        WHEN 'Sun' THEN '日'
      END`, tg.Column.Field)
		formula = fmt.Sprintf("CONCAT(DATE(%s),' ( ', %s, ' ) ')", tg.Column.Field, caseConvert)
	case Tg_Interval_Week:
		dateFormat := "%Y-%m-%d 当周"
		dateField := fmt.Sprintf("DATE_SUB(%s, WEEKDAY(%s))", tg.Column.Field, tg.Column.Field)
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s')", dateField, dateFormat)
	case Tg_Interval_Month:
		dateFormat := "%Y-%m"
		formula = fmt.Sprintf("DATE_FORMAT(%s,'%s')", tg.Column.Field, dateFormat)
	case Tg_Interval_Quarter:
		formula = fmt.Sprintf("CONCAT(YEAR(%s),' ', QUARTER(%s),'季')", tg.Column.Field, tg.Column.Field)
	case Tg_Interval_Year:
		formula = fmt.Sprintf("YEAR(%s)", tg.Column.Field)
	}
	return formula, tg.Column.Alias
}

// 时间窗口期构建
func BuildWindowTimeGrain(tg TimeGrain) (string, string) {
	if tg.WindowNum == 0 {
		tg.WindowNum = 1
	}

	var formula string
	tgField := fmt.Sprintf("%s.%s", EVENT_TABLE_ALIAS, tg.Column.Field)
	switch interval := tg.Interval; interval {
	case Tg_Interval_Minute:
		formula = fmt.Sprintf(`MINUTES_ADD("1970-01-01", (FLOOR(MINUTES_DIFF(%s, "1970-01-01") / %d) * %d))`, tgField, tg.WindowNum, tg.WindowNum)
	case Tg_Interval_Hour:
		formula = fmt.Sprintf(`HOURS_ADD("1970-01-01", (FLOOR(HOURS_DIFF(%s, "1970-01-01") / %d) * %d))`, tgField, tg.WindowNum, tg.WindowNum)
	case Tg_Interval_Day:
		formula = fmt.Sprintf(`DATE_ADD("1970-01-01", (FLOOR(DATEDIFF(%s, "1970-01-01") / %d) * %d))`, tgField, tg.WindowNum, tg.WindowNum)
	case Tg_Interval_Week:
		formula = fmt.Sprintf(`WEEKS_ADD("1970-01-01", (FLOOR(WEEKS_DIFF(%s, "1970-01-01") / %d) * %d))`, tgField, tg.WindowNum, tg.WindowNum)
	case Tg_Interval_Month:
		formula = fmt.Sprintf(`MONTHS_ADD("1970-01-01", (FLOOR(MONTHS_DIFF(%s, "1970-01-01") / %d) * %d))`, tgField, tg.WindowNum, tg.WindowNum)
	case Tg_Interval_Year:
		formula = fmt.Sprintf(`YEARS_ADD("1970-01-01", (FLOOR(YEARS_DIFF(%s, "1970-01-01") / %d) * %d))`, tgField, tg.WindowNum, tg.WindowNum)
	}
	return formula, tg.Column.Alias
}

func GetTimeZone(timeZone string) string {
	tz := "Asia/Shanghai"
	switch timeZone {
	case "+00:00":
		tz = "+00:00"
	case "+03:00":
		tz = "+03:00"
	case "-05:00":
		tz = "-05:00"
	case "-08:00":
		tz = "-08:00"
	case "+08:00", "":
		tz = "+08:00"
	}
	return tz
}

// 各时区统计添加时间 - 解决时间选择时覆盖当前日期
func AddDateForTzTime(t time.Time, tz string) time.Time {
	switch tz {
	case "+00:00":
		t = t.Add(16 * time.Hour)
	case "-08:00":
		t = t.Add(32 * time.Hour)
	case "-05:00":
		t = t.Add(26 * time.Hour)
	case "+03:00":
		t = t.Add(10 * time.Hour)
	}
	return t
}
