package genesis

// ----- ---- --- -- -
// Copyright 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/oneiro-ndev/genesis/pkg/config"
	"github.com/pkg/errors"
)

// ProcessGenesisJSON updates genesis.json for a given chain
//
// currently it only updates the chain name and genesis time
func ProcessGenesisJSON(conf *config.Config, chainName, path string) error {
	// read genesis.json
	gjBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Reading %s", path))
	}

	var genesisI interface{}
	err = json.Unmarshal(gjBytes, &genesisI)
	genesis, isMap := genesisI.(map[string]interface{})
	if !isMap {
		return errors.New("genesis.json doesn't unpack into map[string]interface{}")
	}

	// update genesis data
	genesis["chain_id"] = chainName

	// write back to genesis.json
	gjBytes, err = json.Marshal(genesis)
	if err != nil {
		return errors.Wrap(err, "Couldn't remarshal genesis.json into json")
	}
	err = ioutil.WriteFile(path, gjBytes, 0600)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Writing %s", path))
	}
	return nil
}
