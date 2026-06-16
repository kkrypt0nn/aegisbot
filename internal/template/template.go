package template

import (
	"bytes"
	"text/template"

	"github.com/kkrypt0nn/aegisbot/internal/log"
)

const (
	DefaultAlert = "⚠️ Rule '{{.RuleName}}' matched and triggered by <@{{.UserID}}>"
	DefaultBan   = "Matched rule '{{.RuleName}}', action=ban"
	DefaultKick  = "Matched rule '{{.RuleName}}', action=kick"
)

func Render(tpl string, variables map[string]any, defaultTpl string) string {
	if tpl == "" {
		tpl = defaultTpl
	}

	var buf bytes.Buffer
	t, err := template.New("template").Parse(tpl)
	if err != nil {
		log.Errorf("Failed parsing template: %v", err)
		return tpl
	}
	if err := t.Execute(&buf, variables); err != nil {
		log.Errorf("Failed executing template: %v", err)
		return tpl
	}
	return buf.String()
}
