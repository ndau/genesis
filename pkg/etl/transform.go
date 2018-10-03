package etl

import (
	"github.com/oneiro-ndev/genesis/pkg/config"
	"github.com/oneiro-ndev/ndau/pkg/ndau/backing"
	"github.com/oneiro-ndev/ndaumath/pkg/constants"
	math "github.com/oneiro-ndev/ndaumath/pkg/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// TransformRow creates an AccountData struct given a RawRow from the spreadsheet
func TransformRow(row RawRow, logger logrus.FieldLogger) (ad backing.AccountData, err error) {
	ad.Balance = math.Ndau(constants.QuantaPerUnit * row.QtyPurchased)

	creation, err := math.TimestampFrom(row.PurchaseDate)
	if err != nil {
		return ad, errors.Wrap(err, "Creating ndau timestamp from row.PurchaseDate")
	}
	ad.LastEAIUpdate = creation
	ad.LastWAAUpdate = creation

	if row.UnlockDate != nil {
		unlockDate, err := math.TimestampFrom(*row.UnlockDate)
		if err != nil {
			return ad, errors.Wrap(err, "Creating ndau timestamp from row.UnlockDate")
		}

		lock := backing.Lock{
			NoticePeriod: unlockDate.Since(creation),
			UnlocksOn:    &unlockDate,
		}

		ad.Lock = &lock
	}

	if row.RewardTarget != nil {
		addr, err := addressFrom(*row.RewardTarget, logger.WithField("column", config.RewardTargetS))
		if err != nil {
			return ad, err
		}
		ad.RewardsTarget = &addr
	}

	return
}
