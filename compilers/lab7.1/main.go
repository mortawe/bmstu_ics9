package main

import "fmt"

// A = $AXIOM nterm_list
// N = $NTERM nterm_list
// T = $TERM term_list
// R = $RULE nterm = list
// nterm_list = nterm_list nterm | nterm
// term_list = term_list term | term
// term = \" symbol \"
// nterm = symbol | symbol\'
// list = list \n list | term list | nterm list | $EPS

func main() {
	tokens, err := lex(`$AXIOM E 
$NTERM E' T T' F
$TERM "+" "*" "(" ")" "n"
* asfafas
$RULE E = T E'
$RULE E' = "+" T E'
$EPS
$RULE T = F T'
$RULE T' = "*" F T'
$EPS
$RULE F = "n"
"(" E ")"`)
	fmt.Println(err)
	for _, token := range tokens {
		fmt.Println(token.String())
	}
}
