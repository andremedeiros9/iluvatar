@echo off
rem Double-clickable installer shipped in the Windows release archive: runs
rem install.ps1 from the same folder, which installs the iluvatar.exe next
rem to it. Windows doesn't run .ps1 files on double-click, hence this.
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0install.ps1"
echo.
pause
