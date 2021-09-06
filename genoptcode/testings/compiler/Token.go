package compiler

type Tag string

var TagVar Tag = "var"
var TagIf Tag = "if"
var TagThen Tag = "then"
var TagElse Tag = "else"
var TagAssign Tag = "="
var TagPlus Tag = "+"
var TagID Tag = "id"
var TagNumber Tag = "number"
var TagSC Tag = ";"
var TagMore Tag = ">"
var TagLess Tag = "<"

type Token struct {
	Tag   Tag
	Value string
}
