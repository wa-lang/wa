# LSP 测试

凹语言的 LSP 只支持根工作区有 wa.mod 的场景，目前还在开发的早期阶段。

## Sublime Text 语法高亮

基于 [Go语言配置](https://github.com/DisposaBoy/GoSublime/blob/development/syntax/GoSublime-Go-Recommended.sublime-syntax) 修改得到凹语言高亮配置文件 [sublime/wa.sublime-syntax](wa.sublime-syntax)。

选择 Preferences -> Browse Packages 菜单打开包文件所在目录，将 [sublime/wa.sublime-syntax](wa.sublime-syntax) 文件复制到 User 子目录下。重启动 Sublime Text，打开`*.wa`后缀名的凹语言程序，查看高亮是否工作正常。

## Sublime Text LSP 支持

1. 安装 Sublime Text, 打开命令面板
2. 安装 Package Control 包管理插件
3. 通过包管理安装 LSP 包: Package Control: Install Package, 选择 LSP
4. 通过包管理查看 LSP 包已经成功安装: Package Control: List Packages, 确认 LSP 在列表中
5. 命令面板选择 LSP: Enable Language Server Globally 全局范围打开 LSP 插件
6. 菜单: Preferences -> Package Settings -> LSP -> Settings, 添加配置(参考后面的例子)
7. 菜单: Tools -> LSP -> Toggle Log Panel 打开 LSP 的日志面板
8. 打开凹语言程序, 在 LSP 的日志面板可以看到服务启动的日志
9. 右键上下文菜单: LSP -> Format File 格式化文件
10. 其他 LSP 功能还在开发中...

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
