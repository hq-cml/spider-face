package core

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"regexp"
)

var (
	ViewRoot      string
	ViewExt       string = ".html"
	ViewTemplates map[string]*template.Template
	templateFiles map[string]string
	viewFunc      template.FuncMap //map[string]interface{}
)

type View struct {
	data map[interface{}]interface{}
}

func init() {
	viewFunc = template.FuncMap{}
	viewFunc["date"] = Date
	viewFunc["strtotime"] = StrToTime
	viewFunc["time"] = Time
}

func NewView() *View {
	return &View{}
}

func (this *View) Assign(key interface{}, value interface{}) {
	if this.data == nil {
		this.data = make(map[interface{}]interface{})
		this.data[key] = value
	} else {
		this.data[key] = value
	}
}

func (this *View) Render(viewPathName string) ([]byte, error) {
	if ViewRoot == "" {
		return []byte(""), errors.New("TplPath not set")
	}
	viewPathName = strings.ToLower(viewPathName)
	//TODO
	//if RunMod == "dev" {
	//	t := template.New(view_name).Delims("{{", "}}").Funcs(view_func)
	//	t, err := parseTemplate(t, ViewRoot+"/"+view_name+ViewExt)
	//	if err != nil || t == nil {
	//		return []byte(""), err
	//	}
	//	ViewTemplates[view_name] = t
	//}

	if tpl, ok := ViewTemplates[viewPathName]; ok == false {
		return []byte(""), errors.New("template " + ViewRoot + "/" + viewPathName + ViewExt + " not found or compile failed")
	} else {
		html_content_bytes := bytes.NewBufferString("")
		err := tpl.ExecuteTemplate(html_content_bytes, viewPathName, this.data)
		if err != nil {
			return []byte(""), err
		}
		html_content, _ := ioutil.ReadAll(html_content_bytes)
		return html_content, nil
	}
}

//TODO支持用户自己增加
func AddViewFunc(key string, func_name interface{}) {
	viewFunc[key] = func_name
}

func InitViewTemplate(viewRoot string) error {
	if viewRoot == "" {
		return nil
	}
	ViewRoot = viewRoot
	templateFiles = make(map[string]string)

	filepath.Walk(viewRoot, func(path string, f os.FileInfo, err error) error {
		//忽略目录名和软链
		if f.IsDir() || (f.Mode()&os.ModeSymlink) > 0 {
			return nil
		}

		//忽略非法后缀名
		if strings.HasSuffix(path, ViewExt) == false {
			return nil
		}

		fileName := strings.Trim(strings.Replace(path, ViewRoot, "", 1), "/")
		templateFiles[strings.TrimSuffix(fileName, ViewExt)] = path
		return nil
	})

	//fmt.Println("templateFiles: ", helper.JsonEncode(templateFiles))
	ViewTemplates = make(map[string]*template.Template)
	for name, filePath := range templateFiles {
		if _, err := os.Stat(filePath); err != nil && os.IsNotExist(err) {
			fmt.Printf("parse template %q err : %q", filePath, err)
			continue
		}

		//注册自定义函数
		tpl := template.New(name).Delims("{{", "}}").Funcs(viewFunc)

		//解析模板
		fmt.Println("\n\nQ----------------")
		//t, err := parseTemplate(t, filePath)
		//if err != nil || t == nil {
		//	continue
		//}

		files := []string{}
		err := getAllFiles(filePath, &files)
		if err != nil {
			panic(err)
		}
		tpl, err = tpl.ParseFiles(files...)
		if err != nil {
			panic(err)
		}
		ViewTemplates[name] = tpl
	}

	return nil
}

func getAllFiles(path string, files *[]string) error {
	*files = append(*files, path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	reg := regexp.MustCompile(`{{\s{0,}template\s{0,}"(.*?)".*?}}`)
	matches := reg.FindAllStringSubmatch(string(data), -1)
	for _, v := range matches {
		if v == nil || v[1] == "" {
			continue
		}

		subFile := ViewRoot + "/" + v[1] + ViewExt
		if subFile == path {
			continue
		}

		*files = append(*files, subFile)

		err = getAllFiles(subFile, files)
		if err != nil {
			return err
		}
	}
	return nil
}

//func parseTemplate(tpl *template.Template, file string) (t *template.Template, err error) {
//	data, _ := ioutil.ReadFile(file)
//	fmt.Println("P--------------", file)
//	t, err = tpl.Parse(string(data))
//	if err != nil {
//		return nil, err
//	}
//
//	//TODO 功能重复??
//	reg := regexp.MustCompile(`{{\s{0,}template\s{0,}"(.*?)".*?}}`)
//	match := reg.FindAllStringSubmatch(string(data), -1)
//	for _, v := range match {
//		if v == nil || v[1] == "" {
//			continue
//		}
//
//		fmt.Println("X--------------", v[1], file)
//		tlook := t.Lookup(v[1])
//		if tlook != nil {
//			continue
//		}
//		deep_file := ViewRoot + "/" + v[1] + ViewExt
//		fmt.Println("Y--------------", deep_file, file)
//		if deep_file == file {
//			continue
//		}
//
//		t, err = parseTemplate(t, deep_file)
//		if err != nil {
//			return nil, err
//		}
//	}
//	return t, nil
//}

//go风格的时间格式替换
var goDateReplacePattern = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", //A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

func StrToTime(dateString, format string) (time.Time, error) {
	replacer := strings.NewReplacer(goDateReplacePattern...)
	format = replacer.Replace(format)
	return time.ParseInLocation(format, dateString, time.Local)
}

func Date(format string, t time.Time) string {
	replacer := strings.NewReplacer(goDateReplacePattern...)
	format = replacer.Replace(format)
	return t.Format(format)
}

func Time() int64 {
	return time.Now().Unix()
}