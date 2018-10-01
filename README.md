# `genesis`: initialize all of oneiro's blockchains

All of oneiro's blockchains need to be initialized in a variety of ways. For simple cases, we can simply define configuration files and copy the data in those files into the appropriate places, but for other cases, we need something more complex. `genesis` has the following responsibilities:

## Quick start

```
git clone git@github.com:oneiro-ndev/genesis.git $GOPATH/src/github.com/oneiro-ndev/genesis
glide install
go build -o ./etl cmd/etl/main.go
## Put ndau.xlsx in the project root.
./etl
```

## Basic configuration

Read configuration data about the chaos chain, such as the names, addresses, and public keys of its genesis nodes, and transform this into an appropriate `genesis.json` file. Do the same for the ndau chain and the order chain.

## initialize system variables

Read system variable initial values from a configuration file and write them to a noms database appropriately. Produce a tarball of this noms database for distribution to the chaos nodes.

## initialize accounts

- extract data from the interim spreadsheet
- transform it appropriately
- load it into a noms database appropriately.

Once complete, produce a tarball of this noms database for distribution to the ndau nodes.

Note that this step expects the file `ndau.xlsx` to exist in the root directory, but it is excluded from the repository for security reasons.

## ndau config

Generate an appropriate configuration file for each ndau node.

