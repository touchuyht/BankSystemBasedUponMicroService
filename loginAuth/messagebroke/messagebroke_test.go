// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package messagebroke

import (
	"fmt"
	"io/ioutil"
	"loginAuth/data"
	"testing"
)

func TestMessageBroke_SendMessage(t *testing.T) {
	var c chan int = make(chan int)
	mb := &MessageBroke{}
	mb.MessageBrokeInit([]string{"localhost:9092"},100,"updateRedis")
	go mb.MessageProducer("updateRedis")
	bytes, err := ioutil.ReadFile(data.Path("config/config.json"))
	if err!=nil{
		t.Error("Fail to read config.json")
	}
	//fmt.Println(string(bytes))
	mb.SendMessage(string(bytes))
	<-c
}

func TestMessageBroke_MessageConsumer(t *testing.T) {
	var c chan int = make(chan int)
	mb := &MessageBroke{}
	mb.MessageBrokeInit([]string{"localhost:9092"},100,"updateRedis")
	go mb.MessageProducer("updateRedis")
	go mb.MessageConsumer([]string{"updateRedis"},"group-1")
	bytes, err := ioutil.ReadFile(data.Path("config/config.json"))
	if err!=nil{
		t.Error("Fail to read config.json")
	}
	//fmt.Println(string(bytes))
	go func(){
		s:=<- mb.Value
		if s!=string(bytes){
			t.Error("Fail to consume message, Got ",s)
		}else{
			fmt.Println(string(bytes))
		}
		c<-1
	}()
	//mb.SendMessage(string(bytes))
	mb.SendMessage("hello, kafka")
	<-c
}