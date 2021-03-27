package cmd

import (
	"bytes"
	"os/exec"
)

func gitFetch(branch string) error {
	gitCmd := exec.Command("git", "fetch", "origin", branch)
	return gitCmd.Run()
}

func gitDiff(branch string) (*bytes.Buffer, error) {
	var stdout bytes.Buffer
	gitCmd := exec.Command("git", "diff", "origin/" + branch)
	gitCmd.Stdout = &stdout
	err := gitCmd.Run()
	if err != nil {
		return nil, err
	}
	return &stdout, nil
}