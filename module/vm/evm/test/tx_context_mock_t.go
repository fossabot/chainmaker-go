/*
Copyright (C) BABEC. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	configPb "chainmaker.org/chainmaker/pb-go/config"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"

	acPb "chainmaker.org/chainmaker/pb-go/accesscontrol"
	commonPb "chainmaker.org/chainmaker/pb-go/common"
	storePb "chainmaker.org/chainmaker/pb-go/store"
	"chainmaker.org/chainmaker/protocol"
	wasm "chainmaker.org/chainmaker-go/wasmer/wasmer-go"
)

var testOrgId = "wx-org1.chainmaker.org"

var CertFilePath = "./config/admin1.sing.crt"
var ByteCodeFile = "./token.bin"

var txType = commonPb.TxType_INVOKE_USER_CONTRACT

const (
	ContractNameTest    = "contract01"
	ContractVersionTest = "v1.0.0"
	ChainIdTest         = "chain01"
)

var bytes []byte
var file []byte

// 初始化上下文和wasm字节码
func InitContextTest(runtimeType commonPb.RuntimeType) (*commonPb.ContractId, *TxContextMockTest, []byte) {
	if bytes == nil {
		bytes, _ = wasm.ReadBytes(ByteCodeFile)
		fmt.Printf("byteCode file size=%d\n", len(bytes))
	}

	contractId := commonPb.ContractId{
		ContractName:    ContractNameTest,
		ContractVersion: ContractVersionTest,
		RuntimeType:     runtimeType,
	}
	if file == nil {
		var err error
		file, err = ioutil.ReadFile(CertFilePath)
		if err != nil {
			panic("file is nil" + err.Error())
		}
	}
	sender := &acPb.SerializedMember{
		OrgId:      testOrgId,
		MemberInfo: file,
		IsFullCert: true,
	}

	txContext := TxContextMockTest{
		lock:      &sync.Mutex{},
		vmManager: nil,
		hisResult: make([]*callContractResult, 0),
		creator:   sender,
		sender:    sender,
		cacheMap:  make(map[string][]byte),
	}

	versionKey := []byte(protocol.ContractVersion + ContractNameTest)
	runtimeTypeKey := []byte(protocol.ContractRuntimeType + ContractNameTest)
	versionedByteCodeKey := append([]byte(protocol.ContractByteCode+ContractNameTest), []byte(contractId.ContractVersion)...)

	txContext.Put(commonPb.ContractName_SYSTEM_CONTRACT_STATE.String(), versionedByteCodeKey, bytes)
	txContext.Put(commonPb.ContractName_SYSTEM_CONTRACT_STATE.String(), versionKey, []byte(contractId.ContractVersion))
	txContext.Put(commonPb.ContractName_SYSTEM_CONTRACT_STATE.String(), runtimeTypeKey, []byte(strconv.Itoa(int(runtimeType))))

	return &contractId, &txContext, bytes
}

type TxContextMockTest struct {
	lock          *sync.Mutex
	vmManager     protocol.VmManager
	gasUsed       uint64 // only for callContract
	currentDepth  int
	currentResult []byte
	hisResult     []*callContractResult

	sender   *acPb.SerializedMember
	creator  *acPb.SerializedMember
	cacheMap map[string][]byte
}

func (s *TxContextMockTest) SetStateKvHandle(i int32, iterator protocol.StateIterator) {
	panic("implement me")
}

func (s *TxContextMockTest) GetStateKvHandle(i int32) (protocol.StateIterator, bool) {
	panic("implement me")
}

func (s *TxContextMockTest) PutRecord(contractName string, value []byte, sqlType protocol.SqlType) {
	panic("implement me")
}

func (s *TxContextMockTest) Select(name string, startKey []byte, limit []byte) (protocol.StateIterator, error) {
	panic("implement me")
}

func (s *TxContextMockTest) GetBlockProposer() []byte {
	panic("implement me")
}

func (s *TxContextMockTest) SetStateSqlHandle(i int32, rows protocol.SqlRows) {
	panic("implement me")
}

func (s *TxContextMockTest) GetStateSqlHandle(i int32) (protocol.SqlRows, bool) {
	panic("implement me")
}

type callContractResult struct {
	contractName string
	method       string
	param        map[string]string
	deep         int
	gasUsed      uint64
	result       []byte
}

func (s *TxContextMockTest) Get(name string, key []byte) ([]byte, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	k := string(key)
	if name != "" {
		k = name + "::" + k
	}
	//println("【get】 key:" + k)
	//fms.Println("【get】 key:", k, "val:", cacheMap[k])
	return s.cacheMap[k], nil
	//return nil,nil
	//data := "hello"
	//for i := 0; i < 70; i++ {
	//	for i := 0; i < 100; i++ {//1k
	//		data += "1234567890"
	//	}
	//}
	//return []byte(data), nil
}

func (s *TxContextMockTest) Put(name string, key []byte, value []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	k := string(key)
	//v := string(value)
	if name != "" {
		k = name + "::" + k
	}
	//println("【put】 key:" + k)
	//fmt.Println("【put】 key:", k, "val:", value)
	s.cacheMap[k] = value
	return nil
}

func (s *TxContextMockTest) Del(name string, key []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	k := string(key)
	//v := string(value)
	if name != "" {
		k = name + "::" + k
	}
	//println("【put】 key:" + k)
	s.cacheMap[k] = nil
	return nil
}
func (s *TxContextMockTest) CallContract(contractId *commonPb.ContractId, method string, byteCode []byte,
	parameter map[string]string, gasUsed uint64, refTxType commonPb.TxType) (*commonPb.ContractResult, commonPb.TxStatusCode) {
	s.gasUsed = gasUsed
	s.currentDepth = s.currentDepth + 1
	if s.currentDepth > protocol.CallContractDepth {
		contractResult := &commonPb.ContractResult{
			Code:    commonPb.ContractResultCode_FAIL,
			Result:  nil,
			Message: fmt.Sprintf("CallContract too deep %d", s.currentDepth),
		}
		return contractResult, commonPb.TxStatusCode_CONTRACT_TOO_DEEP_FAILED
	}
	if s.gasUsed > protocol.GasLimit {
		contractResult := &commonPb.ContractResult{
			Code:    commonPb.ContractResultCode_FAIL,
			Result:  nil,
			Message: fmt.Sprintf("There is not enough gas, gasUsed %d GasLimit %d ", gasUsed, int64(protocol.GasLimit)),
		}
		return contractResult, commonPb.TxStatusCode_CONTRACT_FAIL
	}
	r, code := s.vmManager.RunContract(contractId, method, byteCode, parameter, s, s.gasUsed, refTxType)

	result := callContractResult{
		deep:         s.currentDepth,
		gasUsed:      s.gasUsed,
		result:       r.Result,
		contractName: contractId.ContractName,
		method:       method,
		param:        parameter,
	}
	s.hisResult = append(s.hisResult, &result)
	s.currentResult = r.Result
	s.currentDepth = s.currentDepth - 1
	return r, code
}

func (s *TxContextMockTest) GetCurrentResult() []byte {
	return s.currentResult
}

func (s *TxContextMockTest) GetTx() *commonPb.Transaction {
	return &commonPb.Transaction{
		Header: &commonPb.TxHeader{
			ChainId:        ChainIdTest,
			Sender:         s.GetSender(),
			TxType:         txType,
			TxId:           "12345678",
			Timestamp:      0,
			ExpirationTime: 0,
		},
		RequestPayload:   nil,
		RequestSignature: nil,
		Result:           nil,
	}
}

func (*TxContextMockTest) GetBlockHeight() int64 {
	return 0
}
func (s *TxContextMockTest) GetTxResult() *commonPb.Result {
	panic("implement me")
}

func (s *TxContextMockTest) SetTxResult(txResult *commonPb.Result) {
	panic("implement me")
}

func (TxContextMockTest) GetTxRWSet(runVmSuccess bool) *commonPb.TxRWSet {
	return &commonPb.TxRWSet{
		TxId:     "txId",
		TxReads:  nil,
		TxWrites: nil,
	}
}

func (s *TxContextMockTest) GetCreator(namespace string) *acPb.SerializedMember {
	return s.creator
}

func (s *TxContextMockTest) GetSender() *acPb.SerializedMember {
	return s.sender
}

func (*TxContextMockTest) GetBlockchainStore() protocol.BlockchainStore {
	return &mockBlockchainStore{}
}

func (*TxContextMockTest) GetAccessControl() (protocol.AccessControlProvider, error) {
	panic("implement me")
}

func (s *TxContextMockTest) GetChainNodesInfoProvider() (protocol.ChainNodesInfoProvider, error) {
	panic("implement me")
}

func (*TxContextMockTest) GetTxExecSeq() int {
	panic("implement me")
}

func (*TxContextMockTest) SetTxExecSeq(i int) {
	panic("implement me")
}

func (s *TxContextMockTest) GetDepth() int {
	return s.currentDepth
}

func BaseParam(parameters map[string]string) {
	parameters[protocol.ContractTxIdParam] = "TX_ID"
	parameters[protocol.ContractCreatorOrgIdParam] = "org_a"
	parameters[protocol.ContractCreatorRoleParam] = "admin"
	parameters[protocol.ContractCreatorPkParam] = "1234567890abcdef1234567890abcdef"
	parameters[protocol.ContractSenderOrgIdParam] = "org_b"
	parameters[protocol.ContractSenderRoleParam] = "user"
	parameters[protocol.ContractSenderPkParam] = "11223344556677889900aabbccddeeff"
	parameters[protocol.ContractBlockHeightParam] = "1"
}

type mockBlockchainStore struct {
}

func (m mockBlockchainStore) GetHeightByHash(blockHash []byte) (uint64, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetBlockHeaderByHeight(height int64) (*commonPb.BlockHeader, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetLastChainConfig() (*configPb.ChainConfig, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetTxHeight(txId string) (uint64, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetArchivedPivot() uint64 {
	panic("implement me")
}

func (m mockBlockchainStore) ArchiveBlock(archiveHeight uint64) error {
	panic("implement me")
}

func (m mockBlockchainStore) RestoreBlocks(serializedBlocks [][]byte) error {
	panic("implement me")
}

func (m mockBlockchainStore) QuerySingle(contractName, sql string, values ...interface{}) (protocol.SqlRow, error) {
	panic("implement me")
}

func (m mockBlockchainStore) QueryMulti(contractName, sql string, values ...interface{}) (protocol.SqlRows, error) {
	panic("implement me")
}

func (m mockBlockchainStore) ExecDdlSql(contractName, sql string) error {
	panic("implement me")
}

func (m mockBlockchainStore) BeginDbTransaction(txName string) (protocol.SqlDBTransaction, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetDbTransaction(txName string) (protocol.SqlDBTransaction, error) {
	panic("implement me")
}

func (m mockBlockchainStore) CommitDbTransaction(txName string) error {
	panic("implement me")
}

func (m mockBlockchainStore) RollbackDbTransaction(txName string) error {
	panic("implement me")
}

func (m mockBlockchainStore) InitGenesis(genesisBlock *storePb.BlockWithRWSet) error {
	panic("implement me")
}

func (m mockBlockchainStore) PutBlock(block *commonPb.Block, txRWSets []*commonPb.TxRWSet) error {
	panic("implement me")
}

func (m mockBlockchainStore) SelectObject(contractName string, startKey []byte, limit []byte) (protocol.StateIterator, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetHistoryForKey(contractName string, key []byte) (protocol.KeyHistoryIterator, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetAccountTxHistory(accountId []byte) (protocol.TxHistoryIterator, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetContractTxHistory(contractName string) (protocol.TxHistoryIterator, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetBlockByHash(blockHash []byte) (*commonPb.Block, error) {
	panic("implement me")
}

func (m mockBlockchainStore) BlockExists(blockHash []byte) (bool, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetBlock(height int64) (*commonPb.Block, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetLastConfigBlock() (*commonPb.Block, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetBlockByTx(txId string) (*commonPb.Block, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetBlockWithRWSets(height int64) (*storePb.BlockWithRWSet, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetTx(txId string) (*commonPb.Transaction, error) {
	panic("implement me")
}

func (m mockBlockchainStore) TxExists(txId string) (bool, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetTxConfirmedTime(txId string) (int64, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetLastBlock() (*commonPb.Block, error) {
	return &commonPb.Block{
		Header: &commonPb.BlockHeader{
			ChainId:        "",
			BlockHeight:    0,
			PreBlockHash:   nil,
			BlockHash:      nil,
			PreConfHeight:  0,
			BlockVersion:   nil,
			DagHash:        nil,
			RwSetRoot:      nil,
			TxRoot:         nil,
			BlockTimestamp: 0,
			Proposer:       nil,
			ConsensusArgs:  nil,
			TxCount:        0,
			Signature:      nil,
		},
		Dag:            nil,
		Txs:            nil,
		AdditionalData: nil,
	}, nil
}

func (m mockBlockchainStore) ReadObject(contractName string, key []byte) ([]byte, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetTxRWSet(txId string) (*commonPb.TxRWSet, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetTxRWSetsByHeight(height int64) ([]*commonPb.TxRWSet, error) {
	panic("implement me")
}

func (m mockBlockchainStore) GetDBHandle(dbName string) protocol.DBHandle {
	panic("implement me")
}

func (m mockBlockchainStore) Close() error {
	panic("implement me")
}