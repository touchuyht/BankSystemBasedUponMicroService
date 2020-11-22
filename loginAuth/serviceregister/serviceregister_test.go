// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serviceregister

import "testing"

func TestNacosService_CreateClient(t *testing.T) {
	/*nacosService := NacosService{}
	nacosService.NacosServiceInit()
	nacosService.SetServerConfig("console1.nacos.io",80)
	err := nacosService.CreateClient()
	if err != nil{
		t.Error("Failed to create naming client")
	}*/
}

func TestNacosService_RegisterService(t *testing.T) {
	nacosService := NacosService{}
	nacosService.NacosServiceInit()
	nacosService.SetServerConfig("console.nacos.io","/nacos",80)
	err := nacosService.CreateClient()
	if err != nil{
		t.Error("Failed to create naming client")
	}
	success, err := nacosService.RegisterService(2550)
	if !success||err!=nil{
		t.Error("Fail to register service")
	}
}