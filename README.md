# 練習使用 grpc 在 client 和 server 端間傳遞訊息
 
 * .proto  
將需要用到的資料格式以及方法都定義在.proto檔案，接著使用以下兩句，去產生 .pb.go grpc.pb.go 這兩個檔案

```sh
$ protoc --go-grpc_out=. *.proto
$ protoc --go_out=. game.proto
```
需要更新的時候可以使用以下這行

```sh
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative game.proto
```

 * 訊息傳遞  
   * 一對一: 進來一個 message，回傳一個 message 回去
   * 多對一: 用 stream 的方式，傳入多個 message，回傳一個 message 回去
   * 一對多: 進來一個 message，用 stream 的方式，回傳多個 message 回去
   * 多對多: 用 stream 的方式，傳入多個 message，再用 stream 的方式，回傳多個 message 回去
