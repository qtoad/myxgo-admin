package util

import (
	"archive/zip"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	textTmpl "text/template"
	"time"
)

func Random(strings []string) ([]string, error) {
	for i := len(strings) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}

	str := make([]string, 0)
	for i := 0; i < len(strings); i++ {
		str = append(str, strings[i])
	}
	return str, nil
}

func CompressedHTML(h *template.HTML) {
	st := strings.Split(string(*h), "\n")
	var ss []string
	for i := 0; i < len(st); i++ {
		st[i] = strings.TrimSpace(st[i])
		if st[i] != "" {
			ss = append(ss, st[i])
		}
	}
	*h = template.HTML(strings.Join(ss, "\n"))
}

func ReplaceNth(s, old, new string, n int) string {
	i := 0
	for m := 1; m <= n; m++ {
		x := strings.Index(s[i:], old)
		if x < 0 {
			break
		}
		i += x
		if m == n {
			return s[:i] + new + s[i+len(old):]
		}
		i += len(old)
	}
	return s
}

func WrapURL(u string) string {
	uarr := strings.Split(u, "?")
	if len(uarr) < 2 {
		return url.QueryEscape(strings.ReplaceAll(u, "/", "_"))
	}
	v, err := url.ParseQuery(uarr[1])
	if err != nil {
		return url.QueryEscape(strings.ReplaceAll(u, "/", "_"))
	}
	return url.QueryEscape(strings.ReplaceAll(uarr[0], "/", "_")) + "?" +
		strings.ReplaceAll(v.Encode(), "%7B%7B.Id%7D%7D", "{{.Id}}")
}

func JSON(a interface{}) string {
	if a == nil {
		return ""
	}
	b, _ := json.Marshal(a)
	return string(b)
}

//获取source的子串,如果start小于0或者end大于source长度则返回""
//start:开始index，从0开始，包括0
//end:结束index，以end结束，但不包括end
func substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

func ParseBool(s string) bool {
	b1, _ := strconv.ParseBool(s)
	return b1
}

func ReplaceAll(s string, oldnew ...string) string {
	repl := strings.NewReplacer(oldnew...)
	return repl.Replace(s)
}

func PackageName(v interface{}) string {
	if v == nil {
		return ""
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		return val.Elem().Type().PkgPath()
	}
	return val.Type().PkgPath()
}

func ParseFloat32(f string) float32 {
	s, _ := strconv.ParseFloat(f, 32)
	return float32(s)
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func CopyMap(m map[string]string) map[string]string {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		panic(err)
	}
	var cm map[string]string
	err = dec.Decode(&cm)
	if err != nil {
		panic(err)
	}
	return cm
}

func ParseTime(stringTime string) time.Time {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", stringTime, loc)
	return theTime
}

func ParseHTML(name, tmpl string, param interface{}) template.HTML {
	t := template.New(name)
	t, err := t.Parse(tmpl)
	if err != nil {
		fmt.Println("util parseHTML error", err)
		return ""
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, param)
	if err != nil {
		fmt.Println("util parseHTML error", err)
		return ""
	}
	return template.HTML(buf.String())
}

func ParseText(name, tmpl string, param interface{}) string {
	t := textTmpl.New(name)
	t, err := t.Parse(tmpl)
	if err != nil {
		fmt.Println("util parseHTML error", err)
		return ""
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, param)
	if err != nil {
		fmt.Println("util parseHTML error", err)
		return ""
	}
	return buf.String()
}

func CompareVersion(src, toCompare string) bool {
	if toCompare == "" {
		return false
	}

	exp, _ := regexp.Compile(`-(.*)`)
	src = exp.ReplaceAllString(src, "")
	toCompare = exp.ReplaceAllString(toCompare, "")

	srcs := strings.Split(src, "v")
	srcArr := strings.Split(srcs[1], ".")
	op := ">"
	srcs[0] = strings.TrimSpace(srcs[0])
	if InArray([]string{">=", "<=", "=", ">", "<"}, srcs[0]) {
		op = srcs[0]
	}

	toCompare = strings.ReplaceAll(toCompare, "v", "")

	if op == "=" {
		return srcs[1] == toCompare
	}

	if srcs[1] == toCompare && (op == "<=" || op == ">=") {
		return true
	}

	toCompareArr := strings.Split(strings.ReplaceAll(toCompare, "v", ""), ".")
	for i := 0; i < len(srcArr); i++ {
		v, err := strconv.Atoi(srcArr[i])
		if err != nil {
			return false
		}
		vv, err := strconv.Atoi(toCompareArr[i])
		if err != nil {
			return false
		}
		switch op {
		case ">", ">=":
			if v < vv {
				return true
			} else if v > vv {
				return false
			} else {
				continue
			}
		case "<", "<=":
			if v > vv {
				return true
			} else if v < vv {
				return false
			} else {
				continue
			}
		}
	}

	return false
}

const (
	Byte  = 1
	KByte = Byte * 1024
	MByte = KByte * 1024
	GByte = MByte * 1024
	TByte = GByte * 1024
	PByte = TByte * 1024
	EByte = PByte * 1024
)

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+" %s", val, suffix)
}

// FileSize calculates the file size and generate user-friendly string.
func FileSize(s uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(s, 1024, sizes)
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// TimeSincePro calculates the time interval and generate full user-friendly string.
func TimeSincePro(then time.Time, m map[string]string) string {
	now := time.Now()
	diff := now.Unix() - then.Unix()

	if then.After(now) {
		return "future"
	}

	var timeStr, diffStr string
	for {
		if diff == 0 {
			break
		}

		diff, diffStr = computeTimeDiff(diff, m)
		timeStr += ", " + diffStr
	}
	return strings.TrimPrefix(timeStr, ", ")
}

// Seconds-based time units
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

func computeTimeDiff(diff int64, m map[string]string) (int64, string) {
	diffStr := ""
	switch {
	case diff <= 0:
		diff = 0
		diffStr = "now"
	case diff < 2:
		diff = 0
		diffStr = "1 " + m["second"]
	case diff < 1*Minute:
		diffStr = fmt.Sprintf("%d "+m["seconds"], diff)
		diff = 0

	case diff < 2*Minute:
		diff -= 1 * Minute
		diffStr = "1 " + m["minute"]
	case diff < 1*Hour:
		diffStr = fmt.Sprintf("%d "+m["minutes"], diff/Minute)
		diff -= diff / Minute * Minute

	case diff < 2*Hour:
		diff -= 1 * Hour
		diffStr = "1 " + m["hour"]
	case diff < 1*Day:
		diffStr = fmt.Sprintf("%d "+m["hours"], diff/Hour)
		diff -= diff / Hour * Hour

	case diff < 2*Day:
		diff -= 1 * Day
		diffStr = "1 " + m["day"]
	case diff < 1*Week:
		diffStr = fmt.Sprintf("%d "+m["days"], diff/Day)
		diff -= diff / Day * Day

	case diff < 2*Week:
		diff -= 1 * Week
		diffStr = "1 " + m["week"]
	case diff < 1*Month:
		diffStr = fmt.Sprintf("%d "+m["weeks"], diff/Week)
		diff -= diff / Week * Week

	case diff < 2*Month:
		diff -= 1 * Month
		diffStr = "1 " + m["month"]
	case diff < 1*Year:
		diffStr = fmt.Sprintf("%d "+m["months"], diff/Month)
		diff -= diff / Month * Month

	case diff < 2*Year:
		diff -= 1 * Year
		diffStr = "1 " + m["year"]
	default:
		diffStr = fmt.Sprintf("%d "+m["years"], diff/Year)
		diff = 0
	}
	return diff, diffStr
}

func DownloadTo(url, output string) error {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	file, err := os.Create(output)

	if err != nil {
		return err
	}

	_, err = io.Copy(file, res.Body)

	if err != nil {
		return err
	}

	return nil
}

func UnzipDir(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	err = os.MkdirAll(dest, 0750)

	if err != nil {
		return err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(filepath.Dir(path), f.Mode())
			if err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

// NormalizeEOL will convert Windows (CRLF) and Mac (CR) EOLs to UNIX (LF)
func NormalizeEOL(input []byte) []byte {
	var right, left, pos int
	if right = bytes.IndexByte(input, '\r'); right == -1 {
		return input
	}
	length := len(input)
	tmp := make([]byte, length)

	// We know that left < length because otherwise right would be -1 from IndexByte.
	copy(tmp[pos:pos+right], input[left:left+right])
	pos += right
	tmp[pos] = '\n'
	left += right + 1
	pos++

	for left < length {
		if input[left] == '\n' {
			left++
		}

		right = bytes.IndexByte(input[left:], '\r')
		if right == -1 {
			copy(tmp[pos:], input[left:])
			pos += length - left
			break
		}
		copy(tmp[pos:pos+right], input[left:left+right])
		pos += right
		tmp[pos] = '\n'
		left += right + 1
		pos++
	}
	return tmp[:pos]
}

func IsDirectory(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsFile(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func MkdirIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModePerm)
	}
}

func MkFileIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, _ = os.Create(path)
	}
}

func UploadFileTo(fh *multipart.FileHeader, destDirectory string) (int64, error) {
	src, err := fh.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()

	out, err := os.OpenFile(filepath.Join(destDirectory, fh.Filename),
		os.O_WRONLY|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return 0, err
	}
	defer out.Close()

	return io.Copy(out, src)
}

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetParentDirectory(directory string) string {
	return Substr(directory, 0, strings.LastIndex(directory, "/"))
}

func ParseFileContentType(fileName string) string {
	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	if strings.HasPrefix(contentType, "text/") {
		contentType = "text/plain"
	}
	return contentType
}

func IsHiddenFile(name string) bool {
	if strings.TrimSpace(name) == "" {
		return false
	}

	return strings.HasPrefix(name, ".")
}

func ByteCountIEC(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
