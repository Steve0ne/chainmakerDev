/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net_libp2p

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/libp2p/go-libp2p-core/crypto"
)

const pemBegin = "-----BEGIN"

// PrivateKeyFromPEM load data from pem and create private key
func PrivateKeyFromPEM(raw []byte, pwd []byte) (crypto.PrivKey, error) {
	var err error

	if len(raw) <= 0 {
		return nil, fmt.Errorf("PEM is nil")
	}

	if !strings.Contains(string(raw), pemBegin) {
		keyBytes, err := hex.DecodeString(string(raw))
		if err != nil {
			return nil, fmt.Errorf("fail to decode Secp256k1 public key: [%v]", err)
		}
		return PrivateKeyFromDER(keyBytes)
	}

	block, _ := pem.Decode(raw)
	if block == nil {
		return PrivateKeyFromDER(raw)
	}

	plain := block.Bytes
	if x509.IsEncryptedPEMBlock(block) {
		if len(pwd) <= 0 {
			return nil, fmt.Errorf("missing password for encrypted PEM")
		}

		plain, err = x509.DecryptPEMBlock(block, pwd)
		if err != nil {
			return nil, fmt.Errorf("fail to decrypt PEM: [%s]", err)
		}
	}

	return PrivateKeyFromDER(plain)
}

// PrivateKeyFromDER load data from der and create private key
func PrivateKeyFromDER(der []byte) (crypto.PrivKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		if priv, _, err := crypto.KeyPairFromStdKey(key); err == nil {
			return priv, nil
		}
	}

	if key, err := x509.ParseECPrivateKey(der); err == nil {
		if priv, _, err := crypto.KeyPairFromStdKey(key); err == nil {
			return priv, nil
		}
	}

	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		if priv, _, err := crypto.KeyPairFromStdKey(key); err == nil {
			return priv, nil
		}
	}

	Secp256k1Key, _ := btcec.PrivKeyFromBytes(btcec.S256(), der)
	key := Secp256k1Key.ToECDSA()
	if priv, _, err := crypto.KeyPairFromStdKey(key); err == nil {
		return priv, nil
	}
	return nil, fmt.Errorf("get priv from der error")
}
