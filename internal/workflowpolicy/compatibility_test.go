package workflowpolicy

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompatibilityPolicy(t *testing.T) {
	for _, test := range []struct {
		name    string
		field   string
		message string
		wantErr bool
		wire    bool
	}{
		{"additive", "string value = 1; string extra = 2;", "feat: add field", false, false},
		{"unmarked source break", "string renamed = 1;", "feat: rename field", true, false},
		{"marked source break", "string renamed = 1;", "feat!: rename field", false, true},
		{"footer marked source break", "string renamed = 1;", "feat: rename field\n\nBREAKING CHANGE: field renamed", false, true},
		{"marked wire break", "int32 value = 1;", "feat!: change field type", true, true},
	} {
		for _, release := range []bool{false, true} {
			name := test.name + "/explicit-base"
			if release {
				name = test.name + "/release-base"
			}
			t.Run(name, func(t *testing.T) {
				dir, base := compatibilityRepo(t)
				writeSchema(t, dir, test.field)
				commitFixture(t, dir, test.message)
				against := base
				if release {
					against = "release"
				}
				var output bytes.Buffer
				err := CheckCompatibility(context.Background(), dir, against, &output)
				if (err != nil) != test.wantErr {
					t.Fatalf("error = %v; output:\n%s", err, &output)
				}
				if strings.Contains(output.String(), "checking wire compatibility") != test.wire {
					t.Fatalf("wire check differs from expectation; output:\n%s", &output)
				}
			})
		}
	}
}

func TestPrecedingRelease(t *testing.T) {
	dir, base := compatibilityRepo(t)
	gitFixture(t, dir, "tag", "v2.0.0")
	gitFixture(t, dir, "checkout", "--orphan", "unrelated")
	commitFixture(t, dir, "feat: unrelated release")
	gitFixture(t, dir, "tag", "v99.0.0")
	gitFixture(t, dir, "checkout", "main")
	commitFixture(t, dir, "chore: later commit")
	gitFixture(t, dir, "tag", "v1.1.0")
	commitFixture(t, dir, "chore: current commit")
	gitFixture(t, dir, "tag", "v3.0.0")
	gitFixture(t, dir, "tag", "v9.0.0-beta.1", base)
	head := gitFixture(t, dir, "rev-parse", "HEAD")
	tag, err := precedingRelease(context.Background(), dir, head)
	if err != nil || tag != "v2.0.0" {
		t.Fatalf("preceding release = %q, %v; want highest stable ancestor excluding HEAD", tag, err)
	}
}

func TestCompatibilityRejectsMissingBaseline(t *testing.T) {
	dir, _ := compatibilityRepo(t)
	for _, ref := range []string{"", "missing-ref", "release"} {
		var output bytes.Buffer
		if err := CheckCompatibility(context.Background(), dir, ref, &output); err == nil {
			t.Fatalf("baseline %q unexpectedly passed", ref)
		}
	}
}

func compatibilityRepo(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	gitFixture(t, dir, "init", "-b", "main")
	gitFixture(t, dir, "config", "user.name", "Example Developer")
	gitFixture(t, dir, "config", "user.email", "developer@example.com")
	gitFixture(t, dir, "config", "commit.gpgsign", "false")
	gitFixture(t, dir, "config", "tag.gpgsign", "false")
	gitFixture(t, dir, "config", "core.hooksPath", filepath.Join(dir, "no-hooks"))
	if err := os.Mkdir(filepath.Join(dir, "proto"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"buf.yaml", "buf.lock"} {
		payload, err := os.ReadFile(filepath.Join("../..", file))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, file), payload, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeSchema(t, dir, "string value = 1;")
	// A marker already in the baseline must not authorize a later source break.
	commitFixture(t, dir, "feat!: baseline schema")
	gitFixture(t, dir, "tag", "v1.0.0")
	return dir, gitFixture(t, dir, "rev-parse", "HEAD")
}

func writeSchema(t *testing.T, dir, fields string) {
	t.Helper()
	payload := "syntax = \"proto3\";\npackage example.v1;\nmessage Thing { " + fields + " }\n"
	if err := os.WriteFile(filepath.Join(dir, "proto", "example.proto"), []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commitFixture(t *testing.T, dir, message string) {
	t.Helper()
	gitFixture(t, dir, "add", ".")
	gitFixture(t, dir, "commit", "--allow-empty", "-m", message)
}

func gitFixture(t *testing.T, dir string, args ...string) string {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = dir
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
	return strings.TrimSpace(string(output))
}
