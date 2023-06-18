// 版权 @2023 凹语言 作者。保留所有权利。

%{
// 这是 凹语言 yacc 的例子, 用于对表达式进行解析, 为了简化词法部分暂时通过手工录入.
%}

%union {
	num :int
}

%type  <num> expr expr1 expr2 expr3
%token '+' '-' '*' '/' '(' ')'
%token <num> NUM

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
	Kind  :int
	Value :int
}

type exprLexer struct {
	tokens :[]exprToken
	pos    :int 
}

func exprLexer.Lex(yylval: *exprSymType) => int {
	if this.pos >= len(this.tokens) {
		return eof
	}
	tok := this.tokens[this.pos]
	this.pos++

	yylval.num = tok.Value
	return tok.Kind
}

func exprLexer.Error(s: string) {
	println("ERROR:", s)
}

func main {
	print("1+2*(3+4)-10 = ")
	exprParse(&exprLexer{
		tokens: []exprToken{
			{Kind: NUM, Value: 1},
			{Kind: '+'},
			{Kind: NUM, Value: 2},
			{Kind: '*'},
			{Kind: '('},
			{Kind: NUM, Value: 3},
			{Kind: '+'},
			{Kind: NUM, Value: 4},
			{Kind: ')'},
			{Kind: '-'},
			{Kind: NUM, Value: 10},
		},
	})
}
