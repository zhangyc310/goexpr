package goexpr

import (
	"math/rand"
)

// Rand returns a random float64 between 0 (inclusive) and 1 (exclusive).
func Rand() Expr {
	return &randExpr{}
}

type randExpr struct {
}

func (e *randExpr) Eval(params Params) interface{} {
	return rand.Float64()
}

func (e *randExpr) String() string {
	return "RAND"
}
