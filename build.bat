::@echo off
set curDir=%~dp0
echo curDir: %curDir%

set outputDir=%curDir%build
echo outputDir: %outputDir%

:: rd /?
:: 删除 build 目录
rd /s /q %outputDir%
echo rd: %outputDir%

:: copy
xcopy i18n "%outputDir%/i18n" /s /e /h /i /y
xcopy static "%outputDir%/static" /s /e /h /i /y
xcopy templates "%outputDir%/templates" /s /e /h /i /y
xcopy data "%outputDir%/data" /s /e /h /i /y
xcopy script "%outputDir%/script" /s /e /h /i /y

:: pkgName
for /F %%i in ('go env GOOS') do (set os=%%i)
for /F %%i in ('go env GOARCH') do (set arch=%%i)
set pkgName=o_%os%_%arch%.exe
echo pkgName: %pkgName%

:: go
set pkgPath="%outputDir%/%pkgName%"
cd ./src && go build -ldflags="-s -w" -o %pkgPath%
::cd ./src && go build -ldflags="-s -w" -o %pkgPath% && upx -9 --brute %pkgPath%
echo pkgPath: %pkgPath%

:: startup.bat
set startupPath=%outputDir%/startup.bat
echo :: startup.bat > %startupPath%
echo %pkgName% >> %startupPath%
echo pause >> %startupPath%
echo startupPath: %startupPath%

pause