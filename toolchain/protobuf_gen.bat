@echo off



if not exist "../../proto" (
	echo "proto not found"
	pause
)


if not exist "../pb_protocol" (
	md "../pb_protocol"
)

set original=%cd%
set pc=%cd%\protoc.exe
cd ../
set out_dir=%cd%\pb_protocol

if not exist %out_dir% (
	md %out_dir%
)

cd ../proto/
set protoDir=%cd%
cd %original%

protoc.exe --go_out=%out_dir% --proto_path=%protoDir%  %protoDir%\*.proto


echo OK

pause

