package zylog

import (
	"os"
	"time"
	"sync"
	"log"
)

const (
	Fatal    = Level(-4)
	Critical = Level(-3)
	Error    = Level(-2)
	Warn     = Level(-1)
	Info     = Level(0)
	Debug    = Level(1)
	Trace    = Level(2)

	EachIsolation  = LevelStrategy(6) // seven files: trace, debug, info, warn, error, critical, *
	ErrorIsolation = LevelStrategy(5) // six files: trace, debug, info, warn, error, *
	WarnIsolation  = LevelStrategy(4) // five files: trace, debug, info, warn, *
	InfoIsolation  = LevelStrategy(3) // four files: trace, debug, info, *
	DebugIsolation = LevelStrategy(2) // three files: trace, debug, *
	TraceIsolation = LevelStrategy(1) // two files: trace, *
	NoneIsolation  = LevelStrategy(0) // one file only
)

type Level int8

type LevelStrategy uint8

type levelInfo struct {
	level    Level
	filePath string
	file     *os.File
	logger   map[int]*log.Logger
	sync.RWMutex
}

type prefixInfo struct {
	nextRegisterId int
	prefixes       map[int]string
	fallPrefixes   map[string]int
	prefixLen      int
	sync.RWMutex
}

type ChildLogger struct {
	manager  *ZyLogger
	prefix   string
	position string
	id       int
	loggers  []*log.Logger
}

type ZyLogger struct {
	Directory     string        //base directory of files
	FileName      string        //the file prefix name, must not be null
	Level         Level         //Level
	LevelStrategy LevelStrategy /* Means strategy to output.
	For example, zylog.Logger will redirect output into two files
	'$FileName.zylog' and 'FileName-Fatal.zylog' if choose
	zylog.CriticalIsolation. Fatal stores to the former, other Levels
	able to output stores in the latter*/

	MaxKeepDuration uint8         //Set 0 if always keep output
	Duration        time.Duration //Unit of MaxKeepDuration, min is one hour

	prefixInfo *prefixInfo

	children  []*log.Logger
	levelInfo []*levelInfo
	fileTime  time.Time

	sync.RWMutex
}

