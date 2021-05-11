// +build rocksdb

/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package store

import (
	"runtime"

	logImpl "chainmaker.org/chainmaker-go/logger"
	"chainmaker.org/chainmaker-go/protocol"
	"chainmaker.org/chainmaker-go/store/blockdb"
	"chainmaker.org/chainmaker-go/store/blockdb/blockkvdb"
	"chainmaker.org/chainmaker-go/store/cache"
	"chainmaker.org/chainmaker-go/store/contracteventdb/eventmysqldb"
	"chainmaker.org/chainmaker-go/store/dbprovider"
	"chainmaker.org/chainmaker-go/store/dbprovider/rocksdbprovider"
	"chainmaker.org/chainmaker-go/store/historydb"
	"chainmaker.org/chainmaker-go/store/historydb/historykvdb"
	"chainmaker.org/chainmaker-go/store/statedb"
	"chainmaker.org/chainmaker-go/store/statedb/statekvdb"
	"chainmaker.org/chainmaker-go/store/types"
	"golang.org/x/sync/semaphore"
)

// Factory is a factory function to create an instance of the block store
// which commits block into the ledger.
type Factory struct {
}

// NewStore constructs new `protocol.BlockchainStore`
func (m *Factory) NewStore(engineType types.EngineType, chainId string) (protocol.BlockchainStore, error) {
	switch engineType {
	case types.RocksDb:
		blockDB, err := m.NewBlockKvDB(chainId, engineType)
		if err != nil {
			return nil, err
		}
		stateDB, err := m.NewStateKvDB(chainId, engineType)
		if err != nil {
			return nil, err
		}
		historyDB, err := m.NewHistoryKvDB(chainId, engineType)
		if err != nil {
			return nil, err
		}
		contractEventDB, err := eventmysqldb.NewContractEventMysqlDB(chainId)
		if err != nil {
			return nil, err
		}
		return NewBlockStoreImpl(chainId, blockDB, stateDB, historyDB, contractEventDB, NewKvDBProvider(chainId, types.CommonDBDir, engineType))
	default:
		return nil, nil
	}
}

// NewBlockKvDB constructs new `BlockDB`
func (m *Factory) NewBlockKvDB(chainId string, engineType types.EngineType) (blockdb.BlockDB, error) {
	nWorkers := runtime.NumCPU()
	blockDB := &blockkvdb.BlockKvDB{
		WorkersSemaphore: semaphore.NewWeighted(int64(nWorkers)),
		Cache:            cache.NewStoreCacheMgr(chainId),

		Logger: logImpl.GetLoggerByChain(logImpl.MODULE_STORAGE, chainId),
	}
	switch engineType {
	case types.RocksDb:
		blockDB.DbProvider = rocksdbprovider.NewBlockProvider(chainId)
	default:
		return nil, nil
	}
	return blockDB, nil
}

// NewStateKvDB constructs new `StabeKvDB`
func (m *Factory) NewStateKvDB(chainId string, engineType types.EngineType) (statedb.StateDB, error) {
	stateDB := &statekvdb.StateKvDB{
		Logger: logImpl.GetLoggerByChain(logImpl.MODULE_STORAGE, chainId),
		Cache:  cache.NewStoreCacheMgr(chainId),
	}
	switch engineType {
	case types.RocksDb:
		stateDB.DbProvider = rocksdbprovider.NewStateProvider(chainId)
	default:
		return nil, nil
	}
	return stateDB, nil
}

// NewHistoryKvDB constructs new `HistoryKvDB`
func (m *Factory) NewHistoryKvDB(chainId string, engineType types.EngineType) (historydb.HistoryDB, error) {
	historyDB := &historykvdb.HistoryKvDB{
		Cache:  cache.NewStoreCacheMgr(chainId),
		Logger: logImpl.GetLoggerByChain(logImpl.MODULE_STORAGE, chainId),
	}
	switch engineType {
	case types.RocksDb:
		historyDB.DbProvider = rocksdbprovider.NewHistoryProvider(chainId)
	default:
		return nil, nil
	}
	return historyDB, nil
}

// NewStateKvDB constructs new `StabeKvDB`
func NewKvDBProvider(chainId string, dbDir string, engineType types.EngineType) dbprovider.Provider {
	switch engineType {
	case types.RocksDb:
		return rocksdbprovider.NewProvider(chainId, dbDir)
	}
	return nil
}
