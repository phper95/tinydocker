package main

import "time"

/*
计算给定日期是星期几。
参数：
date_string（str）：日期字符串，格式为 'YYYY - MM - DD'。
返回：
int：表示星期几的整数，其中 0 代表星期一，1 代表星期二，依此类推，6 代表星期日。
*/
func GetWeekday(date_string string) int {
	// 解析日期字符串
	layout := "2006-01-02" // Go语言的时间格式化布局
	date, err := time.Parse(layout, date_string)
	if err != nil {
		return -1 // 返回 -1 表示解析错误
	}

	// 获取星期几
	weekday := date.Weekday()

	// 将星期几转换为整数，0 代表星期日，6 代表星期六
	return int(weekday)
}
