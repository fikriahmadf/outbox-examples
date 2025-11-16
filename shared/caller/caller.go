package caller

import (
	"fmt"
	"runtime"
	"strings"
)

type Opt struct {
	Skip int
}

type Option func(*Opt)

const defaultSkip = 1

// FuncName will return caller function name. Use WithSkip to skip frame.
func FuncName(options ...Option) string {
	option := &Opt{Skip: defaultSkip}
	for _, opt := range options {
		opt(option)
	}

	pc, _, _, ok := runtime.Caller(option.Skip)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "?"
	}

	funcNameWithModule := fn.Name()
	funcNameWithModuleSplit := strings.Split(funcNameWithModule, "/")
	funcName := funcNameWithModuleSplit[len(funcNameWithModuleSplit)-1]

	return funcName
}

// FileLine return caller "file.go:123". Use WithSkip to skip frame.
func FileLine(options ...Option) string {
	option := &Opt{Skip: defaultSkip}
	for _, opt := range options {
		opt(option)
	}

	_, file, line, ok := runtime.Caller(option.Skip)
	if !ok {
		return "?"
	}

	return fmt.Sprintf("%s:%d", file, line)
}

// Info return caller file line and func name.
func Info(options ...Option) (fileLine string, funcName string) {
	option := &Opt{Skip: defaultSkip}
	for _, opt := range options {
		opt(option)
	}

	pc, file, line, ok := runtime.Caller(option.Skip)
	if !ok {
		return "?", "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "?", "?"
	}

	funcNameWithModule := fn.Name()
	funcNameWithModuleSplit := strings.Split(funcNameWithModule, "/")
	funcName = funcNameWithModuleSplit[len(funcNameWithModuleSplit)-1]

	fileLine = fmt.Sprintf("%s:%d", file, line)

	return fileLine, funcName
}

func WithSkip(skip int) Option {
	return func(o *Opt) {
		o.Skip = skip + defaultSkip
	}
}
