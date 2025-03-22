# pipeantenna

## Usage

```
pipeantenna ~/go/src/github.com/ashmeet28/pipeantenna/web_pages/index.html /mnt/t/pipeantenna_https_crt /mnt/t/pipeantenna_https_key /mnt/t/pipeantenna_auth_key_file 8080
```

```
openssl req -x509 -newkey rsa:4096 -noenc -keyout /mnt/t/pipeantenna_https_key -out /mnt/t/pipeantenna_https_crt
```

```
openssl pkey -in /mnt/t/pipeantenna_https_key -pubout -outform DER | sha256sum
```
