@ECHO OFF
SET BINARY=logstats.exe
SET VERSION=1.0.0
for /f "delims=" %%a in ('git rev-parse HEAD') do SET BUILD=%%a
SET LDFLAGS=-ldflags "-X main.Version=%VERSION% -X main.Build=%BUILD%"

set ARG1=%1

if "%ARG1%" == "clean" (
    CALL :clean
) ELSE IF "%ARG1%" == "build" (
    CALL :build
) ELSE IF "%ARG1%" == "generate-log" (
    CALL :generate-log
) ELSE IF "%ARG1%" == "run" (
    CALL :run
) ELSE IF "%ARG1%" == "clean-build" (
    CALL :clean
    CALL :build
) ELSE (
    ECHO Incorrect arguments passed [%ARG1%]
)
EXIT /B 0

:clean
if exist bin (
    RD /S /Q "bin"
)
EXIT /B 0

:build
@ECHO ON
cd cmd\ircstats
go get -t -v ./...
go build %LDFLAGS% -o ..\..\bin\%BINARY%
@ECHO OFF
cd ..\..\extra
copy "config.yaml" "..\bin\config.yaml" > nul
copy "template.html" "..\bin\template.html" > nul
cd ..
EXIT /B 0

:run
bin\%BINARY% -d ./bin
EXIT /B 0

:generate-log
php extra\createtestlog.php > .\bin\irctest.log
EXIT /B 0