// File: tool.go
// Author: walker
// Mail: walkerIVI@gmail.com
// ToolList:
//    func LogPrint()
//    func LogFile(stdout bool)
//    func CheckErr(err error)
//    func Request(url string, method string, reqD string, reqH map[string]string) string
//    func Btoi(s string) (int, error)
//    type Base64 struct
//    func Base64New() Base64
//    func (this *Base64) Set(table string, padding string)
//    func (this *Base64) Encode(msg string) string
//    func (this *Base64) Decode(data string) string
//    func IndexSlice(array, element interface{}) int
//    func RandSlice(array []interface}) []interface}
//    func ParseTime(strTime string) time.Duration
//    func ReadYaml(file string) map[string]interface}
//    func DumpsJson(data interface}, indent int) []byte
//    func LoadsJson(data []byte) interface}
//    func FormatMap(data interface}) interface}
//    func CurrentDir(level int) string

package tool

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
)

var outFile *os.File

func LogPrint() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	outFile.Close()
	log.SetOutput(os.Stderr)
}

func LogFile(stdout bool) {
	_, err := os.Stat("logs")
	if err != nil {
		os.Mkdir("logs", os.ModePerm)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	now := time.Now()
	logfile := fmt.Sprintf("logs/%v.log", now.Format("20060102"))
	outFile.Close()
	outFile, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	CheckErr(err)
	if stdout {
		multi := io.MultiWriter(os.Stdout, outFile)
		log.SetOutput(multi)
	} else {
		log.SetOutput(outFile)
	}
}

func CheckErr(err error) {
	if err != nil {
		log.SetFlags(log.Ldate | log.Ltime)

		_, path, line, _ := runtime.Caller(1)
		split := strings.Split(path, "/")
		msg := fmt.Sprintf("%v %v %v", split[len(split)-1], line, err)
		log.Println(msg)

		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		return
	}
}

func Request(url string, method string, reqD string, reqH map[string]string) string {

	req, err := http.NewRequest(method, url, strings.NewReader(reqD))
	CheckErr(err)

	for k, v := range reqH {
		req.Header.Set(k, v)
	}

	rsp, err := http.DefaultClient.Do(req)
	CheckErr(err)
	defer rsp.Body.Close()
	rspD, err := ioutil.ReadAll(rsp.Body)
	CheckErr(err)

	return string(rspD)
}

// Binary String to Int
func Btoi(s string) (int, error) {
	ans := 0
	base := 2
	for n := len(s) - 1; n >= 0; n-- {
		i := len(s) - 1 - n
		p, err := strconv.Atoi(string(s[n]))
		if err != nil {
			return 0, err
		}
		ans += int(math.Pow(float64(base), float64(i))) * p
	}
	return ans, nil
}

// Base64
type Base64 struct {
	table   string
	padding string
}

func Base64New() Base64 {
	obj := Base64{
		table:   "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/",
		padding: "=",
	}
	return obj
}

func (this *Base64) Set(table string, padding string) {
	this.table = table
	this.padding = padding
}

func (this *Base64) Encode(msg string) string {
	bin := ""
	code := ""
	for _, c := range msg {
		s := fmt.Sprintf("%b", c)
		s = strings.Repeat("0", 8-len(s)) + s
		bin += s
	}
	r := len(bin) % 6
	bin += strings.Repeat("0", 6-r)
	for i := 0; i < len(bin)/6; i++ {
		s := bin[6*i : 6*(i+1)]
		pos, _ := Btoi(s)
		code += string(this.table[pos])
	}
	code += strings.Repeat(this.padding, (6-r)/2)
	return code
}

func (this *Base64) Decode(data string) string {
	bin := ""
	msg := ""
	for _, n := range data {
		if string(n) == this.padding {
			bin = bin[0 : len(bin)-2]
		} else {
			index := IndexSlice(this.table, string(n))
			s := fmt.Sprintf("%b", index)
			s = strings.Repeat("0", 6-len(s)) + s
			bin += s
		}
	}
	msg = ""
	for i := 0; i < len(bin)/8; i++ {
		s := bin[8*i : 8*(i+1)]
		char, _ := Btoi(s)
		msg += string(char)
	}
	return msg
}

// Index of List
func IndexSlice(array, element interface{}) int {
	switch array.(type) {
	case string:
		data := array.(string)
		val := element.(string)
		for n, _ := range data {
			flag := true
			for m, _ := range val {
				if m+n > len(data) {
					return -1
				}
				if data[n+m] != val[m] {
					flag = false
					break
				}
			}
			if flag {
				return n
			}
		}
		return -1
	case []string:
		data := array.([]string)
		val := element.(string)
		for n, v := range data {
			if v == val {
				return n
			}
		}
		return -1
	case []int:
		data := array.([]int)
		val := element.(int)
		for n, v := range data {
			if v == val {
				return n
			}
		}
		return -1
	case []int64:
		data := array.([]int64)
		val := element.(int64)
		for n, v := range data {
			if v == val {
				return n
			}
		}
		return -1
	case []float64:
		data := array.([]float64)
		val := element.(float64)
		for n, v := range data {
			if v == val {
				return n
			}
		}
		return -1
	default:
		data := array.([]interface{})
		val := element.(interface{})
		for n, v := range data {
			if v == val {
				return n
			}
		}
		return -1
	}
}

// Random List
func RandSlice(array []interface{}) []interface{} {
	rand.Seed(time.Now().UnixNano())
	c := 4 * len(array)
	for i := 0; i < c; i++ {
		m := rand.Intn(len(array))
		n := rand.Intn(len(array))
		tmp := array[m]
		array[m] = array[n]
		array[n] = tmp
	}
	return array
}

// Change 1s, 1m to 1*time.Second, 1*time.Minute
func ParseTime(strTime string) time.Duration {
	var _time time.Duration
	num := len(strTime) - 1
	val, _ := strconv.ParseInt(string(strTime[:num]), 0, 64)
	str := string(strTime[num])
	switch str {
	case "s":
		_time = time.Duration(val) * time.Second
	case "m":
		_time = time.Duration(val) * time.Minute
	case "h":
		_time = time.Duration(val) * time.Hour
	case "d":
		_time = 24 * time.Duration(val) * time.Hour
	}
	return _time
}

// YAML
func ReadYaml(file string) map[string]interface{} {

	var yml []byte
	var err error
	var dict map[string]interface{}

	yml, err = ioutil.ReadFile(file)
	CheckErr(err)

	err = yaml.Unmarshal(yml, &dict)
	CheckErr(err)

	return FormatMap(dict).(map[string]interface{})
}

// Json
func DumpsJson(data interface{}, indent int) []byte {
	var jon []byte
	var err error
	if indent == 0 {
		jon, err = json.Marshal(FormatMap(data))
		CheckErr(err)
	} else {
		jon, err = json.MarshalIndent(FormatMap(data), "", strings.Repeat(" ", indent))
		CheckErr(err)
	}
	return jon
}

func LoadsJson(data []byte) interface{} {
	var dict interface{}
	err := json.Unmarshal(data, &dict)
	CheckErr(err)
	return dict
}

// Format All Map is map[string]interface{}
func FormatMap(data interface{}) interface{} {
	var rsp interface{}

	switch data.(type) {
	case map[interface{}]interface{}:
		dict := map[string]interface{}{}
		for k, v := range data.(map[interface{}]interface{}) {
			dict[fmt.Sprintf("%v", k)] = FormatMap(v)
		}
		rsp = dict
	case map[string]interface{}:
		dict := map[string]interface{}{}
		for k, v := range data.(map[string]interface{}) {
			dict[k] = FormatMap(v)
		}
		rsp = dict
	case []interface{}:
		array := []interface{}{}
		for _, v := range data.([]interface{}) {
			d := FormatMap(v)
			array = append(array, d)
		}
		rsp = array
	default:
		rsp = data
	}
	return rsp
}

func CurrentDir(level int) string {
	_, file, _, _ := runtime.Caller(level)
	dir := filepath.Dir(file)
	return dir
}
