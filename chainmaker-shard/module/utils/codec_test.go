/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"testing"
)

func TestBase64EncodeToString(t *testing.T) {
	dataString := "i am chainmaker cross chain"
	encodeToString := Base64EncodeToString([]byte(dataString))
	dataBytes, err := Base64DecodeToBytes(encodeToString)
	if err != nil {
		t.Error("decode error")
	}
	if string(dataBytes) != dataString {
		t.Error("decode data error")
	}
}

//func TestMyChannel(t *testing.T) {
//	ch := make(chan string, 100)
//	go func(c chan string) {
//		time.Sleep(time.Second)
//		c <- "zhangsan"
//	}(ch)
//	listen(ch)
//}
//
//func listen(ch chan string) {
//	select {
//	case data := <-ch:
//		fmt.Println("-----" + data)
//	}
//}
