/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blockpool

import (
	"bytes"
	"errors"
	"fmt"

	"chainmaker.org/chainmaker/pb-go/common"
	chainedbftpb "chainmaker.org/chainmaker/pb-go/consensus/chainedbft"
)

//BlockPool store block and qc in memory
type BlockPool struct {
	//mtx                   sync.RWMutex
	idToQC                map[string]*chainedbftpb.QuorumCert // store qc in memory, key is BlockId, value is blockQC
	blockTree             *BlockTree                          // store block in memory
	highestQC             *chainedbftpb.QuorumCert            // highest qc in local node
	highestCertifiedBlock *common.Block                       // highest block with qc in local node
}

//NewBlockPool init a block pool with rootBlock, rootQC and maxPrunedSize
func NewBlockPool(rootBlock *common.Block,
	rootQC *chainedbftpb.QuorumCert, maxPrunedSize int) *BlockPool {
	blockPool := &BlockPool{
		idToQC:                make(map[string]*chainedbftpb.QuorumCert, 0),
		blockTree:             NewBlockTree(rootBlock, rootQC, maxPrunedSize),
		highestQC:             rootQC,
		highestCertifiedBlock: rootBlock,
	}
	blockPool.idToQC[string(rootQC.BlockId)] = rootQC
	return blockPool
}

//InsertBlock insert block to block pool
func (bp *BlockPool) InsertBlock(block *common.Block) error {
	if err := bp.blockTree.InsertBlock(block); err != nil {
		return err
	}
	if _, exist := bp.idToQC[string(block.Header.BlockHash)]; exist {
		if bp.highestCertifiedBlock.Header.BlockHeight < block.Header.BlockHeight {
			bp.highestCertifiedBlock = block
		}
	}
	return nil
}

//InsertQC store qc
func (bp *BlockPool) InsertQC(qc *chainedbftpb.QuorumCert) error {
	if qc == nil {
		return errors.New("qc is nil")
	}
	if _, exist := bp.idToQC[string(qc.BlockId)]; exist {
		return nil
	}
	bp.idToQC[string(qc.BlockId)] = qc

	if qc.Level <= bp.highestQC.Level {
		return nil
	}
	bp.highestQC = qc
	if blk := bp.blockTree.GetBlockByID(string(qc.BlockId)); blk != nil {
		bp.highestCertifiedBlock = blk
	}
	return nil
}

func (bp *BlockPool) GetBlocks(height uint64) []*common.Block {
	return bp.blockTree.GetBlocks(height)
}

//GetRootBlock get root block
func (bp *BlockPool) GetRootBlock() *common.Block {
	return bp.blockTree.GetRootBlock()
}

func (bp *BlockPool) GetRootQC() *chainedbftpb.QuorumCert {
	return bp.blockTree.GetRootQC()
}

//GetBlockByID get block by block hash
func (bp *BlockPool) GetBlockByID(id string) *common.Block {
	return bp.blockTree.GetBlockByID(id)
}

//GetQCByID get qc by block hash
func (bp *BlockPool) GetQCByID(id string) *chainedbftpb.QuorumCert {
	return bp.idToQC[id]
}

//GetHighestQC get highest qc
func (bp *BlockPool) GetHighestQC() *chainedbftpb.QuorumCert {
	return bp.highestQC
}

//GetHighestCertifiedBlock get highest certified block
func (bp *BlockPool) GetHighestCertifiedBlock() *common.Block {
	return bp.highestCertifiedBlock
}

//BranchFromRoot get branch from root to input block
func (bp *BlockPool) BranchFromRoot(block *common.Block) []*common.Block {
	return bp.blockTree.BranchFromRoot(block)
}

//PruneBlock prune block
func (bp *BlockPool) PruneBlock(newRootID string) error {
	newRootQC := bp.idToQC[newRootID]
	prunedBlocks, err := bp.blockTree.PruneBlock(newRootID, newRootQC)
	if err != nil || prunedBlocks == nil {
		return err
	}
	for _, block := range prunedBlocks {
		delete(bp.idToQC, block)
	}
	return nil
}

func (bp *BlockPool) Details() string {
	qcCount := len(bp.idToQC)
	qcContents := bytes.NewBufferString(fmt.Sprintf("BlockPool qcCount: %d\n", qcCount))
	for blkID, qc := range bp.idToQC {
		qcContents.WriteString(fmt.Sprintf("blkID: %s, height: %d, level: %d\n", blkID, qc.Height, qc.Level))
	}
	qcContents.WriteString(bp.blockTree.Details())
	return qcContents.String()
}
