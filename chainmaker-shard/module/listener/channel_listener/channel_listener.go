/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package channel_listener

import (
	"context"
	"errors"
	"time"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/event/coder"
	"chainmaker.org/chainmaker-cross/handler"
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/net"
	libp2p "chainmaker.org/chainmaker-cross/net/net_libp2p"
	"chainmaker.org/chainmaker-cross/utils"
	"github.com/libp2p/go-libp2p-core/protocol"
	"go.uber.org/zap"
)

const (
	MinDataLength = 3
)

// ChannelListener the struct of channel listener
type ChannelListener struct {
	peer       net.Peer               // 点对点网络连接的本地节点
	log        *zap.SugaredLogger     // log
	coders     *coder.EventCoderTools // 编解码器
	cancelFunc context.CancelFunc     // 退出信号函数
}

// NewChannelListener create new channel listener
func NewChannelListener() *ChannelListener {
	channelConf := conf.Config.ListenerConfig.ChannelConfig
	var peer net.Peer
	if net.PeerProvider(channelConf.Provider) == net.LibP2PPeer {
		libp2pChannelConfig := channelConf.LibP2PChannel
		if libp2pChannelConfig == nil {
			panic("channel listener's config is error")
		}
		monitorHost, err := libp2p.NewMonitorHost(libp2pChannelConfig.Address, conf.FinalCfgPath(libp2pChannelConfig.PrivKeyFile),
			protocol.ID(libp2pChannelConfig.ProtocolID), libp2pChannelConfig.GetDelimit())
		if err != nil {
			panic(err)
		}
		peer = monitorHost
	}
	return &ChannelListener{
		peer:   peer,
		log:    logger.GetLogger(logger.ModuleChannelListener),
		coders: coder.GetEventCoderTools(),
	}
}

// ListenStart start this channel listener
func (cl *ChannelListener) ListenStart() error {
	// set stream handler
	ch, err := cl.peer.Listen()
	if err != nil {
		cl.log.Error(err)
		return err
	}
	eventHandler, exist := handler.GetEventHandlerTools().GetHandler(handler.TransactionProcess)
	if !exist {
		cl.log.Errorf("eventHandler[%d] not exist", handler.TransactionProcess)
		return errors.New("can not find handler to handle transaction event")
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	cl.cancelFunc = cancelFunc
	go func(ctx context.Context, ch chan net.Message) {
		for {
			select {
			case <-ctx.Done():
				// 完成跳出
				if err := cl.peer.Stop(); err != nil {
					// 打印错误
					cl.log.Errorf("channel listener exit error", err)
				}
				return
			case msg := <-ch:
				cl.log.Info("channel listener receive data")
				// 启动独立goroutine处理该问题
				go func(msg net.Message) {
					// 获取数据的序列化方式
					if len(msg.GetPayload()) < MinDataLength {
						// 打印错误信息
						cl.log.Error("receive data length is illegal")
					} else {
						receivedData, _ := utils.Base64DecodeToBytes(string(msg.GetPayload()))
						eventTy, marshalTy := eventproto.EventType(receivedData[coder.EventTyIndex]), event.MarshalType(receivedData[coder.MarshalTyIndex])
						if eventTy == eventproto.TransactionCtxEventType {
							if eveCoder, exist := cl.coders.GetDefaultCoder(eventTy); exist {
								// 处理
								if marshalTy == event.BinaryMarshalType {
									eve, err := eveCoder.UnmarshalFromBinary(receivedData)
									if err != nil {
										// 打印错误信息
										cl.log.Error("unmarshal receive data failed, ", err)
									} else {
										result, err := eventHandler.Handle(eve, true)
										if err != nil {
											// 打印错误信息
											cl.log.Error("handle event failed, ", err)
										} else {
											// 需要结果是*event.ProofResponse
											if resp, ok := result.(*event.ProofResponse); ok {
												if respCoder, exist := cl.coders.GetDefaultCoder(eventproto.ProofRespEventType); exist {
													// 进行二进制的序列化
													binary, err := respCoder.MarshalToBinary(resp)
													if err == nil {
														sendData := utils.Base64EncodeToString(binary)
														if m, err := libp2p.NewLibP2pMessage(msg.GetNodeID(), []byte(sendData), false); err != nil {
															cl.log.Error("generate LibP2pMessage failed, ", err)
														} else {
															if err := cl.peer.Write(m); err != nil {
																// 记录
																cl.log.Error("write proof response event to connection failed, ", err)
															} else {
																// todo 暂不处理数据库记录
															}
														}
													} else {
														// 日志打印
														cl.log.Error("marshal proof response event failed, ", err)
													}
												}
											} else {
												// 打印信息
												cl.log.Error("resp result can not convert to ProofResponse")
											}
										}
									}
								}
							} else {
								// 打印错误信息
								cl.log.Error("can not find coder for transaction event ctx type")
							}
						} else {
							// 打印错误信息
							cl.log.Error("received data is not type of transaction event context")
						}
					}
				}(msg)
			case <-time.After(conf.LogWritePeriod):
				// logger
				cl.log.Info("channel listener is running periodically!")
			}
		}
	}(ctx, ch)
	return nil
}

// Stop stop listener server
func (cl *ChannelListener) Stop() error {
	cl.cancelFunc()
	cl.log.Info("Module channel-listener stopped")
	return nil
}
