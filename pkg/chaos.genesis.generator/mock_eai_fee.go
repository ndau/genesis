package generator

import (
	"github.com/oneiro-ndev/ndaumath/pkg/address"
	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/oneiro-ndev/ndaumath/pkg/signature"
	"github.com/oneiro-ndev/ndaumath/pkg/signed"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	sv "github.com/oneiro-ndev/system_vars/pkg/system_vars"
)

func makeMockEAIFeeTable() sv.EAIFeeTable {
	return sv.EAIFeeTable{
		makeMockEAIFee("ndev operations", 40),
		makeMockEAIFee("ntrd operations", 10),
		makeMockEAIFee("rfe account", 1),
		makeMockEAIFee("rewards nomination acct", 1),
		makeMockNodeRewardEAIFee(98),
	}
}

func makeMockEAIFee(_ string, thousandths int64) sv.EAIFee {
	public, _, err := signature.Generate(signature.Ed25519, nil)
	if err != nil {
		panic(err)
	}
	addr, err := address.Generate(address.KindNdau, public.KeyBytes())
	if err != nil {
		panic(err)
	}
	fee, err := signed.MulDiv(thousandths, constants.QuantaPerUnit, 1000)
	if err != nil {
		panic(err)
	}
	return sv.EAIFee{
		Fee: math.Ndau(fee),
		To:  &addr,
	}
}

func makeMockNodeRewardEAIFee(thousandths int64) sv.EAIFee {
	fee, err := signed.MulDiv(thousandths, constants.QuantaPerUnit, 1000)
	if err != nil {
		panic(err)
	}
	return sv.EAIFee{
		Fee: math.Ndau(fee),
		To:  nil,
	}
}
