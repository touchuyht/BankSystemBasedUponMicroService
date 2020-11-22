// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Package config is use to read the environment related variables and the config json
package config

import (
	"sync"
	"loginAuth/data"
	"io/ioutil"
	"encoding/json"
	"github.com/jeanphorn/log4go"
)

var lock = &sync.Mutex{}
var Conf *Config

type RedisRPCServer struct{
	NetProtocol string `json:"netProtocol"`
	ServerAddr string  `json:"serverAddr"`
	CertFile string    `json:"certFile"`
	KeyFile string     `json:"keyFile"`
}

type MessageBroke struct{
	Address []string   `json:"address"`
	BufferSize int     `json:"bufferSize"`
	TopicName string   `json:"topicName"`
}

type RedisHelper struct {
	Protocol string    `json:"protocol"`
	Ip string          `json:"ip"`
	Port string        `json:"port"`
}

type Message struct{
	MessageId uint64   `json:"messageId"`
	Topic string       `json:"topic"`
	Receiver string    `json:"receiver"`
	IpAddr string      `json:"ipAddr"`
	Port uint64        `json:"port"`
	UserName string    `json:"userName"`
	Password string    `json:"password"`
	Token string       `json:"token"`
}

type NacosServer struct{
	IpAddr string      `json:"ipAddr"`
	ContextPath string `json:"contextPath"`
	Port uint64        `json:"port"`
	ServicePort uint64 `json:"servicePort"`
}

//Config saves the config for structs
type Config struct {
	RedisRPCServerConfig RedisRPCServer `json:"RedisRPCServer"`
	MessageBrokeConfig MessageBroke     `json:"MessageBroke"`
	RedisHelperConfig RedisHelper       `json:"RedisHelper"`
	NacosServerConfig NacosServer       `json:"NacosServer"`
	DataPush Message                    `json:"DataPush"`
	UpdateToken Message                 `json:"UpdateToken"`
}

func ConfigInstance()(*Config, error){
	if Conf == nil{
		lock.Lock()
		defer lock.Unlock()
		if Conf==nil{
			Conf = &Config{}
			err := Conf.ConfigInit()
			if err!=nil{
				log4go.Error("Fail to init Config")
				return nil, err
			}
		}else{
			log4go.Info("Config Instance already created")
		}
	}else{
		log4go.Info("Config Instance already created")
	}
	return Conf, nil
}

//ConfigInit reads config from config.json and writes the data into Config struct
func (c *Config)ConfigInit() error {
	bytes, err := ioutil.ReadFile(data.Path("config/config.json"))
	if err != nil{
		log4go.Error("Fail to open file config.json")
		return err
	}
	err = json.Unmarshal(bytes, c)
	if err != nil{
		log4go.Error("Fail to unmarshal from config.json")
		return err
	}
	return err
}