# 游戏区服报警

#### 安装说明
编译Linux,Windows和Mac环境下可执行程序

```go
git clone https://github.com/kyirving/game.git
cd game
go mod tidy 
```

###### linux
```go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/exec cmd/main.go
```
###### windows
```go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows/exec cmd/main.go
```
###### mac
```go
go build -o bin/mac/exec cmd/main.go
```