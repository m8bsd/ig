### Coding my own mailserver
### Install Go on FreeBSD using pkg:
```
sudo pkg update
sudo pkg install go git openssl
```
### Create Project
```
mkdir ~/mailserver
cd ~/mailserver

go mod init mailserver
go get github.com/emersion/go-smtp
```
