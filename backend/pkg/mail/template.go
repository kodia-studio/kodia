package mail

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

// TemplateEngine handles email template rendering
type TemplateEngine struct {
	templateDir string
	funcMap     template.FuncMap
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine(templateDir string) *TemplateEngine {
	return &TemplateEngine{
		templateDir: templateDir,
		funcMap:     createFuncMap(),
	}
}

// createFuncMap creates the function map for templates
func createFuncMap() template.FuncMap {
	return template.FuncMap{
		// Add custom template functions here
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}
}

// RenderHTML renders an HTML template with data
func (te *TemplateEngine) RenderHTML(ctx context.Context, templateName string, data map[string]any) (string, error) {
	templatePath := filepath.Join(te.templateDir, templateName+".html")

	tmpl, err := template.New(templateName).Funcs(te.funcMap).ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, filepath.Base(templatePath), data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// RenderPlainText renders a plain text template with data
func (te *TemplateEngine) RenderPlainText(ctx context.Context, templateName string, data map[string]any) (string, error) {
	templatePath := filepath.Join(te.templateDir, templateName+".txt")

	tmpl, err := template.New(templateName).Funcs(te.funcMap).ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, filepath.Base(templatePath), data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// TemplateMail is a mail class that uses templates
type TemplateMail struct {
	*BaseMail
	templateEngine *TemplateEngine
	templateName   string
	templateData   map[string]any
}

// NewTemplateMail creates a new template-based mail
func NewTemplateMail(engine *TemplateEngine) *TemplateMail {
	return &TemplateMail{
		BaseMail:       NewBaseMail(),
		templateEngine: engine,
		templateData:   make(map[string]any),
	}
}

// WithTemplate sets the template name
func (tm *TemplateMail) WithTemplate(name string) *TemplateMail {
	tm.templateName = name
	return tm
}

// WithTemplateData sets template data
func (tm *TemplateMail) WithTemplateData(data map[string]any) *TemplateMail {
	tm.templateData = data
	return tm
}

// SetTemplateVariable sets a single template variable
func (tm *TemplateMail) SetTemplateVariable(key string, value any) *TemplateMail {
	tm.templateData[key] = value
	return tm
}

// GetTemplate returns the template name
func (tm *TemplateMail) GetTemplate() string {
	return tm.templateName
}

// Build renders the template and sets the body
func (tm *TemplateMail) Build(ctx context.Context, data map[string]any) error {
	if tm.templateName == "" {
		return fmt.Errorf("template name not set")
	}

	// Merge provided data with mail's template data
	finalData := tm.templateData
	for k, v := range data {
		finalData[k] = v
	}

	// Try to render HTML template
	htmlContent, err := tm.templateEngine.RenderHTML(ctx, tm.templateName, finalData)
	if err == nil && htmlContent != "" {
		tm.HTMLBody(htmlContent)
	}

	// Also try to render plain text template
	textContent, err := tm.templateEngine.RenderPlainText(ctx, tm.templateName, finalData)
	if err == nil && textContent != "" {
		tm.Body(textContent)
	}

	return nil
}
