package mq

import "awesome/anet"

/*AMQTopic
 *  ReserveLow part store topic for a message.
 */
type AMQTopic = uint32

type AMQProtocolSubTopic struct {
	Topics []AMQTopic `json:"topics"`
}
type AMQProtocolSubTopicAck struct {
	Status        int        `json:"status"`
	Message       string     `json:"message"`
	SucceedTopics []AMQTopic `json:"topics"`
}

type AmqAckType = uint32

type AMQCallback func(head *anet.PackHead)
