# 凹语言版本 yacc 简介 - 以表达式解析为例

yacc 是用于生成语法解析器的程序，是编译器爱好者的工具。凹语言的 yacc 从 goyacc 移植而来，目前可以初步支持输出 凹语言 版本解析器代码。本文以以表达式解析为例展示下用法。

## 1. yacc 文件结构简介

yacc 文件一般以 `*.y` 格式命名，其格式如下：

```
// *.y 文件本身的注释

%{
// 生成解析器代码的头部，一般是 import 等语句
%}


// yacc 语法对应的词法类似、语法树节点等

%%

// BNF 语法定义

%%

// 生成解析器代码的尾部
```

简单来说，y 文件由两个 `%%` 行分隔为三个部分（类似文章的凤头、猪肚、豹尾）：
- 凤头：对应生成的解析器的头部，其中`%{ ... %}` 包含的部分为原样输出，其他部分是 yacc 语法定义的词法类型和语法树节点等
- 猪肚：是 yacc 文件等核心，通过 BNF 语法定义了语法结构，这里主要是针对 `LALR(1)` 语法
- 豹尾：如果是独立的程序，可以在这个部分引入词法解析器和 main 函数；如果是 package 则是可以省略的

## 2. 定义`expr.y`文件 - 凤头部分

创建表达式语法文件如下：

```yacc
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
```

其中`%union`定义了词法和语法解析器对接的结构体类型。`%type`语句定义了`expr expr1 expr2 expr3`几种语法节点，都是对应`<num>`类型值，而数字的值需要填充到`%union`定义的`num`属性部分。`%token`语句定义的运算符和`NUM`类型的数字。

## 3. 定义`expr.y`文件 - 猪肚部分

猪肚部分对应表达式的语法结构：

```
%%

top:
	expr { println($1) }

expr:
	expr1
	| '+' expr { $$ = $2 }
	| '-' expr { $$ = -$2 }

expr1:
	expr2
	| expr1 '+' expr2 { $$ = $1 + $3 }
	| expr1 '-' expr2 { $$ = $1 - $3 }

expr2:
	expr3
	| expr2 '*' expr3 { $$ = $1 * $3 }
	| expr2 '/' expr3 { $$ = $1 / $3 }

expr3:
	NUM
	| '(' expr ')' { $$ = $2 }

%%
```

当遇到`expr`语法规则是直接输出结果，`expr1`表示加减法、`expr2`表示乘除法、`expr3`表示数字或小括弧。在每个最终后面的`{}`中包含的是动作代码，它们根据不同的语法规则选择不同的计算方式得到结果，结果赋值给`$$`（也就是对应`%type  <num> expr expr1 expr2 expr3`语句中的`<num>`部分类型，也对应`%union`定义的`num`成员）。

## 4. 定义`expr.y`文件 - 豹尾部分 - 01

有了凤头和猪肚部分，yacc就可以生成必要的解析器代码了。默认后生成以下格式的解析器函数`yyParse`：

```
func yyParse(yylex: *yyLexer) => int {
	return yyNewParser().Parse(yylex)
}
```

而`yyLexer`词法解析器则是用户需要自行实现的（词法解析实现相对简单），主要包含以下2个方法：

```
type yyLexer struct {}

func yyLexer.Lex(yylval *yySymType) => int {
	// 返回 Token 类型, 并且将对应的值填充到 yylval 相应的属性中
}

func yyLexer.Error(s string) {
	// 遇到错误
}
```

`yyLexer.Lex` 返回 Token 类型，并且将对应的值填充到 `yylval` 相应的属性中，遇到文件结尾时返回`0`表示文件结束。方法参数对应的`yySymType`类型由yacc工具生成，对应如下的代码：

```
type yySymType struct {
	yys :int
	num :int
}
```

其中`num`对应对应`%union`定义的属性，也就是数字的值。

## 5. 定义`expr.y`文件 - 豹尾部分 - 02

为了简化演示代码，我们先手工构造词法序列，然后通过`yyLexer.Lex` 返回。

```
// Lex 结束标志
const eof = 0

type yyToken struct {
	Kind  :int
	Value :int
}

type yyLexer struct {
	tokens :[]yyToken
	pos    :int 
}

func yyLexer.Lex(yylval *yySymType) => int {
	if this.pos >= len(this.tokens) {
		return eof
	}
	tok := this.tokens[this.pos]
	this.pos++

	yylval.num = tok.Value
	return tok.Kind
}

func yyLexer.Error(s string) {
	println("ERROR:", s)
}
```

首先定义`yyToken`，对应token的类型和值信息。然后`yyLexer`定义全部的token列表和当前的pos信息。`yyLexer.Lex`方法每次从`this.tokens`列表对应的`this.pos`位置返回一个token，如果是结束则返回`eof`。

然后就可以构造main函数启动了：

```
func main {
	print("1+2*3 = ")
	yyParse(&yyLexer{
		tokens: []exprToken{
			{Kind: NUM, Value: 1},
			{Kind: '+'},
			{Kind: NUM, Value: 2},
			{Kind: '*'},
			{Kind: NUM, Value: 3},
		},
	})
}
```

## 6. 生成解析器代码

在生成解析器代码前再准备一个`copyright.txt`文件，比如“保留所有权利”或者“自由使用”之类的。然后通过以下命令生成解析器代码：

```
$ wa yacc -l -p=yy -c="copyright.txt" -o="y.wa" expr.y
```

其中`-l`表示生成的代码禁止映射到`*.y`文件行列号（用生成代码的位置），`-p=yy`表示生成的解析器函数和类型等用`yy`前缀（这也是默认值），`-c="copyright.txt"`为生成代码指定版权信息，`-o="y.wa"`指定输出文件，最后的`expr.y`对熟人的yacc规则文件。

生成代码成功之后可以执行：

```
$ wa y.wa
1+2*3 = 7
```

完整的例子可以参考（这里使用的是`expr`前缀）：https://gitee.com/wa-lang/wa/blob/master/_examples/expr/expr.y

## 7. 下一步

目前的凹语言版 yacc 工具还是Go语言实现的，只是输出的解析器是凹语言代码。我们希望下一步可以将 yacc 工具本身移植到凹语言实现，最终可以通过 wasm 模块执行。

