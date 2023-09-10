# proto actor go sample

[protoactor-go](https://github.com/asynkron/protoactor-go) を使ったサンプルです。

このサンプルはアクターモデル勉強会・およびハンズオンで利用するものです。   
勉強会の内容やハンズオン資料などについてはお問い合わせください。

リポジトリにHTTPとありますが、HTTPサーバ以外のサンプルの方が多いです。

| dir                      |                                    |
|:-------------------------|:-----------------------------------|
| [simple](./simple)       | アクターシステムとHTTPサーバを起動して処理するサンプル      |
| [ticket](./ticket)       | アクターシステムとHTTPサーバを起動してチケットを発行するサンプル |
| [routing](./routing)     | アクターのルーティングサンプル                    |
| [structure](./structure) | アクターの構造パターン サンプル                   |

## 使い方例

```sh
$ cd simple
$ go run main.go
```
