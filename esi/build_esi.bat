echo off
cls
del *.go
java -jar swagger-codegen-cli.jar generate -l go -i https://esi.tech.ccp.is/latest/swagger.json?datasource=tranquility --model-package esi

pause