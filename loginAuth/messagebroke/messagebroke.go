// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
MessageBroke has a string chan which you use it for sending message.
Use it in the following order:MessageBrokeInit->go MessageProducer()->SendMessage
*/
package messagebroke

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/jeanphorn/log4go"
)

//MessageBroke is the type which you can instantiate to send message
type MessageBroke struct{
	address []string      //address is the ip and port of kafka brokers
	message chan string   //message is the channel to receiver string from other goroutine
	bufferSize int        //bufferSize is the size of message buffer size, if set to 0 would cause the sending goroutine
	                      //waiting for other goroutine to read from another side of the channel
	topicName string      //topicName represents the topic name in kafka
	Value chan string     //chan is the channel to send string to other goroutine
	stopConsumer chan int //this chan is use for stopping consumer
}

//MessageBrokeInit would initialize the attributes of MessageBroke's instantiates
func (mb *MessageBroke)MessageBrokeInit(address []string, bufferSize int, topicName string){
	mb.address = address
	mb.bufferSize = bufferSize
	mb.message = make(chan string, mb.bufferSize)
	mb.Value = make(chan string, mb.bufferSize)
	mb.topicName = topicName
	mb.stopConsumer = make(chan int)
}

func (mb *MessageBroke)SetTopicName(topicName string){
	mb.topicName = topicName
}

func (mb *MessageBroke)SetBufferSize(bufferSize int){
	mb.bufferSize = bufferSize
}

//SendMessage receives a string and send it to kafka brokers
func (mb *MessageBroke)SendMessage(message string) {
	mb.message <- message
}

func (mb *MessageBroke)StopConsumer(){
	mb.stopConsumer<-1
}

/*//GetMessage receives message from mb.value channel and return it as a string
func (mb *MessageBroke)GetMessage() string{
	return <-mb.value
}*/

//MessageProducer receives the sending message from the string channel, but when there is no message in the channel , it would pend
//And it will make a goroutine to receive the response of kafka brokers
func (mb *MessageBroke)MessageProducer(topic string){
	config:=sarama.NewConfig()
	config.Producer.RequiredAcks=sarama.WaitForAll //Wait for the signal which indicates the succession in saving all the duplicates
	config.Producer.Partitioner=sarama.NewRandomPartitioner //Send to partitions randomly
	config.Producer.Return.Successes=true //whether wait or not wait for the response from brokers which indicates the sending failure or success
	config.Producer.Return.Errors=true
	config.Version=sarama.V0_10_0_0 //config.Version should be in accordance with kafka version
	var value string = ""

	log4go.Info("start to make produce")
	producer, err:=sarama.NewAsyncProducer(mb.address, config)
	if err!=nil {
		log4go.Error(err)
		return
	}
	defer producer.AsyncClose()

	log4go.Info("start goroutine")
	go func(p sarama.AsyncProducer){
		for{
			select{
			case <-p.Successes():
			case fail:=<-p.Errors():
				log4go.Error("err:",fail.Err)
			}
		}
	}(producer)

	for {
		msg:=&sarama.ProducerMessage{
			Topic: topic,
		}
		value = <-mb.message
		msg.Value = sarama.ByteEncoder(value)
		producer.Input()<-msg
	}
}

//MessageConsumer receives message from kafka and sends it through a string chan to other goroutine and it create two
//goroutines to deal errors and notifications
func (mb *MessageBroke)MessageConsumer(topics []string, groupId string) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// init consumer
	consumer, err := cluster.NewConsumer(mb.address, groupId, topics, config)
	if err != nil {
		log4go.Error(fmt.Sprintf("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err))
		return
	}
	defer consumer.Close()

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log4go.Error(fmt.Sprintf("%s:Error: %s\n", groupId, err))
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log4go.Info(fmt.Sprintf("%s:Rebalanced: %+v \n", groupId, ntf))
		}
	}()

	// consume messages, watch signals
	var successes int
	Loop:
	for{
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				log4go.Info(fmt.Sprintf("%s:%s/%d/%d\t%s\t%s\n", groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value))
				//fmt.Fprintf(os.Stdout, "%s:%s/%d/%d\t%s\t%s\n", groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				mb.Value <- string(msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
				successes++
			}
		case <-mb.stopConsumer:
			break Loop
		}
	}
	log4go.Info(fmt.Sprintf("%s consume %d messages \n", groupId, successes))
}