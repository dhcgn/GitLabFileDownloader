# WireGuard Configuration Auto Update

## Into

I want to manage the WireGuard configuration of multiple servers in one git repository.
So I created this tutorial in which I save all my WireGuard configuration files in my GitLap repository and update my WireGuard instance if the config was changed.

## Prerequisite

1. Install WireGuard
2. Install [GitLabFileDownloader](https://github.com/dhcgn/GitLabFileDownloader/releases)
3. Install pwsh (PowerShell)
4. SetUp GitLap Project
   1. Create New Project
   2. Save a valid WireGuard Configuration like from `wg showconf wg0`

## Update Script

Save this at /root/Update-WgConf.ps1

```pwsh
$projectNumber = 0 # ProjectNumber from GitLab
$token = '' # Personal token from GitLab

$reproFilePath = 'wg0.ini' # Filename from saved config in GitLab
$conf = '/root/wg0.conf'
$wgInterface = 'wg0'

/usr/local/bin/gdown -reproFilePath $reproFilePath -outPath $conf -projectNumber $projectNumber -url https://gitlab.com/api/v4/ -token $token
if ($LASTEXITCODE -eq 0) {
    Write-Host 'Update wg conf'
    wg setconf $wgInterface $conf
}
```

## Crontab configuration

Call every 15 minutes `Update-WgConf.ps1` and log the Output to `/root/log.log`

```crontab
*/15 * * * * /usr/bin/pwsh /root/Update-WgConf.ps1 > /root/log.log
```
