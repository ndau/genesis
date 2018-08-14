package etl

import (
	"fmt"
	"strings"
	"time"

	"github.com/oneiro-ndev/genesis/pkg/config"
	"github.com/tealeg/xlsx"
)

// RawRow encapsulates the raw data of a single row of the ndau spreadsheet
type RawRow struct {
	RowNumber    uint64
	Address      string
	QtyPurchased float64
	PurchaseDate time.Time
	UnlockDate   *time.Time
	RewardTarget *string
}

func (rr RawRow) String() string {
	rt := "nil"
	if rr.RewardTarget != nil {
		rt = *rr.RewardTarget
	}
	return fmt.Sprintf(
		"%d %s: %f ndau on %s, unlocking %s, rewards to %s",
		rr.RowNumber,
		rr.Address,
		rr.QtyPurchased,
		rr.PurchaseDate,
		rr.UnlockDate,
		rt,
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

func extractRow(xrow *xlsx.Row, conf *config.Config, date1904 bool) (rr RawRow, err error) {
	getCell := func(name string) *xlsx.Cell {
		col := conf.Columns[name]
		if col >= len(xrow.Cells) {
			return nil
		}
		return xrow.Cells[col]
	}

	ac := getCell(config.AddressS)
	if ac != nil {
		rr.Address = ac.String()
	}
	qpc := getCell(config.QtyPurchasedS)
	if qpc != nil {
		rr.QtyPurchased, err = qpc.Float()
		if err != nil {
			return RawRow{}, err
		}
	}

	// any cell with 0 purchased should be considered blank
	if rr.QtyPurchased == 0 {
		return RawRow{}, blank{}
	}

	pdc := getCell(config.PurchaseDateS)
	if pdc != nil {
		rr.PurchaseDate, err = pdc.GetTime(date1904)
		if err != nil {
			return RawRow{}, err
		}
	}
	udc := getCell(config.UnlockDateS)
	if udc != nil {
		ud, err := udc.GetTime(date1904)
		if err != nil {
			return RawRow{}, err
		}
		if ud != xlsx.TimeFromExcelTime(0, date1904) {
			rr.UnlockDate = &ud
		}
	}
	rtc := getCell(config.RewardTargetS)
	if rtc != nil {
		rts := rtc.String()
		if len(rts) > 0 && !(strings.EqualFold("false", rts) || rts == "0") {
			rr.RewardTarget = &rts
		}
	}

	return rr, nil
}

// Extract the input spreadsheet into a list of raw rows
func Extract(conf *config.Config) ([]RawRow, error) {
	err := conf.CheckColumns()
	if err != nil {
		return nil, err
	}

	file, err := xlsx.OpenFile(conf.Path)
	if err != nil {
		return nil, err
	}
	sheet, ok := file.Sheet[conf.Sheet]
	if !ok {
		return nil, fmt.Errorf("Sheet '%s' not found in %s", conf.Sheet, conf.Path)
	}

	raws := make([]RawRow, 0, sheet.MaxRow-conf.FirstRow+1)
	for row := conf.FirstRow; row < sheet.MaxRow; row++ {
		raw, err := extractRow(sheet.Rows[row], conf, file.Date1904)
		raw.RowNumber = uint64(row + 1)
		if err != nil {
			if isBlank(err) {
				// if it's just a blank row, it's not a problem
				continue
			}
			return nil, fmt.Errorf("Failure to extract row %d: %s", raw.RowNumber, err.Error())
		}
		raws = append(raws, raw)
	}

	return raws, nil
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
