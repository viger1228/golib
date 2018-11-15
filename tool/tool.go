package tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func LogFile() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	now := time.Now()
	logfile := fmt.Sprintf("logs/%v.log", now.Format("20060102"))
	outFile.Close()
	outFile, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	CheckErr(err)
	log.SetOutput(outFile)
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
