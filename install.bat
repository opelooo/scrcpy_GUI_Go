@echo off
setlocal enabledelayedexpansion

echo =======================================================
echo    Scrcpy GUI - Windows Installer
echo =======================================================

echo [+] Detected Operating System: Windows

set "APP_NAME=Scrcpy GUI"
set "APP_DIR=scrcpy-gui"
set "EXE_NAME=scrcpy-gui.exe"
set "DIST_EXE=dist\%EXE_NAME%"
set "INSTALL_DIR=%LOCALAPPDATA%\%APP_DIR%"

:: 1. Check or build binary
if not exist "%DIST_EXE%" (
    echo [+] Compiling %APP_NAME%...
    where go >nul 2>&1
    if errorlevel 1 (
        echo [-] Error: Go compiler not found and %DIST_EXE% is missing!
        echo     Please install Go or build the binary first.
        pause
        exit /b 1
    )
    if not exist dist mkdir dist
    go build -ldflags="-H=windowsgui -s -w" -o "%DIST_EXE%" .\cmd\scrcpy-gui
    if errorlevel 1 (
        echo [-] Build failed! Check your Go environment.
        pause
        exit /b 1
    )
    echo [+] Build complete: %DIST_EXE%
)

:: 2. Create destination folder
echo [+] Creating installation directory: %INSTALL_DIR%
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

:: 3. Copy binary and icon
echo [+] Copying application files...
copy /Y "%DIST_EXE%" "%INSTALL_DIR%\%EXE_NAME%" >nul
if exist icon.png copy /Y icon.png "%INSTALL_DIR%\icon.png" >nul

:: 4. Create Shortcuts via PowerShell
echo [+] Creating Desktop and Start Menu shortcuts...
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
    "$WshShell = New-Object -ComObject WScript.Shell; " ^
    "$installDir = [Environment]::ExpandEnvironmentVariables('%INSTALL_DIR%'); " ^
    "$targetExe = Join-Path $installDir '%EXE_NAME%'; " ^
    "$desktop = [Environment]::GetFolderPath('Desktop'); " ^
    "$programs = [Environment]::GetFolderPath('Programs'); " ^
    "$sc1 = $WshShell.CreateShortcut((Join-Path $desktop '%APP_NAME%.lnk')); " ^
    "$sc1.TargetPath = $targetExe; " ^
    "$sc1.WorkingDirectory = $installDir; " ^
    "$sc1.Save(); " ^
    "$sc2 = $WshShell.CreateShortcut((Join-Path $programs '%APP_NAME%.lnk')); " ^
    "$sc2.TargetPath = $targetExe; " ^
    "$sc2.WorkingDirectory = $installDir; " ^
    "$sc2.Save(); " ^
    "$userPath = [Environment]::GetEnvironmentVariable('Path', 'User'); " ^
    "if ($userPath -notlike ('*' + $installDir + '*')) { [Environment]::SetEnvironmentVariable('Path', $userPath + ';' + $installDir, 'User') }"

echo.
echo =======================================================
echo    Installation Complete!
echo =======================================================
echo Application:  %APP_NAME%
echo Location:     %INSTALL_DIR%\%EXE_NAME%
echo Shortcuts:    Placed on Desktop and Start Menu
echo.
pause

