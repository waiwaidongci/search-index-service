package domain

import (
	"errors"
	"testing"
	"time"
)

func TestRuleChainDoesNotSkipPercentageMatch(t *testing.T) {
	pct := 50
	template := QueryTemplate{
		Type:         TypeString,
		DefaultValue: "off",
		Rules: []Rule{{
			ID:         "full-rollout",
			Priority:   1,
			Percentage: &pct,
			Value:      "on",
		}},
	}
	value, _, err := Evaluate(template, nil, EvaluationContext{SubjectID: "subject-a"})
	if err != nil {
		t.Fatal(err)
	}
	if value != "on" {
		t.Fatalf("expected full rollout rule to match, got %v", value)
	}
}

func TestSegmentConstraintsDoNotSkipUnknownOperator(t *testing.T) {
	segment := Segment{ID: "s", Name: "S", Constraints: []Constraint{{Key: "env", Operator: "invalid", Value: "prod"}}, Percentage: 0}
	if err := segment.Validate(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid operator error, got %v", err)
	}
}

func TestRuleValidationRejectsInvalidTransition(t *testing.T) {
	start := time.Now().Add(time.Hour)
	end := time.Now().Add(-time.Hour)
	rules := []Rule{{ID: "r", Priority: 1, Value: "on", StartAt: &start, EndAt: &end}}
	if err := ValidateRules(TypeString, rules); !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid time window error, got %v", err)
	}
}

func TestEvaluateDoesNotMutateRuleOrder(t *testing.T) {
	template := QueryTemplate{
		Type:         TypeString,
		DefaultValue: "off",
		Rules: []Rule{
			{ID: "second", Priority: 2, Value: "two"},
			{ID: "first", Priority: 1, Value: "one"},
		},
	}
	originalFirst := template.Rules[0].ID
	_, _, err := Evaluate(template, nil, EvaluationContext{})
	if err != nil {
		t.Fatal(err)
	}
	if template.Rules[0].ID != originalFirst {
		t.Fatalf("Evaluate mutated the caller's rule slice: %#v", template.Rules)
	}
}
