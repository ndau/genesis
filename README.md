# `genesis`: initialize all of oneiro's blockchains

All of oneiro's blockchains need to be initialized in a variety of ways. For simple cases, we can simply define configuration files and copy the data in those files into the appropriate places, but for other cases, we need something more complex. `genesis` has the following responsibilities:

## Existing Features
# `etl`: process spreadsheet .csv file and populate ndau noms database

## To run the ETL process on a fresh spreadsheet (of the form "DashData XX-XX-19 release XX.csv")

- run process_csv.py script on the current spreadsheet to merge the user IDs in the spreadsheet with account addresses in the MongoDB.  This step requires that you be whitelisted on the AWS MongoDB server for the IP address that you're running this script from.  The steps to enable the whitelisting for unsupported IP addresses is documented in the process_csv.py script.  Running this script will create a file called "output.csv" which includes all the merged account addresses:

    ```sh
    $ ./process_csv.py -i DashData XX-XX-19 release XX.csv
    ```

- at this point you can run the etl command to push the spreadsheet data from the "output.csv" file to a ndau noms data directory.  This command requires a "config.toml" in the current directory.  To create a config.toml file, copy the config.template file contained in this repo to config.toml in the directory you plan to run the etl tool. Modify the config.toml file to contain the appropriate arguments for the etl tool (see the config.template file for default values):
    - path to input file (default "./outpout.csv")
    - sheet name in the .csv file that contains the data
    - first row in the spreadsheet that contains data
    - path to the ndau noms data directory that will receive the data
    - time of genesis
    - column definition of data in the CSV
        - address column #
        - notify_immediately column #
        - purchase date column #
        - qty purchased column #
        - reward target column #
        - unlock date column #
        - delegate node column #

    ```sh
    $ ~/go/src/github.com/oneiro-ndev/commands/cmd/etl/etl
    ```

- The spreadsheet data should now be contained in new entries in the noms db data directory specified in the config.toml file referenced above.

### Account ETL

The `etl` program can read an input spreadsheet from a ".xlsx" or ".csv" file and produce an appropriate noms database.

```sh
# clone repos if necessary
git clone git@github.com:oneiro-ndev/ndau.git $GOPATH/src/github.com/oneiro-ndev/ndau
git clone git@github.com:oneiro-ndev/genesis.git $GOPATH/src/github.com/oneiro-ndev/genesis
# set up ndau state
cd $GOPATH/src/github.com/oneiro-ndev/ndau
bin/reset.sh && bin/build.sh && bin/init.sh
# set up ETL
cd ../genesis
# IMPORTANT: Update `config.toml` with the path to the source spreadsheet
./bin/etl.sh
# run the ndau node
cd ../ndau
bin/run.sh
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

