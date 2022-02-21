set GOOS=linux
set SRCPATH=%cd%
cd ../
cd ../
cd ../
cd ../
set WORKINGPATH=%cd%

set GOARCH=amd64

set GOPATH=%cd%
echo %GOPATH%
go build -o %SRCPATH%\amq %SRCPATH%\main.go
cd %WORKINGPATH%
echo @'build finish'
pause