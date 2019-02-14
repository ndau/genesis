#! /bin/bash

# Submit all 10 blocks of post-genesis transactions

for f in ../transactions/[01]*.json; do
    ../bin/submitTx.py --main --input $d --submit
done