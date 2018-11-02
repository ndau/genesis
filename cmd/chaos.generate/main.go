package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"

	cli "github.com/jawher/mow.cli"
	"github.com/oneiro-ndev/chaos/pkg/genesisfile"
	generator "github.com/oneiro-ndev/genesis/pkg/chaos.genesis.generator"
)

func ndauhome() string {
	ndauhome := os.ExpandEnv("$NDAUHOME")
	if len(ndauhome) == 0 {
		home := os.ExpandEnv("$HOME")
		ndauhome = path.Join(home, ".ndau")
	}
	return ndauhome
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(
			os.Stderr,
			err.Error(),
		)
		os.Exit(1)
	}
}

func main() {
	app := cli.App("chaos.generate", "generate chaos mockfile and associated data")

	var (
		verbose = app.BoolOpt("v verbose", false, "emit more detailed information")
		dryRun  = app.BoolOpt("d dry-run", false, "don't actually generate any data")
		gfpath  = app.StringOpt(
			"g genesisfile",
			genesisfile.DefaultPath(ndauhome()),
			"path to genesisfile",
		)
		afpath = app.StringOpt(
			"a associatedfile",
			generator.DefaultAssociated(ndauhome()),
			"path to genesisfile",
		)
	)

	app.Action = func() {
		if *verbose {
			fmt.Printf("%25s: %s\n", "genesisfile path", *gfpath)
			fmt.Printf("%25s: %s\n", "associatedfile path", *afpath)
		}

		if !*dryRun {
			bpc, err := generator.Generate(*gfpath, *afpath)

			if bpc != nil {
				fmt.Printf(
					"%25s: %s\n",
					"bpc public key",
					base64.StdEncoding.EncodeToString(bpc),
				)
			}

			check(err)
		}
	}

	check(app.Run(os.Args))
}
