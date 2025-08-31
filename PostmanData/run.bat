
@echo off
REM Usage: runNewman.bat "FolderName"

IF "%~1"=="" (
    echo Please provide a folder name as argument.
    echo Example: runNewman.bat Register
    exit /b 1
)

set FOLDER=%~1

newman run https://api.postman.com/collections/29129091-bf2b9350-9f28-4f4c-b3e9-f47b356e1391?access_key=PMAT-01K3XV4C4DDW78Z2WV5NPYKKW1 --folder "%FOLDER%" -d ".\dataset.csv"
