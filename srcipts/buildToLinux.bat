@echo off
setlocal

:: 设置目标操作系统和架构
set GOOS=linux
set GOARCH=amd64

:: 设置 Go 项目的入口文件（main.go）和输出文件名
set MAIN_FILE=..\cmd\api\main.go
set OUTPUT_FILE=..\..\bin\hotsearch-linux-amd64

:: 清理之前的构建结果
echo Clean up previous build results
:: del %OUTPUT_FILE%.exe
del %OUTPUT_FILE%

:: 交叉编译
echo Start cross-compiling
set CGO_ENABLED=0
go build -o %OUTPUT_FILE% -ldflags "-s -w" %MAIN_FILE%

:: 检查构建是否成功
if exist %OUTPUT_FILE% (
    echo Build successfully
    echo Binary files have been saved to %OUTPUT_FILE%
) else (
    echo Build failed
)

endlocal
pause