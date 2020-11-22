// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package userauth

import "testing"

func TestCreateToken(t *testing.T) {
	claims := Claims{UserName: "zhangshan"}
	token, success := CreateToken(&claims)
	if !success {
		t.Error("CreateToken failed. Got ", token, "expected ")
	}
}

func TestCheckToken(t *testing.T) {
	claims := Claims{UserName: "zhangshan"}
	token, _ := CreateToken(&claims)
	if CheckToken(token)!=nil{
		t.Error("token check failed")
	}
}