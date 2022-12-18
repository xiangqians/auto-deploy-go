set curDir=%~dp0
set output="%curDir%build\auto-deploy.exe"
cd ./src && go build -ldflags="-s -w" -o %output%
upx -9 --brute %output%
pause