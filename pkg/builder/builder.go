package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/zk3151463/pake-go/pkg/config"
)

// Builder handles the application building process
type Builder struct {
	config *config.Config
}

// NewBuilder creates a new Builder instance
func NewBuilder(config *config.Config) *Builder {
	return &Builder{
		config: config,
	}
}

// Build builds the application
func (b *Builder) Build() error {
	// Create project directory
	projectDir := filepath.Join("build", b.config.Name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Generate main.go
	if err := b.generateMainGo(projectDir); err != nil {
		return fmt.Errorf("failed to generate main.go: %w", err)
	}

	// Generate go.mod
	if err := b.generateGoMod(projectDir); err != nil {
		return fmt.Errorf("failed to generate go.mod: %w", err)
	}

	// Copy icon if provided
	if b.config.Icon != "" {
		if err := b.copyIcon(projectDir); err != nil {
			return fmt.Errorf("failed to copy icon: %w", err)
		}
	}

	// Generate wails.json
	if err := b.generateWailsConfig(projectDir); err != nil {
		return fmt.Errorf("failed to generate wails.json: %w", err)
	}

	// Generate frontend
	if err := b.generateFrontend(projectDir); err != nil {
		return fmt.Errorf("failed to generate frontend: %w", err)
	}

	// Build the application
	if err := b.runWailsBuild(projectDir); err != nil {
		return fmt.Errorf("failed to build application: %w", err)
	}

	// Move the built application to the final location
	builtAppPath := filepath.Join(projectDir, "build", "bin")
	finalAppPath := filepath.Join("build", "bin")
	if err := os.MkdirAll(filepath.Dir(finalAppPath), 0755); err != nil {
		return fmt.Errorf("failed to create final directory: %w", err)
	}

	// Remove existing app if it exists
	if err := os.RemoveAll(finalAppPath); err != nil {
		return fmt.Errorf("failed to remove existing app: %w", err)
	}

	// Move the built app to final location
	if err := os.Rename(builtAppPath, finalAppPath); err != nil {
		return fmt.Errorf("failed to move built app: %w", err)
	}

	// Clean up the temporary project directory
	if err := os.RemoveAll(projectDir); err != nil {
		return fmt.Errorf("failed to clean up project directory: %w", err)
	}

	return nil
}

// generateMainGo generates the main.go file
func (b *Builder) generateMainGo(projectDir string) error {
	// Create webview manager file
	webviewFile, err := os.Create(filepath.Join(projectDir, "webview.go"))
	if err != nil {
		return err
	}
	defer webviewFile.Close()

	webviewTmpl := template.Must(template.New("webview").Parse(webviewManagerTemplate))
	if err := webviewTmpl.Execute(webviewFile, nil); err != nil {
		return err
	}

	// Create main.go file
	mainFile, err := os.Create(filepath.Join(projectDir, "main.go"))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	mainTmpl := template.Must(template.New("main").Parse(mainTemplate))
	return mainTmpl.Execute(mainFile, b.config)
}

// generateGoMod generates the go.mod file
func (b *Builder) generateGoMod(projectDir string) error {
	tmpl := template.Must(template.New("go.mod").Parse(goModTemplate))
	file, err := os.Create(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, b.config)
}

// copyIcon copies the icon file to the project directory
func (b *Builder) copyIcon(projectDir string) error {
	src, err := os.Open(b.config.Icon)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(projectDir, "icon.png"))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	return err
}

// generateWailsConfig generates the wails.json file
func (b *Builder) generateWailsConfig(projectDir string) error {
	tmpl := template.Must(template.New("wails.json").Parse(wailsConfigTemplate))
	file, err := os.Create(filepath.Join(projectDir, "wails.json"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, b.config)
}

// generateFrontend generates the frontend files
func (b *Builder) generateFrontend(projectDir string) error {
	frontendDir := filepath.Join(projectDir, "frontend")
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		return err
	}

	// Generate package.json
	if err := b.generatePackageJSON(frontendDir); err != nil {
		return err
	}

	// Generate vite.config.js
	if err := b.generateViteConfig(frontendDir); err != nil {
		return err
	}

	// Generate src directory
	srcDir := filepath.Join(frontendDir, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return err
	}

	// Generate App.vue
	if err := b.generateAppVue(srcDir); err != nil {
		return err
	}

	// Generate main.js
	if err := b.generateMainJS(srcDir); err != nil {
		return err
	}

	// Generate index.html
	if err := b.generateIndexHTML(frontendDir); err != nil {
		return err
	}

	return nil
}

// generatePackageJSON generates the package.json file
func (b *Builder) generatePackageJSON(frontendDir string) error {
	tmpl := template.Must(template.New("package.json").Parse(packageJSONTemplate))
	file, err := os.Create(filepath.Join(frontendDir, "package.json"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, b.config)
}

// generateViteConfig generates the vite.config.js file
func (b *Builder) generateViteConfig(frontendDir string) error {
	tmpl := template.Must(template.New("vite.config.js").Parse(viteConfigTemplate))
	file, err := os.Create(filepath.Join(frontendDir, "vite.config.js"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, b.config)
}

// generateAppVue generates the App.vue file
func (b *Builder) generateAppVue(srcDir string) error {
	tmpl := template.Must(template.New("App.vue").Parse(appVueTemplate))
	file, err := os.Create(filepath.Join(srcDir, "App.vue"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, b.config)
}

// generateMainJS generates the main.js file
func (b *Builder) generateMainJS(srcDir string) error {
	tmpl := template.Must(template.New("main.js").Parse(mainJSTemplate))
	file, err := os.Create(filepath.Join(srcDir, "main.js"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, b.config)
}

// generateIndexHTML generates the index.html file
func (b *Builder) generateIndexHTML(frontendDir string) error {
	tmpl := template.Must(template.New("index.html").Parse(indexHTMLTemplate))
	file, err := os.Create(filepath.Join(frontendDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, b.config)
}

// runWailsBuild runs the wails build command
func (b *Builder) runWailsBuild(projectDir string) error {
	cmd := exec.Command("wails", "build")
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

const mainTemplate = `package main

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// domReady is called after the front-end dom has been loaded
func (a *App) domReady(ctx context.Context) {
	// 使用 JavaScript 重定向到目标 URL，并处理背景色
	script := fmt.Sprintf(` + "`" + `
		(function() {
			// 设置背景色
			document.body.style.backgroundColor = '#ffffff';
			document.documentElement.style.backgroundColor = '#ffffff';

			// 设置用户代理
			{{if .UserAgent}}
			Object.defineProperty(navigator, 'userAgent', {
				get: function() { return '{{.UserAgent}}'; }
			});
			{{end}}

			// 使用 sessionStorage 来防止循环重定向
			if (!sessionStorage.getItem('hasRedirected') && window.location.href !== "{{.URL}}") {
				sessionStorage.setItem('hasRedirected', 'true');
				window.location.href = "{{.URL}}";
			}
		})();
	` + "`" + `)
	runtime.WindowExecJS(ctx, script)
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "{{.Name}}",
		Width:            {{.Width}},
		Height:           {{.Height}},
		DisableResize:    false,
		Fullscreen:       false,
		WindowStartState: options.Normal,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			TitleBar:            {{if .HideTitleBar}}mac.TitleBarHidden(){{else}}mac.TitleBarDefault(){{end}},
			Appearance:          mac.NSAppearanceNameAqua,
		},
		Frameless:   {{.HideTitleBar}},
		AlwaysOnTop: {{.AlwaysOnTop}},
		{{if .UserAgent}}
		WebviewUserAgent: "{{.UserAgent}}",
		{{end}}
	})

	if err != nil {
		log.Fatal(err)
	}
}
`

const goModTemplate = `module {{if .Name}}github.com/{{.Name}}{{else}}pake-app{{end}}

go 1.21

require github.com/wailsapp/wails/v2 v2.10.1

require (
	github.com/bep/debounce v1.2.1 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jchv/go-winloader v0.0.0-20210711035445-715c2860da7e // indirect
	github.com/labstack/echo/v4 v4.10.2 // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/leaanthony/go-ansi-parser v1.6.0 // indirect
	github.com/leaanthony/gosod v1.0.3 // indirect
	github.com/leaanthony/slicer v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/samber/lo v1.38.1 // indirect
	github.com/tkrajina/go-reflector v0.5.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/wailsapp/go-webview2 v1.0.10 // indirect
	github.com/wailsapp/mimetype v1.4.1 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)`

const wailsConfigTemplate = `{
	"name": "{{.Name}}",
	"outputfilename": "{{.Name}}",
	"frontend:install": "npm install",
	"frontend:build": "npm run build",
	"frontend:dev": "npm run dev",
	"author": {
		"name": "Pake-Go",
		"email": "pake-go@example.com"
	},
	"info": {
		"companyName": "Pake-Go",
		"productName": "{{.Name}}",
		"productVersion": "1.0.0",
		"copyright": "Copyright © 2024 Pake-Go",
		"comments": "Built with Pake-Go"
	}
}
`

const packageJSONTemplate = `{
	"name": "{{.Name}}",
	"version": "1.0.0",
	"description": "Built with Pake-Go",
	"type": "module",
	"scripts": {
		"dev": "vite",
		"build": "vite build",
		"preview": "vite preview"
	},
	"dependencies": {
		"vue": "^3.3.0"
	},
	"devDependencies": {
		"@vitejs/plugin-vue": "^4.5.0",
		"vite": "^4.5.0"
	}
}
`

const viteConfigTemplate = `import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
	plugins: [vue()],
	server: {
		port: 34115
	}
})
`

const appVueTemplate = `<template>
	<div id="app">
		<div
			id="frame"
			:style="{
				width: '100%',
				height: '100%',
				border: 'none',
				opacity: isLoading ? 0 : 1,
				transition: 'opacity 0.3s ease-in-out',
				background: '#ffffff'
			}"
			ref="frame"
		></div>
		<div v-if="isLoading" class="loading">
			<div class="spinner"></div>
		</div>
	</div>
</template>

<script>
export default {
	name: 'App',
	data() {
		return {
			isLoading: true
		}
	},
	methods: {
		handleLoad() {
			const frame = this.$refs.frame;
			if (frame) {
				try {
					// 注入自定义 CSS
					if ('{{.InjectCSS}}') {
						const style = document.createElement('style');
						style.textContent = '{{.InjectCSS}}';
						document.head.appendChild(style);
					}

					// 注入自定义 JS
					if ('{{.InjectJS}}') {
						const script = document.createElement('script');
						script.textContent = '{{.InjectJS}}';
						document.body.appendChild(script);
					}

					// 300ms 后隐藏加载状态
					setTimeout(() => {
						this.isLoading = false;
					}, 300);
				} catch (e) {
					console.error('Failed to inject content:', e);
				}
			}
		}
	},
	mounted() {
		this.handleLoad();
	}
}
</script>

<style>
html, body {
	margin: 0;
	padding: 0;
	width: 100%;
	height: 100%;
	overflow: hidden;
	background: #ffffff !important;
}

#app {
	width: 100%;
	height: 100%;
	margin: 0;
	padding: 0;
	overflow: hidden;
	background: #ffffff;
	position: relative;
}

#frame {
	display: block;
	width: 100%;
	height: 100%;
	border: none;
	background: #ffffff;
}

.loading {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: #ffffff;
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 9999;
}

.spinner {
	width: 40px;
	height: 40px;
	border: 4px solid #f3f3f3;
	border-top: 4px solid #3498db;
	border-radius: 50%;
	animation: spin 1s linear infinite;
}

@keyframes spin {
	0% { transform: rotate(0deg); }
	100% { transform: rotate(360deg); }
}
</style>
`

const mainJSTemplate = `import { createApp } from 'vue'
import App from './App.vue'

const app = createApp(App)
app.mount('#app')`

const indexHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta content="width=device-width, initial-scale=1.0" name="viewport" />
		<title>{{.Name}}</title>
	</head>
	<body>
		<div id="app"></div>
		<script src="/src/main.js" type="module"></script>
	</body>
</html>`

const webviewManagerTemplate = `package main

import (
	"strings"
	"sync"
)

// WebViewManager manages web content customization
type WebViewManager struct {
	mu sync.RWMutex
	rules []struct {
		URL     string
		CSS     []string
		JS      []string
		Headers map[string]string
	}
}

// NewWebViewManager creates a new WebViewManager
func NewWebViewManager() *WebViewManager {
	return &WebViewManager{
		rules: make([]struct {
			URL     string
			CSS     []string
			JS      []string
			Headers map[string]string
		}, 0),
	}
}

// AddRule adds a new injection rule
func (w *WebViewManager) AddRule(url string, css []string, js []string, headers map[string]string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.rules = append(w.rules, struct {
		URL     string
		CSS     []string
		JS      []string
		Headers map[string]string
	}{
		URL:     url,
		CSS:     css,
		JS:      js,
		Headers: headers,
	})
}

// GetRulesForURL returns all matching rules for a given URL
func (w *WebViewManager) GetRulesForURL(url string) ([]string, []string, map[string]string) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var css []string
	var js []string
	headers := make(map[string]string)

	for _, rule := range w.rules {
		if strings.Contains(url, rule.URL) {
			css = append(css, rule.CSS...)
			js = append(js, rule.JS...)
			for k, v := range rule.Headers {
				headers[k] = v
			}
		}
	}

	return css, js, headers
}

// ClearRules clears all injection rules
func (w *WebViewManager) ClearRules() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.rules = make([]struct {
		URL     string
		CSS     []string
		JS      []string
		Headers map[string]string
	}, 0)
}

// GenerateInjectionScript generates the JavaScript code for injecting CSS and JS
func (w *WebViewManager) GenerateInjectionScript(css []string, js []string) string {
	var script strings.Builder

	// Inject CSS
	for _, style := range css {
		// Escape single quotes and backslashes
		escapedStyle := strings.ReplaceAll(style, "\\", "\\\\")
		escapedStyle = strings.ReplaceAll(escapedStyle, "'", "\\'")
		
		script.WriteString("(function() {")
		script.WriteString("var style = document.createElement('style');")
		script.WriteString("style.textContent = '" + escapedStyle + "';")
		script.WriteString("document.head.appendChild(style);")
		script.WriteString("})();")
	}

	// Inject JS
	for _, code := range js {
		script.WriteString("(function() {")
		script.WriteString(code)
		script.WriteString("})();")
	}

	return script.String()
}
`
