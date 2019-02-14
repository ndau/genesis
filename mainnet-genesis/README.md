# Mainnet genesis setup and testing

These procedures replicate components of the mainnet genesis transaction block for testing. This is a simplified process designed to chase down a few bugs.

These accounts and transactions are the real ones to be used at mainnet genesis. These validation keypairs have been generated with YubiHSM hardware security modules and these transactions have been signed by them.

## ETL config

The `etl/config.toml` file **must be updated** to use the correct `noms_path`. Check that the `path` setting is pointing to the correct .csv file. The current genesis file is `output-rel48.csv`.

## Setup

1. Extract  `genesis_files_mainnet.tar` into the appropriate directory *[what is this?]*.
1. Run the standard `./bin/setup.sh` script

At this point, load the ETL data:

1. `cd ./etl`
1. `../../../commands/etl`

Start the blockchain:

1. Run `../bin/run-nofinalize.sh`: this is the standard `./bin/run.sh` script without the ndau and chaos finalize steps.
1. Run `./bin/ndauapi.sh`

You should now have a running blockchain and API with no accounts defined and a complete set of system variables properly loaded.

## Submit post-genesis transactions

The `../bin/submitTx.py` script processes an input JSON file, prevalidating and (optionally) submitting the transactions it specifies. 10 separate JSON files for each set of genesis transactions are included in `../transactions`. The script `../bin/transactions.sh` will run `../bin/submitTx.py` to submit all post-genesis transactions in order.

1. `../bin/transactions.sh`

Submits

1. 10 `ClaimAccount` transactions
2. 1 `Issue` transaction
3. 5 `CreditEAI` transactions
4. 5 `Transfer` transactions
5. 5 `Lock` transactions
6. 5 `SetRewardsDestination` transactions
7. 5 `RegisterNode` transactions
8. 1 `NominateNodeRewards` transactions
9. 1 `ClaimReward` transaction
10. 2 `Delegate` transactions for the ndev and ntrd operations accounts and 5 `Delegate` transactions for the 5 network nodes