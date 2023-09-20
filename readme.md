# proto actor go sample

[protoactor-go](https://github.com/asynkron/protoactor-go) を使ったサンプルです。

Akka in ActionのコードをProto.Actorに置き換えたら、のサンプルになっていますので、　  
翻訳版のAkka実践バイブルなどをお持ちの方は、　  
そちらと合わせてご覧ください。

Proto Actor(Go)は、Akka/Pekkoのようなアクターモデルを実現するためのライブラリで、  
Akka Typedが導入される前のAkkaのようなAPIを提供しています。  
細部は異なりますが、クラスタをはじめとした機能はほぼ同じです。

物理的に異なるサーバ上で動作するアクターも、同一アプリケーションのように扱うことができます。  
また、アクターのメッセージの送受信は、同期的にも非同期的にも行うことができますので、  
アクターモデルによる並列処理などの学習にどうぞ。  

このサンプルはアクターモデル勉強会・およびハンズオンで利用するものです。  
勉強会の内容やハンズオン資料などについてはお問い合わせください。

| dir                          |                                    |
|:-----------------------------|:-----------------------------------|
| [simple](./simple)           | アクターシステムとHTTPサーバを起動して処理するサンプル      |
| [ticket](./ticket)           | アクターシステムとHTTPサーバを起動してチケットを発行するサンプル |
| [routing](./routing)         | アクターのルーティングサンプル                    |
| [structure](./structure)     | アクターの構造パターン サンプル                   |
| [persistence](./persistence) | アクターの永続化 サンプル                      |

## 使い方例

```sh
$ cd calculator
$ go run main.go
```
