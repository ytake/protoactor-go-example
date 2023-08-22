# チケット購入のサンプル

## 概要

Akka in ActionのコードをProto.Actorに置き換えたら、のコードです。  

流れは下記の通りです。  

Root Actor <-> Box Office <-> Ticket Seller 

場合によってはTicket SellerからBoxOfficeを経由せずにRoot Actorに直接メッセージを送ることもあり、  
Ticket Sellerアクターはイベントごとに存在します。    

## 使い方

### チケット作成

好きなイベントのチケットを作成しましょう  
イベント名をパスで指定し、ticketsに作成するチケットの枚数を指定します。  

好きなだけチケットを作ってください。  
ただし同じイベント名のチケットは作成できません。  
アクターが重複をチェックしています。

*チケットはメモリに保存されますので、停止すると消えます*

```bash
$ curl --request POST \
  --url http://127.0.0.1:8080/events/frank_zappa_live1985 \
  --header 'Content-Type: multipart/form-data' \
  --form tickets=40
```

### 全てのチケット取得

作成されたチケットを取得します。  
イベント名とチケット枚数が返ってきます。

```bash
$ curl --request GET \
  --url http://127.0.0.1:8080/events
```

### 単一のチケット取得

作成されたチケットを一件取得します。  

```bash
$ curl --request GET \
  --url http://127.0.0.1:8080/events/frank_zappa_live1985
```
