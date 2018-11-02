package generator

import (
	"encoding/base64"

	"github.com/oneiro-ndev/chaincode/pkg/vm"
	"github.com/oneiro-ndev/chaos/pkg/genesisfile"
	"github.com/oneiro-ndev/ndaumath/pkg/address"
	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/oneiro-ndev/ndaumath/pkg/eai"
	"github.com/oneiro-ndev/ndaumath/pkg/signature"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/oneiro-ndev/system_vars/pkg/svi"
	sv "github.com/oneiro-ndev/system_vars/pkg/system_vars"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
)

// Generate creates a genesisfile and associated data.
//
// Both arguments are paths to the files which should be written.
//
// Both files are written as TOML. In both cases, if there exists a TOML
// file in that path already, new data will be added but existing data will
// not be overwritten.
func Generate(gfilepath, associated string) (bpc []byte, err error) {
	bpc, _, err = ed25519.GenerateKey(nil)
	if err != nil {
		return
	}

	var gfile genesisfile.GFile
	var ma Associated
	gfile, ma, err = generateData(bpc)
	if err != nil {
		return
	}

	// update the associated file
	asscFile := make(AssociatedFile)
	err = Update(associated, asscFile, func() error {
		asscFile[base64.StdEncoding.EncodeToString(bpc)] = ma
		return nil
	})
	if err != nil {
		err = errors.Wrap(err, "updating associated file")
		return
	}

	// update the mockfile
	existingGfile := make(genesisfile.GFile)
	err = Update(gfilepath, existingGfile, func() error {
		for k, v := range gfile {
			existingGfile[k] = v
		}
		return nil
	})
	if err != nil {
		err = errors.Wrap(err, "updating genesis file")
		return
	}

	return
}

// mock up some chaos data
// this is in a separate, private function because we don't want people blowing
// away their old data; they have to go through the public Generate function
// which preserves old data
func generateData(bpc []byte) (mock genesisfile.GFile, ma Associated, err error) {
	mock = make(genesisfile.GFile)
	ma = make(Associated)

	sets := func(key string, val interface{}) {
		mock.Set(svi.Location{Namespace: bpc, Key: []byte(key)}, val)
	}

	sviLoc := svi.Location{Namespace: bpc, Key: []byte("svi")}
	mock.Set(sviLoc, "SVI stub: do not edit")
	mock.Edit(sviLoc, func(v *genesisfile.Value) error {
		tru := true
		v.SVIStub = &tru
		return nil
	})

	// set RFE address
	// generate ownership keys
	ma[sv.ReleaseFromEndowmentOwnershipName], ma[sv.ReleaseFromEndowmentOwnershipPrivateName], err = signature.Generate(signature.Ed25519, nil)
	if err != nil {
		err = errors.Wrap(err, "generating RFE keys")
		return
	}
	// now generate and set the address
	rfeOwnership := ma[sv.ReleaseFromEndowmentOwnershipName].(signature.PublicKey)
	var rfeAddr address.Address
	rfeAddr, err = address.Generate(address.KindNdau, rfeOwnership.KeyBytes())
	if err != nil {
		err = errors.Wrap(err, "generating RFE addr")
		return
	}
	sets(sv.ReleaseFromEndowmentAddressName, rfeAddr)

	// set default rate tables
	sets(sv.UnlockedRateTableName, eai.DefaultUnlockedEAI)
	sets(sv.LockedRateTableName, eai.DefaultLockBonusEAI)

	// make default settlement duration
	ded := sv.DefaultSettlementDuration{Duration: math.Day * 2}
	sets(sv.DefaultSettlementDurationName, ded)

	// make default tx fee script
	// this one is very simple: unconditionally returns numeric 0
	sets(sv.TxFeeScriptName, vm.MiniAsm("zero").Bytes())

	// min stake for an account to be active
	sets(sv.MinStakeName, math.Ndau(1000*constants.QuantaPerUnit))

	// make default node goodness script
	// empty: returns the value on top of the stack
	// as goodness functions have the total stake on top of the stack,
	// that's actually not a terrible default
	sets(sv.NodeGoodnessFuncName, vm.MiniAsm("").Bytes())

	// make eai fee table
	sets(sv.EAIFeeTableName, makeMockEAIFeeTable())

	// set default min duration between node rewards nominations
	sets(sv.MinDurationBetweenNodeRewardNominationsName, math.Duration(1*math.Day))

	// set nominate reward
	// generate ownership keys
	ma[sv.NominateNodeRewardOwnershipName], ma[sv.NominateNodeRewardOwnershipPrivateName], err = signature.Generate(signature.Ed25519, nil)
	if err != nil {
		panic(err)
	}
	// now generate and set the address
	nnrOwnership := ma[sv.NominateNodeRewardOwnershipName].(signature.PublicKey)
	nnrAddr, err := address.Generate(address.KindNdau, nnrOwnership.KeyBytes())
	if err != nil {
		panic(err)
	}
	sets(sv.NominateNodeRewardAddressName, nnrAddr)

	// set node reward nomination timeout
	sets(sv.NodeRewardNominationTimeoutName, math.Duration(30*math.Second))

	// set command validator change keys
	// generate ownership keys
	ma[sv.CommandValidatorChangeOwnershipName], ma[sv.CommandValidatorChangeOwnershipPrivateName], err = signature.Generate(signature.Ed25519, nil)
	if err != nil {
		panic(err)
	}
	// now generate and set the address
	cvcOwnership := ma[sv.CommandValidatorChangeOwnershipName].(signature.PublicKey)
	cvcAddr, err := address.Generate(address.KindNdau, cvcOwnership.KeyBytes())
	if err != nil {
		panic(err)
	}
	sets(sv.CommandValidatorChangeAddressName, cvcAddr)

	return
}