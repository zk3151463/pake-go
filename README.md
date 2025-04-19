# Pake-Go

Pake-Go 是一个用 Go 语言编写的轻量级工具，可以将任何网页打包成桌面应用程序。它使用 [Wails](https://wails.io/) 作为底层框架，提供了原生的桌面应用体验。

## 特性

- 🚀 快速将网页打包成桌面应用
- 🎨 支持自定义窗口样式
- 💉 支持注入自定义 CSS 和 JavaScript
- 🔒 支持自定义请求头
- 📱 支持自定义 User-Agent
- 🖥️ 支持自定义窗口大小
- 🎯 支持窗口置顶
- 🎨 支持隐藏标题栏
- 📦 支持多平台打包（Windows、macOS、Linux）

## 安装

确保你的系统已经安装了 Go 1.21 或更高版本。

```bash
go install github.com/zk3151463/pake-go@latest
```

## 使用方法

### 命令行方式

```bash
pake-go -u https://example.com -n MyApp
```

### 配置文件方式

创建 `config.json` 文件：

```json
{
  "url": "https://example.com",
  "name": "MyApp",
  "width": 1200,
  "height": 800,
  "hideTitleBar": false,
  "transparent": false,
  "alwaysOnTop": false,
  "userAgent": "",
  "icon": "path/to/icon.png",
  "injectCSS": "",
  "injectJS": "",
  "headers": {}
}
```

然后运行：

```bash
pake-go -c config.json
```

## 配置选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| url | 要打包的网页地址 | - |
| name | 应用程序名称 | - |
| width | 窗口宽度 | 1200 |
| height | 窗口高度 | 800 |
| hideTitleBar | 是否隐藏标题栏 | false |
| transparent | 是否透明背景 | false |
| alwaysOnTop | 是否窗口置顶 | false |
| userAgent | 自定义 User-Agent | - |
| icon | 应用图标路径 | - |
| injectCSS | 注入的 CSS 代码 | - |
| injectJS | 注入的 JavaScript 代码 | - |
| headers | 自定义请求头 | {} |

## 开发

1. 克隆仓库：
```bash
git clone https://github.com/zk3151463/pake-go.git
cd pake-go
```

2. 安装依赖：
```bash
go mod download
```

3. 构建项目：
```bash
go build
```

## 许可证

MIT License 