// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipaddress

import "testing"

func TestIPAddress_GetIpAddress(t *testing.T) {
	ipa := &IPAddress{}
	if ip,_ := ipa.GetIpAddress(); ip!="172.20.10.5"{
		t.Error("Got wrong ip")
	}
}