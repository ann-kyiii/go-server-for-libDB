# go-server-for-libDB

図書貸出アプリのサーバプログラムです．  
`golang，echo` で作られています．

## サーバの起動方法
サーバは，以下のコマンドでクイック起動できます．
```
go run main_server.go
```

また，ビルドとサーバのバックグラウンド実行は以下のコマンドで行います．
```
go build main_server.go
./main_server.go &
```

バックグラウンド実行したサーバを終了するには，実行中のプロセスを以下のコマンドで探して，2項目のPIDを指定して終了します．
```
ps aux | grep main_server
kill [PID(.main_serverプロセスの2項目の数字))]
```

###  サーバへのアクセス例
APIサーバへのアクセス方法は，`test_command` を例とします．

例：
```
curl [サーバのグローバルIPアドレス]:1313 
```