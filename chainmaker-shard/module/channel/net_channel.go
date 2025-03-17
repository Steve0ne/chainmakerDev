/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"time"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"errors"
	"fmt"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/event/coder"
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/net"
	"chainmaker.org/chainmaker-cross/net/net_http"
	"chainmaker.org/chainmaker-cross/net/net_libp2p"
	"chainmaker.org/chainmaker-cross/utils"
	"go.uber.org/zap"
)

const (
	MinDataLength              = 2
	HttpCrossTransactionRouter = "/cross?method=transaction"
)

// NetChannel net channel which handle message between two cross-chain proxy
type NetChannel struct {
	connection net.Connection               // 与其他跨链代理节点的连接
	log        *zap.SugaredLogger           // 日志
	coders     *coder.EventCoderTools       // 消息编解码器
	contexts   *event.ProofResponseContexts // 消息证明的Map
}

// NewNetChannel create new net channel
func NewNetChannel(connection net.Connection) *NetChannel {
	return &NetChannel{
		connection: connection,
		log:        logger.GetLogger(logger.ModuleNet),
		coders:     coder.GetEventCoderTools(),
		contexts:   event.GetProofResponseContexts(),
	}
}

// Init init channel connection
func (n *NetChannel) Init() error {
	dataChan, err := n.connection.ReadData()
	if err != nil {
		return err
	}
	// 启动监听，用于读取数据并打印，不做启动事情
	go func() {
		for {
			select {
			case msg, ok := <-dataChan:
				if !ok {
					// 通道关闭
					n.log.Warn("net channel will be closed")
					return
				}
				n.handleReceivedData(msg)
			case <-time.After(conf.LogWritePeriod):
				// 打印日志，表明在正常活着
				n.log.Info("net channel is running periodically!")
			}
		}
	}()
	return nil
}

func (n *NetChannel) handleReceivedData(msg net.Message) {
	// 从通道中读到数据
	if len(msg.GetPayload()) < MinDataLength {
		// 打印错误信息
		n.log.Error("receive data is illegal")
		return
	}
	n.log.Debugf("receive data length = %v", len(msg.GetPayload()))
	receivedData, err := utils.Base64DecodeToBytes(string(msg.GetPayload()))
	if err != nil {
		n.log.Error("base64 decode data failed, ", err)
		return
	}
	eventTy, marshalTy := eventproto.EventType(receivedData[coder.EventTyIndex]), event.MarshalType(receivedData[coder.MarshalTyIndex])
	if eventTy == eventproto.ProofRespEventType {
		if eveCoder, exist := n.coders.GetDefaultCoder(eventTy); exist {
			// 处理
			if marshalTy == event.BinaryMarshalType {
				eve, err := eveCoder.UnmarshalFromBinary(receivedData)
				if err != nil {
					// 打印错误信息
					n.log.Error("unmarshal receive data failed, ", err)
				} else {
					if resp, ok := eve.(*event.ProofResponse); ok {
						// 填充结果
						if resp.Code == event.SuccessResp {
							n.log.Infof("cross[%s]->chain[%s]->key[%s] response is success",
								resp.GetCrossID(), resp.GetChainID(), resp.Key)
							// 操作成功，填充结果
							n.contexts.DoneByProofResp(resp)
						} else {
							n.log.Errorf("cross[%s]->chain[%s]->key[%s] response is failed",
								resp.GetCrossID(), resp.GetChainID(), resp.Key)
							n.contexts.DoneError(resp.Key, resp.Msg)
						}
					}
				}
			}
		}
		return
	}
	n.log.Errorf("the event type = [%v] which id not ProofRespEvent", eventTy)
}

// GetChanType return type of channel
func (n *NetChannel) GetChanType() TransmissionChanType {
	return NetTransmissionChan
}

// Deliver marshal transaction event and transaction it to other cross-chain proxy
func (n *NetChannel) Deliver(eve *event.TransactionEventContext) error {
	var (
		msg net.Message
		err error
	)
	switch n.connection.GetProvider() {
	case net.LibP2PConnection:
		// 需要序列化eve
		binary, err := coder.GetTransactionEventCtxCoder().MarshalToBinary(eve)
		if err != nil {
			n.log.Errorf("cross[%s]->chain[%s]->key[%s] marshal to binary bytes failed, ",
				eve.GetEvent().GetCrossID(), eve.GetEvent().GetChainID(), eve.GetKey(), err)
			return err
		}
		// 将binary转换为Base64
		base64String := utils.Base64EncodeToString(binary)
		n.log.Infof("cross[%s]->chain[%s]->key[%s] begin write to net channel, length = [%v]",
			eve.GetEvent().GetCrossID(), eve.GetEvent().GetChainID(), eve.GetKey(), len(base64String))
		if msg, err = net_libp2p.NewLibP2pMessage(n.connection.PeerID(), []byte(base64String), false); err != nil {
			n.log.Errorf("cross[%s]->chain[%s]->key[%s] generate LibP2p_message error, ",
				eve.GetEvent().GetCrossID(), eve.GetEvent().GetChainID(), eve.GetKey(), err)
		}
	case net.HttpConnection:
		//router, ok := conf.Config.RouterConfigs.RouterConfigs
		if msg, err = net_http.NewRequest(eve, HttpCrossTransactionRouter, event.BinaryMarshalType); err != nil {
			n.log.Errorf("cross[%s]->chain[%s]->key[%s] generate http_message error, ",
				eve.GetEvent().GetCrossID(), eve.GetEvent().GetChainID(), eve.GetKey(), err)
			return err
		}
	default:
		err = errors.New(fmt.Sprintf("unsupported connection provider #{n.connection.Provider()}"))
		n.log.Errorf("cross[%s]->chain[%s]->key[%s] write to channel failed, ",
			eve.GetEvent().GetCrossID(), eve.GetEvent().GetChainID(), eve.GetKey(), err)
		return err
	}

	err = n.connection.WriteData(msg)
	if err != nil {
		n.log.Errorf("cross[%s]->chain[%s]->key[%s] write to net channel failed, ",
			eve.GetEvent().GetCrossID(), eve.GetEvent().GetChainID(), eve.GetKey(), err)
		return err
	} else {
		n.log.Infof("cross[%s]->chain[%s]->key[%s] write to net channel success",
			eve.GetEvent().GetCrossID(), eve.GetEvent().GetChainID(), eve.GetKey())
	}
	return nil
}
