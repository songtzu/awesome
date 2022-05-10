package anet


type InterfaceNet interface {
	IOnInit(*Connection)      //初始化操作，比如心跳的设置...
	IOnProcessPack(pack *PackHead,connection *Connection) //处理消息
	/*
	 * this interface SHOULD NOT CALL close.
	 */
	IOnClose(err error)(tryReconnect bool )                //
	//IWrite(msg interface{}, ph *PackHead)

	IOnConnect(isOk bool)
	IOnNewConnection(connection *Connection)

}

type InterfaceServer interface {
}