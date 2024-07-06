taskkill -im ftp_server.exe -f
rm tmp/ftp_server.exe
go build -o ./tmp/ftp_server.exe .
tmp\\ftp_server.exe