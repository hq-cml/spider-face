package core

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"fmt"
)

//TODO 待测试
type Rewriter struct {
	rewriteRegexp []*RegexRewriteRule
	rewriteStatic map[string]string
}

type RegexRewriteRule struct {
	pattern string
	match   string
	regex   *regexp.Regexp
}

func NewRewriter() *Rewriter {
	return &Rewriter{
		rewriteRegexp : []*RegexRewriteRule{},
		rewriteStatic : map[string]string{},
	}
}
// This most like nginx rewrite module
// Support regex
// This is different from controller rewrite-router
func (rwt *Rewriter)RegRewriteRule(list map[string]string) {
	for p, m := range list {
		if strings.Index(p, "(") < 0 {
			rwt.rewriteStatic[p] = m
			continue
		}
		r := regexp.MustCompile(p)
		if r == nil {
			continue
		}
		reg := &RegexRewriteRule{
			pattern: p,
			match:   m,
			regex:   r,
		}
		rwt.rewriteRegexp = append(rwt.rewriteRegexp, reg)
	}
}

// Match rewrite
// Note that this will change URL.RawQuery
// and URL.Path in http.Request
func (rwt *Rewriter)MatchRewrite(r *http.Request) {
	urlPath := r.URL.Path
	var rewrite_url string = ""
	var ok bool

	fmt.Println("urlPath-----", urlPath)
	if rewrite_url, ok = rwt.rewriteStatic[urlPath]; ok == true {
		goto RESET_URI
	}

	for _, rewrite := range rwt.rewriteRegexp {
		fmt.Println("B--------")
		match := rewrite.regex.FindAllStringSubmatch(urlPath, -1)
		if match == nil {
			continue
		}
		match_cnt := len(match[0])
		if match_cnt == 1 {
			return
		}

		rewrite_url = rewrite.match

		for n := 1; n < match_cnt; n++ {
			replace_val := "[" + strconv.Itoa(n) + "]"
			rewrite_url = strings.Replace(rewrite_url, replace_val, match[0][n], -1)
		}
		break
	}
	if rewrite_url == "" {
		fmt.Println("C--------")
		return
	}

	RESET_URI:
	fmt.Println("rewrite_url:", rewrite_url)
	rewrite_url = strings.Replace(rewrite_url, "[args]", r.URL.RawQuery, -1)
	fmt.Println("rewrite_url:", rewrite_url)
	uri_map := strings.SplitN(rewrite_url, "?", 2)

	if len(uri_map) == 2 {
		fmt.Println("X-----")
		r.URL.Path = uri_map[0]
		r.URL.RawQuery = uri_map[1]
	} else {
		fmt.Println("Y-----")
		r.URL.Path = uri_map[0]
	}
}
