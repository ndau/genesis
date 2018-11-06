package genesis

import (
	"fmt"

	"github.com/oneiro-ndev/chaos/pkg/chaos"
	"github.com/oneiro-ndev/chaos/pkg/chaos/backing"
	"github.com/oneiro-ndev/chaos/pkg/genesisfile"
	metast "github.com/oneiro-ndev/metanode/pkg/meta/state"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Upload the given genesisfile into the chaos chain
func Upload(gfile genesisfile.GFile, nomsPath string, dryrun bool) error {
	app, err := chaos.NewApp(nomsPath)
	if err != nil {
		return errors.Wrap(err, "failed to create chaos app")
	}

	logger := logrus.StandardLogger()
	app.SetLogger(logger)
	logger.WithField("noms_path", nomsPath).Info("Initialized temporary app for state update")

	err = app.UpdateStateImmediately(func(stI metast.State) (metast.State, error) {
		st := stI.(*backing.State)

		return st, gfile.Iter(func(ns, key []byte, val genesisfile.Value) error {
			vb, err := val.IntoBytes()
			if err != nil {
				logger.WithError(err).WithField("value", fmt.Sprintf("%v", val.Data)).Error("attempting to encode value to bytes")
				return errors.Wrap(err, "attempting to encode value as bytes")
			}

			if !dryrun {
				err = st.SetNamespaced(app.GetDB(), 0, ns, key, vb)
				if err != nil {
					logger.WithError(err).Error("failed to set namespaced k-v")
					return errors.Wrap(err, "attempting set namespaced k-v")
				}
			}

			return nil
		})
	})
	if err != nil {
		logger.WithError(err).Error("updating state")
		return errors.Wrap(err, "updating state")
	}
	return errors.Wrap(app.Close(), "closing app")
}
