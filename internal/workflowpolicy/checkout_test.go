package workflowpolicy

import (
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestCheckoutDoesNotPersistCredentials(t *testing.T) {
	for name, source := range workflows(t) {
		t.Run(name, func(t *testing.T) {
			var workflow struct {
				Jobs map[string]struct {
					Steps []struct {
						Uses string
						With map[string]any
					}
				}
			}
			if err := yaml.Unmarshal([]byte(source), &workflow); err != nil {
				t.Fatal(err)
			}
			for job, config := range workflow.Jobs {
				for _, step := range config.Steps {
					if strings.HasPrefix(step.Uses, "actions/checkout@") && step.With["persist-credentials"] != false {
						t.Errorf("%s: checkout must set persist-credentials: false", job)
					}
				}
			}
		})
	}
}
