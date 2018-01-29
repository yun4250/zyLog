package zylog

import (
	"time"
	"os"
	"path/filepath"
	"strings"
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

func NewManagerInDir(fileName string, directory string) *ZyLogger {
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

func (zy *ZyLogger) GetChild(prefix string, position string) *ChildLogger {
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

func (zy *ZyLogger) GetChildWithPid(prefix string, pid int, position string) *ChildLogger {
	prefix = fmt.Sprintf("%s[%d]", prefix, pid)
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

//write lock
func (l *ZyLogger) initLevelInfo() {
	l.Lock()
	defer l.Unlock()
	if l.Directory == "" {
		l.Directory = GetDirectory()
	} else {
		absPath, _ := filepath.Abs(l.Directory)
		l.Directory = absPath
	}
	if e := os.MkdirAll(l.Directory, 0777); e != nil {
		Panic("Make dir %s failed", l.Directory)
	}

	li := make([]*levelInfo, 7)
	independent := l.LevelStrategy >= ErrorIsolation && l.Level >= Error
	li[6] = initLi(Error, l.Directory, l.FileName, independent)
	independent = l.LevelStrategy >= WarnIsolation && l.Level >= Warn
	li[5] = initLi(Warn, l.Directory, l.FileName, independent)
	independent = l.LevelStrategy >= InfoIsolation && l.Level >= Info
	li[4] = initLi(Info, l.Directory, l.FileName, independent)
	independent = l.LevelStrategy >= DebugIsolation && l.Level >= Debug
	li[3] = initLi(Debug, l.Directory, l.FileName, independent)
	independent = l.LevelStrategy >= TraceIsolation && l.Level == Trace
	li[2] = initLi(Trace, l.Directory, l.FileName, independent)
	independent = l.LevelStrategy >= EachIsolation && l.Level >= Critical
	li[1] = initLi(Critical, l.Directory, l.FileName, independent)
	independent = l.LevelStrategy >= NoneIsolation && l.Level >= Fatal
	li[0] = initLi(Fatal, l.Directory, l.FileName, independent)
	l.levelInfo = li
	PrintDebug("levelInfo init")
}

//write lock than read lock
func (l *ZyLogger) getLogger(level Level) *levelInfo {
	l.RLock()
	defer l.RUnlock()
	if l.levelInfo == nil {
		l.RUnlock()
		l.initLevelInfo()
		l.RLock()
	}
	go l.checkAndRotate()
	return l.levelInfo[LevelToIndex(level)]
}

//write lock
func (l *ZyLogger) checkAndRotate() {
	l.RLock()
	defer l.RUnlock()
	if l.MaxKeepDuration > 0 && l.Duration > 0 {
		if l.fileTime.IsZero() {
			l.RUnlock()
			l.Lock()
			l.fileTime = time.Now().Add(l.Duration)
			l.Unlock()
			l.RLock()
		}
		if time.Now().After(l.fileTime) {
			go func() {
				l.Lock()
				defer l.Unlock()
				for _, li := range l.levelInfo {
					bakName := li.filePath + "." + l.fileTime.Format(DurationToFormat(l.Duration))
					li.rotate(bakName)
				}
				l.fileTime = l.fileTime.Add(l.Duration)
				go l.ClearOldFile()
			}()
		}
	}
}

//no lock
func (l *ZyLogger) ClearOldFile() {
	prefix := l.FileName + ".log."
	keep := time.Now().Add(l.Duration * -1 * time.Duration(l.MaxKeepDuration))
	filepath.Walk(l.Directory, func(filename string, fi os.FileInfo, err error) error {
		if i := strings.Index(filename, prefix); i != -1 {
			str := Substr(filename, i+len(prefix), -1)
			if t := Parse(str); !t.IsZero() && t.Before(keep) {
				os.Remove(filename)
				PrintDebug("Clear %s", filename)
			}
		}
		return err
	})
}

func (zy *ZyLogger) OpenPrefixFormat() *ZyLogger {
	zy.Lock()
	defer zy.Unlock()
	if zy.prefixInfo == nil {
		zy.initPi()
	}
	zy.prefixInfo.Open()
	return zy
}

func (zy *ZyLogger) ClosePrefixFormat() *ZyLogger {
	zy.Lock()
	defer zy.Unlock()
	if zy.prefixInfo == nil {
		zy.initPi()
	}
	zy.prefixInfo.Close()
	return zy
}

func (zy *ZyLogger) AddPrefix(prefix string) int {
	if zy.prefixInfo == nil {
		zy.initPi()
	}
	return zy.prefixInfo.addPrefix(prefix)
}

func initLi(level Level, dir string, fileName string, independent bool) *levelInfo {
	li := &levelInfo{
		level: level,
	}
	li.Lock()
	defer li.Unlock()
	if independent && li.level != Fatal {
		li.filePath = filepath.Clean(dir + "/" + fileName + "-" + li.levelStr() + ".log")
	} else {
		li.filePath = filepath.Clean(dir + "/" + fileName + ".log")
	}
	return li
}

func (li *levelInfo) levelStr() string {
	switch li.level {
	case Fatal:
		return "fatal"
	case Critical:
		return "critical"
	case Error:
		return "error"
	case Warn:
		return "warn"
	case Info:
		return "info"
	case Debug:
		return "debug"
	case Trace:
		return "trace"
	default:
		return "info"
	}
}

func (li *levelInfo) openOrCreateFile() {
	if f, e := os.OpenFile(li.filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777); e == nil {
		li.file = f
		PrintDebug("Open file %s", li.filePath)
	} else {
		Print("Open and create file %s failed", li.filePath)
	}
}

func (li *levelInfo) getLogger(id int) *log.Logger {
	li.RLock()
	defer li.RUnlock()
	if li.file == nil && li.logger == nil {
		li.RUnlock()
		li.Lock()
		li.openOrCreateFile()
		li.logger = make(map[int]*log.Logger)
		li.Unlock()
		li.RLock()
	}
	if v, b := li.logger[id]; b {
		return v
	} else {
		return nil
	}
}

func (li *levelInfo) newLogger(id int, prefix string) *log.Logger {
	li.Lock()
	defer li.Unlock()
	if _, b := li.logger[id]; !b {
		li.logger[id] = log.New(li.file, prefix, log.LstdFlags)
	}
	return li.logger[id]
}

func (li *levelInfo) rotate(bakName string) {
	li.Lock()
	defer li.Unlock()
	if IsExists(li.filePath) && li.file != nil {
		li.file.Close()
		if e := os.Rename(li.filePath, bakName); e != nil {
			Print("Rotate failed :%s\n", e)
		} else {
			PrintDebug("Rotate file %s to %s", li.filePath, bakName)
		}
		li.openOrCreateFile()
		for _, loger := range li.logger {
			loger.SetOutput(li.file)
		}
	}
}

func (l *ZyLogger) initPi() {
	PrintDebug("ZyLogger.prefixInfo init")
	pi := &prefixInfo{
		nextRegisterId: 0,
		prefixes:       make(map[int]string),
		fallPrefixes:   make(map[string]int),
		prefixLen:      0,
	}
	l.prefixInfo = pi
}

func (pi *prefixInfo) Open() {
	pi.Lock()
	defer pi.Unlock()
	if pi.prefixLen == 0 {
		pi.prefixLen = 1
		for _, p := range pi.prefixes {
			if len := len(p); len > pi.prefixLen {
				pi.prefixLen = len
			}
		}
		pi.formatAll()
		PrintDebug("Open formatting prefix, now length is %d", pi.prefixLen)
	}
}

func (pi *prefixInfo) Close() {
	pi.Lock()
	defer pi.Unlock()
	clearNum := 0
	if pi.prefixLen != 0 {
		for k, v := range pi.prefixes {
			pi.prefixes[k] = strings.TrimRight(v, " ")
		}
		pi.prefixLen = 0
		PrintDebug("Close formatting prefix, reformat %d prefixes", clearNum)
	}
}

func (pi *prefixInfo) addPrefix(prefix string) (id int) {
	if k, v := pi.fallPrefixes[prefix]; v {
		return k
	}
	id = pi.nextRegisterId
	pi.nextRegisterId += 1
	pi.fallPrefixes[prefix] = id

	if pi.prefixLen == 0 {
		pi.prefixes[id] = prefix
	} else if le := len(prefix); pi.prefixLen < le {
		pi.prefixLen = le
		pi.prefixes[id] = prefix
		pi.formatAll()
	} else {
		pi.prefixes[id] = pi.format(prefix)
	}
	PrintDebug("Add prefix [%s], now len is %d", pi.prefixes[id], pi.prefixLen)
	return id
}

func (pi *prefixInfo) formatAll() {
	for k, v := range pi.prefixes {
		pi.prefixes[k] = pi.format(v)
	}
}

func (pi *prefixInfo) format(prefix string) string {
	s := prefix
	for i := 0; i < pi.prefixLen-len(prefix); i++ {
		s = s + " "
	}
	return s
}
