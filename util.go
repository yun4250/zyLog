package zylog

import (
	"os"
	"path/filepath"
	"time"
	"fmt"
)

const (
	NormalFormat = "2006-01-02.15-04-05"
	MinuteFormat = "2006-01-02.15-04"
	HourFormat   = "2006-01-02.15"
	DayFormat    = "2006-01-02"
	DayFormatO   = "02"
	MonthFormat  = "2006-01"
	YearFormat   = "2006"
)

func DurationToFormat(duration time.Duration) string {
	if duration >= 0 && duration < time.Minute {
		return NormalFormat
	} else if duration >= time.Minute && duration < time.Hour {
		return MinuteFormat
	} else if duration >= time.Hour && duration < time.Hour*24 {
		return HourFormat
	} else if duration >= time.Hour*24 {
		return DayFormat
	} else {
		return DurationToFormat(duration * -1)
	}
}

func Parse(s string) (t time.Time) {
	var e error
	if len(s) == len(NormalFormat) {
		t, e = time.ParseInLocation(NormalFormat, s, time.Local)
	} else if len(s) == len(MinuteFormat) {
		t, e = time.ParseInLocation(MinuteFormat, s, time.Local)
	} else if len(s) == len(DayFormat) {
		t, e = time.ParseInLocation(DayFormat, s, time.Local)
	} else if len(s) == len(HourFormat) {
		t, e = time.ParseInLocation(HourFormat, s, time.Local)
	} else {
		t = time.Time{}
	}
	if e != nil {
		Print("parse %s fail: %s", s, e.Error())
	}
	return
}

func GetDirectory() string {
	absPath, _ := filepath.Abs("?")
	return string([]rune(absPath)[0:len(absPath)-1])
}

func OpenOrCreate(path string) *os.File {
	absPath, _ := filepath.Abs(path)
	if f, e := os.OpenFile(absPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777); e == nil {
		PrintDebug("OpenFile success: %s\n", absPath)
		return f
	} else {
		Print("OpenFile failed: %s\n", e.Error())
		return nil
	}
}

/*
截取字符串 start 起点下标 length 需要截取的长度(-1 到末尾)
 */
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	if length == 0 || start > rl {
		return ""
	}
	end := start + length
	if length < 0 || end > rl {
		end = rl
	}
	if start < 0 {
		start = 0
	}
	return string(rs[start:end])
}

func IsExists(path string) bool {
	_, e := os.Stat(path)
	if os.IsNotExist(e) {
		return false
	} else if e == nil {
		return true
	} else {
		fmt.Println(e)
		return true
	}
}
