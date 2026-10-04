# Installs the latest iluvatar release on Windows.
#
#   irm https://github.com/andremedeiros9/iluvatar/releases/latest/download/install.ps1 | iex
#
# The binary goes into %LOCALAPPDATA%\Programs\iluvatar, which is added to
# your user PATH; set ILUVATAR_INSTALL_DIR to put it somewhere else.
$ErrorActionPreference = 'Stop'

$repo = 'andremedeiros9/iluvatar'

$installDir = $env:ILUVATAR_INSTALL_DIR
if (-not $installDir) {
    $installDir = Join-Path $env:LOCALAPPDATA 'Programs\iluvatar'
}

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

    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
    Expand-Archive -Force -Path (Join-Path $tmp $archive) -DestinationPath $installDir
}
finally {
    Remove-Item -Recurse -Force $tmp
}

Write-Host "Installed iluvatar to $(Join-Path $installDir 'iluvatar.exe')"

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if (($userPath -split ';') -notcontains $installDir) {
    [Environment]::SetEnvironmentVariable('Path', ($userPath.TrimEnd(';') + ';' + $installDir), 'User')
    Write-Host "Added $installDir to your PATH. Open a new terminal to run 'iluvatar'."
}
