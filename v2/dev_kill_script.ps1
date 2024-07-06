Get-ChildItem * -Include *.lock -Recurse | Remove-Item
powershell "taskkill -im air.exe -im node.exe -im ftp_server.exe -im tee -im main.exe -f"

