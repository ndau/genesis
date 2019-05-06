package etl

import (
	"fmt"

	"github.com/oneiro-ndev/genesis/pkg/config"
	"github.com/oneiro-ndev/ndau/pkg/ndau/backing"
	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	"github.com/oneiro-ndev/ndaumath/pkg/eai"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// TransformRow creates an AccountData struct given a RawRow from the spreadsheet
func TransformRow(row RawRow, logger logrus.FieldLogger, bonusTable eai.RateTable) (ad backing.AccountData, err error) {
	ad.Balance = math.Ndau(constants.QuantaPerUnit * row.QtyPurchased)

	creation, err := math.TimestampFrom(row.PurchaseDate)
	if err != nil {
		return ad, errors.Wrap(err, config.PurchaseDateS)
	}
	ad.LastEAIUpdate = creation
	ad.LastWAAUpdate = creation

	if row.QtyPurchased >= 1000 {
		ad.CurrencySeatDate = &creation
	}

	if row.UnlockDate != nil {
		unlockDate, err := math.TimestampFrom(*row.UnlockDate)
		if err != nil {
			return ad, errors.Wrap(err, config.UnlockDateS)
		}

		ad.Lock = backing.NewLock(unlockDate.Since(creation), bonusTable)
		ad.Lock.Notify(creation, 0)
	}

	if row.RewardTarget != nil {
		addr, err := addressFrom(*row.RewardTarget, logger.WithField("column", config.RewardTargetS))
		if err != nil {
			return ad, errors.Wrap(err, config.RewardTargetS)
		}
		ad.RewardsTarget = &addr
	}

	if row.DelegationNode != nil {
		addr, err := addressFrom(*row.DelegationNode, logger.WithField("column", config.DelegationNodeS))
		if err != nil {
			return ad, errors.Wrap(err, config.DelegationNodeS)
		}
		ad.DelegationNode = &addr
	}

	if row.SettlementPeriod > 0 {
		ad.RecourseSettings.Period = math.DurationFrom(row.SettlementPeriod)
	}

	if row.RewardSource != nil {
		addr, err := addressFrom(*row.RewardSource, logger.WithField("column", config.RewardSourceS))
		if err != nil {
			return ad, errors.Wrap(err, config.RewardSourceS)
		}
		ad.IncomingRewardsFrom = append(ad.IncomingRewardsFrom, addr)
	}

	if len(row.ValidationPublic) == 2 && len(row.ValidationScript) > 0 {
		ad.ValidationKeys = row.ValidationPublic
		ad.ValidationScript = row.ValidationScript.Bytes()
	} else if len(row.ValidationPublic) == 0 && len(row.ValidationScript) == 0 {
		// we don't have to do anything in this case, but we _do_ want to distinguish
		// it from the error case where some but not all of these fields were set
	} else {
		err = fmt.Errorf("either all or none of %s, %s, and %s must be set", config.ValidationPublic1S, config.ValidationPublic2S, config.ValidationScriptS)
		return
	}

	return
}
