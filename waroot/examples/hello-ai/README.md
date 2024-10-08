# 凹语言试验性 AI 示例

先参考 [Device Model: Chrome AI Gemini Nano](https://chrome-ai.edgeone.app/), 确保本地 Chrome AI 正常工作, 确保开发者控制台可成功执行 `await ai.assistant.create();` 程序。

示例代码:

```wa
import "ai"

func main {
	ai.RequestSession(func(session: ai.Session){
		session.PromptAsync("Who are you?", func(res: string) {
			println(res)
		})
	})
}
```

然后本地命令行环境执行`wa run`, 然后在打开的页面的开发者控制台可以看到以下输出:

```
 I am a large language model, trained by Google.
```


<!--
- https://github.com/lightning-joyce/chromeai

BUGS:

- [The model took too long too many times for this version.](https://issues.chromium.org/issues/356649889)
-->

