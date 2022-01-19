# GitLabFileDownloader

![Go](https://github.com/dhcgn/GitLabFileDownloader/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/dhcgn/GitLabFileDownloader/branch/master/graph/badge.svg)](https://codecov.io/gh/dhcgn/GitLabFileDownloader)
[![Go Report Card](https://goreportcard.com/badge/github.com/dhcgn/GitLabFileDownloader)](https://goreportcard.com/report/github.com/dhcgn/GitLabFileDownloader)

Download a file from a GitLab server and save it to disk if file is different.

## Latest

See Releases

## Scopes for personal access tokens

- `read_repository`: Allows read-access to the repository files.
- `api`: Allows read-write access to the repository files.

## TLS Security

**Will be a config switch, soon.**

```go
tr := &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
```

## Using

```plain
gdown.exe -h
2020/01/30 20:49:44 GitLab File Downloader Version: 2.0.2
2020/01/30 20:49:44 Project: https://github.com/dhcgn/GitLabFileDownloader/
Usage of gitlabfiledownloader.exe:
  -branch string
        Branch (default "master")
  -exclude string
        Exclude these regex pattern
  -includeonly string
        Include only these regex pattern
  -outFolder string
        Folder to write file to disk
  -outPath string
        Path to write file to disk
  -projectNumber int
        The Project ID from your project
  -repoFilePath string
        File path in repo, like src/main.go
  -repoFolder string
        Folder to write file to disk
  -token string
        Private-Token with access right for "api" and "read_repository"
  -url string
        Url to Api v4, like https://my-git-lab-server.local/api/v4/
```

## Use Case

### Download file from your gitlab repository 

You want to have the benefits from git to manage your config files.
With this (windows and linux) tool you can now download theses config files from an on-promise GitLab instance and save them to disk.

The file will be **only** replaced if the hash is different (from disk to git).


**Working example!**

See https://gitlab.com/gitLabFileDownloader/test-project for file.

```bat
gdown.exe -outPath settings.json -projectNumber 16447351 -repoFilePath settings.json -token 5BUJpxdVx9fyq5KrXJx6 -url https://gitlab.com/api/v4/
```

```log
2022/01/19 21:16:18 GitLab File Downloader Version: undef Commit: undef
2022/01/19 21:16:18 Project: https://github.com/dhcgn/GitLabFileDownloader/
2022/01/19 21:16:18 Mode: File
2022/01/19 21:16:18 Wrote file: settings.json , because is new or changed
```

### Download folder from your gitlab repository 

**Working example!**

See https://gitlab.com/gitLabFileDownloader/test-project for folder structure.

```bat
gdown.exe -outFolder my_local_dir -projectNumber 16447351 -repoFolder test_dir -token 5BUJpxdVx9fyq5KrXJx6 -url https://gitlab.com/api/v4/ -exclude .gitkeep
```

```log
2022/01/19 21:14:05 GitLab File Downloader Version: undef Commit: undef
2022/01/19 21:14:05 Project: https://github.com/dhcgn/GitLabFileDownloader/
2022/01/19 21:14:05 Mode: Folder
2022/01/19 21:14:06 Sync 5 files, from remote folder test_dir
2022/01/19 21:14:06 Sync 4 files, from remote folder test_dir/test_space_special +ä$
2022/01/19 21:14:06 Sync 2 files, from remote folder test_dir/test_space_special +ä$/deep
2022/01/19 21:14:07 Sync 1 files, from remote folder test_dir/test_space_special +ä$/deep/deep in deep
2022/01/19 21:14:07 Wrote file: test_dir/test_space_special +ä$/deep/deep in deep/another deep file , because is new or changed
2022/01/19 21:14:08 Wrote file: test_dir/test_space_special +ä$/deep/deep_file , because is new or changed
2022/01/19 21:14:08 Wrote file: test_dir/test_space_special +ä$/test#1.txt , because is new or changed
2022/01/19 21:14:08 Wrote file: test_dir/test_space_special +ä$/test#2 #.txt , because is new or changed
2022/01/19 21:14:09 Wrote file: test_dir/test_space_special +ä$/test#3 äöü.bin , because is new or changed
2022/01/19 21:14:09 Skip: .gitkeep because exclude rule: .gitkeep
2022/01/19 21:14:09 Wrote file: test_dir/file1.txt , because is new or changed
2022/01/19 21:14:09 Wrote file: test_dir/file2.txt , because is new or changed
2022/01/19 21:14:10 Wrote file: test_dir/file_space[ ].txt , because is new or changed
```
