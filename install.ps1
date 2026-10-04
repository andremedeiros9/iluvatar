# Installs iluvatar on Windows.
#
# Run straight from the web, it downloads the latest release:
#
#   irm https://github.com/andremedeiros9/iluvatar/releases/latest/download/install.ps1 | iex
#
# Run from an extracted release archive (double-click install.cmd), it
# installs the iluvatar.exe sitting next to it instead, without
# downloading anything.
#
# The binary goes into %LOCALAPPDATA%\Programs\iluvatar, which is added to
# your user PATH; set ILUVATAR_INSTALL_DIR to put it somewhere else.
$ErrorActionPreference = 'Stop'

$repo = 'andremedeiros9/iluvatar'

$installDir = $env:ILUVATAR_INSTALL_DIR
if (-not $installDir) {
    $installDir = Join-Path $env:LOCALAPPDATA 'Programs\iluvatar'
}

# $bundled is the binary shipped alongside this script in a release
# archive. $PSScriptRoot is empty when the script is piped into iex.
$bundled = $null
if ($PSScriptRoot -and (Test-Path (Join-Path $PSScriptRoot 'iluvatar.exe'))) {
    $bundled = Join-Path $PSScriptRoot 'iluvatar.exe'
}

if ($bundled) {
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
    Copy-Item -Force $bundled (Join-Path $installDir 'iluvatar.exe')
}
else {
    switch ($env:PROCESSOR_ARCHITECTURE) {
        'AMD64' { $arch = 'amd64' }
        'ARM64' { $arch = 'arm64' }
        default { throw "iluvatar: unsupported architecture $env:PROCESSOR_ARCHITECTURE" }
    }

    $archive = "iluvatar_windows_$arch.zip"
    $baseUrl = "https://github.com/$repo/releases/latest/download"

    $tmp = Join-Path ([IO.Path]::GetTempPath()) ("iluvatar-" + [Guid]::NewGuid())
    New-Item -ItemType Directory -Path $tmp | Out-Null

    try {
        Write-Host "Downloading $archive..."
        Invoke-WebRequest -UseBasicParsing -Uri "$baseUrl/$archive" -OutFile (Join-Path $tmp $archive)
        Invoke-WebRequest -UseBasicParsing -Uri "$baseUrl/checksums.txt" -OutFile (Join-Path $tmp 'checksums.txt')

        $expected = $null
        foreach ($line in Get-Content (Join-Path $tmp 'checksums.txt')) {
            $parts = $line -split '\s+'
            if ($parts.Count -ge 2 -and $parts[1] -eq $archive) {
                $expected = $parts[0]
            }
        }
        $actual = (Get-FileHash -Algorithm SHA256 (Join-Path $tmp $archive)).Hash
        if (-not $expected -or $actual -ne $expected) {
            throw "iluvatar: checksum mismatch for $archive"
        }

        Expand-Archive -Force -Path (Join-Path $tmp $archive) -DestinationPath (Join-Path $tmp 'extracted')
        New-Item -ItemType Directory -Force -Path $installDir | Out-Null
        Copy-Item -Force (Join-Path $tmp 'extracted\iluvatar.exe') (Join-Path $installDir 'iluvatar.exe')
    }
    finally {
        Remove-Item -Recurse -Force $tmp
    }
}

# A browser download is marked as coming from the internet; the user chose
# to install it, so lift that.
Unblock-File (Join-Path $installDir 'iluvatar.exe')

Write-Host "Installed iluvatar to $(Join-Path $installDir 'iluvatar.exe')"

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if (($userPath -split ';') -notcontains $installDir) {
    [Environment]::SetEnvironmentVariable('Path', ($userPath.TrimEnd(';') + ';' + $installDir), 'User')
    Write-Host "Added $installDir to your PATH. Open a new terminal to run 'iluvatar'."
}
