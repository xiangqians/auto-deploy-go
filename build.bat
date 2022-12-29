::@echo off
set curDir=%~dp0
set outputDir=%curDir%build

:: copy
xcopy i18n "%outputDir%\i18n" /s /e /h /i /y
xcopy static "%outputDir%\static" /s /e /h /i /y
xcopy templates "%outputDir%\templates" /s /e /h /i /y
xcopy data "%outputDir%\data" /s /e /h /i /y
xcopy Dockerfile "%outputDir%" /y
xcopy docker_build.sh "%outputDir%" /y
xcopy docker_run.sh "%outputDir%" /y

:: go
::set outputPkg="%outputDir%/o"
::cd ./src && go build -ldflags="-s -w" -o %outputPkg%
::cd ./src && go build -ldflags="-s -w" -o %outputPkg% && upx -9 --brute %outputPkg%

:: gox
set outputPkg="%outputDir%/o_{{.OS}}_{{.Arch}}"
::cd ./src && gox -os="windows linux" -output %outputPkg%
::cd ./src && gox -osarch "windows/amd64 linux/amd64" -output %outputPkg%
cd ./src && gox -osarch "windows/amd64 linux/amd64" -ldflags="-s -w" -output %outputPkg%

pause