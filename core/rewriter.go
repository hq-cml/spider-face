package core

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Rewriter struct {
	logger        SpiderLogger
	RegexpRewrite []*RegexRewriteRule
	StaticRewrite map[string]string
}

type RegexRewriteRule struct {
	pattern    string
	rewriteUrl string
	regex      *regexp.Regexp
}

func NewRewriter(logger SpiderLogger) *Rewriter {
	return &Rewriter{
		logger: logger,
		RegexpRewrite : []*RegexRewriteRule{},
		StaticRewrite : map[string]string{},
	}
}

//注册rewrite规则
// This most like nginx rewrite module
// Support regex
// This is different from controller rewrite-routerManger
func (rwt *Rewriter) RegisterRewriteRule(list map[string]string) {
	for p, m := range list {
		if strings.Index(p, "(") < 0 {
			rwt.StaticRewrite[p] = m
		} else {
			r := regexp.MustCompile(p)
			if r == nil {
				continue
			}
			reg := &RegexRewriteRule{
				pattern:    p,
				rewriteUrl: m,
				regex:      r,
			}
			rwt.RegexpRewrite = append(rwt.RegexpRewrite, reg)
		}
	}
}

// Try to Match rewrite
// Note :
// that this will change URL.RawQuery
// and URL.Path in http.Request !!!
func (rwt *Rewriter)TryMatchRewrite(r *http.Request) {
	urlPath := r.URL.Path
	var rewriteUrl string = ""
	var exist bool

	_, exist = rwt.StaticRewrite[urlPath]
	if exist {
		rewriteUrl = rwt.StaticRewrite[urlPath]
	} else {
		for _, rewrite := range rwt.RegexpRewrite {
			matches := rewrite.regex.FindAllStringSubmatch(urlPath, -1)
			if matches == nil {
				continue
			}
			matchCnt := len(matches[0])
			if matchCnt == 1 { //完全匹配, 没什么要改写的了
				return
			}

			//获取改写后的值
			rewriteUrl = rewrite.rewriteUrl
			//改写参数，逐个迁移到改写后的值上面
			for i := 1; i < matchCnt; i++ {
				replaceVal := "$" + strconv.Itoa(i)
				rewriteUrl = strings.Replace(rewriteUrl, replaceVal, matches[0][i], -1)
			}
			break
		}
	}

	//没有改写必要
	if rewriteUrl == "" {
		return
	}

	//实施改写: 直接变更request.URL的值!!
	rwt.logger.Infof("Before Rewrite. UrlPath: '%s', RawQuery: '%s'", urlPath, r.URL.RawQuery)
	uriMap := strings.SplitN(rewriteUrl, "?", 2)
	if len(uriMap) == 2 {
		r.URL.Path = uriMap[0]
		r.URL.RawQuery = uriMap[1]
	} else {
		r.URL.Path = uriMap[0]
	}
	rwt.logger.Infof("After Rewrite. UrlPath: '%s', RawQuery: '%s'", r.URL.Path, r.URL.RawQuery)
}
