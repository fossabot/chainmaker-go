/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"chainmaker.org/chainmaker-go/utils"

	"chainmaker.org/chainmaker/pb-go/syscontract"

	commonPb "chainmaker.org/chainmaker/pb-go/common"
	"github.com/spf13/cobra"

	"github.com/gogo/protobuf/proto"
)

func SubscribeTxCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribeTx",
		Short: "Subscribe Tx",
		Long:  "Subscribe Tx",
		RunE: func(_ *cobra.Command, _ []string) error {
			return subscribeTx()
		},
	}
	return cmd
}

func subscribeTx() error {
	start, _ := utils.Int64ToBytes(startBlock)
	end, _ := utils.Int64ToBytes(endBlock)
	payload := &commonPb.Payload{
		Parameters: []*commonPb.KeyValuePair{
			{Key: syscontract.SubscribeTx_START_BLOCK.String(), Value: start},
			{Key: syscontract.SubscribeTx_END_BLOCK.String(), Value: end},
			{Key: syscontract.SubscribeTx_TX_IDS.String(), Value: []byte(txIds)},
		},
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf("marshal payload failed, %s", err.Error())
	}

	_, err = subscribeRequest(sk3, client, syscontract.SubscribeFunction_SUBSCRIBE_TX.String(), chainId, payloadBytes)
	if err != nil {
		return err
	}

	return nil
}
