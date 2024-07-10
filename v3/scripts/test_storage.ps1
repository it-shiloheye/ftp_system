mkdir test_tmp/test_storage
cd test_tmp/test_storage
Get-ChildItem * -Include *.lock,data,*/data/*,*.exe -Recurse | Remove-Item -Force -Recurse
mkdir data
cd ../../peer
go build -o ../test_tmp/test_storage/ftp_server.exe .
cd ../test_tmp/test_storage 
./ftp_server.exe
