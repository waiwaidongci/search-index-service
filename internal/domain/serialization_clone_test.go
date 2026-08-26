package domain

import (
	"testing"
	"time"
)

func TestCloneRulesOwnsPointersAndMaps(t *testing.T) {
	pct := 50
	start := time.Now()
	end := start.Add(time.Hour)
	in := []Rule{{
		ID:         "r",
		Priority:   1,
		Tags:       map[string]string{"env": "prod"},
		Percentage: &pct,
		StartAt:    &start,
		EndAt:      &end,
		Value:      "on",
	}}
	out := CloneRules(in)
	*in[0].Percentage = 90
	*in[0].StartAt = start.Add(time.Minute)
	*in[0].EndAt = end.Add(time.Hour)
	in[0].Tags["env"] = "staging"
	if out[0].Tags["env"] != "prod" {
		t.Fatalf("tags map was shared: %#v", out[0].Tags)
	}
	if *out[0].Percentage != 50 || out[0].StartAt.Equal(start) || out[0].EndAt.Equal(end) {
		t.Fatalf("pointer fields were shared: %#v", out[0])
	}
}

func TestCloneRulesOwnsJSONValue(t *testing.T) {
	in := []Rule{{
		ID:    "r",
		Value: map[string]any{"nested": map[string]any{"x": "one"}},
	}}
	out := CloneRules(in)
	root := in[0].Value.(map[string]any)
	root["nested"].(map[string]any)["x"] = "changed"
	if out[0].Value.(map[string]any)["nested"].(map[string]any)["x"] != "one" {
		t.Fatalf("JSON value was shared: %#v", out[0].Value)
	}
}
