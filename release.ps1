function Get-ScriptDirectory {
    Split-Path -parent $PSCommandPath
}

$equinox = (Resolve-Path (Join-Path (Get-ScriptDirectory) "build\tools\equinox.exe")).Path
$key = (Resolve-Path (Join-Path (Get-ScriptDirectory) "build\secrets\equinox.key")).Path
$token = Get-Content (Resolve-Path (Join-Path (Get-ScriptDirectory) "build\secrets\token.txt")).Path
$version = "1.0.6"

. $equinox release --version=$version `
    --platforms="windows_amd64 linux_amd64" `
    --channel "stable" `
    --signing-key=$key `
    --app="app_dTSgRj85fgP" `
    --token=$token `
    -- -ldflags ("-s -w -X main.version={0}" -f $version) -trimpath `
    github.com/GitLabFileDownloader/src