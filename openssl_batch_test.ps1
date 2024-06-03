# winget install -e --id ShiningLight.OpenSSL
# set PATH=%PATH%;"C:\Program Files\OpenSSL-Win64\bin"

openssl genrsa -out tmp/yourdomain.key 2048 

openssl req -new -key tmp/yourdomain.key -out tmp/yourdomain.csr -subj "/C=US/ST=Utah/L=Lehi/O=Your Company, Inc./OU=IT/CN=yourdomain.com"