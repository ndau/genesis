# `genesis`: initialize all of oneiro's blockchains

All of oneiro's blockchains need to be initialized in a variety of ways. For simple cases, we can simply define configuration files and copy the data in those files into the appropriate places, but for other cases, we need something more complex. `genesis` has the following responsibilities:

## Existing Features

### Account ETL

The `etl` program can read an input spreadsheet from a ".xlsx" or ".csv" file and produce an appropriate noms database.

```sh
git clone git@github.com:oneiro-ndev/genesis.git $GOPATH/src/github.com/oneiro-ndev/genesis
cd $GOPATH/src/github.com/oneiro-ndev/genesis
glide install
go build ./cmd/etl
## Update `config.toml` with the path to the source spreadsheet
./etl
```

Once ETL is complete, the noms database can be examined directly to see the results:

```sh
noms show /Users/prgn/.ndau/ndau/noms::ndau
```

## Planned Future Features

### Basic configuration

Read configuration data about the chaos chain, such as the names, addresses, and public keys of its genesis nodes, and transform this into an appropriate `genesis.json` file. Do the same for the ndau chain and the order chain.

### initialize system variables

Read system variable initial values from a configuration file and write them to a noms database appropriately. Produce a tarball of this noms database for distribution to the chaos nodes.

### initialize accounts

Produce a tarball of the noms database for distribution to the ndau nodes.

### ndau config

Generate an appropriate configuration file for each ndau node.

