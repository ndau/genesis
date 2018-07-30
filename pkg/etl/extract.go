package etl

import (
	"fmt"
	"time"

	"github.com/tealeg/xlsx"
)

// each field of a RawRow needs a const string identifier
const (
	AddressS           = "address"
	QtyPurchasedS      = "qty_purchased"
	PurchaseDateS      = "purchase_date"
	UnlockDateS        = "unlock_date"
	NotifyImmediatelyS = "notify_immediately"
)

// RawRow encapsulates the raw data of a single row of the ndau spreadsheet
type RawRow struct {
	Address      string
	QtyPurchased float64
	PurchaseDate time.Time
	UnlockDate   *time.Time
}

type blank struct{}

func (blank) Error() string {
	return "Row is blank"
}

func isBlank(err error) bool {
	_, ok := err.(blank)
	return ok
}

func extractRow(xrow *xlsx.Row, conf *Config, date1904 bool) (rr RawRow, err error) {
	getCell := func(name string) *xlsx.Cell {
		col := conf.Columns[name]
		if col >= len(xrow.Cells) {
			return nil
		}
		return xrow.Cells[col]
	}

	ac := getCell(AddressS)
	if ac != nil {
		rr.Address = ac.String()
	}
	qpc := getCell(QtyPurchasedS)
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

	pdc := getCell(PurchaseDateS)
	if pdc != nil {
		rr.PurchaseDate, err = pdc.GetTime(date1904)
		if err != nil {
			return RawRow{}, err
		}
	}
	udc := getCell(UnlockDateS)
	if udc != nil {
		ud, err := udc.GetTime(date1904)
		if err != nil {
			return RawRow{}, err
		}
		if ud != xlsx.TimeFromExcelTime(0, date1904) {
			rr.UnlockDate = &ud
		}
	}

	return rr, nil
}

// Extract the input spreadsheet into a list of raw rows
func Extract(conf *Config) ([]RawRow, error) {
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
		if err != nil {
			if isBlank(err) {
				// if it's just a blank row, it's not a problem
				continue
			}
			return nil, fmt.Errorf("Failure to extract row %d: %s", row, err.Error())
		}
		raws = append(raws, raw)
	}

	return raws, nil
}
