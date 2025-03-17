/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

import "strconv"

// 安装合约时会执行此方法，必须
//export init_contract
func initContract() {}

// 升级合约时会执行此方法，必须
//export upgrade
func upgrade() {}

const (
	KeyLocal     = "Local"
	FieldBalance = "Balance"
)

//export Show
func Show() {
	if b, resultCode := GetStateByte(KeyLocal, FieldBalance); resultCode != SUCCESS {
		ErrorResult("failed to get state")
		return
	} else {
		SuccessResult(string(b))
		return
	}
}

//export Plus
func Plus() {
	// get number
	if num, resultCode := Arg("number"); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to get number")
		return
	} else {
		if b, resultCode := GetStateByte(KeyLocal, FieldBalance); resultCode != SUCCESS {
			ErrorResult("failed to get state")
			return
		} else {
			number, _ := strconv.Atoi(num)
			balance, _ := strconv.Atoi(string(b))
			balance += number
			if balance < 0 {
				ErrorResult("balance less than 0")
				return
			}
			if balance > 10000 {
				ErrorResult("balance greater than 10000")
				return
			}
			b := strconv.Itoa(balance)
			PutStateByte(KeyLocal, FieldBalance, []byte(b))
			SuccessResult(b)
			return
		}
	}
}

//export Minus
func Minus() {
	// get number
	if num, resultCode := Arg("number"); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to get number")
		return
	} else {
		if b, resultCode := GetStateByte(KeyLocal, FieldBalance); resultCode != SUCCESS {
			ErrorResult("failed to get state")
			return
		} else {
			number, _ := strconv.Atoi(num)
			balance, _ := strconv.Atoi(string(b))
			balance -= number
			if balance < 0 {
				ErrorResult("balance less than 0")
				return
			}
			if balance > 10000 {
				ErrorResult("balance greater than 10000")
				return
			}
			b := strconv.Itoa(balance)
			PutStateByte(KeyLocal, FieldBalance, []byte(b))
			SuccessResult(b)
			return
		}
	}
}

//export Reset
func Reset() {
	PutStateByte(KeyLocal, FieldBalance, []byte("1000"))
	SuccessResult("1000")
	return
}

func main() {}
