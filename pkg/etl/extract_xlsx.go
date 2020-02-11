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
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/oneiro-ndev/chaincode/pkg/vm"
	"github.com/oneiro-ndev/genesis/pkg/config"
	"github.com/oneiro-ndev/ndaumath/pkg/signature"
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
)

func extractXlsxRow(xrow *xlsx.Row, conf *config.Config, date1904 bool) (rr RawRow, err error) {
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
			return RawRow{}, errors.Wrap(err, config.QtyPurchasedS)
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
			return RawRow{}, errors.Wrap(err, config.PurchaseDateS)
		}
	}
	udc := getCell(config.UnlockDateS)
	if udc != nil {
		ud, err := udc.GetTime(date1904)
		if err != nil {
			return RawRow{}, errors.Wrap(err, config.UnlockDateS)
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
	dnc := getCell(config.DelegationNodeS)
	if dnc != nil {
		dns := dnc.String()
		if len(dns) > 0 && !(strings.EqualFold("false", dns) || dns == "0") {
			rr.DelegationNode = &dns
		}
	}
	vp1 := getCell(config.ValidationPublic1S)
	if vp1 != nil {
		vp1k, err := signature.ParsePublicKey(vp1.String())
		if err != nil {
			return RawRow{}, errors.Wrap(err, config.ValidationPublic1S)
		}
		rr.ValidationPublic = append(rr.ValidationPublic, *vp1k)
	}
	vp2 := getCell(config.ValidationPublic2S)
	if vp2 != nil {
		vp2k, err := signature.ParsePublicKey(vp2.String())
		if err != nil {
			return RawRow{}, errors.Wrap(err, config.ValidationPublic2S)
		}
		rr.ValidationPublic = append(rr.ValidationPublic, *vp2k)
	}
	vs := getCell(config.ValidationScriptS)
	if vs != nil {
		scriptBytes, err := base64.StdEncoding.DecodeString(vs.String())
		if err != nil {
			return RawRow{}, errors.Wrap(err, config.ValidationScriptS)
		}
		rr.ValidationScript = vm.ConvertToOpcodes(scriptBytes)
		if err = rr.ValidationScript.IsValid(); err != nil {
			return RawRow{}, errors.Wrap(err, config.ValidationScriptS)
		}
	}

	return rr, nil
}

func extractXlsx(conf *config.Config) ([]RawRow, error) {
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
		raw, err := extractXlsxRow(sheet.Rows[row], conf, file.Date1904)
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
