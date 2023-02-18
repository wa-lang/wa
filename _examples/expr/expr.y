// 版权 @2023 凹语言 作者。保留所有权利。

%{
package main
%}

%union {
	num int
}

%type	<num>	expr expr1 expr2 expr3

%token '+' '-' '*' '/' '(' ')'

%token	<num>	NUM

%%

top:
	expr {
		println($1)
	}

expr:
	expr1
	| '+' expr {
		$$ = $2
	}
	| '-' expr {
		$$ = -$2
	}

expr1:
	expr2
	| expr1 '+' expr2 {
		$$ = $1 + $3
	}
	| expr1 '-' expr2 {
		$$ = $1 - $3
	}

expr2:
	expr3
	| expr2 '*' expr3 {
		$$ = $1 * $3
	}
	| expr2 '/' expr3 {
		$$ = $1 / $3
	}

expr3:
	NUM
	| '(' expr ')' {
		$$ = $2
	}


%%

// Lex 结束标志
const eof = 0

type exprToken struct {
	Kind  int
	Value int
}

type exprLexer struct {
	tokens []exprToken
	pos    int 
}

func (p *exprLexer) Lex(yylval *exprSymType) int {
	if p.pos >= len(p.tokens) {
		return eof
	}
	tok := p.tokens[p.pos]
	p.pos++

	yylval.num = tok.Value
	return tok.Kind
}

func (x *exprLexer) Error(s string) {
	println("ERROR:", s)
}

func main() {
	exprParse(&exprLexer{
		tokens: []exprToken{
			// 1+2*3
			{Kind: NUM, Value: 1},
			{Kind: '+'},
			{Kind: NUM, Value: 2},
			{Kind: '*'},
			{Kind: NUM, Value: 3},
		},
	})
}
