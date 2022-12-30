::@echo off
set curDir=%~dp0
set outputDir=%curDir%build

:: rd /?
:: 删除 build 目录
rd /s /q %outputDir%

:: copy
xcopy i18n "%outputDir%/i18n" /s /e /h /i /y
xcopy static "%outputDir%/static" /s /e /h /i /y
xcopy templates "%outputDir%/templates" /s /e /h /i /y
xcopy data "%outputDir%/data" /s /e /h /i /y
xcopy script "%outputDir%/script" /s /e /h /i /y

:: go
::set outputPkg="%outputDir%/o"
::cd ./src && go build -ldflags="-s -w" -o %outputPkg%
::cd ./src && go build -ldflags="-s -w" -o %outputPkg% && upx -9 --brute %outputPkg%

:: gox
set outputPkg="%outputDir%/o_{{.OS}}_{{.Arch}}"
::cd ./src && gox -os="windows linux" -output %outputPkg%
cd ./src && gox -osarch "windows/amd64 linux/amd64" -ldflags="-s -w" -output %outputPkg%

:: startup.bat
set startupName=%outputDir%/startup.bat
echo :: startup.bat > %startupName%
echo o_windows_amd64.exe >> %startupName%
echo pause >> %startupName%

:: startup.sh
set startupName=%outputDir%/startup.sh
echo # startup.sh > %startupName%
echo ./o_linux_amd64 >> %startupName%

pause