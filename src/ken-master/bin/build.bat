@echo off
set GOROOT=D:\Downloads\go1.9.windows-amd64\go
set GOPATH=D:\Codes\Golang

set BUILDPATH=..\build
set date=%Date:~0,4%%Date:~5,2%%Date:~8,2%%Time:~0,2%%Time:~3,2%%Time:~6,2%
set curPath=%BUILDPATH%\ken-master-%date%

call:getReady
call:buildStart

::准备编译
:buildStart
    echo 开始编译 ...
    set goosType=windows
    set ext=.exe
    if "%1"=="linux" (
        set goosType=linux
        set ext=
    )
    set CGO_ENABLED=0
    set GOOS=%goosType%
    set GOARCH=amd64
    %GOROOT%/bin/go build -i -o %curPath%\bin\ken-master%ext% ..\src\main.go
    %GOROOT%/bin/go build -i -o %curPath%\bin\ken-master-cli%ext% ..\src\client.go
GOTO:EOF

::准备工作
:getReady
    echo 创建临时目录 ...  %curPath%
    md %curPath%
    md %curPath%\bin
    md %curPath%\conf
    md %curPath%\logs
    md %curPath%\certs
    copy ..\conf\agent.conf.exp %curPath%\conf\agent.conf
GOTO:EOF