package web

import "html/template"

type TemplateProvider interface {
	MetricsViewTemplate() *template.Template
}
