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

```go
tr := &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
```

## Using

```plain
gitlabfiledownloader.exe -h
2020/01/30 20:49:44 GitLab File Downloader Version: 2.0.0
2020/01/30 20:49:44 Project: https://github.com/dhcgn/GitLabFileDownloader/
Usage of gitlabfiledownloader.exe:
  -branch string
        Branch (default "master")
  -outPath string
        Path to write file to disk
  -projectNumber int
        The Project ID from your project
  -repoFilePath string
        File path in repo, like src/main.go
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

```bat
gitlabfiledownloader.exe -outPath settings.json -projectNumber 16447351 -repoFilePath settings.json -token 5BUJpxdVx9fyq5KrXJx6 -url https://gitlab.com/api/v4/
```

### Download folder from your gitlab repository 

```bat
gitlabfiledownloader.exe -outFolder test_dir -projectNumber 16447351 -repoFolder test_dir -token 5BUJpxdVx9fyq5KrXJx6 -url https://gitlab.com/api/v4/
```