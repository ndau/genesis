import csv
import json


def ClaimAccount(d):
    return dict(
        txtype="ClaimAccount",
        tx=dict(
            target=d["target"],
            ownership=d["ownership"],
            validation_keys=[d["validation_keys"]],
            validation_script="",
            sequence=d["sequence"],
            signature="",
        ),
    )


def Issue(d):
    return dict(
        txtype="Issue", tx=dict(qty=d["qty"], sequence=d["sequence"], signatures=[""])
    )


def Delegate(d):
    return dict(
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
        txtype="CreditEAI",
        tx=dict(node=d["target"], sequence=d["sequence"], signatures=[""]),
    )


def Lock(d):
    return dict(
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
        txtype="NominateNodeRewards",
        tx=dict(
            random=0, sequence=d["sequence"], signatures=[""]  # nominate the 0 node
        ),
    )


def ClaimReward(d):
    return dict(
        txtype="ClaimReward",
        tx=dict(node=d["target"], sequence=d["sequence"], signatures=[""]),
    )


if __name__ == "__main__":
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
        print(json.dumps(txs, indent=2))
