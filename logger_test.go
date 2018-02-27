package zylog

import (
	"testing"
)

func TestLogger(t *testing.T) {
	dz = true
	Print("test")
	//PrintDebug("debug")

	//position := "loading"
	//fmt.Println(strings.ToUpper(position[0:1]) + strings.ToLower(position[1:]))

	//fe := NewError(fmt.Errorf("test"))
	//fmt.Println(reflect.TypeOf(fmt.Errorf("test")).Name())
	//fmt.Println(reflect.TypeOf(fe).Name())
	//
	//fmt.Println(NewError(fe).String())
	//fmt.Println(NewError(fmt.Errorf("error")).String())
	//fmt.Println(NewError("string").String())
	//fmt.Println(NewError(123).String())

	//fmt.Printf("%s",NewError(fmt.Errorf("test error")))
	//fmt.Printf("%s",NewError(fmt.Errorf("test error")).Error())
	//var LoggerManager = &ZyLogger{
	//	Directory:       "logs",
	//	FileName:        "default",
	//	Level:           Info,
	//	LevelStrategy:   NoneIsolation,
	//	Duration:        time.Hour * 24,
	//	MaxKeepDuration: 1,
	//}
	//fmt.Println("Manager Init")
	//var Logger = LoggerManager.GetChild("test")
	//fmt.Println("ChildLogger Init")
	//
	//Logger.Manager.MaxKeepDuration = 5
	//Logger.Manager.Duration = time.Second * 1
	//Logger.Manager.LevelStrategy = EachIsolation
	//LoggerManager.OpenPrefixFormat()
	//other := LoggerManager.GetChild("abcdefgh")
	//other2 := LoggerManager.GetChild("abc  defg")
	//other3 := LoggerManager.GetChild("a")
	//for i := 0; i < 20; i++ {
	//	Logger.Info(i)
	//	Logger.Warn(i)
	//	other.Info(i)
	//	other2.Info(i)
	//	other3.Info(i)
	//	time.Sleep(time.Second)
	//}
	//Logger.Manager.ClosePrefixFormat()
	//for i := 20; i < 25; i++ {
	//	Logger.Info(i)
	//	Logger.Warn(i)
	//	other.Info(i)
	//	other2.Info(i)
	//	other3.Info(i)
	//	time.Sleep(time.Second)
	//}
}
