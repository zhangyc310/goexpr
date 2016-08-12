package goexpr

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	zeroTime = time.Time{}
)

func Binary(operator string, left Expr, right Expr) (Expr, error) {
	// Normalize equals and not equals
	switch operator {
	case "=":
		operator = "=="
	case "<>":
		operator = "!="
	}
	o, found := ops[operator]
	if !found {
		return nil, fmt.Errorf("No op for %v", operator)
	}
	// Fill in the blanks
	return &binaryExpr{operator, buildOpFN(o), left, right}, nil
}

type opFN func(a interface{}, b interface{}) interface{}

type binaryExpr struct {
	operatorString string
	operator       opFN
	left           Expr
	right          Expr
}

func (e *binaryExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", e.left, e.operatorString, e.right)
}

func (e *binaryExpr) Eval(params Params) interface{} {
	return e.operator(e.left.Eval(params), e.right.Eval(params))
}

func buildOpFN(o op) opFN {
	return func(a interface{}, b interface{}) interface{} {
		return applyOp(o, a, b)
	}
}

type coercedTo int

var (
	coerceNotAllowed = coercedTo(0)
	coercedToNil     = coercedTo(1)
	coercedToBool    = coercedTo(2)
	coercedToUInt    = coercedTo(3)
	coercedToInt     = coercedTo(4)
	coercedToFloat   = coercedTo(5)
	coercedToString  = coercedTo(6)
	coercedToTime    = coercedTo(7)
)

type op struct {
	nl   func(a interface{}, b interface{}) interface{}
	bl   func(a bool, b bool) interface{}
	uin  func(a uint64, b uint64) interface{}
	sin  func(a int64, b int64) interface{}
	fl   func(a float64, b float64) interface{}
	st   func(a string, b string) interface{}
	ts   func(a time.Time, b time.Time) interface{}
	dflt interface{}
}

var ops = map[string]op{
	"==": op{
		nl: func(a interface{}, b interface{}) interface{} {
			return a == nil && b == nil
		},
		bl: func(a bool, b bool) interface{} {
			return a == b
		},
		uin: func(a uint64, b uint64) interface{} {
			return a == b
		},
		sin: func(a int64, b int64) interface{} {
			return a == b
		},
		fl: func(a float64, b float64) interface{} {
			return a == b
		},
		st: func(a string, b string) interface{} {
			return a == b
		},
		ts: func(a time.Time, b time.Time) interface{} {
			return a == b
		},
		dflt: false,
	},
	"LIKE": op{
		nl: func(a interface{}, b interface{}) interface{} {
			return a == nil && b == nil
		},
		bl: func(a bool, b bool) interface{} {
			return a == b
		},
		uin: func(a uint64, b uint64) interface{} {
			return a == b
		},
		sin: func(a int64, b int64) interface{} {
			return a == b
		},
		fl: func(a float64, b float64) interface{} {
			return a == b
		},
		st: func(a string, b string) interface{} {
			lb := len(b)
			last := lb - 1
			startWildcard := b[0] == '%'
			endWildcard := b[last] == '%'
			if endWildcard {
				b = b[:last]
			}
			if startWildcard {
				b = b[1:]
			}
			if lb == 0 {
				return true
			}
			if startWildcard && endWildcard {
				return strings.Contains(a, b)
			}
			la := len(a)
			if la < lb {
				return false
			}
			if endWildcard {
				return a[:lb] == b
			}
			return a[lb-la:] == b
		},
		ts: func(a time.Time, b time.Time) interface{} {
			return a == b
		},
		dflt: false,
	},
	"!=": op{
		nl: func(a interface{}, b interface{}) interface{} {
			return a == nil && b != nil || a != nil && b == nil
		},
		bl: func(a bool, b bool) interface{} {
			return a != b
		},
		uin: func(a uint64, b uint64) interface{} {
			return a != b
		},
		sin: func(a int64, b int64) interface{} {
			return a != b
		},
		fl: func(a float64, b float64) interface{} {
			return a != b
		},
		st: func(a string, b string) interface{} {
			return a != b
		},
		ts: func(a time.Time, b time.Time) interface{} {
			return a != b
		},
		dflt: false,
	},
	"<": op{
		nl: func(a interface{}, b interface{}) interface{} {
			return a == nil && b != nil
		},
		bl: func(a bool, b bool) interface{} {
			return !a && b
		},
		uin: func(a uint64, b uint64) interface{} {
			return a < b
		},
		sin: func(a int64, b int64) interface{} {
			return a < b
		},
		fl: func(a float64, b float64) interface{} {
			return a < b
		},
		st: func(a string, b string) interface{} {
			return a < b
		},
		ts: func(a time.Time, b time.Time) interface{} {
			return a.Before(b)
		},
		dflt: false,
	},
	"<=": op{
		nl: func(a interface{}, b interface{}) interface{} {
			return a == nil
		},
		bl: func(a bool, b bool) interface{} {
			return true
		},
		uin: func(a uint64, b uint64) interface{} {
			return a <= b
		},
		sin: func(a int64, b int64) interface{} {
			return a <= b
		},
		fl: func(a float64, b float64) interface{} {
			return a <= b
		},
		st: func(a string, b string) interface{} {
			return a <= b
		},
		ts: func(a time.Time, b time.Time) interface{} {
			return !a.After(b)
		},
		dflt: false,
	},
	">": op{
		nl: func(a interface{}, b interface{}) interface{} {
			return a != nil && b == nil
		},
		bl: func(a bool, b bool) interface{} {
			return a && !b
		},
		uin: func(a uint64, b uint64) interface{} {
			return a > b
		},
		sin: func(a int64, b int64) interface{} {
			return a > b
		},
		fl: func(a float64, b float64) interface{} {
			return a > b
		},
		st: func(a string, b string) interface{} {
			return a > b
		},
		ts: func(a time.Time, b time.Time) interface{} {
			return a.After(b)
		},
		dflt: false,
	},
	">=": op{
		nl: func(a interface{}, b interface{}) interface{} {
			return b == nil
		},
		bl: func(a bool, b bool) interface{} {
			return true
		},
		uin: func(a uint64, b uint64) interface{} {
			return a >= b
		},
		sin: func(a int64, b int64) interface{} {
			return a >= b
		},
		fl: func(a float64, b float64) interface{} {
			return a >= b
		},
		st: func(a string, b string) interface{} {
			return a >= b
		},
		ts: func(a time.Time, b time.Time) interface{} {
			return !a.Before(b)
		},
		dflt: false,
	},
	"+": op{
		nl: func(a interface{}, b interface{}) interface{} {
			return nil
		},
		bl: func(a bool, b bool) interface{} {
			return nil
		},
		uin: func(a uint64, b uint64) interface{} {
			return a + b
		},
		sin: func(a int64, b int64) interface{} {
			return a + b
		},
		fl: func(a float64, b float64) interface{} {
			return a + b
		},
		st: func(a string, b string) interface{} {
			return nil
		},
		ts: func(a time.Time, b time.Time) interface{} {
			return nil
		},
		dflt: nil,
	},
	"-": op{
		nl: func(a interface{}, b interface{}) interface{} {
			return nil
		},
		bl: func(a bool, b bool) interface{} {
			return nil
		},
		uin: func(a uint64, b uint64) interface{} {
			return a - b
		},
		sin: func(a int64, b int64) interface{} {
			return a - b
		},
		fl: func(a float64, b float64) interface{} {
			return a - b
		},
		st: func(a string, b string) interface{} {
			return nil
		},
		ts: func(a time.Time, b time.Time) interface{} {
			return nil
		},
		dflt: nil,
	},
	"*": op{
		nl: func(a interface{}, b interface{}) interface{} {
			return nil
		},
		bl: func(a bool, b bool) interface{} {
			return nil
		},
		uin: func(a uint64, b uint64) interface{} {
			return a * b
		},
		sin: func(a int64, b int64) interface{} {
			return a * b
		},
		fl: func(a float64, b float64) interface{} {
			return a * b
		},
		st: func(a string, b string) interface{} {
			return nil
		},
		ts: func(a time.Time, b time.Time) interface{} {
			return nil
		},
		dflt: nil,
	},
	"/": op{
		nl: func(a interface{}, b interface{}) interface{} {
			return nil
		},
		bl: func(a bool, b bool) interface{} {
			return nil
		},
		uin: func(a uint64, b uint64) interface{} {
			if b == 0 {
				return nil
			}
			return a / b
		},
		sin: func(a int64, b int64) interface{} {
			if b == 0 {
				return nil
			}
			return a / b
		},
		fl: func(a float64, b float64) interface{} {
			if b == 0 {
				return nil
			}
			return a / b
		},
		st: func(a string, b string) interface{} {
			return nil
		},
		ts: func(a time.Time, b time.Time) interface{} {
			return nil
		},
		dflt: nil,
	},
}

func strToBool(str string) (bool, bool) {
	if str == "" {
		return false, true
	}
	r, err := strconv.ParseBool(str)
	return r, err == nil
}

func strToFloat(str string) (float64, bool) {
	r, err := strconv.ParseFloat(str, 64)
	return r, err == nil
}

func boolToUInt(v bool) uint64 {
	if v {
		return 1
	}
	return 0
}

func boolToInt(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

func boolToFloat(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

func boolToString(v bool) string {
	return strconv.FormatBool(v)
}

func floatToString(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func div(x interface{}, y interface{}) interface{} {
	switch v1 := x.(type) {
	case uint64:
		v2 := y.(uint64)
		if y == 0 {
			return 0
		}
		return v1 + v2
	case int64:
		v2 := y.(int64)
		if y == 0 {
			return 0
		}
		return v1 + v2
	case float64:
		v2 := y.(float64)
		if y == 0 {
			return 0
		}
		return v1 + v2
	}
	return 0
}

// func coerce(a interface{}, b interface{}) (c coercedTo, a interface{}, b interface{}) {
// 	switch x := a.(type) {
// 	case nil:
// 		return coercedToNil, a, b
// 	case bool:
// 		switch y := b.(type) {
// 		case nil:
// 			return coercedToNil, a, b
// 		case bool:
// 			return coercedToBool, x, y
// 		case byte:
// 			return coercedToUInt, boolToUInt(x), uint64(y)
// 		case uint16:
// 			return coercedToUInt, boolToUInt(x), uint64(y)
// 		case uint32:
// 			return coercedToUInt, boolToUInt(x), uint64(y)
// 		case uint64:
// 			return coercedToUInt, boolToUInt(x), uint64(y)
// 		case uint:
// 			return coercedToUInt, boolToUInt(x), uint64(y)
// 		case int8:
// 			return coercedToInt, boolToInt(x), uint64(y)
// 		case int16:
// 			return coercedToInt, boolToInt(x), uint64(y)
// 		case int32:
// 			return coercedToInt, boolToInt(x), uint64(y)
// 		case int64:
// 			return coercedToInt, boolToInt(x), uint64(y)
// 		case int:
// 			return coercedToInt, boolToInt(x), uint64(y)
// 		case float32:
// 			return coercedToFloat, boolToFloat(x), uint64(y)
// 		case float64:
// 			return coercedToFloat, boolToFloat(x), uint64(y)
// 		case string:
// 			yb, ok := strToBool(y)
// 			if ok {
// 				return coercedToBool, x, yb
// 			}
// 		}
// 	case byte:
// 		switch y := b.(type) {
// 		case nil:
// 			return coercedToNil, a, b
// 		case bool:
// 			return coercedToBool, x, y
// 		case byte:
// 			return coercedToUInt, boolToUInt(x), uint64(y)
// 		case uint16:
// 			return coercedToUInt, boolToUInt(x), uint64(y)
// 		case uint32:
// 			return coercedToUInt, boolToUInt(x), uint64(y)
// 		case uint64:
// 			return coercedToUInt, boolToUInt(x), uint64(y)
// 		case uint:
// 			return coercedToUInt, boolToUInt(x), uint64(y)
// 		case int8:
// 			return coercedToInt, boolToInt(x), uint64(y)
// 		case int16:
// 			return coercedToInt, boolToInt(x), uint64(y)
// 		case int32:
// 			return coercedToInt, boolToInt(x), uint64(y)
// 		case int64:
// 			return coercedToInt, boolToInt(x), uint64(y)
// 		case int:
// 			return coercedToInt, boolToInt(x), uint64(y)
// 		case float32:
// 			return coercedToFloat, boolToFloat(x), uint64(y)
// 		case float64:
// 			return coercedToFloat, boolToFloat(x), uint64(y)
// 		case string:
// 			yb, ok := strToBool(y)
// 			if ok {
// 				return coercedToBool, x, yb
// 			}
// 		}
// 	}
//
// 	return coerceNotAllowed, a, b
// }
