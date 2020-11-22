// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
serviceregister contains NacosService struct and some functions.When try to register a service,
follow these steps:
	1.Initalize the NacosService's object by calling the NacosServiceInit method.
	2.If you need to change customize the value of the NacosService, call SetServerConfig to change
	the default setting of serverConfig. Call AppendServerConfig to add additional server to the
	serverConfig. Call SetClientConfig to change the clientConfig.
	3.Call CreateClient to instantiate the client attribute of NacosService.
	4.Call RegisterService with a parameter indicating the port to register the redisrpc service
*/
package serviceregister

import (
	"fmt"
	"loginAuth/ipaddress"

	"github.com/jeanphorn/log4go"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

//NacosService have three attributes necessary for registering a service
type NacosService struct{
	clientConfig constant.ClientConfig
	serverConfig []constant.ServerConfig //serverConfig is used for finding the nacos server
	client naming_client.INamingClient
}

//Initialize the NacosService object
func (ns *NacosService)NacosServiceInit(){
	ns.clientConfig = constant.ClientConfig{
		NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", //namespace id
		TimeoutMs:           5000,     //http request out time
		ListenInterval:      10000,    //only valid in config client
		BeatInterval:        5000,     //only valid in naming client
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}
	ns.serverConfig = []constant.ServerConfig{
		{
			IpAddr: "nacos-headless",
			ContextPath: "/nacos",
			Port:   8848,
		},
	}
}

//SetServerConfig changes the value of NacosService.serverConfig
func (ns *NacosService)SetServerConfig(serverDomain, contextPath string, port uint64){
	ns.serverConfig = []constant.ServerConfig{
		{
			IpAddr: serverDomain,
			ContextPath: contextPath,
			Port: port,
		},
	}
}

//AppendServerConfig append a new serverConfig to the NacosService.serverConfig
func (ns *NacosService)AppendServerConfig(serverDomain string, port uint64){
	ns.serverConfig = append(ns.serverConfig, constant.ServerConfig{
		IpAddr: serverDomain,
		Port: port,
	})
}

//SetClientConfig change the value of NacosService.clientConfig
func (ns *NacosService)SetClientConfig(sc constant.ClientConfig){
	ns.clientConfig = sc
}

//CreateClient creates a naming client
func (ns *NacosService)CreateClient() (err error){
	ns.client, err = clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": ns.serverConfig,
		"clientConfig": ns.clientConfig,
	})
	return
}

//RegisterService register the service to nacos
func (ns *NacosService)RegisterService(port uint64) (success bool, err error) {
	ipa := &ipaddress.IPAddress{}
	ip, err := ipa.GetIpAddress()
	if err!=nil{
		log4go.Error(err)
		return
	}
	param := vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: "loginAuth",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true, //true will enable heart beat
		Metadata:    map[string]string{"zju":"ningbo"},
		//ClusterName: "loginAuth", // default value is DEFAULT
	}
	success, err = ns.client.RegisterInstance(param)
	log4go.Info(fmt.Sprintf("RegisterServiceInstance,param:%+v,result:%+v \n\n",param, success))
	if !success||err!=nil{
		log4go.Error("Fail to register service")
	}
	return
}