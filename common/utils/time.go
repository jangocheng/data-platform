package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetPreTime(timeFlag string) string {
	monthDays := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	today := time.Now()
	year := today.Year()
	month := today.Month()
	day := today.Day()
	if strings.HasSuffix(timeFlag, "d") {
		dayStr := strings.TrimRight(timeFlag, "d")
		if day, err := strconv.Atoi(dayStr); err != nil {
			return ""
		} else {
			return time.Unix(today.Unix()-int64(day)*86400, 0).Format("2006-01-02 15:04:05")
		}
	} else if strings.HasSuffix(timeFlag, "m") {
		monStr := strings.TrimRight(timeFlag, "m")
		if mon, err := strconv.Atoi(monStr); err != nil {
			return ""
		} else {
			for i:=0; i<mon; i++ {
				if month - 1 <= 0 {
					month = 12
					year -= 1
				} else {
					month -= 1
				}
			}
			days := monthDays[month-1]
			if day > days {
				day = days
			}

			return time.Date(year, month, day, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
		}
	} else {
		return ""
	}
}


func TimeFormat(timeStrParam string) string {
	timeStr := strings.Trim(timeStrParam, " ")
	timeStrSplit := strings.Split(timeStr, " ")
	if len(timeStrSplit) != 2 {
		timeStrSplit = append(timeStrSplit, "00:00:00")
	}
	_date := timeStrSplit[0]
	_time := timeStrSplit[1]
	ymd := ""
	if len(_date) == 8 && !strings.Contains(_date, "-") {
		ymd = fmt.Sprintf("%s-%s-%s", _date[0:4], _date[4:6], _date[6:8])
	} else {
		ymd = strings.Replace(_date, "年", "-", -1)
		ymd = strings.Replace(ymd, "月", "-", -1)
		ymd = strings.Replace(ymd, "日", "", -1)
		ymd = strings.Replace(ymd, "/", "-", -1)
		ymdSplit := strings.Split(ymd, "-")
		if len(ymdSplit) == 3 {
			if len(ymdSplit[1]) == 1 {
				ymdSplit[1] = fmt.Sprintf("0%s", ymdSplit[1] )
			}
			if len(ymdSplit[2]) == 1 {
				ymdSplit[2] = fmt.Sprintf("0%s", ymdSplit[2] )
			}
			ymd = strings.Join(ymdSplit, "-")
		}

	}
	hms := ""
	if !strings.Contains(_time, ":") {
		hms = "00:00:00"
	} else {
		hms = _time
	}
	return fmt.Sprintf("%s %s", ymd, hms)
}

func Timestamp(timeStrParams ...string) int64 {
	if len(timeStrParams) == 0 {
		return time.Now().Unix()
	}
	timeStr := timeStrParams[0]
	timeFormat := TimeFormat(timeStr)
	timeTemplate := "2006-01-02 15:04:05"
	stamp, err := time.ParseInLocation(timeTemplate, timeFormat, time.Local)
	if err != nil {
		return time.Now().Unix()
	}
	return stamp.Unix()
}