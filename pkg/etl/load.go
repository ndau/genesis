package etl

import (
	"crypto/sha256"
	"os"

	"github.com/oneiro-ndev/genesis/pkg/config"
	metast "github.com/oneiro-ndev/metanode/pkg/meta/state"
	"github.com/oneiro-ndev/ndau/pkg/ndau"
	"github.com/oneiro-ndev/ndau/pkg/ndau/backing"
	nconfig "github.com/oneiro-ndev/ndau/pkg/ndau/config"
	"github.com/oneiro-ndev/ndaumath/pkg/address"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func addressFrom(addrS string, rlogger logrus.FieldLogger) (address.Address, error) {
	addr, err := address.Validate(addrS)
	if err != nil {
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

// Load the given rows into the noms configuration
func Load(conf *config.Config, rows []RawRow, ndauhome string) error {
	nconf, err := nconfig.LoadDefault(nconfig.DefaultConfigPath(ndauhome))
	if err != nil {
		return errors.Wrap(err, "Failed to load ndau config")
	}

	nomsPath := os.ExpandEnv(conf.NomsPath)

	app, err := ndau.NewAppSilent(nomsPath, "", -1, *nconf)
	if err != nil {
		return errors.Wrap(err, "Failed to create ndau app")
	}

	logger := logrus.StandardLogger()
	app.SetLogger(logger)

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

			ad, err := TransformRow(row, logger.WithField("column", config.AddressS))
			if err != nil {
				rlogger.WithError(err).Error("failed to transform row")
				return st, errors.Wrap(err, "failed to transform row")
			}

			st.Accounts[addr.String()] = ad

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
