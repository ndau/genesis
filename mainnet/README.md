# Mainnet genesis setup and testing

These procedures replicate components of the mainnet genesis transaction block for testing. This is a simplified process designed to chase down a few bugs.

These accounts and transactions are the real ones to be used at mainnet genesis. These validation keypairs have been generated with YubiHSM hardware security modules and these transactions have been signed by them.

## Setup

1. Extract  `genesis_files_mainnet.tar` into an empty `.localnet` directory.
2. Run the standard `./bin/setup.sh 1` script
3. Run `./run-nofinalize.sh` in this directory: this is the standard `./bin/run.sh` script without the ndau and chaos finalize steps.
4. Run `./bin/ndauapi.sh`

You should now have a running blockchain and API with no accounts defined and a complete set of system variables properly loaded.

## Testing

The `./submitTx.py` script processes an input JSON file, prevalidating and (optionally) submitting the transactions it specifies. Separate JSON files for each set of genesis transactions are included here. The script `./transactions.sh` will run `./submitTx.py` to submit four of the five JSON blocks (see below for block 3). Each block can be submitted independently, in order.

1. Submit 10 `ClaimAccount` transactions.
2. Submit 5 `CreditEAI` transactions, one for each node account.
3. Submit 5 `Lock` transactions, one for each node account.
4. Submit 5 `SetRewardsDestination` transactions, one for each node account.

## Using with ETL data

The ETL script may be run between steps 2 and 3 of the **Setup** process above to load real account data. This should be done for testing the effect of the `CreditEAI` transactions in step 2 of the **Testing** procedure, as otherwise there will be no EAI to credit.

## Current Problems

1. The _SetRewardsDestination_ transactions (in the file **5-SetRewardsDestination.json**) do not work, giving invalid signature errors. I don't think this is correct: these were signed through the same procedure as the **3-Lock.json** transactions, which work. The validation script handles each of those two transactions identically.

2. If the ETL process is inserted between setup steps 2 and 3, EAI is correctly credited to some accounts but not to others, even for accounts delegated to the same node.

3. In addition, EAI fees are not credited to the various accounts which are to receive them. The _Transfer_ transactions in **3-Transfer.json** cannot transfer accrued EAI fees to the 5 node accounts so they can stake the minimum 1,000 ndau required. This set of JSON transactions is omitted from the `./transactions.sh` script.