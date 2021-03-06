// Package goexpr provides basic expression evaluation of Go. It supports
// values of the following types: bool, byte, uint16, uint32, uint64, int8,
// int16, int32, int64, int, float32, float64, string and time.Time.
package goexpr

import (
	"fmt"
)

type Params interface {
	Get(key string) interface{}
}

type Expr interface {
	Eval(Params) interface{}

	String() string
}

func Param(name string) Expr {
	return &param{name}
}

type param struct {
	name string
}

func (e *param) Eval(params Params) interface{} {
	return params.Get(e.name)
}

func (e *param) String() string {
	return e.name
}

func Constant(val interface{}) Expr {
	return &constant{val}
}

type constant struct {
	val interface{}
}

func (e *constant) Eval(params Params) interface{} {
	return e.val
}

func (e *constant) String() string {
	return fmt.Sprint(e.val)
}

func Not(wrapped Expr) Expr {
	return &notExpr{wrapped}
}

type notExpr struct {
	wrapped Expr
}

func (e *notExpr) Eval(params Params) interface{} {
	return !e.wrapped.Eval(params).(bool)
}

func (e *notExpr) String() string {
	return fmt.Sprintf("NOT %v", e.wrapped)
}

type MapParams map[string]interface{}

func (p MapParams) Get(name string) interface{} {
	return p[name]
}
