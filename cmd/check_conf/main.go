package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	path := etl.DefaultConfigPath(getNdauhome())
	err := etl.WithConfig(path, func(c *etl.Config) error {
		return c.CheckColumns()
	})
	check(err)
	fmt.Println(path)
}
