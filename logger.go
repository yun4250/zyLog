package zylog

import (
	"fmt"
	"strings"
	"log"
)

func (cl *ChildLogger) Prefix() string {
	return cl.prefix
}

func (cl *ChildLogger) Manager() *ZyLogger {
	return cl.manager
}

func (cl *ChildLogger) Position(position string) *ChildLogger {
	position = strings.ToUpper(position[0:1]) + strings.ToLower(position[1:])
	ano := &ChildLogger{
		manager:  cl.manager,
		prefix:   cl.prefix,
		position: position,
		id:       cl.id,
		loggers:  cl.loggers,
	}
	return ano
}

func (cl *ChildLogger) positionStr() string {
	if cl.position == "" {
		return ""
	}
	return cl.position + ": "
}

func (cl *ChildLogger) get(i Level) *log.Logger {
	li := cl.manager.getLogger(i)
	if v := li.getLogger(cl.id); v != nil {
		return v
	} else {
		return li.newLogger(cl.id, cl.manager.prefixInfo.prefixes[cl.id]+" ")
	}
}

//panic an error of zylog.FatalError,
//you can use Catch or CatchAndThrow to handle it
func (l *ChildLogger) Fatalf(format string, v ...interface{}) {
	if l.manager.Level >= Fatal {
		err := NewError(fmt.Errorf(format, v...))
		l.get(Fatal).Fatalf("FATAL " + l.positionStr() + strings.Replace(err.Error(), "\n", "\n\t", -1))
		panic(err)
	}
}

//panic an error of zylog.FatalError,
//you can use Catch or CatchAndThrow to handle it
func (l *ChildLogger) Fatal(v ...interface{}) {
	if l.manager.Level >= Fatal {

		err := NewError(fmt.Sprint(v...))
		l.get(Fatal).Println("FATAL " + l.positionStr() + strings.Replace(err.Error(), "\n", "\n\t", -1))
		panic(err)
	}
}

func (l *ChildLogger) Criticalf(format string, v ...interface{}) {
	if l.manager.Level >= Critical {
		body := fmt.Sprintf(format, v...)
		l.get(Critical).Printf("CRITICAL " + l.positionStr() + strings.Replace(body, "\n", "\n\t", -1))
	}
}
func (l *ChildLogger) Critical(v ...interface{}) {
	if l.manager.Level >= Critical {
		body := fmt.Sprint(v...)
		l.get(Critical).Println("CRITICAL " + l.positionStr() + strings.Replace(body, "\n", "\n\t", -1))
	}
}

func (l *ChildLogger) Errorf(format string, v ...interface{}) {
	if l.manager.Level >= Error {
		body := fmt.Sprintf(format, v...)
		l.get(Error).Printf("ERROR " + l.positionStr() + strings.Replace(body, "\n", "\n\t", -1))
	}
}
func (l *ChildLogger) Error(v ...interface{}) {
	if l.manager.Level >= Error {
		body := fmt.Sprint(v...)
		l.get(Error).Println("ERROR " + l.positionStr() + strings.Replace(body, "\n", "\n\t", -1))
	}
}

func (l *ChildLogger) Warnf(format string, v ...interface{}) {
	if l.manager.Level >= Warn {
		l.get(Warn).Printf("WARN "+l.positionStr()+format, v...)
	}
}
func (l *ChildLogger) Warn(v ...interface{}) {
	if l.manager.Level >= Warn {
		l.get(Warn).Println("WARN " + l.positionStr() + fmt.Sprint(v...))
	}
}

func (l *ChildLogger) Infof(format string, v ...interface{}) {
	if l.manager.Level >= Info {
		l.get(Info).Printf("INFO "+l.positionStr()+format, v...)
	}
}
func (l *ChildLogger) Info(v ...interface{}) {
	if l.manager.Level >= Info {
		l.get(Info).Println("INFO " + l.positionStr() + fmt.Sprint(v...))
	}
}

func (l *ChildLogger) Debugf(format string, v ...interface{}) {
	if l.manager.Level >= Debug {
		l.get(Debug).Printf("DEBUG "+l.positionStr()+format, v...)
	}
}
func (l *ChildLogger) Debug(v ...interface{}) {
	if l.manager.Level >= Debug {
		l.get(Debug).Println("DEBUG " + l.positionStr() + fmt.Sprint(v...))
	}
}

func (l *ChildLogger) Tracef(format string, v ...interface{}) {
	if l.manager.Level >= Trace {
		l.get(Trace).Printf("TRACE "+l.positionStr()+format, v...)
	}
}
func (l *ChildLogger) Trace(v ...interface{}) {
	if l.manager.Level >= Trace {
		l.get(Trace).Println("TRACE " + l.positionStr() + fmt.Sprint(v...))
	}
}
