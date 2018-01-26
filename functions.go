package zylog

import (
	"fmt"
	"strings"
)

//panic an error of zylog.FatalError,
//you can use Catch or CatchAndThrow to handle it
func (l *ChildLogger) Fatalf(format string, v ...interface{}) {
	if l.Manager.Level >= Fatal {
		err := NewError(fmt.Errorf(format, v...))
		l.get(Fatal).Fatalf("Fatal - " + strings.Replace(err.Error(), "\n", "\n\t", -1))
		panic(err)
	}
}

//panic an error of zylog.FatalError,
//you can use Catch or CatchAndThrow to handle it
func (l *ChildLogger) Fatal(v ...interface{}) {
	if l.Manager.Level >= Fatal {

		err := NewError(fmt.Sprint(v...))
		l.get(Fatal).Println("Fatal - " + strings.Replace(err.Error(), "\n", "\n\t", -1))
		panic(err)
	}
}

func (l *ChildLogger) Criticalf(format string, v ...interface{}) {
	if l.Manager.Level >= Critical {
		body := fmt.Sprintf(format, v...)
		l.get(Critical).Printf("Critical - " + strings.Replace(body, "\n", "\n\t", -1))
	}
}
func (l *ChildLogger) Critical(v ...interface{}) {
	if l.Manager.Level >= Critical {
		body := fmt.Sprint(v...)
		l.get(Critical).Println("Critical - " + strings.Replace(body, "\n", "\n\t", -1))
	}
}

func (l *ChildLogger) Errorf(format string, v ...interface{}) {
	if l.Manager.Level >= Error {
		body := fmt.Sprintf(format, v...)
		l.get(Error).Printf("Error - " + strings.Replace(body, "\n", "\n\t", -1))
	}
}
func (l *ChildLogger) Error(v ...interface{}) {
	if l.Manager.Level >= Error {
		body := fmt.Sprint(v...)
		l.get(Error).Println("Error - " + strings.Replace(body, "\n", "\n\t", -1))
	}
}

func (l *ChildLogger) Warnf(format string, v ...interface{}) {
	if l.Manager.Level >= Warn {
		l.get(Warn).Printf("Warn - "+format, v...)
	}
}
func (l *ChildLogger) Warn(v ...interface{}) {
	if l.Manager.Level >= Warn {
		l.get(Warn).Println("Warn - " + fmt.Sprint(v...))
	}
}

func (l *ChildLogger) Infof(format string, v ...interface{}) {
	if l.Manager.Level >= Info {
		l.get(Info).Printf("Info - "+format, v...)
	}
}
func (l *ChildLogger) Info(v ...interface{}) {
	if l.Manager.Level >= Info {
		l.get(Info).Println("Info - " + fmt.Sprint(v...))
	}
}

func (l *ChildLogger) Debugf(format string, v ...interface{}) {
	if l.Manager.Level >= Debug {
		l.get(Debug).Printf("Debug - "+format, v...)
	}
}
func (l *ChildLogger) Debug(v ...interface{}) {
	if l.Manager.Level >= Debug {
		l.get(Debug).Println("Debug - " + fmt.Sprint(v...))
	}
}

func (l *ChildLogger) Tracef(format string, v ...interface{}) {
	if l.Manager.Level >= Trace {
		l.get(Trace).Printf("Trace - "+format, v...)
	}
}
func (l *ChildLogger) Trace(v ...interface{}) {
	if l.Manager.Level >= Trace {
		l.get(Trace).Println("Trace - " + fmt.Sprint(v...))
	}
}
