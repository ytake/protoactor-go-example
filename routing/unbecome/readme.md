# UnBecome Router

HTTPのルーターではなく、アクターの状態によりルーティングするサンプルです。  

[Become](../become) と動作は同じですが、 `unbecome` で状態を戻す点が異なります。
アクターの状態でルーティングする場合にBecomeStackedでルーティングを追加すると、  
UnbecomeStackedで一個状態を戻すことができます。  
追加されたルーティングを削除するということになるため、挙動としては Become と同じです。

