package action

import (
	"reflect"
	"runtime"
)

const ERR_INVALID_PARAMS = 1001 // invalid params
const ERR_SYS_ERROR = 5001 // system error

var errorMap map[int]string = make(map[int]string)

func init() {
	errorMap[ERR_INVALID_PARAMS] = "invalid-params"
}

func ErrorMessage(code int) string {
	return errorMap[code]
}

//-----------

type ActionParam map[string]string

type ActionResult struct {
	Meta ActionResultMeta `json:"meta"`
	Data interface{}      `json:"data"`
}

type ActionResultMeta struct {
	Status       bool   `json:"status"`
	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func NewResult(data interface{}) ActionResult {
	return ActionResult{
		Meta: ActionResultMeta{
			Status: true,
		},
		Data: data,
	}
}

func NewResultError(errorCode int) ActionResult {
	return ActionResult{
		Meta: ActionResultMeta{
			Status:       false,
			ErrorCode:    errorCode,
			ErrorMessage: ErrorMessage(errorCode),
		},
		Data: nil,
	}
}

func NewSystemError(err error) ActionResult {
    return ActionResult{
        Meta:ActionResultMeta{
            Status:false,
            ErrorCode:ERR_SYS_ERROR,
            ErrorMessage:err.Error(),
        },
        Data:nil,
    }
}

//------------

type ActionFunc func(ActionParam) ActionResult
type ActionBeforeFunc func(*ActionParam)
type ActionAfterFunc func(*ActionResult)

var (
	beforeFuncs = make(map[string][]ActionBeforeFunc)
	afterFuncs  = make(map[string][]ActionAfterFunc)
)

func Before(fn ActionFunc, before ActionBeforeFunc) {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	if len(beforeFuncs[name]) == 0 {
		beforeFuncs[name] = []ActionBeforeFunc{before}
		return
	}
	beforeFuncs[name] = append(beforeFuncs[name], before)
}

func After(fn ActionFunc, after ActionAfterFunc) {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	if len(afterFuncs[name]) == 0 {
		afterFuncs[name] = []ActionAfterFunc{after}
		return
	}
	afterFuncs[name] = append(afterFuncs[name], after)
}

func Call(fn ActionFunc, params ActionParam) ActionResult {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	befores := beforeFuncs[name]
	if len(befores) > 0 {
		for _, f := range befores {
			f(&params)
		}
	}
	result := fn(params)
	afters := afterFuncs[name]
	if len(afters) > 0 {
		for _, f := range afters {
			f(&result)
		}
	}
	return result
}
