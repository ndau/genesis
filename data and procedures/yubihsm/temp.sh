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


# Delegate - all node accounts delegate to themselves - 2 ndev signatures required

for ORDINAL in first SECOND; do
    DELEGATE="07-Delegate-0"
    echo "7. Delegate nodes to themselves - Insert the" $ORDINAL "ndev key -" $NDEV
    for n in $(seq 0 $MAXNODE); do
        echo \"`$SIGN $DELEGATE$n$SIGNABLE $NDEVKEY | $B64TONDAU`\", >> $DELEGATE$n$SIGTEMP
    done
done
for n in $(seq 0 $MAXNODE); do
    sed '$s/,//' $DELEGATE$n$SIGTEMP > $DELEGATE$n$SIGS
    rm $DELEGATE$n$SIGTEMP
done
