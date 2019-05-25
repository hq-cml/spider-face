package core

type SpiderConfig struct {
	BindAddr      string
	StaticPath     string //静态文件根目录
	Logger        SpiderLogger

	TplPath string

	//ReadTimeout   int
	//WriteTimeout  int
	//MaxHeaderByte int

	//HttpErrorHtml map[int]string
}

type SpiderLogger interface {
	log()
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
