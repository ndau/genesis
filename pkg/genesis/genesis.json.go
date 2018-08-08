package genesis

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
	genesis["genesis_time"] = conf.Genesis
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
