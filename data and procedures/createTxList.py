import sys
import csv
import copy
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
            sequence=int(d["sequence"]),
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
        tx=dict(
            qty=int(d["qty"]) * 100_000_000,
            sequence=int(d["sequence"]),
            signatures=[""],
        ),
    )


def Delegate(d):
    return dict(
        comment=d["header"],
        txtype="Delegate",
        tx=dict(
            target=d["target"],
            node=d["ownership"],
            sequence=int(d["sequence"]),
            signatures=[""],
        ),
    )


def CreditEAI(d):
    return dict(
        comment=d["header"],
        txtype="CreditEAI",
        tx=dict(node=d["target"], sequence=int(d["sequence"]), signatures=[""]),
    )


def Lock(d):
    return dict(
        comment=d["header"],
        txtype="Lock",
        tx=dict(
            target=d["target"],
            period=d["period"],
            sequence=int(d["sequence"]),
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
            sequence=int(d["sequence"]),
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
            qty=int(d["qty"]) * 100_000_000,
            sequence=int(d["sequence"]),
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
            sequence=int(d["sequence"]),
            signatures=[""],
        ),
    )


def NominateNodeReward(d):
    return dict(
        comment=d["header"],
        txtype="NominateNodeReward",
        tx=dict(
            random=0,
            sequence=int(d["sequence"]),
            signatures=[""],  # nominate the 0 node
        ),
    )


def ClaimNodeReward(d):
    return dict(
        comment=d["header"],
        txtype="ClaimNodeReward",
        tx=dict(node=d["target"], sequence=int(d["sequence"]), signatures=[""]),
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
        NominateNodeReward=[NominateNodeReward],
        ClaimNodeReward=[ClaimNodeReward],
        Transfer=[Transfer],
    )
    with open("Post-Genesis Transaction Block.csv") as csvfile:
        rdr = csv.DictReader(csvfile)
        rows = [r for r in rdr if r["txtype"] != ""]
        txs = []
        for row in rows:
            txs.extend([tx(row) for tx in txmap[row["txtype"]]])
        newtxs = []

        if len(sys.argv) > 1:
            for t in txs:
                tx = copy.deepcopy(t["tx"])
                if "signature" in tx:
                    del tx["signature"]
                if "signatures" in tx:
                    del tx["signatures"]

                j = json.dumps(tx, indent=2)
                r = subprocess.run(
                    [ndautool, "signable-bytes", t["txtype"]],
                    input=j,
                    text=True,
                    capture_output=True,
                )
                if r.returncode > 0:
                    t["signable_bytes"] = f"ERROR: {r.stderr}"
                else:
                    t["signable_bytes"] = r.stdout
                newtxs.append(t)
            txs = newtxs

        print(json.dumps(txs, indent=2))
