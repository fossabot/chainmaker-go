/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package resultdb

import (
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	"chainmaker.org/chainmaker-go/store/serialization"
)

// ResultDB provides handle to rwSets instances
type ResultDB interface {
	InitGenesis(genesisBlock *serialization.BlockWithSerializedInfo) error
	// CommitBlock commits the block rwsets in an atomic operation
	CommitBlock(blockInfo *serialization.BlockWithSerializedInfo) error

	// GetTxRWSet returns an txRWSet for given txId, or returns nil if none exists.
	GetTxRWSet(txid string) (*commonPb.TxRWSet, error)

	// GetLastSavepoint returns the last block height
	GetLastSavepoint() (uint64, error)

	// Close is used to close database
	Close()
}