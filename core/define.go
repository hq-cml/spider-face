package core

type SpiderConfig struct {
	BindAddr            string
	StaticPath          string //静态文件根目录
	LogPath             string
	LogLevel            string
	TplPath             string

	ReadTimeout   		int64
	WriteTimeout  		int64
	MaxHeaderByte 		int64

	CustomHttpErrorHtml map[int]string    //定制化的错误页面 httpCode => customErr.html
	CustomRewriteRule   map[string]string

	Mime                bool    //配置Mime
	Gzip                bool    //开启gzip压缩，默认对：.css、.js、.html、.jpg、.png进行压缩
	CustomGzipExt       string  //用户自己指定的压缩文件后缀，用|分隔
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

var GlobalConf *SpiderConfig

const (
	DEFAULT_CONTROLLER     string       = "index"
	DEFAULT_ACTION         string       = "index"
	CONTROLLER_SUFFIX      string       = "Controller"
	ACTION_SUFFIX          string       = "Action"
)

const (
	PATH_INFO_IDENTITY string = "***"
)
