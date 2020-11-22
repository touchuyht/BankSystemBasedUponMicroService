package config

import (
	"fmt"
	"testing"
)

func TestConfig_ConfigInit(t *testing.T) {
	c, err := ConfigInstance()
	if err!=nil{
		t.Error("Fail to init config")
	}
	fmt.Println(c)
	if c.RedisHelperConfig.Protocol!="tcp" || c.MessageBrokeConfig.Address[0]!="kafka-headless:9092" || c.DataPush.UserName!="zhangshan"{
		t.Error("Fail to read config.json as struct")
	}
}