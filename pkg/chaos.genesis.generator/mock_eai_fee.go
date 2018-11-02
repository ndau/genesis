package generator

import (
	"github.com/oneiro-ndev/ndaumath/pkg/address"
	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/oneiro-ndev/ndaumath/pkg/signature"
	"github.com/oneiro-ndev/ndaumath/pkg/signed"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	sv "github.com/oneiro-ndev/system_vars/pkg/system_vars"
	"github.com/pkg/errors"
)

func makeEAIFeeTable() (table sv.EAIFeeTable, err error) {
	add := func(thousandths int64, maker func(int64) (sv.EAIFee, error)) {
		var fee sv.EAIFee
		fee, err = maker(thousandths)
		table = append(table, fee)
	}

	// ndev operations
	add(40, makeEAIFee)
	if err != nil {
		return
	}

	// ntrd operations
	add(10, makeEAIFee)
	if err != nil {
		return
	}

	// rfe acct
	add(1, makeEAIFee)
	if err != nil {
		return
	}

	// rewards nomination acct
	add(1, makeEAIFee)
	if err != nil {
		return
	}

	// node rewards
	add(98, makeNodeRewardEAIFee)
	if err != nil {
		return
	}

	return
}

func makeEAIFee(thousandths int64) (ef sv.EAIFee, err error) {
	var public signature.PublicKey
	public, _, err = signature.Generate(signature.Ed25519, nil)
	if err != nil {
		return
	}
	*ef.To, err = address.Generate(address.KindNdau, public.KeyBytes())
	if err != nil {
		return
	}
	var fee int64
	fee, err = signed.MulDiv(thousandths, constants.QuantaPerUnit, 1000)
	if err != nil {
		return
	}
	ef.Fee = math.Ndau(fee)

	return
}

func makeNodeRewardEAIFee(thousandths int64) (sv.EAIFee, error) {
	fee, err := signed.MulDiv(thousandths, constants.QuantaPerUnit, 1000)
	return sv.EAIFee{
		Fee: math.Ndau(fee),
		To:  nil,
	}, errors.Wrap(err, "making node rewards eai fee")
}
