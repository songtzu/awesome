@echo off



if not exist "../../proto" (
	echo "proto not found"
	pause
)


if not exist "../../protocol" (
	md "../../protocol"
)


set pc=%cd%/protoc.exe

cd ../../proto/

%pc% --go_out=../protocol *.proto

echo OK

pause