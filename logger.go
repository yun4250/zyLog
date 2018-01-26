package zylog

import (
	"log"
	"fmt"
)

func NewManager(fileName string) *ZyLogger {
	zy := &ZyLogger{
		FileName:        fileName,
		Directory:       "",
		Level:           Info,
		LevelStrategy:   NoneIsolation,
		MaxKeepDuration: 0,
		Duration:        0,
	}
	zy.initPi()
	zy.initLevelInfo()
	return zy
}

func NewManager2(fileName string, directory string) *ZyLogger {
	zy := &ZyLogger{
		FileName:        fileName,
		Directory:       directory,
		Level:           Info,
		LevelStrategy:   NoneIsolation,
		MaxKeepDuration: 0,
		Duration:        0,
	}
	zy.initPi()
	zy.initLevelInfo()
	return zy
}

func (zy *ZyLogger) SetLevel(l Level) *ZyLogger {
	zy.Lock()
	defer zy.Unlock()
	zy.Level = l
	return zy
}

func (zy *ZyLogger) SetLevelStrategy(l LevelStrategy) *ZyLogger {
	zy.Lock()
	zy.LevelStrategy = l
	zy.Unlock()
	zy.initLevelInfo()
	return zy
}

func (zy *ZyLogger) GetChild(prefix string) *ChildLogger {
	zy.Lock()
	defer zy.Unlock()
	if zy.FileName == "" {
		Panic("zyLogger.FileName can not be empty")
	}
	return &ChildLogger{
		Manager: zy,
		Prefix:  prefix,
		id:      zy.AddPrefix(prefix),
	}
}

func (zy *ZyLogger) GetChildWithPid(prefix string, pid int) *ChildLogger {
	prefix = fmt.Sprintf("[%d]%s", pid, prefix)
	zy.Lock()
	defer zy.Unlock()
	if zy.FileName == "" {
		Panic("zyLogger.FileName can not be empty")
	}
	return &ChildLogger{
		Manager: zy,
		Prefix:  prefix,
		id:      zy.AddPrefix(prefix),
	}
}

func (l *ChildLogger) get(i Level) *log.Logger {
	li := l.Manager.getLogger(i)
	if v := li.getLogger(l.id); v != nil {
		return v
	} else {
		return li.newLogger(l.id, l.Manager.prefixInfo.prefixes[l.id]+" ")
	}
}
