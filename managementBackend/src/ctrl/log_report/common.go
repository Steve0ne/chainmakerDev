/*
   Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
     SPDX-License-Identifier: Apache-2.0
*/
package log_report

import (
	loggers "management_backend/src/logger"
	"management_backend/src/sync"
	"time"
)

var (
	log = loggers.GetLogger(loggers.ModuleWeb)
)

const (
	NO_AUTO = iota
	AUTO
)

var TickerMap = map[string]*ticker{}

type ticker struct {
	stopCh     chan struct{}
	tickerTime int
}

func NewTicker(tickerTime int) *ticker {
	return &ticker{
		stopCh:     make(chan struct{}),
		tickerTime: tickerTime,
	}
}

func (tickerUp *ticker) Start(chainId string) {
	TickerMap[chainId] = tickerUp
	go func() {
		ticker := time.NewTicker(time.Hour * time.Duration(tickerUp.tickerTime))
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				//定时上报链信息
				err := sync.ReportChainData(chainId)
				log.Error(err)

			case <-tickerUp.stopCh:
				//停止定时任务，并杀死进程
				return
			}
		}
	}()
}

func (tickerUp *ticker) StopTicker(chainId string) {
	delete(TickerMap, chainId)
	close(tickerUp.stopCh)
}
