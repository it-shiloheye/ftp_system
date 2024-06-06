powershell "taskkill -im ftp_server.exe  -im ftp_client.exe  -im air.exe -f"
start powershell "npm run dev:client"
start powershell "npm run dev:server"
