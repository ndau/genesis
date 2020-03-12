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
	"bufio"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ndau/chaincode/pkg/vm"
	"github.com/ndau/genesis/pkg/config"
	"github.com/ndau/ndaumath/pkg/signature"
	"github.com/pkg/errors"
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
			return RawRow{}, errors.Wrap(err, config.QtyPurchasedS)
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
			return RawRow{}, errors.Wrap(err, config.PurchaseDateS)
		}
	}
	udc := getCell(config.UnlockDateS)
	if cellHasValue(udc) {
		ud, err := time.Parse("1/2/06", *udc)
		if err != nil {
			return RawRow{}, errors.Wrap(err, config.UnlockDateS)
		}
		rr.UnlockDate = &ud
	}
	rtc := getCell(config.RewardTargetS)
	if cellHasValue(rtc) {
		rr.RewardTarget = rtc
	}
	dnc := getCell(config.DelegationNodeS)
	if cellHasValue(dnc) {
		rr.DelegationNode = dnc
	}
	spc := getCell(config.RecourseS)
	if cellHasValue(spc) {
		rr.RecoursePeriod, err = time.ParseDuration(*spc)
		if err != nil {
			return RawRow{}, errors.Wrap(err, config.RecourseS)
		}
	}
	rsc := getCell(config.RewardSourceS)
	if cellHasValue(rsc) {
		rr.RewardSource = rsc
	}
	vp1 := getCell(config.ValidationPublic1S)
	if cellHasValue(vp1) {
		vp1k, err := signature.ParsePublicKey(*vp1)
		if err != nil {
			return RawRow{}, errors.Wrap(err, config.ValidationPublic1S)
		}
		rr.ValidationPublic = append(rr.ValidationPublic, *vp1k)
	}
	vp2 := getCell(config.ValidationPublic2S)
	if cellHasValue(vp2) {
		vp2k, err := signature.ParsePublicKey(*vp2)
		if err != nil {
			return RawRow{}, errors.Wrap(err, config.ValidationPublic2S)
		}
		rr.ValidationPublic = append(rr.ValidationPublic, *vp2k)
	}
	vs := getCell(config.ValidationScriptS)
	if cellHasValue(vs) {
		scriptBytes, err := base64.StdEncoding.DecodeString(*vs)
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
