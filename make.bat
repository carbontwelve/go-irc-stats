@ECHO OFF

REM How to use:
REM To begin with execute make build && make copy && make generate-log
REM This will initiate the logstats folder in the bin directory. Afterwards you can then
REM run make build && make run to build and then execute from source.

REM Set executable name
SET PACKAGENAME=main.exe

REM Set GOBIN
SET GOBIN=%GOPATH%\bin\logstats

REM Set ld flags containing current version and git hash
SET VERSION=1.0.0
for /f "delims=" %%a in ('git rev-parse HEAD') do SET BUILD=%%a
SET LDFLAGS=-ldflags "-X main.Version=%VERSION% -X main.Build=%BUILD%"

REM Set current working directory
SET CWD=%cd%

REM Get first argument passed by command line
set ARG1=%1

REM Do action if found
if "%ARG1%" == "clean" (
    echo Cleaning Build Directory
    CALL :clean
) ELSE IF "%ARG1%" == "copy" (
    echo Copying extra files
    CALL :copy
) ELSE IF "%ARG1%" == "build" (
    echo Building Project
    CALL :build
    echo Copying extra files
    CALL :copy
) ELSE IF "%ARG1%" == "update" (
    echo Updating dependencies
    CALL :update
) ELSE IF "%ARG1%" == "run" (
    CALL :run
) ELSE IF "%ARG1%" == "build-test-log" (
    echo Building Test Log
    CALL :generate-log
) ELSE (
    ECHO Incorrect arguments passed [%ARG1%]
)

EXIT /B 0

:run
%GOBIN%\%PACKAGENAME% -d %GOBIN%
EXIT /B 0

:update
go get -v ./...
EXIT /B 0

:clean
if exist %GOBIN% (
    RD /S /Q %GOBIN%
)
EXIT /B 0

:build
@ECHO ON
go install %LDFLAGS% %CWD%\cmd\ircstats\main.go
@ECHO OFF
EXIT /B 0

:copy
cd extra
copy "config.yaml" "%GOBIN%\config.yaml" > nul
copy "template.html" "%GOBIN%\template.html" > nul
cd ..
EXIT /B 0

:generate-log
php %CWD%\extra\createtestlog.php > %GOBIN%\irctest.log
EXIT /B 0