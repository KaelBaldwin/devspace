package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	goansi "github.com/k0kubun/go-ansi"
	"github.com/mgutz/ansi"

	"github.com/sirupsen/logrus"
)

var stdout = goansi.NewAnsiStdout()
var stderr = goansi.NewAnsiStderr()

type stdoutLogger struct {
	logMutex sync.Mutex
	level    logrus.Level

	loadingText *loadingText
	fileLogger  Logger
}

type fnTypeInformation struct {
	tag      string
	color    string
	logLevel logrus.Level
	stream   io.Writer
}

var fnTypeInformationMap = map[logFunctionType]*fnTypeInformation{
	debugFn: {
		tag:      "[debug]  ",
		color:    "green+b",
		logLevel: logrus.DebugLevel,
		stream:   stdout,
	},
	infoFn: {
		tag:      "[info]   ",
		color:    "cyan+b",
		logLevel: logrus.InfoLevel,
		stream:   stdout,
	},
	warnFn: {
		tag:      "[warn]   ",
		color:    "166+b",
		logLevel: logrus.WarnLevel,
		stream:   stdout,
	},
	errorFn: {
		tag:      "[error]  ",
		color:    "red+b",
		logLevel: logrus.ErrorLevel,
		stream:   stdout,
	},
	fatalFn: {
		tag:      "[fatal]  ",
		color:    "red+b",
		logLevel: logrus.FatalLevel,
		stream:   stdout,
	},
	panicFn: {
		tag:      "[panic]  ",
		color:    "red+b",
		logLevel: logrus.PanicLevel,
		stream:   stderr,
	},
	doneFn: {
		tag:      "[done] √ ",
		color:    "green+b",
		logLevel: logrus.InfoLevel,
		stream:   stdout,
	},
	failFn: {
		tag:      "[fail] X ",
		color:    "red+b",
		logLevel: logrus.ErrorLevel,
		stream:   stdout,
	},
}

func (s *stdoutLogger) writeMessage(fnType logFunctionType, message string) {
	fnInformation := fnTypeInformationMap[fnType]

	if s.level >= fnInformation.logLevel {
		if s.loadingText != nil {
			s.loadingText.Stop()
		}

		fnInformation.stream.Write([]byte(ansi.Color(fnInformation.tag, fnInformation.color)))
		// ct.Foreground(fnInformation.color, false)
		// fnInformation.stream.Write([]byte(fnInformation.tag))
		// ct.ResetColor()

		fnInformation.stream.Write([]byte(message))

		if s.loadingText != nil && fnType != fatalFn {
			s.loadingText.Start()
		}
	}
}

func (s *stdoutLogger) writeMessageToFileLogger(fnType logFunctionType, args ...interface{}) {
	fnInformation := fnTypeInformationMap[fnType]

	if s.level >= fnInformation.logLevel && s.fileLogger != nil {
		switch fnType {
		case doneFn:
			s.fileLogger.Done(args...)
		case infoFn:
			s.fileLogger.Info(args...)
		case debugFn:
			s.fileLogger.Debug(args...)
		case warnFn:
			s.fileLogger.Warn(args...)
		case failFn:
			s.fileLogger.Fail(args...)
		case errorFn:
			s.fileLogger.Error(args...)
		case panicFn:
			s.fileLogger.Panic(args...)
		case fatalFn:
			s.fileLogger.Fatal(args...)
		}
	}
}

func (s *stdoutLogger) writeMessageToFileLoggerf(fnType logFunctionType, format string, args ...interface{}) {
	fnInformation := fnTypeInformationMap[fnType]

	if s.level >= fnInformation.logLevel && s.fileLogger != nil {
		switch fnType {
		case doneFn:
			s.fileLogger.Donef(format, args...)
		case infoFn:
			s.fileLogger.Infof(format, args...)
		case debugFn:
			s.fileLogger.Debugf(format, args...)
		case warnFn:
			s.fileLogger.Warnf(format, args...)
		case failFn:
			s.fileLogger.Failf(format, args...)
		case errorFn:
			s.fileLogger.Errorf(format, args...)
		case panicFn:
			s.fileLogger.Panicf(format, args...)
		case fatalFn:
			s.fileLogger.Fatalf(format, args...)
		}
	}
}

// StartWait prints a wait message until StopWait is called
func (s *stdoutLogger) StartWait(message string) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.loadingText != nil {
		if s.loadingText.Message == message {
			return
		}

		s.loadingText.Stop()
		s.loadingText = nil
	}

	s.loadingText = &loadingText{
		Message: message,
		Stream:  goansi.NewAnsiStdout(),
	}

	s.loadingText.Start()
}

// StartWait prints a wait message until StopWait is called
func (s *stdoutLogger) StopWait() {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.loadingText != nil {
		s.loadingText.Stop()
		s.loadingText = nil
	}
}

// PrintTable implements logger interface
func (s *stdoutLogger) PrintTable(header []string, values [][]string) {
	columnLengths := make([]int, len(header))

	for k, v := range header {
		columnLengths[k] = len(v)
	}

	// Get maximum column length
	for _, v := range values {
		for key, value := range v {
			if len(value) > columnLengths[key] {
				columnLengths[key] = len(value)
			}
		}
	}

	s.Write([]byte("\n"))

	// Print Header
	for key, value := range header {
		WriteColored(" "+value+"  ", "green+b")

		padding := columnLengths[key] - len(value)

		if padding > 0 {
			s.Write([]byte(strings.Repeat(" ", padding)))
		}
	}

	s.Write([]byte("\n"))

	if len(values) == 0 {
		s.Write([]byte(" No entries found\n"))
	}

	// Print Values
	for _, v := range values {
		for key, value := range v {
			s.Write([]byte(" " + value + "  "))

			padding := columnLengths[key] - len(value)

			if padding > 0 {
				s.Write([]byte(strings.Repeat(" ", padding)))
			}
		}

		s.Write([]byte("\n"))
	}

	s.Write([]byte("\n"))
}

func (s *stdoutLogger) Debug(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(debugFn, fmt.Sprintln(args...))
	s.writeMessageToFileLogger(debugFn, args...)
}

func (s *stdoutLogger) Debugf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(debugFn, fmt.Sprintf(format, args...)+"\n")
	s.writeMessageToFileLoggerf(debugFn, format, args...)
}

func (s *stdoutLogger) Info(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(infoFn, fmt.Sprintln(args...))
	s.writeMessageToFileLogger(infoFn, args...)
}

func (s *stdoutLogger) Infof(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(infoFn, fmt.Sprintf(format, args...)+"\n")
	s.writeMessageToFileLoggerf(infoFn, format, args...)
}

func (s *stdoutLogger) Warn(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(warnFn, fmt.Sprintln(args...))
	s.writeMessageToFileLogger(warnFn, args...)
}

func (s *stdoutLogger) Warnf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(warnFn, fmt.Sprintf(format, args...)+"\n")
	s.writeMessageToFileLoggerf(warnFn, format, args...)
}

func (s *stdoutLogger) Error(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(errorFn, fmt.Sprintln(args...))
	s.writeMessageToFileLogger(errorFn, args...)
}

func (s *stdoutLogger) Errorf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(errorFn, fmt.Sprintf(format, args...)+"\n")
	s.writeMessageToFileLoggerf(errorFn, format, args...)
}

func (s *stdoutLogger) Fatal(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(fatalFn, fmt.Sprintln(args...))
	s.writeMessageToFileLogger(fatalFn, args...)

	if s.fileLogger == nil {
		os.Exit(1)
	}
}

func (s *stdoutLogger) Fatalf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(fatalFn, fmt.Sprintf(format, args...)+"\n")
	s.writeMessageToFileLoggerf(fatalFn, format, args...)

	if s.fileLogger == nil {
		os.Exit(1)
	}
}

func (s *stdoutLogger) Panic(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(panicFn, fmt.Sprintln(args...))
	s.writeMessageToFileLogger(panicFn, args...)

	if s.fileLogger == nil {
		panic(fmt.Sprintln(args...))
	}
}

func (s *stdoutLogger) Panicf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(panicFn, fmt.Sprintf(format, args...)+"\n")
	s.writeMessageToFileLoggerf(panicFn, format, args...)

	if s.fileLogger == nil {
		panic(fmt.Sprintf(format, args...))
	}
}

func (s *stdoutLogger) Done(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(doneFn, fmt.Sprintln(args...))
	s.writeMessageToFileLogger(doneFn, args...)

}

func (s *stdoutLogger) Donef(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(doneFn, fmt.Sprintf(format, args...)+"\n")
	s.writeMessageToFileLoggerf(doneFn, format, args...)
}

func (s *stdoutLogger) Fail(args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(failFn, fmt.Sprintln(args...))
	s.writeMessageToFileLogger(failFn, args...)
}

func (s *stdoutLogger) Failf(format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(failFn, fmt.Sprintf(format, args...)+"\n")
	s.writeMessageToFileLoggerf(failFn, format, args...)
}

func (s *stdoutLogger) Print(level logrus.Level, args ...interface{}) {
	switch level {
	case logrus.InfoLevel:
		s.Info(args...)
	case logrus.DebugLevel:
		s.Debug(args...)
	case logrus.WarnLevel:
		s.Warn(args...)
	case logrus.ErrorLevel:
		s.Error(args...)
	case logrus.PanicLevel:
		s.Panic(args...)
	case logrus.FatalLevel:
		s.Fatal(args...)
	}
}

func (s *stdoutLogger) Printf(level logrus.Level, format string, args ...interface{}) {
	switch level {
	case logrus.InfoLevel:
		s.Infof(format, args...)
	case logrus.DebugLevel:
		s.Debugf(format, args...)
	case logrus.WarnLevel:
		s.Warnf(format, args...)
	case logrus.ErrorLevel:
		s.Errorf(format, args...)
	case logrus.PanicLevel:
		s.Panicf(format, args...)
	case logrus.FatalLevel:
		s.Fatalf(format, args...)
	}
}

func (s *stdoutLogger) With(obj interface{}) *LoggerEntry {
	return &LoggerEntry{
		logger: s,
		context: map[string]interface{}{
			"context-1": obj,
		},
	}
}

func (s *stdoutLogger) WithKey(key string, obj interface{}) *LoggerEntry {
	return &LoggerEntry{
		logger: s,
		context: map[string]interface{}{
			key: obj,
		},
	}
}

func (s *stdoutLogger) SetLevel(level logrus.Level) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.level = level
}

func (s *stdoutLogger) printWithContext(fnType logFunctionType, context map[string]interface{}, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(fnType, fmt.Sprintln(args...))

	if s.fileLogger != nil && s.level >= fnTypeInformationMap[fnType].logLevel {
		s.fileLogger.printWithContext(fnType, context, args...)
	}
}

func (s *stdoutLogger) printWithContextf(fnType logFunctionType, context map[string]interface{}, format string, args ...interface{}) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	s.writeMessage(fnType, fmt.Sprintf(format, args...)+"\n")

	if s.fileLogger != nil && s.level >= fnTypeInformationMap[fnType].logLevel {
		s.fileLogger.printWithContextf(fnType, context, format, args...)
	}
}

func (s *stdoutLogger) Write(message []byte) (int, error) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.loadingText != nil {
		s.loadingText.Stop()
	}

	n, err := fnTypeInformationMap[infoFn].stream.Write(message)

	if s.loadingText != nil {
		s.loadingText.Start()
	}

	return n, err
}

func (s *stdoutLogger) WriteString(message string) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.loadingText != nil {
		s.loadingText.Stop()
	}

	fnTypeInformationMap[infoFn].stream.Write([]byte(message))

	if s.loadingText != nil {
		s.loadingText.Start()
	}
}
