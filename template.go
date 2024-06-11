package gonertia

import (
	"fmt"
	"html/template"
	"net/http"
)

// TemplateData are the data that will be transferred
// and will be available in the root template.
type TemplateData map[string]any

func (i *Inertia) buildTemplateData(r *http.Request, page *page) (TemplateData, error) {
	pageJSON, err := i.marshallJSON(page)
	if err != nil {
		return nil, fmt.Errorf("marshal page into json: %w", err)
	}

	// Get template data from context.
	ctxTemplateData, err := TemplateDataFromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting template data from context: %w", err)
	}

	// Defaults.
	result := TemplateData{
		"inertiaHead": "", // todo reserved for SSR.
		"inertia":     i.inertiaContainerHTML(pageJSON),
	}

	// Add the shared template data to the result.
	for key, val := range i.sharedTemplateData {
		result[key] = val
	}

	// Add template data from context to the result.
	for key, val := range ctxTemplateData {
		result[key] = val
	}

	return result, nil
}

func (i *Inertia) buildSharedTemplateFuncs() template.FuncMap {
	// Defaults.
	result := template.FuncMap{
		"mix": func(path string) (string, error) {
			if val, ok := i.mixManifestData[path]; ok {
				return val, nil
			}
			return path, fmt.Errorf("file %q not found in mix manifest file", path)
		},
	}

	// Add the shared template funcs to the result.
	for key, val := range i.sharedTemplateFuncs {
		result[key] = val
	}

	return result
}
