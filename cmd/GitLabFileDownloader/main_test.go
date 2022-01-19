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

	"github.com/dhcgn/GitLabFileDownloader/internal"
	"github.com/dhcgn/GitLabFileDownloader/internal/api"
)

func Test_main_no_arguments(t *testing.T) {

	output := captureOutput(func() {
		main()
	})

	fmt.Println(output)

	expected := "Arguments are missing"
	if !strings.Contains(output, expected) {
		t.Errorf("main() got console output = \"%v\", want \"%v\"", output, expected)
	}
}

func Test_main_mode_file(t *testing.T) {
	err, filePath := getTempFilePath()
	if err != nil {
		t.Error(err)
	}

	setFlags(filePath)

	api.HttpGetFunc = func(url string, s internal.Settings) ([]byte, error) {
		return []byte(`{
			"file_name": "settings.json",
			"file_path": "settings.json",
			"size": 66,
			"encoding": "base64",
			"content_sha256": "3de0a34a2cd8d60061f9ac2feda73053b0b8de80995d3fd167c2c225f73817a4",
			"ref": "master",
			"blob_id": "3bb802a168cc02233c337503990b8d906619583b",
			"commit_id": "726a84679597812d8085085f742fb5ddba8a0299",
			"last_commit_id": "4005048b4c3d556ebcdb40bd7dc471fd2216d635",
			"content": "ewogICAgImZydWl0IjogIkFwcGxlIiwKICAgICJzaXplIjogIkxhcmdlIiwKICAgICJjb2xvciI6ICJSZWQiCn0K"
		}`), nil
	}

	var output string

	tests := []struct {
		name        string
		prepare     func()
		wantContent []string
	}{
		{
			name: "New File",
			prepare: func() {

			},
			wantContent: []string{"Wrote"},
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
			wantContent: []string{"Wrote"},
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
			wantContent: []string{"Skip"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()

			output = captureOutput(func() {
				main()
			})

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
