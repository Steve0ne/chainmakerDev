/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package entity

type SdkConfig struct {
	ChainId     string
	OrgId       string
	UserName    string
	Tls         bool
	TlsHost     string
	Remote      string
	CaCert      []byte
	UserPrivKey []byte
	UserCert    []byte
}
