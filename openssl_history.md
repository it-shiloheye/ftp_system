## Generating Your Private Key
After deciding on a key algorithm, key size, and whether to use a passphrase, you are ready to generate your private key.

Use the following command to generate your private key using the RSA algorithm:

```
    openssl genrsa -out yourdomain.key 2048
```

This command generates a private key in your current directory named yourdomain.key (-out yourdomain.key) using the RSA algorithm (genrsa) with a key length of 2048 bits (2048). The generated key is created using the OpenSSL format called PEM.

Use the following command to view the raw, encoded contents (PEM format) of the private key:

```
    cat yourdomain.key
```
Even though the contents of the file might look like a random chunk of text, it actually contains important information about the key.

Use the following command to decode the private key and view its contents:

```
    openssl rsa -text -in yourdomain.key -noout
```
The -noout switch omits the output of the encoded version of the private key.

## Extracting Your Public Key
The private key file contains both the private key and the public key. You can extract your public key from your private key file if needed.

Use the following command to extract your public key:

```
    openssl rsa -in yourdomain.key -pubout -out yourdomain_public.key
```

## Creating Your CSR
After generating your private key, you are ready to create your CSR. The CSR is created using the PEM format and contains the public key portion of the private key as well as information about you (or your company).

Use the following command to create a CSR using your newly generated private key:

```
    openssl req -new -key yourdomain.key -out yourdomain.csr
```
After entering the command, you will be asked series of questions. Your answers to these questions will be embedded in the CSR. Answer the questions as described below:

```
    Country Name (2 letter code)	The two-letter country code where your company is legally located.
    State or Province Name (full name)	The state/province where your company is legally located.
    Locality Name (e.g., city)	The city where your company is legally located.
    Organization Name (e.g., company)	Your company's legally registered name (e.g., YourCompany, Inc.).
    Organizational Unit Name (e.g., section)	The name of your department within the organization. (You can leave this option blank; simply press Enter.)
    Common Name (e.g., server FQDN)	The fully-qualified domain name (FQDN) (e.g., www.example.com).
    Email Address	Your email address. (You can leave this option blank; simply press Enter.)
    A challenge password	Leave this option blank (simply press Enter).
    An optional company name	Leave this option blank (simply press Enter).
```
Some of the above CSR questions have default values that will be used if you leave the answer blank and press Enter. These default values are pulled from the OpenSSL configuration file located in the OPENSSLDIR (see Checking Your OpenSSL Version). If you want to leave a question blank without using the default value, type a "." (period) and press Enter.

## Using the -subj Switch
Another option when creating a CSR is to provide all the necessary information within the command itself by using the -subj switch.

Use the following command to disable question prompts when generating a CSR:

```
    openssl req -new -key yourdomain.key -out yourdomain.csr \
        -subj "/C=US/ST=Utah/L=Lehi/O=Your Company, Inc./OU=IT/CN=yourdomain.com"
```
This command uses your private key file (-key yourdomain.key) to create a new CSR (-out yourdomain.csr) and disables question prompts by providing the CSR information (-subj).

## Creating Your CSR with One Command
Instead of generating a private key and then creating a CSR in two separate steps, you can actually perform both tasks at once.

Use the following command to create both the private key and CSR:

```
    openssl req -new \
        -newkey rsa:2048 -nodes -keyout yourdomain.key \
        -out yourdomain.csr \
        -subj "/C=US/ST=Utah/L=Lehi/O=Your Company, Inc./OU=IT/CN=yourdomain.com"
```
This command generates a new private key (-newkey) using the RSA algorithm with a 2048-bit key length (rsa:2048) without using a passphrase (-nodes) and then creates the key file with a name of yourdomain.key (-keyout yourdomain.key).

The command then generates the CSR with a filename of yourdomain.csr (-out yourdomain.csr) and the information for the CSR is supplied (-subj).

Note: While it is possible to add a subject alternative name (SAN) to a CSR using OpenSSL, the process is a bit complicated and involved. If you do need to add a SAN to your certificate, this can easily be done by adding them to the order form when purchasing your DigiCert certificate.

## Verifying CSR Information
After creating your CSR using your private key, we recommend verifying that the information contained in the CSR is correct and that the file hasn't been modified or corrupted.

Use the following command to view the information in your CSR before submitting it to a CA (e.g., DigiCert):

```
    openssl req -text -in yourdomain.csr -noout -verify
```
The -noout switch omits the output of the encoded version of the CSR. The -verify switch checks the signature of the file to make sure it hasn't been modified.

Running this command provides you with the following output:
```
verify OK
Certificate Request:
    Data:
        Version: 0 (0x0)
        Subject: C=US, ST=Utah, L=Lehi, O=Your Company, Inc., OU=IT, CN=yourdomain.com
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (2048 bit)
                Modulus:
                    00:bb:31:71:40:81:2c:8e:fb:89:25:7c:0e:cb:76:
                    [...17 lines removed]
                Exponent: 65537 (0x10001)
        Attributes:
            a0:00
    Signature Algorithm: sha256WithRSAEncryption
         0b:9b:23:b5:1f:8d:c9:cd:59:bf:b7:e5:11:ab:f0:e8:b9:f6:
         [...14 lines removed]
```
On the first line of the above output, you can see that the CSR was verified (verify OK). On the fourth line, the Subject: field contains the information you provided when you created the CSR. Make sure this information is correct.

If any of the information is wrong, you will need to create an entirely new CSR to fix the errors. This is because CSR files are digitally signed, meaning if even a single character is changed in the file it will be rejected by the CA.