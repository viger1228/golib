package tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
)

var outFile *os.File

func init() {
	LogPrint()
}

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

	return dict
}

// Json
func DumpJson(data interface{}, indent int) string {

	var jon []byte
	fmt.Printf("%T\n", data)
	jon, err := json.MarshalIndent(data, "", strings.Repeat(" ", indent))
	CheckErr(err)

	return string(jon)

}

// Format All submap is map[string]interface{}
func FormatMap(data interface{}) interface{} {
	switch i := data.(type) {
	case array:

	}
}
