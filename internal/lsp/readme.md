# LSP 测试

凹语言的 LSP 只支持根工作区有 wa.mod 的场景，目前还在开发的早期阶段。

## Sublime Text 语法高亮

TODO

## Sublime Text LSP 支持

以下是在 Sublime Text 测试的说明：

- 安装 Sublime Text, 打开命令面板
- 安装 Package Control 包管理插件
- 通过包管理安装 LSP 包: Package Control: Install Package, 选择 LSP
- 通过包管理查看 LSP 包已经成功安装: Package Control: List Packages, 确认 LSP 在列表中
- 命令面板选择 LSP: Enable Language Server Globally 全局范围打开 LSP 插件
- 菜单: Preferences -> Package Settings -> LSP -> Settings, 添加配置(参考后面的例子)

Sublime Text 的 凹语言 LSP 配置:

```json
{
  // General settings
  "show_diagnostics_panel_on_save": 0,

  // Language server configurations
  "clients": {
    "phpactor": {
      // enable this configuration
      "enabled": true,
      // the startup command -- what you would type in a terminal
      "command": ["wa", "lsp"],
      // the selector that selects which type of buffers this language server attaches to
      "selector": "source.wa"
    }
  }
}
```

## VS Code 支持

TODO
