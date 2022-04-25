package main

import (
	"strings"
	"os/exec"
)

// pulls a line-by-line block of all git
// tags, in creation (revision) order
// we'll use this to find the most recent
// revision if it was unspecified
func raw_git_revisions() (string, bool) {
	cmd := exec.Command("git", "tag", "-l", "-n1", "--sort=-creatordate")

	result, err := cmd.Output()

	if err != nil {
		return "", false
	}

	return strings.TrimSpace(string(result)), true
}

// this is the "same" function as load_file_normalise, just
// from a specific git tag.  if, for a given file, the tag
// has no reference (newer file, for example) it falls back
// to load_file_normalise and one-stop-shops the process
func load_file_git_tag(file_name, revision_tag string) (string, bool) {
	cmd := exec.Command("git", "diff", "-U999999", revision_tag, file_name)

	result, err := cmd.Output()

	if err != nil {
		return "", false
	}

	text := string(result)

	// if there's nothing to diff,
	// fallback to regular old nonsense
	if len(text) == 0 {
		text, ok := load_file_normalise(file_name)

		if ok {
			text = artificial_diff(text)

			switch text[0] {
			case '\\', '+', '-':
			default: text = " " + text
			}
		}

		return text, ok
	}

	for i := 0; i < 2; i++ {
		n := rune_pair(text, '@', '@')

		if n < 0 {
			break
		}

		text = text[n + 2:]
	}

	text = strings.TrimSpace(text)

	switch text[0] {
	case '\\', '+', '-':
	default: text = " " + text
	}

	return normalise_text(text), true
}

// adds the additional byte-per-line for diff aware mode
// so that we don't have to track pieces individually
func artificial_diff(input string) string {
	buffer := strings.Builder {}
	buffer.Grow(len(input) * 2)

	for _, c := range input {
		buffer.WriteRune(c)
		if c == '\n' {
			buffer.WriteRune(' ')
		}
	}

	return buffer.String()
}