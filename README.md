# GitLabFileDownloader
Download a file from a GitLab server

```
GitLabDownloader_windows_amd64.exe -h

Version: 0.0.0
Usage of GitLabDownloader_windows_amd64.exe:
  -branch string
        Branch (default "master")
  -outPath string
        Path to write the file (default "my_file.json")
  -projectNumber int
        Url to Api v4 (default 34)
  -reproFilePath string
        gitLabFilePathInReproPtr (default "my_config.json")
  -token string
        Private-Token (default "xxxxxxxxxxxxxxxxxxxx")
  -url string
        Url to Api v4 (default "https://my-git-lap-server.local/api/v4/")
```

## Use Case

You want to save for production configs in a on-promise GitLab instance.
With this tool you can sabe the config to the desired location (on windows and linux).

The file will be only replaced if the hash is different.
