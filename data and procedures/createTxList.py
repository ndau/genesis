import csv
import json
import subprocess


def ClaimAccount(d):
    tx = dict(
        comment=d["header"],
        txtype="ClaimAccount",
        tx=dict(
            target=d["target"],
            ownership=d["ownership"],
            validation_script=d["validation_script"],
            sequence=d["sequence"],
            signature="",
        ),
    )
    keys = [d.get(f"validation_keys_{i}", "") for i in range(9)]
    keys = [k for k in keys if k != ""]
    tx["tx"]["validation_keys"] = keys
    return tx


def Issue(d):
    return dict(
        comment=d["header"],
        txtype="Issue",
        tx=dict(qty=d["qty"], sequence=d["sequence"], signatures=[""]),
    )


def Delegate(d):
    return dict(
        comment=d["header"],
        txtype="Delegate",
        tx=dict(
            target=d["source"],
            node=d["destination"],
            sequence=d["sequence"],
            signatures=[""],
        ),
    )


def CreditEAI(d):
    return dict(
        comment=d["header"],
        txtype="CreditEAI",
        tx=dict(node=d["target"], sequence=d["sequence"], signatures=[""]),
    )


def Lock(d):
    return dict(
        comment=d["header"],
        txtype="Lock",
        tx=dict(
            target=d["target"],
            period=d["period"],
            sequence=d["sequence"],
            signatures=[""],
        ),
    )


def SetRewardsDestination(d):
    return dict(
        comment=d["header"],
        txtype="SetRewardsDestination",
        tx=dict(
            source=d["source"],
            destination=d["destination"],
            sequence=d["sequence"],
            signatures=[""],
        ),
    )


def Transfer(d):
    return dict(
        comment=d["header"],
        txtype="Transfer",
        tx=dict(
            source=d["source"],
            destination=d["destination"],
            qty=d["qty"],
            sequence=d["sequence"],
            signatures=[""],
        ),
    )


def RegisterNode(d):
    return dict(
        comment=d["header"],
        txtype="RegisterNode",
        tx=dict(
            node=d["target"],
            distribution_script="",
            rpc_address="",
            sequence=d["sequence"],
            signatures=[""],
        ),
    )


def NominateNodeRewards(d):
    return dict(
        comment=d["header"],
        txtype="NominateNodeRewards",
        tx=dict(
            random=0, sequence=d["sequence"], signatures=[""]  # nominate the 0 node
        ),
    )


def ClaimReward(d):
    return dict(
        comment=d["header"],
        txtype="ClaimReward",
        tx=dict(node=d["target"], sequence=d["sequence"], signatures=[""]),
    )


if __name__ == "__main__":
    ndautool = "/Users/kentquirk/go/src/github.com/oneiro-ndev/commands/ndau"
    txmap = dict(
        ClaimAccount=[ClaimAccount],
        Issue=[Issue],
        Delegate=[Delegate],
        CreditEAI=[CreditEAI],
        RegisterNode=[RegisterNode],
        Lock=[Lock, SetRewardsDestination],
        NominateNodeRewards=[NominateNodeRewards],
        ClaimReward=[ClaimReward],
        Transfer=[Transfer],
    )
    with open("Post-Genesis Transaction Block.csv") as csvfile:
        rdr = csv.DictReader(csvfile)
        rows = [r for r in rdr if r["txtype"] != ""]
        txs = []
        for row in rows:
            txs.extend([tx(row) for tx in txmap[row["txtype"]]])
        newtxs = []
        for t in txs[10]:
            j = json.dumps(t["tx"])
            r = subprocess.run(
                ndautool, "signable-bytes", t["txtype"], input=j, text=True
            )
            sb = r.stdout
            t["signable_bytes"] = sb
            newtxs.append(t)
        print(json.dumps(newtxs, indent=2))
