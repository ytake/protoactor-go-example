# persistence event / snapshot

いわゆるイベントソーシングの実装です。  

障害が起きた場合や、システムの再起動時に  
イベントの復元と、スナップショットの復元を行います。  

イベントは永続化され、スナップショットは定期的に最新状態が永続化されます。

## イベントの永続化

永続化には例として下記のライブラリを使用しています。  

[goleveldb](https://github.com/syndtr/goleveldb)

protoactor-goの`persistence`パッケージを使用して、独自のドライバとして実装しています。  
こちらも利用例などがあるので、参考にしてください。

