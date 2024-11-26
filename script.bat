@echo off
setlocal

REM Define the path to the Goose executable if it's not in PATH
set "goosePath=goose.exe"

REM Define the connection string and migration directory
set "connectionString=sqlserver://@localhost:1433?database=jethur_test&trustServerCertificate=true"
set "migrationDir=migrations"
set "seederDir=seeders"

REM Capture the operation argument
set "operation=%1"
set "name=%2"

if "%operation%"=="" (
    echo No operation provided.
    echo Usage:
    echo   %0 createmigration [name] - Create a new migration file
    echo   %0 migrateup   - Apply all up migrations
    echo   %0 migratedown - Apply all down migrations
    echo   %0 createseeder [name] - Create a new seeder file
    echo   %0 seed - Seed the database
    echo   %0 unseed - Remove seed data from the database
    exit /b 1
)

echo Operation: %operation%
echo Connection String: %connectionString%
echo Migration Directory: %migrationDir%
echo Seeder Directory: %seederDir%

REM Run the Goose migration command based on the operation
if /I "%operation%"=="createmigration" (
    echo Creating new migration: %name%
    %goosePath% -dir=%migrationDir% create %name% sql 
)else if /I "%operation%"=="migrateup" (
    echo Running migrations up...
    %goosePath% -dir=%migrationDir% mssql "%connectionString%" up > %migrationDir%/output.txt 2>&1
) else if /I "%operation%"=="migratedown" (
    echo Rolling back migrations...
    %goosePath% -dir=%migrationDir% mssql "%connectionString%" down
) else if /I "%operation%"=="createseeder" (
    echo Creating new seeder: %name%
    %goosePath% -dir=%seederDir% create %name% sql
) else if /I "%operation%"=="seed" (
    echo Running seeders...
    %goosePath%  -table=goose_seeder_versions -dir=%seederDir% mssql "%connectionString%"  up > %seederDir%/output.txt 2>&1
) else if /I "%operation%"=="unseed" (
    REM Run the seeder down migrations to remove the seeded data from the database.
    echo Removing seeded data...
    %goosePath% -table=goose_seeder_versions -dir=%seederDir% mssql "%connectionString%" down
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