#! /usr/bin/env python3

import csv
import copy
import json
import argparse
import subprocess


def ClaimAccount(d):
    tx = dict(
        comment=d["header"],
        txtype="ClaimAccount",
        tx=dict(
            target=d["target"],
            ownership=d["ownership"],
            validation_script=d["validation_script"],
            sequence=int(d["sequence"]),
            pvt_key=d["pvt_key"],
            signature="",
        ),
    )
    keys = [d.get(f"validation_keys_{i}", "") for i in range(9)]
    keys = [k for k in keys if k != ""]
    tx["tx"]["validation_keys"] = keys
    return [tx]


def Issue(d):
    tx = dict(
        comment=d["header"],
        txtype="Issue",
        tx=dict(
            qty=int(d["qty"]),
            sequence=int(d["sequence"]),
            signatures=[""],
        ),
    )
    return [tx]


def Delegate(d):
    tx = dict(
        comment=d["header"],
        txtype="Delegate",
        tx=dict(
            target=d["target"],
            node=d["ownership"],
            sequence=int(d["sequence"]),
            pvt_key=d["pvt_key"],
            signatures=[""],
        ),
    )
    return [tx]


def CreditEAI(d):
    tx = dict(
        comment=d["header"],
        txtype="CreditEAI",
        tx=dict(
            node=d["target"],
            sequence=int(d["sequence"]),
            pvt_key=d["pvt_key"],
            signatures=[""],
        ),
    )
    return [tx]


def Lock(d):
    tx = dict(
        comment=d["header"],
        txtype="Lock",
        tx=dict(
            target=d["target"],
            period=d["period"],
            sequence=int(d["sequence"]),
            pvt_key=d["pvt_key"],
            signatures=[""],
        ),
    )
    return [tx]


def SetRewardsDestination(d):
    tx = dict(
        comment=d["header"],
        txtype="SetRewardsDestination",
        tx=dict(
            target=d["target"],
            destination=d["destination"],
            sequence=int(d["sequence"]),
            pvt_key=d["pvt_key"],
            signatures=[""],
        ),
    )
    return [tx]


def Transfer(d):
    tx = dict(
        comment=d["header"],
        txtype="Transfer",
        tx=dict(
            source=d["source"],
            destination=d["destination"],
            qty=int(d["qty"]) * 100_000_000,
            sequence=int(d["sequence"]),
            pvt_key=d["pvt_key"],
            signatures=[""],
        ),
    )
    return [tx]


def RegisterNode(d):
    tx = dict(
        comment=d["header"],
        txtype="RegisterNode",
        tx=dict(
            node=d["target"],
            distribution_script=d["distribution"],
            rpc_address=d["rpc_address"],
            sequence=int(d["sequence"]),
            pvt_key=d["pvt_key"],
            signatures=[""],
        ),
    )
    return [tx]


def NominateNodeReward(d):
    tx = dict(
        comment=d["header"],
        txtype="NominateNodeReward",
        tx=dict(
            random=1,
            sequence=int(d["sequence"]),
            pvt_key=d["pvt_key"],
            signatures=[""],  # nominate the 0 node
        ),
    )
    return [tx]


def ClaimNodeReward(d):
    tx = dict(
        comment=d["header"],
        txtype="ClaimNodeReward",
        tx=dict(
            node=d["target"],
            sequence=int(d["sequence"]),
            pvt_key=d["pvt_key"],
            signatures=[""],
        ),
    )
    return [tx]


def generateSignableBytes(obj, ndautool):
    tx = copy.deepcopy(obj)
    if "signature" in tx:
        del tx["signature"]
    if "signatures" in tx:
        del tx["signatures"]
    if "pvt_key" in tx:
        del tx["pvt_key"]

    j = json.dumps(tx, indent=2)
    r = subprocess.run(
        [ndautool, "signable-bytes", t["txtype"]],
        input=j,
        text=True,
        capture_output=True,
    )
    if r.returncode > 0:
        return f"ERROR: {r.stderr}"
    return r.stdout.strip()


def tryToSign(t, keytool):
    pk = t["tx"].get("pvt_key", None)
    tx = copy.deepcopy(t["tx"])
    if pk is not None:
        del tx["pvt_key"]
    if not pk:
        t["tx"] = tx
        return t

    sb = t["signable_bytes"]

    if not pk.startswith("npvt"):
        sig = pk
    else:
        args = [keytool, "sign", pk, sb, "-b"]
        r = subprocess.run(args, text=True, capture_output=True)
        if r.returncode > 0:
            sig = f"ERROR: {r.stderr}"
        else:
            sig = r.stdout.strip()

    if "signature" in tx:
        tx["signature"] = sig
    else:
        tx["signatures"] = [sig]
    t["tx"] = tx
    return t


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--sign",
        action="store_true",
        help="attempt to sign all transactions for which keys are present",
    )
    parser.add_argument(
        "--keytool",
        action="store",
        default="../../commands/cmd/keytool/keytool",
        help="location of keytool",
    )
    parser.add_argument(
        "--ndautool",
        action="store",
        default="../../commands/ndau",
        help="location of ndautool",
    )
    args = parser.parse_args()

    txmap = dict(
        ClaimAccount=ClaimAccount,
        Issue=Issue,
        Delegate=Delegate,
        CreditEAI=CreditEAI,
        RegisterNode=RegisterNode,
        Lock=Lock,
        SetRewardsDestination=SetRewardsDestination,
        NominateNodeReward=NominateNodeReward,
        ClaimNodeReward=ClaimNodeReward,
        Transfer=Transfer,
    )
    with open("Post-Genesis Transaction Block.csv") as csvfile:
        rdr = csv.DictReader(csvfile)
        rows = [r for r in rdr if r["txtype"] != ""]
        txs = []
        for row in rows:
            txs.extend(txmap[row["txtype"]](row))
        newtxs = []

        if args.sign:
            for t in txs:
                sb = generateSignableBytes(t["tx"], args.ndautool)
                t["signable_bytes"] = sb
                if not sb.startswith("ERROR"):
                    t = tryToSign(t, args.keytool)
                newtxs.append(t)
            txs = newtxs

        print(json.dumps(txs, indent=2))
