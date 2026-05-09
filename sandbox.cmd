@echo off
setlocal
rem Command Prompt wrapper; the PowerShell script contains the Windows logic.
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%~dp0sandbox.ps1" %*
exit /b %ERRORLEVEL%
