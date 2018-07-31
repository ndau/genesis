package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/oneiro-ndev/genesis/pkg/config"
	"github.com/oneiro-ndev/genesis/pkg/etl"
)

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}

func getNdauhome() string {
	nh := os.ExpandEnv("$NDAUHOME")
	if len(nh) > 0 {
		return nh
	}
	return filepath.Join(os.ExpandEnv("$HOME"), ".ndau")
}

func main() {
	ndauhome := getNdauhome()
	path := config.DefaultConfigPath(ndauhome)
	var rows []etl.RawRow
	var err error
	err = config.WithConfig(path, func(conf *config.Config) error {
		rows, err = etl.Extract(conf)
		if err != nil {
			return err
		}
		duplicates := etl.DuplicateAddresses(rows)
		if len(duplicates) > 0 {
			fmt.Println("ERROR: duplicate addresses present:")
			for addr, rows := range duplicates {
				fmt.Printf("  %s:\n", addr)
				fmt.Printf("    ")
				for _, row := range rows {
					fmt.Printf("%d ", row)
				}
				fmt.Println()
			}
			return errors.New("duplicate addresses")
		}

		return etl.Load(conf, rows, ndauhome)
	})
	check(err)
}
