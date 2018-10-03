package etl

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/oneiro-ndev/genesis/pkg/config"
)

func cellHasValue(cell *string) bool {
	return cell != nil && !(len(*cell) == 0 || *cell == "0" || strings.EqualFold(*cell, "false"))
}

func extractCSVRow(record []string, conf *config.Config) (rr RawRow, err error) {
	getCell := func(name string) *string {
		col := conf.Columns[name]
		if col >= len(record) {
			return nil
		}
		return &record[col]
	}

	ac := getCell(config.AddressS)
	if cellHasValue(ac) {
		rr.Address = *ac
	}
	qpc := getCell(config.QtyPurchasedS)
	if cellHasValue(qpc) {
		rr.QtyPurchased, err = strconv.ParseFloat(*qpc, 64)
		if err != nil {
			return RawRow{}, err
		}
	}

	// any cell with 0 purchased should be considered blank
	if rr.QtyPurchased == 0 {
		return RawRow{}, blank{}
	}

	pdc := getCell(config.PurchaseDateS)
	if cellHasValue(pdc) {
		rr.PurchaseDate, err = time.Parse("1/2/06", *pdc)
		if err != nil {
			return RawRow{}, err
		}
	}
	udc := getCell(config.UnlockDateS)
	if cellHasValue(udc) {
		ud, err := time.Parse("1/2/06", *udc)
		if err != nil {
			return RawRow{}, err
		}
		rr.UnlockDate = &ud
	}
	rtc := getCell(config.RewardTargetS)
	if cellHasValue(rtc) {
		rr.RewardTarget = rtc
	}

	return rr, nil
}

func extractCSV(conf *config.Config) ([]RawRow, error) {
	fp, err := os.Open(conf.Path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	csvReader := csv.NewReader(bufio.NewReader(fp))

	raws := make([]RawRow, 0)
	row := 0
	for {
		row++
		record, err := csvReader.Read()
		if row <= conf.FirstRow {
			continue
		}
		if record != nil {
			raw, err := extractCSVRow(record, conf)
			raw.RowNumber = uint64(row)
			if err != nil {
				if isBlank(err) {
					// if it's just a blank row, it's not a problem
					continue
				}
				return nil, fmt.Errorf("Failure to extract row %d: %s", raw.RowNumber, err.Error())
			}
			raws = append(raws, raw)
		}
		if err == io.EOF {
			break
		}
	}

	return raws, nil
}
