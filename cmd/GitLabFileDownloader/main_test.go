package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_main_use_cases(t *testing.T) {
	err, filePath := getTempFilePath()
	if err != nil {
		t.Error(err)
	}

	setFlags(filePath)

	var output string

	tests := []struct {
		name         string
		prepare      func()
		wantContent  []string
		wantExitCode int
	}{
		{
			name: "New File",
			prepare: func() {

			},
			wantContent:  []string{"File from disk differs", "New file was copied"},
			wantExitCode: 0,
		},
		{
			name: "Diff File",
			prepare: func() {
				f, err := os.Create(filePath)
				if err != nil {
					t.Error(err)
				}
				f.Write(getContent())
				f.Write([]byte("Add some content to file to create a different hash"))
				f.Close()
			},
			wantContent:  []string{"File from disk differs", "New file was copied"},
			wantExitCode: 0,
		},
		{
			name: "No Diff File",
			prepare: func() {
				f, err := os.Create(filePath)
				if err != nil {
					t.Error(err)
				}
				f.Write(getContent())
				f.Close()
			},
			wantContent:  []string{"No diff, nothing to do."},
			wantExitCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()

			output = captureOutput(func() {
				main()
			})

			fmt.Println(output)

			if exitCode != tt.wantExitCode {
				t.Errorf("main() got exitCode = \"%v\", want \"%v\"", exitCode, tt.wantExitCode)
			}

			for _, line := range tt.wantContent {
				if !strings.Contains(output, line) {
					t.Errorf("main() got console output = \"%v\", want \"%v\"", output, line)
				}
			}

			log.SetOutput(nil)

		})
	}
}

func getTempFilePath() (error, string) {
	tmpfileTarget, _ := ioutil.TempFile("", "golang-test.*")
	filePath := tmpfileTarget.Name()

	if err := tmpfileTarget.Close(); err != nil {
		return err, ""
	}

	if err := os.Remove(filePath); err != nil {
		return err, ""
	}

	return nil, filePath
}

func setFlags(path string) {
	filePath := "settings.json"
	flagRepoFilePathPar = &filePath

	outPath := path
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

func getContent() []byte {
	// file: https://gitlab.com/gitLabFileDownloader/test-project/blob/master/settings.json
	const contentBase64 = "ewogICAgImZydWl0IjogIkFwcGxlIiwKICAgICJzaXplIjogIkxhcmdlIiwKICAgICJjb2xvciI6ICJSZWQiCn0K"
	data, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		panic(err)
	}

	return data
}
