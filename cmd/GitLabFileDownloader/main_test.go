package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	tmpfileTarget, _ := ioutil.TempFile("", "golang-test.*")

	setFlags(tmpfileTarget)

	_ = captureOutput(func() {
		main()
	})

	outputGitLab, _ := ioutil.ReadFile(tmpfileTarget.Name())
	actualoutputGitLab := string(outputGitLab)

	tests := []struct {
		name        string
		wantContent string
	}{
		{
			name: "Integration Test",
			wantContent: `{
    "fruit": "Apple",
    "size": "Large",
    "color": "Red"
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actualoutputGitLab != tt.wantContent {
				t.Errorf("main() gotExists = %v, want %v", actualoutputGitLab, tt.wantContent)
			}
		})
	}
}

func Test_main_console_output(t *testing.T) {
	tmpfileTarget, _ := ioutil.TempFile("", "golang-test.*")
	setFlags(tmpfileTarget)

	output := captureOutput(func() {
		main()
	})

	tests := []struct {
		name        string
		wantContent []string
	}{
		{
			name:        "Integration Test",
			wantContent: []string{"File from disk differs", "New file was copied"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, line := range tt.wantContent {
				if !strings.Contains(output, line) {
					t.Errorf("main() gotExists = \"%v\", want \"%v\"", output, line)
				}
			}
		})
	}
}

func setFlags(tmpfileTarget *os.File) {
	filePath := "settings.json"
	flagRepoFilePathPar = &filePath

	outPath := tmpfileTarget.Name()
	flagOutPathPtr = &outPath

	projectNumber := 16447351
	flagProjectNumberPtr = &projectNumber

	url := "https://gitlab.com/api/v4/"
	flagUrlPtr = &url

	// This token for https://gitlab.com/gitLabFileDownloader/test-project/blob/master/settings.json
	token := "5BUJpxdVx9fyq5KrXJx6"
	flagTokenPtr = &token

	branch := "master"
	flagBranchPtr = &branch
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
