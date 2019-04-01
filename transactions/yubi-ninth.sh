#! /bin/bash

SIGN=~/yubi/sign.py
B64TONDAU=~/go/src/github.com/oneiro-ndev/commands/keytool\ ed\ raw\ signature\ --stdin\ --b64

MAXNODE=0           # Node numbers start at 0
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

# SetRewardDestination - all node accounts send EAI to ndev operations - 2 ndev signatures required

for ORDINAL in first SECOND; do
    SETREWARDDEST="04-SetRewardsDestination-0"
    echo "4. Set node reward destinations - Insert the" $ORDINAL "ndev key -" $NDEV
    for n in $(seq 0 $MAXNODE); do
        echo \"`$SIGN $SETREWARDDEST$n$SIGNABLE $NDEVKEY | $B64TONDAU`\", >> $SETREWARDDEST$n$SIGTEMP
    done
done
for n in $(seq 0 $MAXNODE); do
    sed '$s/,//' $SETREWARDDEST$n$SIGTEMP > $SETREWARDDEST$n$SIGS
    rm $SETREWARDDEST$n$SIGTEMP
done