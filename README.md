# go-server-for-libDB

図書貸出アプリのサーバプログラムです．  
`golang，echo` で作られています．

## サーバの起動方法
サーバは，以下のコマンドでクイック起動できます( `...` はリンクする他の `.go` プログラム)．
```
go run main_server.go ... 
```

また，ビルドとサーバのバックグラウンド実行は以下のコマンドで行います( `...` はリンクする他の `.go` プログラム)．
```
go build main_server.go ...
./main_server.go &
```

バックグラウンド実行したサーバを終了するには，実行中のプロセスを以下のコマンドで探して，2項目のPIDを指定して終了します．
```
ps aux | grep main_server
kill [PID(.main_serverプロセスの2項目の数字))]
```

`master` ブランチは，公開サーバ用のブランチです．
公開サーバでは以下のスクリプトを使用して，サーバの立ち上げと停止をしてください．
- `exec_server_background.sh` を実行して，サーバを常時バックグラウンド実行
- `kill_background_server.sh` を実行して，バックグラウンドしていたサーバを停止

###  サーバへのアクセス例
APIサーバへのアクセス方法は，`test_request.http` を例とします( `localhost` は適宜変えてください)．
