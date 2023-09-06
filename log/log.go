package log

import (
	"errors"
	"fmt"
	"io"
	"runtime"

	"github.com/sirupsen/logrus"
)

type ILogger interface {
	SetLevel(lvl string)
	SetOutput(writer io.Writer)
	SetPrintCaller(printCaller bool)
	SetCaller(caller func() (file string, line int, funcName string, err error))
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	WithFields(fields map[string]interface{}) ILogger
}

type Logger struct {
	entry *logrus.Entry
	// panic fatal error warn warning info debug trace
	level       string
	printCaller bool
	caller      func() (file string, line int, funcName string, err error)
}

func NewLogger() ILogger {
	return newLogger()
}

func defaultCaller() (file string, line int, funcName string, err error) {
	pc, f, l, ok := runtime.Caller(4)
	if !ok {
		err = errors.New("caller failure")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	file, line = f, l
	return
}

// 设置日志打印级别
func (l *Logger) SetLevel(lvl string) {
	if lvl == "" {
		return
	}
	level, err := logrus.ParseLevel(lvl)
	if err == nil {
		l.level = lvl
		l.entry.Logger.Level = level
	}
}

// 设置日志输出位置
func (l *Logger) SetOutput(writer io.Writer) {
	l.entry.Logger.SetOutput(writer)
}

// 设置是否打印调用信息
func (l *Logger) SetPrintCaller(printCaller bool) {
	l.printCaller = printCaller
}

func (l *Logger) SetCaller(caller func() (file string, line int, funcName string, err error)) {
	l.caller = caller
}

// 获取caller信息
func (l *Logger) getCallerInfo(level logrus.Level) map[string]interface{} {
	mp := make(map[string]interface{})
	if l.printCaller || level != logrus.InfoLevel {
		file, line, funcName, err := l.caller()
		if err == nil {
			mp["file"] = fmt.Sprintf("%s:%d", file, line)
			mp["func"] = funcName
		}
	}
	return mp
}

func (l *Logger) log(level logrus.Level, args ...interface{}) {
	l.entry.WithFields(l.getCallerInfo(level)).Log(level, args...)
}

func (l *Logger) logf(level logrus.Level, format string, args ...interface{}) {
	l.entry.WithFields(l.getCallerInfo(level)).Logf(level, format, args...)
}

func (l *Logger) Trace(args ...interface{}) {
	l.log(logrus.TraceLevel, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(logrus.DebugLevel, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log(logrus.InfoLevel, args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.log(logrus.WarnLevel, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log(logrus.ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log(logrus.FatalLevel, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.log(logrus.PanicLevel, args...)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.logf(logrus.TraceLevel, format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(logrus.DebugLevel, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(logrus.InfoLevel, format, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.logf(logrus.WarnLevel, format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(logrus.ErrorLevel, format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logf(logrus.FatalLevel, format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logf(logrus.PanicLevel, format, args...)
}

func (l *Logger) WithFields(fields map[string]interface{}) ILogger {
	entry := l.entry.WithFields(fields)
	return &Logger{
		entry:       entry,
		level:       l.level,
		printCaller: l.printCaller,
		caller:      l.caller,
	}
}

//***********//

var log *Logger

func init() {
	log = newLogger()
}

func newLogger() *Logger {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.AddHook(&errorHook{})
	logger := &Logger{
		entry:  logrus.NewEntry(log),
		caller: defaultCaller,
	}
	return logger
}

func SetLevel(lvl string) {
	if lvl == "" {
		return
	}
	level, err := logrus.ParseLevel(lvl)
	if err == nil {
		log.level = lvl
		log.entry.Logger.Level = level
	}
}

func SetOutput(writer io.Writer) {
	log.entry.Logger.SetOutput(writer)
}

func Trace(args ...interface{}) {
	log.log(logrus.TraceLevel, args...)
}

func Debug(args ...interface{}) {
	log.log(logrus.DebugLevel, args...)
}

func Info(args ...interface{}) {
	log.log(logrus.InfoLevel, args...)
}

func Warning(args ...interface{}) {
	log.log(logrus.WarnLevel, args...)
}

func Error(args ...interface{}) {
	log.log(logrus.ErrorLevel, args...)
}

func Fatal(args ...interface{}) {
	log.log(logrus.FatalLevel, args...)
}

func Panic(args ...interface{}) {
	log.log(logrus.PanicLevel, args...)
}

func Tracef(format string, args ...interface{}) {
	log.logf(logrus.TraceLevel, format, args...)
}

func Debugf(format string, args ...interface{}) {
	log.logf(logrus.DebugLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	log.logf(logrus.InfoLevel, format, args...)
}

func Warningf(format string, args ...interface{}) {
	log.logf(logrus.WarnLevel, format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.logf(logrus.ErrorLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.logf(logrus.FatalLevel, format, args...)
}

func Panicf(format string, args ...interface{}) {
	log.logf(logrus.PanicLevel, format, args...)
}

func WithFields(fields map[string]interface{}) *Logger {
	entry := log.entry.WithFields(fields)
	return &Logger{
		entry:       entry,
		level:       log.level,
		printCaller: log.printCaller,
		caller:      log.caller,
	}
}

func SetPrintCaller(printCaller bool) {
	log.printCaller = printCaller
}

func SetCaller(caller func() (file string, line int, funcName string, err error)) {
	log.caller = caller
}
