::@echo off
set curDir=%~dp0
set outputDir=%curDir%build
set outputPkg="%outputDir%/auto-deploy.exe"
xcopy i18n "%outputDir%\i18n" /s /e /h /i /y && xcopy static "%outputDir%\static" /s /e /h /i /y && xcopy templates "%outputDir%\templates" /s /e /h /i /y && xcopy data "%outputDir%\data" /s /e /h /i /y && cd ./src && go build -ldflags="-s -w" -o %outputPkg% && upx -9 --brute %outputPkg%
::xcopy i18n "%outputDir%\i18n" /s /e /h /i /y && xcopy static "%outputDir%\static" /s /e /h /i /y && xcopy templates "%outputDir%\templates" /s /e /h /i /y && xcopy data "%outputDir%\data" /s /e /h /i /y && cd ./src && go build -ldflags="-s -w" -o %outputPkg%
pause