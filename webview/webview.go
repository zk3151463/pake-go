package webview

import (
	"strings"
	"sync"
)

// WebViewManager manages web content customization
type WebViewManager struct {
	mu    sync.RWMutex
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
