package domain

import (
	"testing"
)

func TestFilterResultDoesNotAliasInputRules(t *testing.T) {
	input := []QueryTemplate{{
		ID:   "q",
		Key:  "flag",
		Type: TypeString,
		Rules: []Rule{{
			ID:    "r",
			Value: "one",
			Tags:  map[string]string{"env": "prod"},
		}},
	}}
	filtered := FilterQueryTemplates(input, QueryTemplateFilter{})
	if len(filtered) != 1 {
		t.Fatalf("expected one filtered template, got %d", len(filtered))
	}
	input[0].Rules[0].Value = "changed"
	input[0].Rules[0].Tags["env"] = "staging"
	if filtered[0].Rules[0].Value != "one" || filtered[0].Rules[0].Tags["env"] != "prod" {
		t.Fatalf("filter result aliases input rules: %#v", filtered[0].Rules[0])
	}
}

func TestCompareReportsBothRuleLengthDirections(t *testing.T) {
	extra := Rule{ID: "extra", Value: "on"}
	a := TemplateRevision{Value: "on"}
	b := TemplateRevision{Value: "on", Rules: []Rule{extra}}
	changes := CompareTemplateRevisions(a, b)
	found := false
	for _, c := range changes {
		if c.Path == "rules[0]" && c.Kind == "added" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected added rules[0] change, got %#v", changes)
	}
}

func TestCompareReportsRemovedRuleDirection(t *testing.T) {
	extra := Rule{ID: "extra", Value: "on"}
	a := TemplateRevision{Value: "on", Rules: []Rule{extra}}
	b := TemplateRevision{Value: "on"}
	changes := CompareTemplateRevisions(a, b)
	found := false
	for _, c := range changes {
		if c.Path == "rules[0]" && c.Kind == "removed" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected removed rules[0] change, got %#v", changes)
	}
}
