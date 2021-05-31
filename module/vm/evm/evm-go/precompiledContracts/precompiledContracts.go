/*
 * Copyright 2020 The SealEVM Authors
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */
package precompiledContracts

import (
	"chainmaker.org/chainmaker/common/evmutils"
)

type PrecompiledContract interface {
	//SetValue(v string)
	GasCost(input []byte) uint64
	Execute(input []byte) ([]byte, error)
}

const ContractsMaxAddress = 16

var Contracts = [ContractsMaxAddress]PrecompiledContract{
	1:  &ecRecover{},
	2:  &sha256hash{},
	3:  &ripemd160hash{},
	4:  &dataCopy{},
	5:  &bigModExp{},
	6:  &bn256AddIstanbul{},
	7:  &bn256ScalarMulIstanbul{},
	8:  &bn256PairingIstanbul{},
	9:  &blake2F{},
	10: &senderOrgId{},
	11: &senderRole{},
	12: &senderPk{},
	13: &creatorOrgId{},
	14: &creatorRole{},
	15: &creatorPk{},
}

func IsPrecompiledContract(address *evmutils.Int) bool {
	if address.IsUint64() {
		addr := address.Uint64()
		return addr < ContractsMaxAddress
	}

	return false
}
