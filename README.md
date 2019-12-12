# GitLabFileDownloader

Download a file from a GitLab server and save it to disk if file is different.

## Latest

https://dl.equinox.io/dhcgn/gitlabfiledownloader/stable

## Using

```
GitLabDownloader_windows_amd64.exe -h

Version: 0.0.0
Usage of GitLabDownloader_windows_amd64.exe:
  -branch string
        Branch (default "master")
  -outPath string
        Path to write the file (default "my_file.json")
  -projectNumber int
        Url to Api v4
  -reproFilePath string
        gitLabFilePathInReproPtr (default "my_config.json")
  -token string
        Private-Token (default "xxxxxxxxxxxxxxxxxxxx")
  -url string
        Url to Api v4 (default "https://my-git-lap-server.local/api/v4/")
```

## Use Case

You want to have the benefits from git to manage your config files.
With this (windows and linux) tool you can now download theses config files from an on-promise GitLab instance and save them to disk.

The file will be **only** replaced if the hash is different (from disk to git).

```bat
GitLabDownloader_windows_amd64_1.0.0.exe -reproFilePath myconfig.xml -outPath c:\App\myconfig.xml -projectNumber 547 -url https://my-git-server.com/api/v4/ -token jd32dwEH2FS42342Sdf32
```
