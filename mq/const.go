package mq

const defaultAMQChanSize = 10000

const (
	AmqAckTypeSuccess AmqAckType = 0
	AmqAckTypeTimeout AmqAckType = 1
)

const (
	AMQCmdDefPub         = 0
	AMQCmdDefSubTopic    = 1
	AMQCmdDefSubTopicAck = 2

	/*
	 *AmqCmdDefUnreliable2All
	 *	不可靠发布（无回包），所有订阅者都会收到此message.
	 *	仅仅使用类似日志，统计类业务，不同的sub订阅同一个topic，加工整理成不同的统计结果。
	 *	业务不依赖此消息队列是否有启动进程处理业务。
	 *	此业务的MQ会缓存Mn 条消息，如有注册可用
	***********/
	AmqCmdDefUnreliable2All = 10

	/*AmqCmdDefUnreliable2RandomOne
	 *	不可靠发布（无回包），N个订阅者中某一个（随机选择）会收到此message，适用场合不多。
	***/
	AmqCmdDefUnreliable2RandomOne = 11

	/*AmqCmdDefReliable2RandomOne
	 *	n个订阅者，所有订阅者都能处理此业务，但是需要有且仅有一个sub处理。
	 *	N个sub中某一个（随机）会收到此message，订阅者需要在n秒内回复处理结果.
	 *	否则，订阅者会把消息转发给剩余订阅者处理。
	 *	直到n秒内有回复的订阅消息，或者返回超时错误给puber。此业务模型面向无状态业务。
	 */
	AmqCmdDefReliable2RandomOne = 12

	/*AmqCmdDefReliable2SpecOne
	 *	n个订阅者，但是其中仅有一个订阅者能处理此业务.
	 *	而发布方和mq proxy节点均不清楚哪个sub节点能处理此业务。
	 *	mq proxy把此msg扇出给所有sub，处理此业务的sub需在处理完业务之后回复结果，或mq proxy超时返回。
	 *	注意使用此业务模型必须自行确保至多唯一一个sub能处理此业务。此业务模型不保证不被重复处理。
	******/
	AmqCmdDefReliable2SpecOne = 13
)
