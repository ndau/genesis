package etl

import (
	metast "github.com/oneiro-ndev/metanode/pkg/meta/state"
	"github.com/oneiro-ndev/ndau/pkg/ndau"
	"github.com/oneiro-ndev/ndau/pkg/ndau/backing"
	nconfig "github.com/oneiro-ndev/ndau/pkg/ndau/config"
	"github.com/oneiro-ndev/ndaumath/pkg/address"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Load the given rows into the noms configuration
func Load(conf *Config, rows []RawRow, ndauhome string) error {
	nconf, err := nconfig.LoadDefault(nconfig.DefaultConfigPath(ndauhome))
	if err != nil {
		return errors.Wrap(err, "Failed to load ndau config")
	}

	app, err := ndau.NewApp(conf.NomsPath, *nconf)
	logger := app.GetLogger()
	logger = logger.WithField("bin", "genesis")
	app.SetLogger(logger)

	err = app.UpdateStateImmediately(func(stI metast.State) (metast.State, error) {
		st := stI.(*backing.State)

		for _, row := range rows {
			rlogger := logger.WithFields(logrus.Fields{
				"row":     row.RowNumber,
				"address": row.Address,
			})
			addr, err := address.Validate(row.Address)
			if err != nil {
				rlogger.Warn("invalid address in spreadsheet; generating one")
				addr, err = address.Generate(address.KindUser, []byte(row.Address))
				if err != nil {
					rlogger.WithError(err).Error("failed to generate address")
					return st, errors.Wrap(err, "failed to generate address")
				}
			}

			ad, err := TransformRow(row)
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
	return errors.Wrap(err, "Updating state")
}
