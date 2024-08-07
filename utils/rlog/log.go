package rlog

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

const Lshortfile = log.Lshortfile
const LstdFlags = log.LstdFlags

const (
	debugLvl = 0
	infoLvl  = 1
	errorLvl = 2
	fatalLvl = 3
)

func getColor(level int) (fg FontStyle, bg FontStyle) {
	if level == debugLvl {
		fg = ForegroundBlack
		bg = BackgroundBlue
	} else if level == infoLvl {
		fg = ForegroundBlack
		bg = BackgroundYellow
	} else if level == errorLvl {
		fg = ForegroundWhite
		bg = BackgroundRed
	} else if level == fatalLvl {
		fg = ForegroundWhite
		bg = BackgroundBlack
	} else {
		fg = ForegroundDefault
		bg = BackgroundDefault
	}
	return
}

var loggers = New(Lshortfile | LstdFlags)

func GetLogger() *Logger {
	return loggers
}

func (logger *Logger) GetWriter() io.Writer {
	return logger.baseLogger.Writer()
}

type Logger struct {
	level      int
	baseLogger *log.Logger
}

func New(flag int) *Logger {
	baseLogger := log.New(os.Stdout, "", flag)

	logger := new(Logger)
	logger.baseLogger = baseLogger

	return logger
}

func (logger *Logger) printfLog(level int, format string, a ...interface{}) {
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	fgColor, bgColor := getColor(level)
	format = fmt.Sprint(PrintWithColor(format, Reset, fgColor, bgColor))
	_ = logger.baseLogger.Output(3, fmt.Sprintf(format, a...))

	if level == fatalLvl {
		os.Exit(1)
	}
}

func (logger *Logger) printLog(level int, a ...interface{}) {
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	p := make([]interface{}, 0, 8)
	for _, b := range a {
		if len(a) > 1 {
			p = append(p, b, " ")
		} else {
			p = append(p, b)
		}
	}
	msg := fmt.Sprint(p...)
	fgColor, bgColor := getColor(level)
	msg = fmt.Sprint(PrintWithColor(msg, Reset, fgColor, bgColor))
	_ = logger.baseLogger.Output(3, msg)

	if level == fatalLvl {
		os.Exit(1)
	}
}

func Debug(a ...any) {
	loggers.printLog(debugLvl, a...)
}

func Info(a ...any) {
	loggers.printLog(infoLvl, a...)
}

func Error(a ...any) {
	loggers.printLog(errorLvl, a...)
}

func Fatal(a ...any) {
	loggers.printLog(fatalLvl, a...)
}

func Debugf(format string, a ...any) {
	loggers.printfLog(debugLvl, format, a...)
}

func Infof(format string, a ...any) {
	loggers.printfLog(infoLvl, format, a...)
}

func Errorf(format string, a ...any) {
	loggers.printfLog(errorLvl, format, a...)
}

func Fatalf(format string, a ...any) {
	loggers.printfLog(fatalLvl, format, a...)
}

func PrettyPrintJSON(data interface{}) {
	// Marshal the data with indentation
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to generate pretty JSON: %v", err)
	}
	// Print the formatted JSON
	fmt.Println(string(prettyJSON))
}
