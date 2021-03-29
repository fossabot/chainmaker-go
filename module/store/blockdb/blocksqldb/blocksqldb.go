/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blocksqldb

import (
	logImpl "chainmaker.org/chainmaker-go/logger"
	commonPb "chainmaker.org/chainmaker-go/pb/protogo/common"
	storePb "chainmaker.org/chainmaker-go/pb/protogo/store"
	"chainmaker.org/chainmaker-go/protocol"
	"chainmaker.org/chainmaker-go/store/serialization"
	"chainmaker.org/chainmaker-go/utils"
	"encoding/hex"
	"golang.org/x/sync/semaphore"
	"runtime"
)

// BlockSqlDB provider a implementation of `blockdb.BlockDB`
// This implementation provides a mysql based data model
type BlockSqlDB struct {
	db               protocol.SqlDBHandle
	workersSemaphore *semaphore.Weighted
	Logger           protocol.Logger
}

// NewBlockSqlDB constructs a new `BlockSqlDB` given an chainId and engine type
func NewBlockSqlDB(chainId string, db protocol.SqlDBHandle, logger protocol.Logger) (*BlockSqlDB, error) {
	nWorkers := runtime.NumCPU()
	if logger == nil {
		logger = logImpl.GetLoggerByChain(logImpl.MODULE_STORAGE, chainId)
	}
	//if err := db.AutoMigrate(&BlockInfo{}); err != nil {
	//	panic(fmt.Sprintf("failed to migrate blockinfo:%s", err))
	//}
	//if err := db.AutoMigrate(&TxInfo{}); err != nil {
	//	panic(fmt.Sprintf("failed to migrate txinfo:%s", err))
	//}
	blockDB := &BlockSqlDB{
		db:               db,
		workersSemaphore: semaphore.NewWeighted(int64(nWorkers)),
		Logger:           logger,
	}
	return blockDB, nil
}
func (b *BlockSqlDB) SaveBlockHeader(header *commonPb.BlockHeader) error {
	blockInfo := ConvertHeader2BlockInfo(header)
	_, err := b.db.Save(blockInfo)
	return err
}

// CommitBlock commits the block and the corresponding rwsets in an atomic operation
func (b *BlockSqlDB) CommitBlock(blocksInfo *serialization.BlockWithSerializedInfo) error {
	block := blocksInfo.Block
	blockHashStr := hex.EncodeToString(block.Header.BlockHash)
	startCommitTxs := utils.CurrentTimeMillisSeconds()
	//save txs
	txInfos := make([]*TxInfo, 0, len(block.Txs))
	for index, tx := range block.Txs {
		txinfo, err := NewTxInfo(tx, block.Header.BlockHeight, int32(index))
		if err != nil {
			b.Logger.Errorf("failed to init txinfo, err:%s", err)
			return err
		}
		txInfos = append(txInfos, txinfo)
	}
	tx := b.db.BeginDbTransaction(blockHashStr)

	for _, txInfo := range txInfos {
		//res := b.db.Clauses(clause.OnConflict{DoNothing: true}).Create(txInfo)
		_, err := tx.Save(txInfo)
		if err != nil {
			b.Logger.Errorf("faield to commit txinfo info, height:%d, tx:%s,err:%s",
				block.Header.BlockHeight, txInfo.TxId, err)
			b.db.RollbackDbTransaction(blockHashStr) //rollback tx
			return err
		}
	}

	elapsedCommitTxs := utils.CurrentTimeMillisSeconds() - startCommitTxs

	//save block info
	startCommitBlockInfo := utils.CurrentTimeMillisSeconds()
	blockInfo, err := NewBlockInfo(block)
	if err != nil {
		b.Logger.Errorf("failed to init blockinfo, err:%s", err)
		return err
	}
	_, err = tx.Save(blockInfo)
	if err != nil {
		b.Logger.Errorf("faield to commit block info, height:%d, err:%s",
			block.Header.BlockHeight, err)
		b.db.RollbackDbTransaction(blockHashStr) //rollback tx
		return err
	}
	err = b.db.CommitDbTransaction(blockHashStr)
	if err != nil {
		b.Logger.Errorf("failed to commit tx, err:%s", err)
		return err
	}
	elapsedCommitBlockInfos := utils.CurrentTimeMillisSeconds() - startCommitBlockInfo
	b.Logger.Infof("chain[%s]: commit block[%d] time used (commit_txs:%d commit_block:%d, total:%d)",
		block.Header.ChainId, block.Header.BlockHeight,
		elapsedCommitTxs, elapsedCommitBlockInfos,
		utils.CurrentTimeMillisSeconds()-startCommitTxs)
	return nil
}

// HasBlock returns true if the block hash exist, or returns false if none exists.
func (b *BlockSqlDB) BlockExists(blockHash []byte) (bool, error) {
	var count int64
	sql := "select count(*) from block_infos where block_hash = ?"
	res, err := b.db.QuerySql(sql, blockHash)
	if err != nil {
		return false, err
	}
	res.ScanColumns(&count)
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// GetBlock returns a block given it's hash, or returns nil if none exists.
func (b *BlockSqlDB) GetBlockByHash(blockHash []byte) (*commonPb.Block, error) {

	return b.getFullBlockBySql("select * from block_infos where block_hash = ?", blockHash)
}
func (b *BlockSqlDB) getBlockInfoBySql(sql string, values ...interface{}) (*BlockInfo, error) {
	//get block info from mysql
	var blockInfo BlockInfo
	res, err := b.db.QuerySql(sql, values...)
	if err != nil {
		return nil, err
	}

	err = res.ScanObject(&blockInfo)
	if err != nil {
		return nil, err
	}
	return &blockInfo, nil
}
func (b *BlockSqlDB) getFullBlockBySql(sql string, values ...interface{}) (*commonPb.Block, error) {
	blockInfo, err := b.getBlockInfoBySql(sql, values...)
	if err != nil {
		return nil, err
	}
	block, err := blockInfo.GetBlock()
	if err != nil {
		return nil, err
	}
	txs, err := b.getTxsByBlockHeight(blockInfo.BlockHeight)
	block.Txs = txs
	return block, nil
}

// GetBlockAt returns a block given it's block height, or returns nil if none exists.
func (b *BlockSqlDB) GetBlock(height int64) (*commonPb.Block, error) {
	return b.getFullBlockBySql("select * from block_infos where block_height =?", height)
}

// GetLastBlock returns the last block.
func (b *BlockSqlDB) GetLastBlock() (*commonPb.Block, error) {
	return b.getFullBlockBySql("select * from block_infos where block_height = (select max(block_height) from block_infos)")
}

// GetLastConfigBlock returns the last config block.
func (b *BlockSqlDB) GetLastConfigBlock() (*commonPb.Block, error) {
	lastBlock, err := b.GetLastBlock()
	if err != nil {
		return nil, err
	}
	if utils.IsConfBlock(lastBlock) {
		return lastBlock, nil
	}
	return b.GetBlock(lastBlock.Header.PreConfHeight)
}

// GetFilteredBlock returns a filtered block given it's block height, or return nil if none exists.
func (b *BlockSqlDB) GetFilteredBlock(height int64) (*storePb.SerializedBlock, error) {
	blockInfo, err := b.getBlockInfoBySql("select * from block_infos where block_height = ?", height)
	if err != nil {
		return nil, err
	}
	return blockInfo.GetFilterdBlock()
}

// GetLastSavepoint reurns the last block height
func (b *BlockSqlDB) GetLastSavepoint() (uint64, error) {
	sql := "select max(block_height) from block_infos"
	row, err := b.db.QuerySql(sql)
	if err != nil {
		return 0, err
	}
	var height uint64
	err = row.ScanColumns(&height)
	if err != nil {
		return 0, err
	}
	return height, nil
}

// GetBlockByTx returns a block which contains a tx.
func (b *BlockSqlDB) GetBlockByTx(txId string) (*commonPb.Block, error) {
	sql := "select * from block_infos where block_height=(select block_height from tx_infos where tx_id=?)"
	return b.getFullBlockBySql(sql, txId)
}

// GetTx retrieves a transaction by txid, or returns nil if none exists.
func (b *BlockSqlDB) GetTx(txId string) (*commonPb.Transaction, error) {
	var txInfo TxInfo
	res, err := b.db.QuerySql("select * from tx_infos where tx_id = ?", txId)
	if err != nil {
		return nil, err
	}
	res.ScanObject(&txInfo)
	return txInfo.GetTx()
}

// HasTx returns true if the tx exist, or returns false if none exists.
func (b *BlockSqlDB) TxExists(txId string) (bool, error) {
	var count int64
	sql := "select count(*) from tx_infos where tx_id = ?"
	res, err := b.db.QuerySql(sql, txId)
	if err != nil {
		return false, err
	}
	res.ScanColumns(&count)
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

//获得某个区块高度下的所有交易
func (b *BlockSqlDB) getTxsByBlockHeight(blockHeight int64) ([]*commonPb.Transaction, error) {
	res, err := b.db.QueryTableSql("select * from tx_infos where block_height = ? order by offsset", blockHeight)
	if err != nil {
		return nil, err
	}
	result := []*commonPb.Transaction{}
	for res.Next() {
		var txInfo TxInfo
		res.ScanObject(&txInfo)
		if err != nil {
			return nil, err
		}
		tx, err := txInfo.GetTx()
		if err != nil {
			return nil, err
		}
		result = append(result, tx)
	}
	return result, nil
}
func (b *BlockSqlDB) GetTxConfirmedTime(txId string) (int64, error) {
	panic("implement me")
}

// Close is used to close database
func (b *BlockSqlDB) Close() {
	b.db.Close()
}

//
//func (b *BlockSqlDB) getBlockByInfo(blockInfo *BlockInfo) (*commonPb.Block, error) {
//	//get txinfos form mysql
//	var txInfos []TxInfo
//	//res = b.db.Debug().Find(&txInfos, txList)
//	res := b.db.Where("block_height = ?",
//		blockInfo.BlockHeight).Order("offset asc").Find(&txInfos)
//	if res.Error == gorm.ErrRecordNotFound {
//		return nil, nil
//	} else if res.Error != nil {
//		b.Logger.Errorf("failed to get tx from tx_info, height:%s, err:%s", blockInfo.BlockHeight, res.Error)
//		return nil, res.Error
//	}
//
//	block, err := blockInfo.GetBlock()
//	if err != nil {
//		b.Logger.Errorf("failed to transform blockinfo to block, chain:%s, block:%d, err:%s",
//			blockInfo.ChainId, blockInfo.BlockHeight, err)
//		return nil, err
//	}
//	for _, txInfo := range txInfos {
//		tx, err := txInfo.GetTx()
//		if err != nil {
//			b.Logger.Errorf("failed to transform txinfo to tx, chain:%s, txid:%s, err:%s",
//				block.Header.ChainId, txInfo.TxId, err)
//		}
//		block.Txs = append(block.Txs, tx)
//	}
//	return block, nil
//}
