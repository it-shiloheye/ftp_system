Get-ChildItem * -Include *.lock -Recurse | Remove-Item
powershell "taskkill -im air.exe -im node.exe -im ftp_server.exe -im ftp_client.exe -im tee -f"

