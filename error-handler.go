package zylog

import (
	"runtime/debug"
	"fmt"
	"reflect"
	"sync"
)

type FatalError struct {
	Type  string
	Err   error
	Stack string
}

func (f FatalError) Error() string {
	return f.Type + ":" + f.Err.Error() + "\n" + f.Stack
}

func (f FatalError) String() string {
	return fmt.Sprintf("error - %s: %s \n%s", f.Type, f.Err.Error(), f.Stack)
}

func NewError(i interface{}) FatalError {
	if i == nil {
		return FatalError{
			Type:  "NonePointer",
			Err:   fmt.Errorf("error is nil"),
			Stack: string(debug.Stack()),
		}
	}
	var t string
	var e error
	n := reflect.TypeOf(i).Name()
	switch n {
	case "FatalError":
		return i.(FatalError)
	case "":
		t = "error"
		e = i.(error)
	default:
		t = n
		e = fmt.Errorf(fmt.Sprint(i))
	}
	if t == "" {
		t = "error"
	}
	if e.Error() == "" {
		e = fmt.Errorf("runtime error")
	}
	return FatalError{
		Type:  t,
		Err:   e,
		Stack: string(debug.Stack()),
	}
}

func IsFatalError(err interface{}) (FatalError, bool) {
	if err == nil {
		return *new(FatalError), false
	}
	switch err.(type) {
	case FatalError:
		return err.(FatalError), true
	default:
		return *new(FatalError), false
	}
}

func Catch(catch func(error)) {
	if err := recover(); err != nil {
		if catch != nil {
			catch(NewError(err))
		}
	}
}

func CatchB(before func(), catch func(error)) {
	if before != nil {
		before()
	}
	if err := recover(); err != nil {
		if catch != nil {
			catch(NewError(err))
		}
	}
}

func GoCatchAndThrow(err error) {
	if err != nil {
		panic(NewError(err))
	}
}

func GoCatchAndThrowB(before func(), err error) {
	if before != nil {
		before()
	}
	panic(NewError(err))
}

func CatchAndThrow() {
	if err := recover(); err != nil {
		panic(NewError(err))
	}
}

func CatchAndThrowAfterUnLock(mutex sync.Locker) {
	mutex.Unlock()
	if err := recover(); err != nil {
		panic(NewError(err))
	}
}

func CatchAndThrowAfterRUnLock(mutex sync.RWMutex) {
	mutex.RUnlock()
	if err := recover(); err != nil {
		panic(NewError(err))
	}
}

func CatchAndThrowA(after func()) {
	if err := recover(); err != nil {
		if after != nil {
			after()
		}
		panic(NewError(err))
	}
}

func CatchAndThrowB(before func()) {
	if before != nil {
		before()
	}
	if err := recover(); err != nil {
		panic(NewError(err))
	}
}

func CatchAndThrowAB(before func(), after func()) {
	if before != nil {
		before()
	}
	if err := recover(); err != nil {
		if after != nil {
			after()
		}
		panic(NewError(err))
	}
}
