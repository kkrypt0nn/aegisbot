package rules

import (
	"os"
	"testing"
)

func createTempRuleFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	_ = tmpFile.Close()
	return tmpFile.Name()
}

func TestParseRuleInvalidYAML(t *testing.T) {
	tmpFile := createTempRuleFile(t, `meta:`)
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	_, err := Parse(tmpFile)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseRulePhishingLinkYAML(t *testing.T) {
	ruleContent := `
- rule:
    name: "PhishingLink"
    meta:
      action: "alert"
      context: "message"
      ignoreBots: true
    strings:
      - name: "link"
        value: "https://badsite.com"
      - name: "scam"
        value: "free nitro"
    expression: |
      message.content.contains(link) && message.content.contains(scam)
`

	tmpFile := createTempRuleFile(t, ruleContent)
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	rules, err := Parse(tmpFile)
	if err != nil {
		t.Fatalf("failed to parse rule: %v", err)
	}

	rule := rules[0]

	if rule.Name != "PhishingLink" {
		t.Errorf("expected name 'PhishingLink', got '%s'", rule.Name)
	}

	if rule.Action != "alert" {
		t.Errorf("expected action 'alert', got '%s'", rule.Action)
	}

	if rule.Context != "message" {
		t.Errorf("expected context 'message', got '%s'", rule.Action)
	}

	if rule.IgnoreBots != true {
		t.Errorf("expected ignoreBots 'true', got '%v'", rule.IgnoreBots)
	}

	if len(rule.Strings) != 2 {
		t.Fatalf("expected 2 strings, got %d", len(rule.Strings))
	}

	link, ok := rule.Strings["link"]
	if !ok {
		t.Errorf("missing string 'link'")
	} else {
		if link.Value != "https://badsite.com" {
			t.Errorf("expected link value 'https://badsite.com', got '%s'", link.Value)
		}
	}

	scam, ok := rule.Strings["scam"]
	if !ok {
		t.Errorf("missing string 'scam'")
	} else {
		if scam.Value != "free nitro" {
			t.Errorf("expected scam value 'free nitro', got '%s'", scam.Value)
		}
	}

	expectedCondition := "message.content.contains(link) && message.content.contains(scam)"
	if rule.Expression != expectedCondition {
		t.Errorf("expected condition '%s', got '%s'", expectedCondition, rule.Expression)
	}
}
