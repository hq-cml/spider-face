package core

import (
	"os"
	"errors"
	"io"
	"strings"
	"path"
)

const (
	CompressMinSize = 200
)

var (
	GzipExt = map[string]bool {
		".css":  true,
		".js":   true,
		".html": true,
		".png":  true,
		".jpg":  true,
	}
)

func InitGzipExt(s string) {
	if len(s) == 0 {
		return
	}
	//split即便没有成功，也会返回一个长度为1的slice，作为整体。。。
	t := strings.Split(s, "|")
	GzipExt = map[string]bool{} //清空默认，以用户输入为准
	for _, ext := range t {
		GzipExt[strings.ToLower(ext)] = true
	}
}

func CheckNeedGzip(fileSize int64, request *Request, filePath string) bool {
	if GlobalConf.Gzip != true ||
		fileSize < int64(CompressMinSize) || //太小，没必要压缩
		strings.Index(request.GetHeader("Accept-Encoding"), "gzip") < 0  {
		return false
	}

	ext := strings.ToLower(path.Ext(filePath))
	_, exist := GzipExt[ext]

	if exist {
		return true
	}
	return false
}


type memFile struct {
	fi     *memFileInfo
	offset int64
}

// Close memfile.
func (f *memFile) Close() error {
	return nil
}

// Get os.FileInfo of memfile.
func (f *memFile) Stat() (os.FileInfo, error) {
	return f.fi, nil
}

// read os.FileInfo of files in directory of memfile.
// it returns empty slice.
func (f *memFile) Readdir(count int) ([]os.FileInfo, error) {
	infos := []os.FileInfo{}

	return infos, nil
}

func (f *memFile) Read(p []byte) (n int, err error) {
	if len(f.fi.content)-int(f.offset) >= len(p) {
		n = len(p)
	} else {
		n = len(f.fi.content) - int(f.offset)
		err = io.EOF
	}
	copy(p, f.fi.content[f.offset:f.offset+int64(n)])
	f.offset += int64(n)
	return
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

func (f *memFile) Seek(offset int64, whence int) (ret int64, err error) {
	switch whence {
	default:
		return 0, errWhence
	case os.SEEK_SET:
	case os.SEEK_CUR:
		offset += f.offset
	case os.SEEK_END:
		offset += int64(len(f.fi.content))
	}
	if offset < 0 || int(offset) > len(f.fi.content) {
		return 0, errOffset
	}
	f.offset = offset
	return f.offset, nil
}

type memFileInfo struct {
	os.FileInfo
	content     []byte
}
