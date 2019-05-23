package core
var (
	DEFAULT_CONTROLLER     string       = "index"
	DEFAULT_ACTION         string       = "index"
	ACTION_SUFFIX          string       = "Action"
	HTTP_METHOD_PARAM_NAME string       = "m"
)

type SpiderLogger interface {
	log()
}

type DefaultLogger struct {

}