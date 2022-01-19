<#

.SYNOPSIS
Build a Go application to various platforms

.DESCRIPTION


.EXAMPLE


.NOTES
Don't move this script, is must be in the root folder.

.LINK
https://github.com/dhcgn/GoTemplate

#>

if ((Get-Command Go -ErrorAction Ignore) -eq $null) {
    Write-Error "Couldn't find Go, is PATH to Go missing?"
    return
}

$appName = "gdown"
$version = "2.1.0"
$publishFolder = "publish"
$debugFolder = "debug"

$commitID = Invoke-Expression "git rev-list -1 HEAD"
if ($LASTEXITCODE -ne 0) {
    Write-Error "Couldn't get commit ID"
    EXIT
}

$rootFolder = Split-Path -parent $PSCommandPath

# Just uncomment the platfoms you don't need
$platforms = @()
$platforms += @{GOOS = "windows"; GOARCH = "amd64"; }
#$platforms += @{GOOS = "windows"; GOARCH = "386"; }
$platforms += @{GOOS = "linux"; GOARCH = "amd64"; }
#$platforms += @{GOOS = "linux"; GOARCH = "386"; }
#$platforms += @{GOOS = "linux"; GOARCH = "arm"; }
$platforms += @{GOOS = "linux"; GOARCH = "arm64"; }
$platforms += @{GOOS = "darwin"; GOARCH = "amd64"; }
$platforms += @{GOOS = "darwin"; GOARCH = "arm64"; }

# Clean Up

Remove-Item -Path ([System.IO.Path]::Combine($rootFolder, "build", $publishFolder)) -Recurse -ErrorAction Ignore
Remove-Item -Path ([System.IO.Path]::Combine($rootFolder, "build", $debugFolder)) -Recurse -ErrorAction Ignore

# Build
$count = 0
$maxCount = $platforms.Count * 2
if($compressPublish)
{
    $maxCount += $platforms.Count
}

# Save GO envs for restore
$savedGOOS = $env:GOOS
$savedGOARCH = $env:GOARCH

foreach ($item in $platforms ) {
    # Write-Host "Build" $item.GOOS $item.GOARCH  -ForegroundColor Green
    Write-Progress -Activity ("Build $($item.GOOS) $($item.GOARCH)")  -PercentComplete ([Double]$count / $maxCount * 100)
    

    $env:GOOS = $item.GOOS
    $env:GOARCH = $item.GOARCH

    if ($item.GOOS -eq "windows") {
        $extension = ".exe"
    }
    else {
        $extension = ".bin"
    }
        
    $buildCode = (Join-Path -Path $rootFolder "cmd\GitLabFileDownloader")
   
    $count += 1
    Write-Progress -Activity ("Build $($item.GOOS) $($item.GOARCH)") -Status "Build publish" -PercentComplete ([Double]$count / $maxCount * 100)

    $buildOutput = ([System.IO.Path]::Combine( $rootFolder, "build", $publishFolder, ("{0}_{1}_{2}_{3}{4}" -f $appName, $item.GOOS, $item.GOARCH, $version, $extension)))
    $executeExpression = "go build -ldflags ""-s -w -X main.version={0} -X main.commitID={1}"" -trimpath -o {2} {3}" -f $version, $commitID, $buildOutput, $buildCode 
    Write-Host "Execute", $executeExpression -ForegroundColor Green
    Invoke-Expression $executeExpression
    
    Compress-Archive -Path $buildOutput -DestinationPath ($buildOutput+".zip")

    if (-not (Test-Path $buildOutput)) {
        Write-Host "ERROR - Build result is missing!" -ForegroundColor Red
        continue
    }

    Start-Sleep -Seconds 1 # Because of stupid AV-Shit 

    $count += 1
    Write-Progress -Activity ("Build $($item.GOOS) $($item.GOARCH)") -Status "Build debug" -PercentComplete ([Double]$count / $maxCount * 100)

    $buildOutput = ([System.IO.Path]::Combine( $rootFolder, "build", $debugFolder, ("{0}_{1}_{2}{3}" -f $appName, $item.GOOS, $item.GOARCH, $extension)))
    $executeExpression = "go build -ldflags ""-X main.version={0}"" -o {1} {2}" -f $version, $buildOutput, $buildCode 
    Write-Host "Execute", $executeExpression -ForegroundColor Green
    Invoke-Expression $executeExpression

    Compress-Archive -Path $buildOutput -DestinationPath ($buildOutput+".zip")
}

# Restore GO envs
$env:GOOS = $savedGOOS
$env:GOARCH = $savedGOARCH

Write-Host "Done!" -ForegroundColor Green