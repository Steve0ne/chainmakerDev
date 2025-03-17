/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"fmt"
	"time"

	"chainmaker.org/chainmaker-cross/utils"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/gogo/protobuf/proto"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	pkgRetry "github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	pkgContext "github.com/hyperledger/fabric-sdk-go/pkg/context"
)

const (
	SUCCESS = 0
	ERROR   = 1
)

type Response struct {
	Code   int
	Result string
}

func SendProposal(client *channel.Client, req *TxRequest) (*channel.Response, error) {
	// get context
	refValue := utils.GetPtrUnExportField(client, "context").Interface()
	ctx, ok := refValue.(*pkgContext.Channel)
	if !ok {
		return nil, fmt.Errorf("type convert error")
	}
	// get target
	targets, err := getTargetFromConfig(ctx.Client, ctx.ChannelID())
	if err != nil {
		return nil, fmt.Errorf("calculate targets from config error")
	}

	// get timer ctx
	reqCtx, cancel := pkgContext.NewRequest(ctx, pkgContext.WithTimeoutType(fab.PeerResponse))
	defer cancel()

	proposal := fab.ProcessProposalRequest{
		SignedProposal: &pb.SignedProposal{
			ProposalBytes: req.Payload,
			Signature:     req.Signature,
		},
	}

	// send proposal
	responses := make([]*fab.TransactionProposalResponse, 0)
	for _, target := range targets {
		proposalResp, err := target.ProcessTransactionProposal(reqCtx, proposal)
		if err != nil {
			return nil, err
		}
		responses = append(responses, proposalResp)
	}
	// construct response
	txProposal := fab.TransactionProposal{Proposal: &pb.Proposal{}}
	err = proto.Unmarshal(req.Payload, txProposal)
	if err != nil {
		return nil, err
	}

	// Note: txProposal unmarshal will loss txID
	if txProposal.TxnID == "" {
		txProposal.TxnID = fab.TransactionID(req.Header.TxId)
	}

	resp := &channel.Response{
		Proposal:         &txProposal,
		Responses:        responses,
		TransactionID:    fab.TransactionID(req.Header.TxId),
		TxValidationCode: pb.TxValidationCode_VALID,
		ChaincodeStatus:  responses[0].ChaincodeStatus,
		Payload:          responses[0].Payload,
	}
	return resp, nil
}

func CommitResp(client *channel.Client, request *channel.Request, resp *channel.Response) error {
	// get context
	refValue := utils.GetPtrUnExportField(client, "context").Interface()
	ctx, ok := refValue.(*pkgContext.Channel)
	if !ok {
		return fmt.Errorf("type convert error")
	}

	// create commit handler
	commitHandler := invoke.NewCommitHandler(nil)

	// get target
	targets, err := getTargetFromConfig(ctx.Client, ctx.ChannelID())
	if err != nil {
		return fmt.Errorf("calculate targets from config error")
	}
	peers := make([]fab.Peer, len(targets))
	for _, target := range targets {
		peer, ok := target.(fab.Peer)
		if !ok {
			return fmt.Errorf("convert peer error")
		}
		peers = append(peers, peer)
	}
	// get timer ctx
	reqCtx, cancel := pkgContext.NewRequest(ctx, pkgContext.WithTimeoutType(fab.PeerResponse))
	defer cancel()

	if request.ChaincodeID == "" || request.Fcn == "" {
		return fmt.Errorf("ChaincodeID and Fcn are required")
	}
	transactor, err := ctx.ChannelService().Transactor(reqCtx)
	if err != nil {
		return fmt.Errorf("failed to create transactor, err: %s", err)
	}
	selection, err := ctx.ChannelService().Selection()
	if err != nil {
		return fmt.Errorf("failed to create selection service, err: %s", err)
	}
	discovery, err := ctx.ChannelService().Discovery()
	if err != nil {
		return fmt.Errorf("failed to create discovery service, err: %s", err)
	}
	peerFilter := func(peer fab.Peer) bool {
		return true
	}

	// get membership
	refValue = utils.GetPtrUnExportField(client, "membership").Interface()
	membership, ok := refValue.(fab.ChannelMembership)
	if !ok {
		return fmt.Errorf("get context membership error, type convert error")
	}

	// get membership
	refValue = utils.GetPtrUnExportField(client, "eventService").Interface()
	eventService, ok := refValue.(fab.EventService)
	if !ok {
		return fmt.Errorf("get context membership error, type convert error")
	}

	// set timeout
	timeOut := make(map[fab.TimeoutType]time.Duration)
	timeOut[fab.Execute] = time.Minute * 5

	// construct options
	opts := invoke.Opts{
		Targets:       peers,
		TargetFilter:  nil,
		TargetSorter:  nil,
		Retry:         pkgRetry.Opts{}, // TODO retry method
		BeforeRetry:   nil,
		Timeouts:      timeOut,
		ParentContext: nil,
		CCFilter:      nil,
	}

	clientContext := &invoke.ClientContext{
		Selection:    selection,
		Discovery:    discovery,
		Membership:   membership,
		Transactor:   transactor,
		EventService: eventService,
	}

	requestContext := &invoke.RequestContext{
		Request:         invoke.Request(*request),
		Opts:            opts,
		Response:        invoke.Response(*resp),
		RetryHandler:    nil,
		Ctx:             reqCtx,
		SelectionFilter: peerFilter,
		PeerSorter:      nil,
	}

	commitHandler.Handle(requestContext, clientContext)
	if requestContext.Error != nil {
		return requestContext.Error
	}
	return nil
}

func getTargetFromConfig(ctx context.Client, channelID string) ([]fab.ProposalProcessor, error) {
	var targets []fab.ProposalProcessor
	chPeers := ctx.EndpointConfig().ChannelPeers(channelID)
	if len(chPeers) == 0 {
		return nil, fmt.Errorf("no channel peers configured for channel [%s]", channelID)
	}

	for _, p := range chPeers {
		newPeer, err := ctx.InfraProvider().CreatePeerFromConfig((&p.NetworkPeer))
		if err != nil || newPeer == nil {
			return nil, fmt.Errorf("new peer failed")
		}

		// Pick peers in the same MSP as the context since only they can query system chaincode
		//if newPeer.MSPID() == ctx.Identifier().MSPID {
		targets = append(targets, newPeer)
		//}
	}
	//target := randomMaxTargets(targets, 1)[0]
	//targets = randomMaxTargets(targets, c.opts.MaxTargets)
	return targets, nil
}

func RetryQueryTx(lc *ledger.Client, txId fab.TransactionID, peersUrls []string) (uint64, error) {
	// send request and handle response
	reqPeers := ledger.WithTargetEndpoints(peersUrls...)
	err := retry.Retry(func(uint) error {
		// send request and handle response
		resp, err := lc.QueryTransaction(txId, reqPeers)
		if err != nil {
			return err
		}
		if resp.ValidationCode == 0 {
			return nil
		} else {
			return fmt.Errorf("try again")
		}
	},
		strategy.Limit(RetryCount),
		strategy.Wait(RetryTimePeriod), // 指定超时等待
	)
	if err != nil {
		return 0, err
	}
	block, err := lc.QueryBlockByTxID(txId, reqPeers)
	if err != nil {
		return 0, err
	}
	return block.Header.Number, nil
}
