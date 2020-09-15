package errs

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/fatih/color"

	"github.com/spf13/cast"
)

const (
	ErrServExists    = 461
	ErrServNotExist  = 462
	ErrDsExists      = 466
	ErrDsNotExist    = 467
	ErrParamBad      = 496
	ErrParamNotFound = 497
	ErrOther         = 499
)

var ErrorStyleFunc = color.New(color.FgHiWhite, color.BgRed).SprintFunc()

var errors = map[int]string{
	// Service related errors
	ErrServExists:   "The service already exists",
	ErrServNotExist: "Service does not exist",

	// DataSource related errors
	ErrDsExists:   "The datasource already exists",
	ErrDsNotExist: "DataSource does not exist",

	// Parameter related errors
	ErrParamBad:      "Bad parameter %s",
	ErrParamNotFound: "Parameter %s not found",
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Cause   string `json:"cause,omitempty"`
	Index   uint64 `json:"index"`
}

func New(code int, msg ...interface{}) *Error {
	err := &Error{
		Code:    code,
		Message: errors[code],
	}
	if msg != nil {
		err.Message = fmt.Sprintf(errors[code], msg...)
	}
	return err
}

func NewError(err error) *Error {
	return &Error{
		Code:    ErrOther,
		Message: err.Error(),
	}
}

// Error is for the error interface
func (e Error) Error() string {
	return e.Message + " (" + e.Cause + ")"
}

func (e Error) toJsonString() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func Join(elems []interface{}, sep string) string {
	var buffer bytes.Buffer
	for i, elem := range elems {
		if i > 0 {
			buffer.WriteString(sep)
		}
		buffer.WriteString(cast.ToString(elem))
	}
	return buffer.String()
}
