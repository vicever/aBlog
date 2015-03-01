package action

import (
	"reflect"
	"runtime"
)

const ERR_INVALID_PARAMS = 1001 // invalid params
const ERR_SYS_ERROR = 5001      // system error

// error codes and messages
var errorMap map[int]string = make(map[int]string)

func init() {
	errorMap[ERR_INVALID_PARAMS] = "invalid-params"
}

// get error message by code
func ErrorMessage(code int) string {
	return errorMap[code]
}

//-----------

// action input params
type ActionParam map[string]string

// action returning result
type ActionResult struct {
	Meta ActionResultMeta `json:"meta"` // meta data shows the result is ok or failed with error
	Data interface{}      `json:"data"` // returning data
}

// action result meta
type ActionResultMeta struct {
	Status       bool   `json:"status"`
	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// new correct result with data
func NewResult(data interface{}) ActionResult {
	return ActionResult{
		Meta: ActionResultMeta{
			Status: true,
		},
		Data: data,
	}
}

// new failed result with error code
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

// new system error with error value. it's used to returning core error in basic libraries, not logic processes
func NewSystemError(err error) ActionResult {
	return ActionResult{
		Meta: ActionResultMeta{
			Status:       false,
			ErrorCode:    ERR_SYS_ERROR,
			ErrorMessage: err.Error(),
		},
		Data: nil,
	}
}

//------------

// action function, solve input params to result
type ActionFunc func(ActionParam) ActionResult

// action before filter, to change params
type ActionBeforeFunc func(*ActionParam)

// action after filter, to change result
type ActionAfterFunc func(*ActionResult)

var (
	beforeFuncs = make(map[string][]ActionBeforeFunc)
	afterFuncs  = make(map[string][]ActionAfterFunc)
)

// set before filter to action func
func Before(fn ActionFunc, before ActionBeforeFunc) {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	if len(beforeFuncs[name]) == 0 {
		beforeFuncs[name] = []ActionBeforeFunc{before}
		return
	}
	beforeFuncs[name] = append(beforeFuncs[name], before)
}

// set after filter to action func
func After(fn ActionFunc, after ActionAfterFunc) {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	if len(afterFuncs[name]) == 0 {
		afterFuncs[name] = []ActionAfterFunc{after}
		return
	}
	afterFuncs[name] = append(afterFuncs[name], after)
}

// call action, run filters together. before -> action -> after, then return result
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
