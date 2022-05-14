package anet

type InterfaceNet interface {
	IOnInit(*Connection)                                   //初始化操作，比如心跳的设置...
	IOnProcessPack(pack *PackHead, connection *Connection) //处理消息
	/*
	 * this interface SHOULD NOT CALL close.
	 */
	IOnClose(err error) (tryReconnect bool) //
	//IWrite(msg interface{}, ph *PackHead)

	IOnConnect(isReconnect bool)
	/*IOnNewConnection
	 * 此接口是服务端响应一个新的客户端请求，创建出一个新的客户端连接，回调给服务端实现使用的。
	 */
	IOnNewConnection(connection *Connection)
}

type InterfaceServer interface {
}
