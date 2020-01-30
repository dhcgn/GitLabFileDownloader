# GitLabFileDownloader

[![CircleCI](https://circleci.com/gh/dhcgn/GitLabFileDownloader.svg?style=svg)](https://circleci.com/gh/dhcgn/GitLabFileDownloader)
[![Go Report Card](https://goreportcard.com/badge/github.com/dhcgn/GitLabFileDownloader)](https://goreportcard.com/report/github.com/dhcgn/GitLabFileDownloader)

Download a file from a GitLab server and save it to disk if file is different.

## Latest

https://dl.equinox.io/dhcgn/gitlabfiledownloader/stable

## Using

```plain
GitLabDownloader_windows_amd64.exe -h

Version: 0.0.0
Usage of GitLabDownloader_windows_amd64.exe:
  -branch string
        Branch (default "master")
  -outPath string
        Path to write file to disk
  -projectNumber int
        The Project ID from your project
  -reproFilePath string
        File path in repro, like src/main.go
  -token string
        Private-Token with access right for api and read_repository
  -url string
        Url to Api v4, like https://my-git-lab-server.local/api/v4/
```

## Use Case

You want to have the benefits from git to manage your config files.
With this (windows and linux) tool you can now download theses config files from an on-promise GitLab instance and save them to disk.

The file will be **only** replaced if the hash is different (from disk to git).

```bat
gitlabfiledownloader.exe -outPath settings.json -projectNumber 16447351 -repoFilePath settings.json -token 5BUJpxdVx9fyq5KrXJx6 -url https://gitlab.com/api/v4/
```
