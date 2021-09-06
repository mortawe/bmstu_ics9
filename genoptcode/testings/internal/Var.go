package internal

import (
	"fmt"
	"strings"
)

type Var struct {
	name     string
	version  int
	sign     string
	rightPar string
}

func NewVar(name string, sign string, version int) *Var {
	v := Var{}
	if strings.Contains(name, ")") {
		v.rightPar = ")+"
		v.name = strings.ReplaceAll(name, ")", "")
	} else {
		v.rightPar = ""
		v.name = name
	}
	v.version = version
	v.sign = sign

	return &v
}

func (v *Var) Equals(o Var) bool {
	return o.name == v.name
}

func (v *Var) String() string {
	return fmt.Sprintf("name : %s, version : %d, rightPar : %s, sign : %s", v.name, v.version, v.rightPar, v.sign)
}
