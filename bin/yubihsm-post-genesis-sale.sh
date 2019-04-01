#! /bin/bash

SIGN=~/yubi/sign.py
B64TONDAU=~/go/src/github.com/oneiro-ndev/commands/cmd/keytool/keytool\ ed\ raw\ signature\ --stdin\ --b64

MAXNODE=4           # Node numbers start at 0
SIGNABLE=".sb"
SIGTEMP=".sigs.temp"
SIGS=".sigs"

BPC="A, B, C"       # Ken (A), Rob (B), Chris (C)
BPCKEY=101          # 101, 102, 103 are BPC - we don't use the other two
AXIOM="B, C, D, F"  # Rob (B), Chris (C), Ciarán (D), Ed (F)
AXIOMKEY=104
NDEV="C, E, F"      # Chris (C), Kent (E), Ed (F)
NDEVKEY=105
NDAU="C, E, F"      # Chris (C), Kent (E), Ed (F)
NDAUKEY=106
NTRD="B, D, G"      # Rob (B), Ciarán (D), Steve (F)
NTRDKEY=107

# Start from scratch - remove all signatures

rm -f *.sigs

# ReleaseFromEndowment - Axiom Foundation releases newly-sold ndau to its own account - 2 Axiom signatures required

for ORDINAL in first SECOND; do
    RFE="01-RFE"
    echo "1. Release From Endowment - Insert the" $ORDINAL "Axiom key -" $AXIOM
    echo  \"`$SIGN $RFE$SIGNABLE $AXIOMKEY | $B64TONDAU`\", >> $ISSUE$SIGTEMP
done
sed '$s/,//' $RFE$SIGTEMP > $RFE$SIGS
rm $RFE$SIGTEMP

# TransferAndLock - Axiom Foundation transfers newly-released ndau to purchaser, locking it

for ORDINAL in first SECOND; do
    TRANSFERANDLOCK="02-TransferAndLock"
    echo "2. Transfer and Lock to purchaser account - Insert the" $ORDINAL "Axiom key -" $NDEV
    for n in $(seq 0 $MAXNODE); do
        echo \"`$SIGN $TRANSFERANDLOCK$SIGNABLE $NDEVKEY | $B64TONDAU`\", >> $TRANSFERANDLOCK$SIGTEMP
    done
done
for n in $(seq 0 $MAXNODE); do
    sed '$s/,//' $TRANSFERANDLOCK$SIGTEMP > $TRANSFERANDLOCK$SIGS
    rm $TRANSFERANDLOCK$SIGTEMP
done

# Issue - Axiom reports number of pre-sale ndau issued - 2 Axiom signatures required

for ORDINAL in first SECOND; do
    ISSUE="03-Issue"
    echo "1. Issue - Insert the" $ORDINAL "Axiom key -" $AXIOM
    echo  \"`$SIGN $ISSUE$SIGNABLE $AXIOMKEY | $B64TONDAU`\", >> $ISSUE$SIGTEMP
done
sed '$s/,//' $ISSUE$SIGTEMP > $ISSUE$SIGS
rm $ISSUE$SIGTEMP