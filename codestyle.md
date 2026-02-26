# 代码风格(仅参考)

为了便于调试，不同的模块可以在 config 包增加 Trace 开关，然后针对每个返回 err 的地方输出日志。
比如 `config.EnableTrace_loader` 标注 loader 模块：

```go
	logger.Trace(&config.EnableTrace_loader, "import "+manifest.MainPkg)
	if _, err := p.Import(manifest.MainPkg); err != nil {
		logger.Tracef(&config.EnableTrace_loader, "err: %v", err)
		return nil, err
	}
```
