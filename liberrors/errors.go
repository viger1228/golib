package liberrors

import (
	"context"
	"fmt"
	"runtime"
	"strings"
)

type errorObject struct {
	errCtx     context.Context // 保存錯誤發生當下訊息
	errContent string          // 錯誤內容
	errCode    int             // 錯誤碼(返回用戶端使用)
	errMessage string          // 錯誤訊息(返回用戶端使用)
}

func (e *errorObject) Error() string {
	return e.errContent
}

// ErrorFormat 错误链
func (e *errorObject) ErrorFormat(currentFunc string) string {
	format := ""
	if funcLinks, ok := e.errCtx.Value("funcLinks").([]string); ok {
		flag := false
		for i := len(funcLinks) - 1; i >= 0; i-- {
			link := funcLinks[i]
			if !flag && link == currentFunc {
				flag = true
			}
			if flag {
				split := strings.Split(link, ".")
				name := split[len(split)-1]
				format += fmt.Sprintf("[%s]", name)
			}
		}
	}
	if format != "" {
		format += ": %v"
	} else {
		format = "%v"
	}
	return format
}

func New(ctx context.Context, text string) error {
	funcLinks := getLinks(3)
	ctx = context.WithValue(ctx, "funcLinks", funcLinks)
	return &errorObject{
		errCtx:     ctx,
		errContent: text,
	}
}

// NewWithMessage 包含返回用戶端的錯誤碼及訊息
func NewWithMessage(ctx context.Context, text string, code int, message string) error {
	funcLinks := getLinks(3)
	ctx = context.WithValue(ctx, "funcLinks", funcLinks)
	return &errorObject{
		errCtx:     ctx,
		errContent: text,
		errCode:    code,
		errMessage: message,
	}
}

// Load 载入一般error
func Load(ctx context.Context, err error) error {
	funcLinks := getLinks(3)
	ctx = context.WithValue(ctx, "funcLinks", funcLinks)
	return &errorObject{
		errCtx:     ctx,
		errContent: err.Error(),
	}
}

// LoadWithMessage 载入一般error，包含返回用戶端的Error Code及訊息
func LoadWithMessage(ctx context.Context, err error, code int, message string) error {
	funcLinks := getLinks(3)
	ctx = context.WithValue(ctx, "funcLinks", funcLinks)
	return &errorObject{
		errCtx:     ctx,
		errContent: err.Error(),
		errCode:    code,
		errMessage: message,
	}
}

// ErrorFormat 錯誤格式化
func ErrorFormat(err interface{}) (format string) {
	format = "%v"
	currentFunc, _ := getFuncName(2)
	if errObject, ok := err.(*errorObject); ok {
		format = errObject.ErrorFormat(currentFunc)
	}
	return format
}

// ErrorCode 錯誤碼
func ErrorCode(err interface{}) (code int) {
	code = 0
	if errObject, ok := err.(*errorObject); ok {
		code = errObject.errCode
	}
	return code
}

// ErrorMessage 錯誤訊息(返回客戶端用)
func ErrorMessage(err interface{}) (message string) {
	message = ""
	if errObject, ok := err.(*errorObject); ok {
		message = errObject.errMessage
	}
	return message
}

// getLinks 获取函数链
func getLinks(start int) (parentFunction []string) {
	for i := start; i <= 100; i++ {
		name, ok := getFuncName(i)
		if !ok {
			break
		}
		parentFunction = append(parentFunction, name)
	}
	return parentFunction
}

// getFuncName 获取层级函数名
func getFuncName(skip int) (string, bool) {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "", false
	}
	f := runtime.FuncForPC(pc)
	name := f.Name()
	return name, true
}
