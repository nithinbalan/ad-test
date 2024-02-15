@echo off
setlocal

REM Define the path to the Goose executable if it's not in PATH
set "goosePath=goose.exe"

REM Define the connection string and migration directory
set "connectionString=sqlserver://@localhost:1433?database=blog&trustServerCertificate=true"
set "migrationDir=migrations"

REM Capture the operation argument
set "operation=%1"
set "migrationName=%2"

if "%operation%"=="" (
    echo No operation provided.
    echo Usage:
    echo   %0 createmigration [name] - Create a new migration file
    echo   %0 migrateup   - Apply all up migrations
    echo   %0 migratedown - Apply all down migrations
    exit /b 1
)

echo Operation: %operation%
echo Connection String: %connectionString%
echo Migration Directory: %migrationDir%

REM Run the Goose migration command based on the operation
if /I "%operation%"=="createmigration" (
    echo Creating new migration: %migrationName%
    %goosePath% -dir=%migrationDir% create %migrationName% sql 
)else if /I "%operation%"=="migrateup" (
    echo Running migrations up...
    %goosePath% -dir=%migrationDir% mssql "%connectionString%" up
) else if /I "%operation%"=="migratedown" (
    echo Rolling back migrations...
    %goosePath% -dir=%migrationDir% mssql "%connectionString%" down
) else (
    echo Invalid operation specified.
    exit /b 1
)

REM Check the result of the Goose command
if %ERRORLEVEL% NEQ 0 (
    echo Error: Migration %operation% failed with error level %ERRORLEVEL%.
) else (
    echo Success: Migration %operation% completed successfully.
)

endlocal