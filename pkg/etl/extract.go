package etl

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/oneiro-ndev/genesis/pkg/config"
)

// RawRow encapsulates the raw data of a single row of the ndau spreadsheet
type RawRow struct {
	RowNumber      uint64
	Address        string
	QtyPurchased   float64
	PurchaseDate   time.Time
	UnlockDate     *time.Time
	RewardTarget   *string
	DelegationNode *string
}

func (rr RawRow) String() string {
	rt := "nil"
	if rr.RewardTarget != nil {
		rt = *rr.RewardTarget
	}
	rd := "nil"
	if rr.DelegationNode != nil {
		rd = *rr.DelegationNode
	}
	return fmt.Sprintf(
		"%d %s: %f ndau on %s, unlocking %s, rewards to %s, delegated to %s",
		rr.RowNumber,
		rr.Address,
		rr.QtyPurchased,
		rr.PurchaseDate,
		rr.UnlockDate,
		rt,
		rd,
	)
}

type blank struct{}

func (blank) Error() string {
	return "Row is blank"
}

func isBlank(err error) bool {
	_, ok := err.(blank)
	return ok
}

// Extract the input spreadsheet into a list of raw rows
func Extract(conf *config.Config) ([]RawRow, error) {
	err := conf.CheckColumns()
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(conf.Path))
	switch ext {
	case ".xlsx":
		return extractXlsx(conf)
	case ".csv":
		return extractCSV(conf)
	default:
		return nil, fmt.Errorf("unknown file extension: %q", ext)
	}
}

// DuplicateAddresses returns a map of every address referred to by more than one row
func DuplicateAddresses(rows []RawRow) (duplicates map[string][]uint64) {
	// fill the duplicates map
	duplicates = make(map[string][]uint64)
	for _, row := range rows {
		duplicates[row.Address] = append(duplicates[row.Address], row.RowNumber)
	}
	// remove entries which are distinct
	for addr := range duplicates {
		if len(duplicates[addr]) <= 1 { // 0 shouldn't exist, but just in case
			delete(duplicates, addr)
		}
	}
	return
}
