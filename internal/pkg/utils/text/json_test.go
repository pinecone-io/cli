package text

import (
	"strings"
	"testing"
)

func TestIndentJSON_DoesNotEscapeHTMLChars(t *testing.T) {
	input := struct {
		URL string `json:"url"`
	}{URL: "https://example.com/auth?foo=1&bar=2"}

	result := IndentJSON(input)

	if strings.Contains(result, `\u0026`) {
		t.Errorf("IndentJSON escaped & as \\u0026; got: %s", result)
	}
	if !strings.Contains(result, "&") {
		t.Errorf("IndentJSON did not preserve literal &; got: %s", result)
	}
}

func TestInlineJSON_DoesNotEscapeHTMLChars(t *testing.T) {
	input := struct {
		URL string `json:"url"`
	}{URL: "https://example.com/auth?foo=1&bar=2"}

	result := InlineJSON(input)

	if strings.Contains(result, `\u0026`) {
		t.Errorf("InlineJSON escaped & as \\u0026; got: %s", result)
	}
	if !strings.Contains(result, "&") {
		t.Errorf("InlineJSON did not preserve literal &; got: %s", result)
	}
}

func TestIndentJSON_IsIndented(t *testing.T) {
	input := struct {
		Key string `json:"key"`
	}{Key: "value"}

	result := IndentJSON(input)

	if !strings.Contains(result, "\n    ") {
		t.Errorf("IndentJSON output is not indented; got: %s", result)
	}
}

func TestInlineJSON_IsCompact(t *testing.T) {
	input := struct {
		Key string `json:"key"`
	}{Key: "value"}

	result := InlineJSON(input)

	if strings.Contains(result, "\n") {
		t.Errorf("InlineJSON output contains newlines; got: %s", result)
	}
}

func TestIndentJSON_NoTrailingNewline(t *testing.T) {
	result := IndentJSON(struct{ K string }{K: "v"})
	if strings.HasSuffix(result, "\n") {
		t.Errorf("IndentJSON has trailing newline")
	}
}

func TestInlineJSON_NoTrailingNewline(t *testing.T) {
	result := InlineJSON(struct{ K string }{K: "v"})
	if strings.HasSuffix(result, "\n") {
		t.Errorf("InlineJSON has trailing newline")
	}
}
