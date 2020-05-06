function Get-ScriptDirectory {
    Split-Path -parent $PSCommandPath
}

$equinox = (Resolve-Path (Join-Path (Get-ScriptDirectory) "build\tools\equinox.exe")).Path
$key = (Resolve-Path (Join-Path (Get-ScriptDirectory) "build\secrets\equinox.key")).Path
$token = Get-Content (Resolve-Path (Join-Path (Get-ScriptDirectory) "build\secrets\token.txt")).Path
$version = "2.0.1"

. $equinox release --version=$version `
    --platforms="windows_amd64 windows_386 linux_amd64 linux_386" `
    --channel "beta" `
    --signing-key=$key `
    --app="app_dTSgRj85fgP" `
    --token=$token `
    -- -ldflags ("-s -w -X main.version={0}" -f $version) -trimpath `
    github.com/dhcgn/GitLabFileDownloader/cmd/GitLabFileDownloader