@echo off



if not exist "../../proto" (
	echo "proto not found"
	pause
)


if not exist "../../protocol" (
	md "../../protocol"
)


set pc=%cd%\protoc.exe
cd ../
set out_dir=%cd%\protocol

if not exist %out_dir% (
	md %out_dir%
)


cd ../proto/

%pc% --go_out=%out_dir% *.proto

echo OK

pause