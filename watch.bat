@echo off

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0;%OLDGOPATH%

gin --bin bin/server.exe --path src --build src/server --appPort 10001