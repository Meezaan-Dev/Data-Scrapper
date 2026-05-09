@echo off
setlocal
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%~dp0sandbox.ps1" %*
exit /b %ERRORLEVEL%
