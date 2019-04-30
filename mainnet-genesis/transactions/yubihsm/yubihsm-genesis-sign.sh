#! /bin/bash

SIGN=~/go/src/github.com/oneiro-ndev/genesis/bin/yubihsm-sign.py
B64TONDAU=~/go/src/github.com/oneiro-ndev/commands/keytool\ ed\ raw\ signature\ --stdin\ --b64

MAXNODE=4           # Node numbers start at 0
SIGNABLE=".sb"
SIGTEMP=".sigs.temp"
SIGS=".sigs"

AUTHKEY=100      # Original genesis keys were set up with authkey number 100

BPC="A, B, C"       # Ken (A), Rob (B), Chris (C)
BPCKEY=101          # 101, 102, 103 are BPC - we don't use the other two
AXIOM="B, C, D, F"  # Rob (B), Chris (C), Ciarán (D), Ed (F)
AXIOMKEY=104
NDEV="C, E, F"      # Chris (C), Kent (E), Ed (F)
NDEVKEY=105
NDAU="C, E, F"      # Chris (C), Kent (E), Ed (F)
NDAUKEY=106
NTRD="B, C, E"      # Rob (B), Chris (C), Kent (E)
NTRDKEY=107

# Start from scratch - remove all signatures

rm -f *.sigs

# Issue - Axiom reports number of pre-sale ndau issued - 2 Axiom signatures required

for ORDINAL in first SECOND; do
    ISSUE="01-Issue"
    echo "1. Issue - Insert the" $ORDINAL "Axiom key -" $AXIOM
    echo  \"`$SIGN $AUTHKEY $ISSUE$SIGNABLE $AXIOMKEY | $B64TONDAU`\", >> $ISSUE$SIGTEMP
done
sed '$s/,//' $ISSUE$SIGTEMP > $ISSUE$SIGS
rm $ISSUE$SIGTEMP

# Transfer - ndev transfers EAI fees to all active node accounts - 2 ndev signatures required

for ORDINAL in first SECOND; do
    TRANSFER="02-Transfer-0"
    echo "2. Transfer to nodes - Insert the" $ORDINAL "ndev key -" $NDEV
    for n in $(seq 0 $MAXNODE); do
        echo \"`$SIGN $AUTHKEY $TRANSFER$n$SIGNABLE $NDEVKEY | $B64TONDAU`\", >> $TRANSFER$n$SIGTEMP
    done
done
for n in $(seq 0 $MAXNODE); do
    sed '$s/,//' $TRANSFER$n$SIGTEMP > $TRANSFER$n$SIGS
    rm $TRANSFER$n$SIGTEMP
done

# Lock - all node accounts lock themselves to become active - 2 ndev signatures required

for ORDINAL in first SECOND; do
    LOCK="03-Lock-0"
    echo "3. Lock nodes - Insert the" $ORDINAL "ndev key -" $NDEV
    for n in $(seq 0 $MAXNODE); do
        echo \"`$SIGN $AUTHKEY $LOCK$n$SIGNABLE $NDEVKEY | $B64TONDAU`\", >> $LOCK$n$SIGTEMP
    done
done
for n in $(seq 0 $MAXNODE); do
    sed '$s/,//' $LOCK$n$SIGTEMP > $LOCK$n$SIGS
    rm $LOCK$n$SIGTEMP
done

# SetRewardDestination - all node accounts send EAI to ndev operations - 2 ndev signatures required

for ORDINAL in first SECOND; do
    SETREWARDDEST="04-SetRewardsDestination-0"
    echo "4. Set node reward destinations - Insert the" $ORDINAL "ndev key -" $NDEV
    for n in $(seq 0 $MAXNODE); do
        echo \"`$SIGN $AUTHKEY $SETREWARDDEST$n$SIGNABLE $NDEVKEY | $B64TONDAU`\", >> $SETREWARDDEST$n$SIGTEMP
    done
done
for n in $(seq 0 $MAXNODE); do
    sed '$s/,//' $SETREWARDDEST$n$SIGTEMP > $SETREWARDDEST$n$SIGS
    rm $SETREWARDDEST$n$SIGTEMP
done

# RegisterNode - all node accounts become active - 2 ndev signatures required

for ORDINAL in first SECOND; do
    REGISTER="05-RegisterNode-0"
    echo "5. Register nodes - Insert the" $ORDINAL "ndev key -" $NDEV
    for n in $(seq 0 $MAXNODE); do
        echo \"`$SIGN $AUTHKEY $REGISTER$n$SIGNABLE $NDEVKEY | $B64TONDAU`\", >> $REGISTER$n$SIGTEMP
    done
done
for n in $(seq 0 $MAXNODE); do
    sed '$s/,//' $REGISTER$n$SIGTEMP > $REGISTER$n$SIGS
    rm $REGISTER$n$SIGTEMP
done

# NominateNodeRewards - nominate node 0 - 1 ndau network signature required

for ORDINAL in first; do
    NNR="06-NNR"
    echo "6. Nominate node reward - Insert the" $ORDINAL "ndau network key -" $NDAU
    echo \"`$SIGN $AUTHKEY $NNR$SIGNABLE $NDAUKEY | $B64TONDAU`\", >> $NNR$SIGTEMP
done
sed '$s/,//' $NNR$SIGTEMP > $NNR$SIGS
rm $NNR$SIGTEMP

# Delegate - all node accounts delegate to each other - 2 ndev signatures required

for ORDINAL in first SECOND; do
    DELEGATE="07-Delegate-0"
    echo "7. Delegate nodes to themselves - Insert the" $ORDINAL "ndev key -" $NDEV
    for n in $(seq 0 $MAXNODE); do
        echo \"`$SIGN $AUTHKEY $DELEGATE$n$SIGNABLE $NDEVKEY | $B64TONDAU`\", >> $DELEGATE$n$SIGTEMP
    done
done
for n in $(seq 0 $MAXNODE); do
    sed '$s/,//' $DELEGATE$n$SIGTEMP > $DELEGATE$n$SIGS
    rm $DELEGATE$n$SIGTEMP
done

# Delegate - ndev delegates its operating account to node 0 - 2 ndev signatures required

for ORDINAL in first SECOND; do
    DELEGATE="07-Delegate-ndev"
    echo "7. Delegate ndev operating account to node 0 - Insert the" $ORDINAL "ndev key -" $NDEV
    echo \"`$SIGN $AUTHKEY $DELEGATE$SIGNABLE $NDEVKEY | $B64TONDAU`\", >> $DELEGATE$SIGTEMP
done
sed '$s/,//' $DELEGATE$SIGTEMP > $DELEGATE$SIGS
rm $DELEGATE$SIGTEMP

# Delegate - ntrd delegates its operating account to node 0 - 2 ntrd signatures required

for ORDINAL in first SECOND; do
DELEGATE="07-Delegate-ntrd"
    echo "7. Delegate ntrd operating account to node 0 - Insert the" $ORDINAL "ntrd key -" $NTRD
    echo \"`$SIGN $AUTHKEY $DELEGATE$SIGNABLE $NTRDKEY | $B64TONDAU`\", >> $DELEGATE$SIGTEMP
done
sed '$s/,//' $DELEGATE$SIGTEMP > $DELEGATE$SIGS
rm $DELEGATE$SIGTEMP

# RecordPrice - Axiom reports nmarket price - 2 Axiom signatures required

for ORDINAL in first SECOND; do
    RECORDPRICE="08-RecordPrice"
    echo "8. RecordPrice - Insert the" $ORDINAL "Axiom key -" $AXIOM
    echo  \"`$SIGN $AUTHKEY $RECORDPRICE$SIGNABLE $AXIOMKEY | $B64TONDAU`\", >> $RECORDPRICE$SIGTEMP
done
sed '$s/,//' $RECORDPRICE$SIGTEMP > $RECORDPRICE$SIGS
rm $RECORDPRICE$SIGTEMP

# SetSysvar - BPC sets the TransactionFeeScript system variable to install transaction fees - 2 BPC signatures required

for ORDINAL in first SECOND; do
    SETSYSVAR="09-SetSysvar"
    echo "9. SetSysvar - Insert the" $ORDINAL "BPC key -" $BPC
    echo  \"`$SIGN $AUTHKEY $SETSYSVAR$SIGNABLE $AXIOMKEY | $B64TONDAU`\", >> $SETSYSVAR$SIGTEMP
done
sed '$s/,//' $SETSYSVAR$SIGTEMP > $SETSYSVAR$SIGS
rm $SETSYSVAR$SIGTEMP
