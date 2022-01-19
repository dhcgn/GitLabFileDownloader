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

	setFlagsFile(filePath)

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
func Test_main_mode_folder(t *testing.T) {
	err, folder := getTempFolderPath()
	if err != nil {
		t.Error(err)
	}

	setFlagsFolder(folder)

	api.HttpGetFunc = func(url string, s internal.Settings) ([]byte, error) {
		if strings.Contains(url, `/repository/tree/?ref=master`) {
			return []byte(`[
				{
					"id": "1e85ff777250e0d0ba1dd079ff562e40784307e1",
					"name": "file1.txt",
					"type": "blob",
					"path": "test_dir/file1.txt",
					"mode": "100644"
				}
			]`), nil
		}
		if strings.Contains(url, `repository/files/test_dir%2Ffile1.txt?ref=master`) {
			return []byte(`{
				"file_name": "file1.txt",
				"file_path": "test_dir/file1.txt",
				"size": 12,
				"encoding": "base64",
				"content_sha256": "11c014f2e9aa58bb56e6a489298ea61a3903c3e632c5aaec5d135996cab0b24e",
				"ref": "master",
				"blob_id": "1e85ff777250e0d0ba1dd079ff562e40784307e1",
				"commit_id": "726a84679597812d8085085f742fb5ddba8a0299",
				"last_commit_id": "9bc24ea56f8862e5964c9f4ee71dab7396902b9f",
				"content": "VGVzdCBGaWxlIDEK"
			}`), nil
		}
		return nil, fmt.Errorf("Unknown URL %v", url)
	}

	var output string

	tests := []struct {
		name        string
		prepare     func()
		wantContent []string
	}{
		{
			name: "No Folder",
			prepare: func() {

			},
			wantContent: []string{"Wrote"},
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

			err := os.RemoveAll(folder)
			if err != nil {
				t.Error(err)
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

func getTempFolderPath() (error, string) {
	filePath, _ := ioutil.TempDir("", "golang-test.*")
	return nil, filePath
}

func setFlagsFile(path string) {
	filePath := "settings.json"
	flagRepoFilePathPar = &filePath

	outPath := path
	flagOutPathPtr = &outPath

	filefolder := ""
	flagRepoFolderPathPtr = &filefolder

	outFolder := ""
	flagOutFolderPtr = &outFolder

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

func setFlagsFolder(folder string) {
	filePath := ""
	flagRepoFilePathPar = &filePath

	outPath := ""
	flagOutPathPtr = &outPath

	filefolder := "test_dir"
	flagRepoFolderPathPtr = &filefolder

	outFolder := folder
	flagOutFolderPtr = &outFolder

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
