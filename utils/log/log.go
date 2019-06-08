package log

/*
 * 日志相关
 */
import (
	"log"
	"os"
	"runtime"
	"strings"
	"fmt"
)

type SpiderLog struct {
	*log.Logger
	level   int
}

const (
	SPIDER_LOG_LEVEL_DEBUG = iota
	SPIDER_LOG_LEVEL_INFO
	SPIDER_LOG_LEVEL_WARN
	SPIDER_LOG_LEVEL_ERR
	SPIDER_LOG_LEVEL_FATAL
)

var DefaultLogger *SpiderLog

func init() {
	DefaultLogger = NewLog("", "debug")
}

func getFileLenPrefix() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	tmp := strings.Split(file, "/")
	sFile := tmp[len(tmp) - 1]

	return fmt.Sprintf("[%s, %v] ", sFile, line)
}

func NewLog(path string, level string) *SpiderLog {
	var f *os.File
	var err error
	if path == "" { //如果没有文件路径，则用标准错误Stderr
		f = os.Stderr
	} else {
		f, err = os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND , 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	levelInt := SPIDER_LOG_LEVEL_DEBUG
	switch level {
	case "debug", "Debug", "DEBUG":
		levelInt = SPIDER_LOG_LEVEL_DEBUG
	case "info", "Info", "INFO":
		levelInt = SPIDER_LOG_LEVEL_INFO
	case "warn", "Warn", "WARN":
		levelInt = SPIDER_LOG_LEVEL_WARN
	case "err", "Err", "ERR":
		levelInt = SPIDER_LOG_LEVEL_ERR
	case "fatal", "Fatal", "FATAL":
		levelInt = SPIDER_LOG_LEVEL_FATAL
	}
	spiderLog := SpiderLog {
		Logger:log.New(f, "", log.LstdFlags),
		level: levelInt,
	}
	return &spiderLog
}

func (spiderLog *SpiderLog)Debugf(format string, v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_DEBUG {return}
	prefix := fmt.Sprintf("[DEBUG]%v", getFileLenPrefix())
	spiderLog.Printf(prefix + format, v...)
}

func (spiderLog *SpiderLog)Debugln(v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_DEBUG {return}
	prefix := fmt.Sprintf("[DEBUG]%v", getFileLenPrefix())
	v1 := []interface{}{prefix}
	v1 = append(v1, v...)
	spiderLog.Println(v1...)
}

func (spiderLog *SpiderLog)Debug(v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_DEBUG {return}
	prefix := fmt.Sprintf("[DEBUG]%v", getFileLenPrefix())
	v1 := []interface{}{prefix}
	v1 = append(v1, v...)
	spiderLog.Print(v1...)
}


func (spiderLog *SpiderLog)Infof(format string, v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_INFO {return}
	prefix := fmt.Sprintf("[ INFO]%v", getFileLenPrefix())
	spiderLog.Printf(prefix + format, v...)
}

func (spiderLog *SpiderLog)Infoln(v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_INFO {return}
	prefix := fmt.Sprintf("[ INFO]%v", getFileLenPrefix())
	v1 := []interface{}{prefix}
	v1 = append(v1, v...)
	spiderLog.Println(v1...)
}

func (spiderLog *SpiderLog)Info(v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_INFO {return}
	prefix := fmt.Sprintf("[ INFO]%v", getFileLenPrefix())
	v1 := []interface{}{prefix}
	v1 = append(v1, v...)
	spiderLog.Print(v1...)
}


func (spiderLog *SpiderLog)Warnf(format string, v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_WARN {return}
	prefix := fmt.Sprintf("[ WARN]%v", getFileLenPrefix())
	spiderLog.Printf(prefix + format, v...)
}

func (spiderLog *SpiderLog)Warnln(v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_WARN {return}
	prefix := fmt.Sprintf("[ WARN]%v", getFileLenPrefix())
	v1 := []interface{}{prefix}
	v1 = append(v1, v...)
	spiderLog.Println(v1...)
}

func(spiderLog *SpiderLog) Warn(v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_WARN {return}
	prefix := fmt.Sprintf("[ WARN]%v", getFileLenPrefix())
	v1 := []interface{}{prefix}
	v1 = append(v1, v...)
	spiderLog.Print(v1...)
}


func (spiderLog *SpiderLog)Errf(format string, v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_ERR {return}
	prefix := fmt.Sprintf("[ERROR]%v", getFileLenPrefix())
	spiderLog.Printf(prefix + format, v...)
}

func (spiderLog *SpiderLog)Errln(v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_ERR {return}
	prefix := fmt.Sprintf("[ERROR]%v", getFileLenPrefix())
	v1 := []interface{}{prefix}
	v1 = append(v1, v...)
	spiderLog.Println(v1...)
}

func (spiderLog *SpiderLog)Err(v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_ERR {return}
	prefix := fmt.Sprintf("[ERROR]%v", getFileLenPrefix())
	v1 := []interface{}{prefix}
	v1 = append(v1, v...)
	spiderLog.Print(v1...)
}


func (spiderLog *SpiderLog)Fatalf(format string, v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_FATAL {return}
	prefix := fmt.Sprintf("[FATAL]%v", getFileLenPrefix())
	spiderLog.Printf(prefix + format, v...)
}

func (spiderLog *SpiderLog)Fatalln(v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_FATAL {return}
	prefix := fmt.Sprintf("[FATAL]%v", getFileLenPrefix())
	v1 := []interface{}{prefix}
	v1 = append(v1, v...)
	spiderLog.Println(v1...)
}

func (spiderLog *SpiderLog)Fatal(v ...interface{}) {
	if spiderLog.level > SPIDER_LOG_LEVEL_FATAL {return}
	prefix := fmt.Sprintf("[FATAL]%v", getFileLenPrefix())
	v1 := []interface{}{prefix}
	v1 = append(v1, v...)
	spiderLog.Print(v1...)
}
