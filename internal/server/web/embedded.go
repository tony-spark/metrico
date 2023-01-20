package web

import (
	"html/template"

	"github.com/rs/zerolog/log"

	"github.com/tony-spark/metrico/assets"
)

type EmbeddedTemplatesProvider struct {
	metricsView *template.Template
}

func NewEmbeddedTemplates() EmbeddedTemplatesProvider {
	metricsViewTemplate, err := template.ParseFS(assets.EmbeddedAssets, "templates/metrics.html")
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not load template %v", err)
	}
	return EmbeddedTemplatesProvider{
		metricsView: metricsViewTemplate,
	}
}

func (e EmbeddedTemplatesProvider) MetricsViewTemplate() *template.Template {
	return e.metricsView
}
