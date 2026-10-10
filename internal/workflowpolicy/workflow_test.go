package workflowpolicy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const workflowsDir = "../../.github/workflows"

func workflows(t *testing.T) map[string]string {
	t.Helper()
	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yml") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(workflowsDir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		out[e.Name()] = string(b)
	}
	if len(out) == 0 {
		t.Fatal("no workflows found")
	}
	return out
}

// zizmor in `mise run actions` owns action pinning, checkout credential
// persistence, and dangerous triggers.
func TestOrgWorkflowInvariants(t *testing.T) {
	for name, w := range workflows(t) {
		t.Run(name, func(t *testing.T) {
			if strings.Contains(w, "\n  schedule:\n") {
				t.Fatal("declares a GitHub schedule; Cloudflare owns recurring dispatch")
			}
			if !strings.Contains(w, "\npermissions:") {
				t.Fatal("no workflow-level permissions block")
			}
		})
	}
}

func TestBothWorkflowsRequireCompatibility(t *testing.T) {
	all := workflows(t)
	for name, base := range map[string]string{
		"verify.yml": "${{ github.event.pull_request.base.sha }}",
		"main.yml":   "release",
	} {
		verification, _, _ := strings.Cut(all[name], "\n  release:")
		if !stepContains(verification, "run: mise run compatibility:check", "CONTRACTS_BASE_REF: "+base) {
			t.Fatalf("%s must run the shared compatibility policy with baseline %s", name, base)
		}
		if strings.Contains(verification, "buf breaking") {
			t.Fatalf("%s duplicates the shared compatibility policy", name)
		}
	}
	_, release, _ := strings.Cut(all["main.yml"], "\n  release:\n")
	release, _, _ = strings.Cut(release, "\n    steps:\n")
	for _, want := range []string{"needs: [verify]", "needs.verify.result == 'success'"} {
		if !strings.Contains(release, want) {
			t.Fatalf("release must require successful verification: missing %q", want)
		}
	}
}

// stepContains reports whether one workflow step holds every line.
func stepContains(workflow string, lines ...string) bool {
	for _, step := range strings.Split(workflow, "\n      - ") {
		have := map[string]bool{}
		for _, line := range strings.Split(step, "\n") {
			have[strings.TrimSpace(line)] = true
		}
		found := true
		for _, line := range lines {
			found = found && have[line]
		}
		if found {
			return true
		}
	}
	return false
}
