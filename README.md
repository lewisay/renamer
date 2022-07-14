### renamer

一个用于批量修改文件名称的命令行工具。

#### 安装
```shell
GOPROXY=https://goproxy.cn/,direct go install github.com/lewisay/renamer@latest
```

#### 用法
```
renamer --help
Batch naming of files

Usage:
  renamer [command]

Available Commands:
  append      append -d . --new=xxx
  forward     forward  -d . --new=xxx
  help        Help about any command
  replace     replace -d file-directory --old=find-string --new=new-string

Flags:
  -d, --dir string   directory path (required)
  -h, --help         help for renamer
  -n, --new string   new string
  -o, --old string   old string

Use "renamer [command] --help" for more information about a command.
```

### 示例
```shell
renamer replace -d . --old=find-string --new=new-string
```
- 在当前目录进行替换
- 原始字符串：find-string
- 替换字符串：new-string



