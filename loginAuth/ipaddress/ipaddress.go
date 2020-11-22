// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
IPAddress can obtain the first usable network device's ip address
*/
package ipaddress

import (
	"net"

	"github.com/jeanphorn/log4go"
)

//IPAddress used for obtaining ipaddress
type IPAddress struct{
	IpAddress string
	IpValid bool
}

//GetIpAddress use the net package to obtain the first workable net device's ipaddress
func (ip *IPAddress)GetIpAddress() (ipAddress string, err error){
	netInterfaces, err := net.Interfaces()
	if err != nil{
		log4go.Error(err)
		return
	}
	Loop:
	for i := 0; i < len(netInterfaces); i++{
		if(netInterfaces[i].Flags & net.FlagUp) != 0{
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs{
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback(){
					if ipnet.IP.To4()!=nil{
						ipAddress = (ipnet.IP.String())
						ip.IpAddress = ipAddress
						ip.IpValid = true
						break Loop
					}
				}
			}
		}
	}
	return
}