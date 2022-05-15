package mq

/*
 * topic
 *	订阅者ID
 * 	订阅时间
 *	订阅者的地址
 *	topic处理数
 */
type statusTopic struct {
	topic      AMQTopic
	subId      int
	subTime    int64
	subAddress string
	count      int64
}

/*MQStatus
 *	当前topic列表，
 *		每个topic的订阅者（id,订阅时间，）。
 *		每个topic的总订阅数。
 *		每个topic，最近100个消息。
 */
type statusData struct {
}
