package etl

// ----- ---- --- -- -
// Copyright 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"crypto/sha256"
	"fmt"
	"os"
	"reflect"

	"github.com/mitchellh/go-homedir"
	"github.com/ndau/genesis/pkg/config"
	metast "github.com/ndau/metanode/pkg/meta/state"
	"github.com/ndau/ndau/pkg/ndau"
	"github.com/ndau/ndau/pkg/ndau/backing"
	nconfig "github.com/ndau/ndau/pkg/ndau/config"
	"github.com/ndau/ndaumath/pkg/address"
	"github.com/ndau/ndaumath/pkg/eai"
	"github.com/ndau/system_vars/pkg/genesisfile"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func addressFrom(addrS string, rlogger logrus.FieldLogger) (address.Address, error) {
	addr, err := address.Validate(addrS)
	if err != nil {
		if rlogger != nil {
			rlogger = rlogger.WithField("spreadsheet address", addrS)
		}
		addrHash := sha256.Sum256([]byte(addrS)) // the address data may be too short
		addr, err = address.Generate(address.KindUser, addrHash[:])
		if err != nil {
			if rlogger != nil {
				rlogger.WithError(err).Error("failed to generate address")
			}
			return addr, errors.Wrap(err, "failed to generate address")
		}
		if rlogger != nil {
			rlogger.WithField("generated_address", addr.String()).Warn("invalid address in spreadsheet")
		}
	}
	return addr, nil
}

func rateTable(config *config.Config) (*eai.RateTable, error) {
	gfilePath, err := homedir.Expand(os.ExpandEnv(config.GenesisToml))
	if err != nil {
		return nil, errors.Wrap(err, "expanding genesisfile path")
	}
	gfile, err := genesisfile.Load(gfilePath)
	if err != nil {
		return nil, err
	}

	var value genesisfile.Value
	for k, v := range gfile {
		if k == "LockedRateTable" {
			value = v
			break
		}
	}

	rti, err := value.Unpack()
	if err != nil {
		return nil, errors.Wrap(err, "unpacking LockedRateTable")
	}

	vals, ok := rti.([]genesisfile.Valuable)
	if !ok {
		return nil, fmt.Errorf(
			"LockedRateTable didn't unpack to []genesisfile.Value, instead: %s",
			reflect.ValueOf(rti).Type(),
		)
	}
	rows := make([]eai.RTRow, len(vals))
	for i := range vals {
		rtr, ok := vals[i].(*eai.RTRow)
		if !ok {
			return nil, fmt.Errorf(
				"LockedRateTable row %d didn't unpack to *eai.RTRow, instead: %s",
				i,
				reflect.ValueOf(vals[i]).Type(),
			)
		}
		rows[i] = *rtr
	}
	rt := eai.RateTable(rows)
	return &rt, nil
}

// Load the given rows into the noms configuration
func Load(conf *config.Config, rows []RawRow, ndauhome string) error {
	nconf, err := nconfig.LoadDefault(nconfig.DefaultConfigPath(ndauhome))
	if err != nil {
		return errors.Wrap(err, "Failed to load ndau config")
	}

	nomsPath := os.ExpandEnv(conf.NomsPath)

	app, err := ndau.NewApp(nomsPath, "", -1, *nconf)
	if err != nil {
		return errors.Wrap(err, "Failed to create ndau app")
	}

	logger := logrus.StandardLogger()
	app.SetLogger(logger)

	rt, err := rateTable(conf)
	if err != nil || rt == nil {
		logger.WithError(err).Error("failed to load rate table")
		return errors.Wrap(err, "failed to load rate table")
	}

	logger.WithField("noms_path", nomsPath).Info("Initialized temporary app for state update")
	err = app.UpdateStateImmediately(func(stI metast.State) (metast.State, error) {
		st := stI.(*backing.State)

		for _, row := range rows {
			rlogger := logger.WithFields(logrus.Fields{
				"row":     row.RowNumber,
				"address": row.Address,
			})
			addr, err := addressFrom(row.Address, rlogger)
			if err != nil {
				return st, err
			}

			ad, err := TransformRow(row, logger.WithField("column", config.AddressS), *rt)
			if err != nil {
				rlogger.WithError(err).Error("failed to transform row")
				return st, errors.Wrap(err, "failed to transform row")
			}

			st.Accounts[addr.String()] = ad
			st.TotalRFE += ad.Balance

			// update the state's delegated nodes
			if ad.DelegationNode != nil {
				dest := st.Delegates[ad.DelegationNode.String()]
				if dest == nil {
					dest = make(map[string]struct{})
				}
				dest[addr.String()] = struct{}{}
				st.Delegates[ad.DelegationNode.String()] = dest
			}

			// we could manually compute EAI at this point, but it's
			// better to wait for the actual delegated node to compute it.
			// Two reasons for this:
			//   1.  We ensure that the calculation is precisely what it would
			//       normally be; there's no need to keep two codebases in sync
			//       in case the semantics of EAI end up changing
			//   2.  We don't have a special case where an account without a
			//       delegated node still gets EAI, somehow.
			//
			// In any case, the total EAI won't change if there is a delay
			// before first calculation; that's the point of using continuous
			// interest.
		}

		return st, nil
	})
	if err != nil {
		return errors.Wrap(err, "Updating state")
	}
	return errors.Wrap(app.Close(), "Closing app")
}
