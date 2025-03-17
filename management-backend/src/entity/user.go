/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package entity

type User struct {
	Id   int64
	Name string
}

func NewUser(id int64, name string) *User {
	return &User{
		Id:   id,
		Name: name,
	}
}

func (user *User) GetName() string {
	return user.Name
}

func (user *User) GetId() int64 {
	return user.Id
}
