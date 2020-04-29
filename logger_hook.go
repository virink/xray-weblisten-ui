

package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// LogHook to send error logs via bot.
type LogHook struct {
	Field  string
	Skip   int
	levels []logrus.Level
}

// 对caller进行递归查询, 直到找到非logrus包产生的第一个调用.
// 因为filename我获取到了上层目录名, 因此所有logrus包的调用的文件名都是 logrus/...
// 因此通过排除logrus开头的文件名, 就可以排除所有logrus包的自己的函数调用
func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// 这里其实可以获取函数名称的: fnName := runtime.FuncForPC(pc).Name()
// 但是我觉得有 文件名和行号就够定位问题, 因此忽略了caller返回的第一个值:pc
// 在标准库log里面我们可以选择记录文件的全路径或者文件名, 但是在使用过程成并发最合适的,
// 因为文件的全路径往往很长, 而文件名在多个包中往往有重复, 因此这里选择多取一层, 取到文件所在的上层目录那层.
func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}

// Fire - LogHook::Fire
func (hook LogHook) Fire(entry *logrus.Entry) error {
	if entry.Level == logrus.ErrorLevel || entry.Level == logrus.DebugLevel {
		entry.Data[hook.Field] = findCaller(hook.Skip)
	}
	return nil
}

// Levels - LogHook::Levels
func (hook LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// NewLogHook - Create a hook to be added to an instance of logger
func NewLogHook(levels ...logrus.Level) logrus.Hook {
	hook := LogHook{
		Field:  "line",
		Skip:   5,
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}
