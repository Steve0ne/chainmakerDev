/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net_libp2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/net"
	"go.uber.org/zap"
)

// LibP2pSteam is stream of p2p connection
type LibP2pSteam struct {
	sync.Mutex
	rw    *bufio.ReadWriter
	log   *zap.SugaredLogger
	delim byte
}

// NewLibP2PSteam create new libp2p stream
func NewLibP2PSteam(rw *bufio.ReadWriter, delim byte) *LibP2pSteam {
	return &LibP2pSteam{
		rw:    rw,
		log:   logger.GetLogger(logger.ModuleP2P),
		delim: delim,
	}
}

// ReadStream read stream from connection
func (s *LibP2pSteam) ReadStream(ch chan net.Message) error {
	delim := s.delim
	go func(c chan net.Message) {
		for {
			str, err := s.rw.ReadString(delim)
			if err != nil {
				if err.Error() == "stream reset" {
					s.log.Debug(err)
				} else if err.Error() == "EOF" {
					return
				} else {
					s.log.Error("read string error: ", err)
				}
				return
			}
			// ignore empty string without delimit
			if str == "" {
				s.log.Warn("get empty string")
				continue
			}
			// ignore empty data with delimit
			if str == string(delim) {
				s.log.Debug("get delimit byte")
				continue
			}

			if str != string(delim) {
				// chunk delimit
				str = strings.Replace(str, string(delim), "", -1)
				s.log.Debug("new incoming stream data, put into read stream channel")
				var msg LibP2pMessage
				err := json.Unmarshal([]byte(str), &msg)
				if err != nil {
					s.log.Error("unmarshal libp2p message error ", err)
				}
				c <- &msg
			}
		}
	}(ch)

	return nil
}

// WriteStream write message to stream
func (s *LibP2pSteam) WriteStream(msg *LibP2pMessage) error {
	s.Lock()
	defer s.Unlock()
	var err error
	// marshal message
	bz, err := json.Marshal(msg)
	if err != nil {
		s.log.Error("marshall libp2p message error", err)
		return err
	}
	// write message
	_, err = s.rw.WriteString(fmt.Sprintf("%s%s", bz, string(s.delim)))
	if err != nil {
		if err.Error() == "stream reset" {
			// ignore stream reset
			s.log.Debug(err)
		} else {
			s.log.Error("write stream error: ", err)
		}
		return err
	}
	err = s.rw.Flush()
	if err != nil {
		s.log.Debug("write stream flush error: ", err)
		return err
	}
	return err
}
