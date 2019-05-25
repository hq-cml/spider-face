package core

type SpiderConfig struct {
	BindAddr      string
	StaticPath     string //静态文件根目录
	LogPath       string
	LogLevel       string

	logger        SpiderLogger

	TplPath string

	//ReadTimeout   int
	//WriteTimeout  int
	//MaxHeaderByte int

	//HttpErrorHtml map[int]string
}

type SpiderLogger interface {
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
	Debug(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Info(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})
	Warn(v ...interface{})
	Errf(format string, v ...interface{})
	Errln(v ...interface{})
	Err(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Fatal(v ...interface{})
}

type DefaultLogger struct {

}

var GlobalConf *SpiderConfig

var (
	DEFAULT_CONTROLLER     string       = "index"
	DEFAULT_ACTION         string       = "index"
	ACTION_SUFFIX          string       = "Action"
	HTTP_METHOD_PARAM_NAME string       = "m"
)
