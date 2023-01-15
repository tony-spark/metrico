// Package web contains helpers for application's web pages rendering
package web

import "html/template"

type TemplateProvider interface {
	MetricsViewTemplate() *template.Template
}
