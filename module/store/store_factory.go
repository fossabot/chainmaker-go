// +build !rocksdb

/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package store

import (
	"chainmaker.org/chainmaker-go/localconf"
	logImpl "chainmaker.org/chainmaker-go/logger"
	"chainmaker.org/chainmaker-go/protocol"
	"chainmaker.org/chainmaker-go/store/blockdb"
	"chainmaker.org/chainmaker-go/store/blockdb/blockkvdb"
	"chainmaker.org/chainmaker-go/store/blockdb/blocksqldb"
	"chainmaker.org/chainmaker-go/store/cache"
	"chainmaker.org/chainmaker-go/store/dbprovider"
	"chainmaker.org/chainmaker-go/store/dbprovider/leveldbprovider"
	"chainmaker.org/chainmaker-go/store/dbprovider/sqldbprovider"
	"chainmaker.org/chainmaker-go/store/historydb"
	"chainmaker.org/chainmaker-go/store/historydb/historykvdb"
	"chainmaker.org/chainmaker-go/store/historydb/historysqldb"
	"chainmaker.org/chainmaker-go/store/statedb"
	"chainmaker.org/chainmaker-go/store/statedb/statekvdb"
	"chainmaker.org/chainmaker-go/store/statedb/statesqldb"
	"chainmaker.org/chainmaker-go/store/types"
	"errors"
	"golang.org/x/sync/semaphore"
	"runtime"
)

// Factory is a factory function to create an instance of the block store
// which commits block into the ledger.
type Factory struct {
}

// NewStore constructs new BlockStore
func (m *Factory) NewStore(engineType types.EngineType, chainId string, logger protocol.Logger) (protocol.BlockchainStore, error) {
	if logger == nil {
		logger = logImpl.GetLoggerByChain(logImpl.MODULE_STORAGE, chainId)
	}
	switch engineType {
	case types.LevelDb, types.RocksDb:
		blockDB, err := m.NewBlockKvDB(chainId, engineType, logger)
		if err != nil {
			return nil, err
		}
		stateDB, err := m.NewStateKvDB(chainId, engineType, logger)
		if err != nil {
			return nil, err
		}
		historyDB, err := m.NewHistoryKvDB(chainId, engineType, logger)
		if err != nil {
			return nil, err
		}
		return NewBlockStoreImpl(chainId, blockDB, stateDB, historyDB, NewKvDBProvider(chainId, types.CommonDBDir, engineType, logger), logger)
	case types.MySQL, types.Sqlite:
		dbprovider := sqldbprovider.NewSqlDBProvider(chainId, localconf.ChainMakerConfig)
		blockDB, err := blocksqldb.NewBlockSqlDB(chainId, dbprovider, logger)
		if err != nil {
			return nil, err
		}
		stateDB, err := statesqldb.NewStateSqlDB(chainId, dbprovider, logger)
		if err != nil {
			return nil, err
		}
		historyDB, err := historysqldb.NewHistorySqlDB(chainId, dbprovider, logger)
		if err != nil {
			return nil, err
		}

		return NewBlockStoreImpl(chainId, blockDB, stateDB, historyDB, dbprovider, logger)
	default:
		return nil, errors.New("invalid engine type")
	}
	return nil, errors.New("invalid engine type")
}

// NewBlockKvDB constructs new `BlockDB`
func (m *Factory) NewBlockKvDB(chainId string, engineType types.EngineType, logger protocol.Logger) (blockdb.BlockDB, error) {
	nWorkers := runtime.NumCPU()
	if logger == nil {
		logger = logImpl.GetLoggerByChain(logImpl.MODULE_STORAGE, chainId)
	}
	blockDB := &blockkvdb.BlockKvDB{
		WorkersSemaphore: semaphore.NewWeighted(int64(nWorkers)),
		Cache:            cache.NewStoreCacheMgr(chainId, logger),

		Logger: logger,
	}
	switch engineType {
	case types.LevelDb:
		blockDB.DbProvider = leveldbprovider.NewBlockProvider(chainId, logger)
	default:
		return nil, nil
	}
	return blockDB, nil
}

// NewStateKvDB constructs new `StabeKvDB`
func (m *Factory) NewStateKvDB(chainId string, engineType types.EngineType, logger protocol.Logger) (statedb.StateDB, error) {
	if logger == nil {
		logger = logImpl.GetLoggerByChain(logImpl.MODULE_STORAGE, chainId)
	}
	stateDB := &statekvdb.StateKvDB{
		Logger: logger,
		Cache:  cache.NewStoreCacheMgr(chainId, logger),
	}
	switch engineType {
	case types.LevelDb:
		stateDB.DbProvider = leveldbprovider.NewStateProvider(chainId, logger)
	default:
		return nil, nil
	}
	return stateDB, nil
}

// NewHistoryKvDB constructs new `HistoryKvDB`
func (m *Factory) NewHistoryKvDB(chainId string, engineType types.EngineType, logger protocol.Logger) (historydb.HistoryDB, error) {
	if logger == nil {
		logger = logImpl.GetLoggerByChain(logImpl.MODULE_STORAGE, chainId)
	}
	historyDB := &historykvdb.HistoryKvDB{
		Cache:  cache.NewStoreCacheMgr(chainId, logger),
		Logger: logger,
	}
	switch engineType {
	case types.LevelDb:
		historyDB.DbProvider = leveldbprovider.NewHistoryProvider(chainId, logger)
	default:
		return nil, nil
	}
	return historyDB, nil
}

// NewKvDBProvider constructs new kv database
func NewKvDBProvider(chainId string, dbDir string, engineType types.EngineType, logger protocol.Logger) dbprovider.Provider {
	switch engineType {
	case types.LevelDb:
		return leveldbprovider.NewLevelDBProvider(chainId, dbDir, logger)
	}
	return nil
}
