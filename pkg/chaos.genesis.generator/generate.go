package generator

import (
	"encoding/base64"
	"strings"

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
	err = Update(associated, &asscFile, func() error {
		asscFile[base64.StdEncoding.EncodeToString(bpc)] = ma
		return nil
	})
	if err != nil {
		err = errors.Wrap(err, "updating associated file")
		return
	}

	// update the mockfile
	// we can't use the handy "Update" function here because it needs
	// some custom load/dump logic
	var existingGfile genesisfile.GFile
	existingGfile, err = genesisfile.Load(gfilepath)
	// can't use os.IsNotExist here because it can't handle error wrapping
	if err != nil && !strings.HasSuffix(err.Error(), "no such file or directory") {
		err = errors.Wrap(err, "loading existing genesis file")
		return
	}
	if existingGfile == nil {
		existingGfile = make(genesisfile.GFile)
	}

	for k, v := range gfile {
		existingGfile[k] = v
	}

	err = existingGfile.Dump(gfilepath)
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

	// this is dumb, but required because there is no such thing as
	// a bool pointer literal
	tru := true
	fals := false

	sets := func(key string, val interface{}) (loc svi.Location, err error) {
		loc = svi.Location{Namespace: bpc, Key: []byte(key)}
		err = mock.Set(loc, val)
		if err != nil {
			return
		}
		err = mock.Edit(loc, func(v *genesisfile.Value) error {
			v.System = &tru
			return nil
		})
		if err != nil {
			return
		}

		return
	}

	var sviLoc svi.Location
	sviLoc, err = sets("svi", "SVI stub: will be automatically filled in")
	if err != nil {
		err = errors.Wrap(err, "make svi stub")
		return
	}
	err = mock.Edit(sviLoc, func(v *genesisfile.Value) error {
		v.SVIStub = &tru
		v.System = &fals
		return nil
	})
	if err != nil {
		err = errors.Wrap(err, "fix svi stub")
		return
	}

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
	_, err = sets(sv.ReleaseFromEndowmentAddressName, rfeAddr)
	if err != nil {
		err = errors.Wrap(err, "setting RFE addr")
		return
	}

	// set default rate tables
	_, err = sets(sv.UnlockedRateTableName, eai.DefaultUnlockedEAI)
	if err != nil {
		err = errors.Wrap(err, "setting unlocked eai table")
		return
	}
	_, err = sets(sv.LockedRateTableName, eai.DefaultLockBonusEAI)
	if err != nil {
		err = errors.Wrap(err, "setting locked rate table")
		return
	}

	// make default settlement duration
	ded := sv.DefaultSettlementDuration{Duration: math.Day * 2}
	_, err = sets(sv.DefaultSettlementDurationName, ded)
	if err != nil {
		err = errors.Wrap(err, "setting default settlement duration")
		return
	}

	// make default tx fee script
	// this one is very simple: unconditionally returns numeric 0
	_, err = sets(sv.TxFeeScriptName, vm.MiniAsm("zero").Bytes())
	if err != nil {
		err = errors.Wrap(err, "setting tx fee script")
		return
	}

	// min stake for an account to be active
	_, err = sets(sv.MinStakeName, math.Ndau(1000*constants.QuantaPerUnit))
	if err != nil {
		err = errors.Wrap(err, "setting min stake")
		return
	}

	// make default node goodness script
	// empty: returns the value on top of the stack
	// as goodness functions have the total stake on top of the stack,
	// that's actually not a terrible default
	_, err = sets(sv.NodeGoodnessFuncName, vm.MiniAsm("").Bytes())
	if err != nil {
		err = errors.Wrap(err, "setting goodness func")
		return
	}

	// make eai fee table
	var eaiFeeTable sv.EAIFeeTable
	eaiFeeTable, err = makeEAIFeeTable()
	_, err = sets(sv.EAIFeeTableName, eaiFeeTable)
	if err != nil {
		err = errors.Wrap(err, "setting eai fee table")
		return
	}

	// set default min duration between node rewards nominations
	_, err = sets(sv.MinDurationBetweenNodeRewardNominationsName, math.Duration(1*math.Day))
	if err != nil {
		err = errors.Wrap(err, "setting min duration between nnr txs")
		return
	}

	// set nominate reward
	// generate ownership keys
	ma[sv.NominateNodeRewardOwnershipName], ma[sv.NominateNodeRewardOwnershipPrivateName], err = signature.Generate(signature.Ed25519, nil)
	if err != nil {
		return
	}
	// now generate and set the address
	nnrOwnership := ma[sv.NominateNodeRewardOwnershipName].(signature.PublicKey)
	var nnrAddr address.Address
	nnrAddr, err = address.Generate(address.KindNdau, nnrOwnership.KeyBytes())
	if err != nil {
		return
	}
	_, err = sets(sv.NominateNodeRewardAddressName, nnrAddr)
	if err != nil {
		err = errors.Wrap(err, "setting nnr address")
		return
	}

	// set node reward nomination timeout
	_, err = sets(sv.NodeRewardNominationTimeoutName, math.Duration(30*math.Second))
	if err != nil {
		err = errors.Wrap(err, "setting nnr timeout")
		return
	}

	// set command validator change keys
	// generate ownership keys
	ma[sv.CommandValidatorChangeOwnershipName], ma[sv.CommandValidatorChangeOwnershipPrivateName], err = signature.Generate(signature.Ed25519, nil)
	if err != nil {
		return
	}
	// now generate and set the address
	cvcOwnership := ma[sv.CommandValidatorChangeOwnershipName].(signature.PublicKey)
	var cvcAddr address.Address
	cvcAddr, err = address.Generate(address.KindNdau, cvcOwnership.KeyBytes())
	if err != nil {
		return
	}
	_, err = sets(sv.CommandValidatorChangeAddressName, cvcAddr)
	if err != nil {
		err = errors.Wrap(err, "setting cvc address")
		return
	}

	return
}
