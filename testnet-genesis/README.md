# Testnet Genesis

These procedures can be used for any non-mainnet application: they're all called *testnet* here for convenience.

The complete set of post-genesis transactions is defined in **Post-Genesis Transaction List - testnet.csv** and the `../bin/createTxList.py` script is used to generate a JSON file with a signed list of transactions from it. That CSV file contains all public and private keys required to sign these (and any other) transactions. The account addresses and keypairs are the only things that differ between this file and the mainnet version.

```
cd transactions
../../bin/createTxList.py --sign --input "Post-Genesis Transaction Block - testnet.csv" > testnet-genesis.json
```
was used to create `transactions/testnet-genesis.json`. That file of 45 transactions can be submitted as:

```
cd transactions
../../bin/submitTx.py --[main|test|dev|local] --action=[submit|prevalidate|both] --delay 0 --input testnet-genesis.json
```

## Post-Genesis Keypairs

All testnet keypairs are stored as plaintext files in the *keys* directory.

At genesis, seven types of validation keys are needed. For mainnet they are assigned in groups to each of seven people. For testnet/devnet/localnet purposes, this structure is mirrored: there are seven groups, each with seven keys.

The YubiHSM device used for mainnet key management identifies keys by an ID number. Those numbers, 101 - 17, are used here. They are used as described in [BPC Genesis Network Values](https://paper.dropbox.com/doc/BPC-Genesis-Network-Values--AYaA0XDGbeshlcw2Fw~Yn4xKAg-U5qFm5bqpvATFAJj75B6b):

- 101 - First BPC Key
- 102 - Second BPC Key
- 103 - Third BPC Key
- 104 - Axiom Foundation Key
- 105 - ndev Operations Key
- 106 - ndau Network Operations Key
- 107 - ntrd Operations Key

These are grouped in groups 1 - 7 corresponding to the hardware keys held by each of the seven people listed in that document. Not every person has an instance of every key type! The first YubiHSM key is assigned to Ken, and he has three BPC keys numbered 101 - 103. He has no other keys. Steve's YubiHSM key only holds key 107, an ntrd Operations Key. Of the 49 possible group/key combinations, only 22 are actually used. Although nine BPC Keys (three instances of 101 - 103) exist, only the three keys labeled 101 are currently used in validation rules for simplicity. Each of Ken, Rob, and Chris holds an instance of keys 101 - 103 and there's little point in making each of them sign each transaction three times.

These keys are assigned to the 10 accounts claimed in the first 10 transactions after genesis according to the validation rules appropriate for each account.

1. *BPC Operations* - 3 BPC keys:
   1. Group 1 - Key 101
   2. Group 2 - Key 101
   3. Group 3 - Key 101

1. *Axiom Foundation* - 4 Axiom keys and 3 BPC keys:
   1. Group 2 - Key 104
   2. Group 3 - Key 104
   3. Group 4 - Key 104
   4. Group 7 - Key 104
   5. Group 1 - Key 101
   6. Group 2 - Key 101
   7. Group 3 - Key 101

1. *ntrd Operations* - 3 ntrd keys:
   1. Group 2 - Key 107
   2. Group 4 - Key 107
   3. Group 7 - Key 107

1. *ndev Operations* - 3 ndev keys:
   1. Group 3 - Key 105
   2. Group 5 - Key 105
   3. Group 6 - Key 105

1. *ndau Network Operations* - 3 ndau Network keys and 3 BPC keys:
   1. Group 3 - Key 106
   2. Group 5 - Key 106
   3. Group 6 - Key 106
   4. Group 1 - Key 101
   5. Group 2 - Key 101
   6. Group 3 - Key 101

1. *ndev Network Nodes - A unique key for each node and 3 ndev Operations keys:
   1. Node [N] Key
   2. Group 3 - Key 105
   3. Group 5 - Key 105
   4. Group 6 - Key 105