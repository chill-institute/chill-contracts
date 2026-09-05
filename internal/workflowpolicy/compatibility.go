package workflowpolicy

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

var stableRelease = regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)
var breakingMarker = regexp.MustCompile(`(?m)(^[a-z]+(\([^)]+\))?!:)|(^BREAKING CHANGE:)`)

const wireConfig = `version: v2
modules:
  - path: proto
deps:
  - buf.build/grpc-ecosystem/grpc-gateway
breaking:
  use:
    - WIRE
`

// CheckCompatibility checks the working schema against an explicit Git ref, or
// the preceding stable release when against is "release".
func CheckCompatibility(ctx context.Context, dir, against string, output io.Writer) error {
	head, err := gitOutput(ctx, dir, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return err
	}
	if against == "release" {
		against, err = precedingRelease(ctx, dir, head)
		if err != nil {
			return err
		}
	}
	if against == "" {
		return fmt.Errorf("a baseline ref or release is required")
	}
	base, err := gitOutput(ctx, dir, "rev-parse", "--verify", "--end-of-options", against+"^{commit}")
	if err != nil {
		return err
	}
	if _, err := gitOutput(ctx, dir, "merge-base", "--is-ancestor", base, head); err != nil {
		return fmt.Errorf("baseline must be an ancestor of HEAD: %w", err)
	}
	if _, err := fmt.Fprintf(output, "Checking compatibility against %s (%s)\n", against, base); err != nil {
		return err
	}
	args := []string{"breaking", "--against", ".git#ref=" + base}
	sourceErr := runBuf(ctx, dir, output, args...)
	if sourceErr == nil {
		return nil
	}
	messages, err := gitOutput(ctx, dir, "log", "--format=%B", base+".."+head)
	if err != nil {
		return err
	}
	if !breakingMarker.MatchString(messages) {
		return fmt.Errorf("source breaking changes require a Conventional Commit breaking marker in %s..HEAD: %w", against, sourceErr)
	}
	if _, err := fmt.Fprintln(output, "Marked major change; checking wire compatibility"); err != nil {
		return err
	}
	return runBuf(ctx, dir, output, append(args, "--config", wireConfig)...)
}

func precedingRelease(ctx context.Context, dir, head string) (string, error) {
	tags, err := gitOutput(ctx, dir, "tag", "--merged", head, "--sort=-version:refname", "--list", "v*")
	if err != nil {
		return "", err
	}
	for _, tag := range strings.Fields(tags) {
		if !stableRelease.MatchString(tag) {
			continue
		}
		commit, err := gitOutput(ctx, dir, "rev-parse", "--verify", tag+"^{commit}")
		if err != nil {
			return "", err
		}
		if commit != head {
			return tag, nil
		}
	}
	return "", fmt.Errorf("no preceding stable release is reachable from HEAD")
}

func gitOutput(ctx context.Context, dir string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, "git", args...)
	command.Dir = dir
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

func runBuf(ctx context.Context, dir string, output io.Writer, args ...string) error {
	command := exec.CommandContext(ctx, "buf", args...)
	command.Dir = dir
	command.Stdout = output
	command.Stderr = output
	if err := command.Run(); err != nil {
		return fmt.Errorf("buf compatibility check: %w", err)
	}
	return nil
}
