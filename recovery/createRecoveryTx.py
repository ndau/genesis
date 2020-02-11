#  ----- ---- --- -- -
#  Copyright 2020 The Axiom Foundation. All Rights Reserved.
# 
#  Licensed under the Apache License 2.0 (the "License").  You may not use
#  this file except in compliance with the License.  You can obtain a copy
#  in the file LICENSE in the source distribution or at
#  https://www.apache.org/licenses/LICENSE-2.0.txt
#  - -- --- ---- -----


#! /usr/bin/env python3

import csv
import copy
import json
import argparse
import subprocess

def getPvtKeys(d):
    pvtkeys = [d.get(f"pvt_key{i}", "") for i in range(9)]
    pvtkeys = [k for k in pvtkeys if k != ""]
    return pvtkeys

def SetValidation(d):
    tx = dict(
        comment=d["header"],
        txtype="SetValidation",
        tx=dict(
            target=d["target"],
            ownership=d["ownership"],
            validation_script=d["validation_script"],
            sequence=int(d["sequence"]),
            pvt_keys=getPvtKeys(d),
            signatures=[],
        ),
    )
    keys = [d.get(f"validation_keys_{i}", "") for i in range(9)]
    keys = [k for k in keys if k != ""]
    tx["tx"]["validation_keys"] = keys

    return [tx]

def ChangeValidation(d):
    tx = dict(
        comment=d["header"],
        txtype="ChangeValidation",
        tx=dict(
            target=d["target"],
            validation_script=d["validation_script"],
            sequence=int(d["sequence"]),
            pvt_keys=getPvtKeys(d),
            signatures=[],
        ),
    )
    keys = [d.get(f"new_keys_{i}", "") for i in range(9)]
    keys = [k for k in keys if k != ""]
    tx["tx"]["new_keys"] = keys

    return [tx]

def ReleaseFromEndowment(d):
    tx = dict(
        comment=d["header"],
        txtype="ReleaseFromEndowment",
        tx=dict(
            destination=d["destination"],
            qty=int(d["qty"]),
#            qty=int(float(d["qty"]) * 100_000_000),
            sequence=int(d["sequence"]),
            pvt_keys=getPvtKeys(d),
            signatures=[],
        ),
    )
    return [tx]
  
def Issue(d):
    tx = dict(
        comment=d["header"],
        txtype="Issue",
        tx=dict(
            qty=int(float(d["qty"]) * 100_000_000),
            sequence=int(d["sequence"]),
            pvt_keys=getPvtKeys(d),
            signatures=[],
        ),
    )
    return [tx]

def RecordPrice(d):
    tx = dict(
        comment=d["header"],
        txtype="RecordPrice",
        tx=dict(
            market_price=int(d["market_price"]),
#            market_price=int(float(d["market_price"]) * 100_000_000_000),
            sequence=int(d["sequence"]),
            pvt_keys=getPvtKeys(d),
            signatures=[],
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
            pvt_keys=getPvtKeys(d),
            signatures=[],
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
            pvt_keys=getPvtKeys(d),
            signatures=[],
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
            pvt_keys=getPvtKeys(d),
            signatures=[],
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
            pvt_keys=getPvtKeys(d),
            signatures=[],
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
            qty=int(float(d["qty"]) * 100_000_000),
            sequence=int(d["sequence"]),
            pvt_keys=getPvtKeys(d),
            signatures=[],
        ),
    )
    return [tx]

def TransferAndLock(d):
    tx = dict(
        comment=d["header"],
        txtype="TransferAndLock",
        tx=dict(
            source=d["source"],
            destination=d["destination"],
            qty=int(float(d["qty"]) * 100_000_000),
            period=d["period"],
            sequence=int(d["sequence"]),
            pvt_keys=getPvtKeys(d),
            signatures=[],
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
            pvt_keys=getPvtKeys(d),
            signatures=[],
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
            pvt_keys=getPvtKeys(d),
            signatures=[],  # nominate the 0 node
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
            pvt_keys=getPvtKeys(d),
            signatures=[],
        ),
    )
    return [tx]

def SetSysvar(d):
    tx = dict(
        comment=d["header"],
        txtype="SetSysvar",
        tx=dict(
            name=d["sysvar_name"],
            value=d["sysvar_value"],
            sequence=int(d["sequence"]),
            pvt_keys=getPvtKeys(d),
            signatures=[],
        ),
    )
    return [tx]

def generateSignableBytes(obj, ndautool):
    tx = copy.deepcopy(obj)
    if "signatures" in tx:
        del tx["signatures"]
    if "pvt_keys" in tx:
        del tx["pvt_keys"]

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
    pks = t["tx"].get("pvt_keys", None)
    tx = copy.deepcopy(t["tx"])
    if pks is not None:
        del tx["pvt_keys"]
    if not pks:
        t["tx"] = tx
        return t

    sb = t["signable_bytes"]

    for pk in pks:
        if not pk.startswith("npvt"):
            if len(tx["signatures"]) == 0:

                args = ["/usr/local/bin/yubihsm-shell", "-p", "lvh2$*BmIi*A2Mm3qmLL", "-a", "sign-eddsa", "-i", pks[0], "-A", "ed25519", "--informat", "base64", "--authkey", "101"]
                r = subprocess.run(args, input=sb, text=True, capture_output=True)

                if r.returncode > 0:
                    sig = f"ERROR: {r.stderr}"
                else:
                    sig = r.stdout.strip()
                    args = [keytool, "ed", "raw", "signature", sig, "-b"]
                    r = subprocess.run(args, text=True, capture_output=True)
                    if r.returncode > 0:
                        sig = f"ERROR: {r.stderr}"
                    else:
                        sig = r.stdout.strip()
                
                tx["signatures"].append(sig)

        else:
            args = [keytool, "sign", pk, sb, "-b"]
            r = subprocess.run(args, text=True, capture_output=True)
            if r.returncode > 0:
                sig = f"ERROR: {r.stderr}"
            else:
                sig = r.stdout.strip()
            tx["signatures"].append(sig)

 #       tx["signature"] = tx["signatures"][0]
 #       del tx["signatures"]

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
        default="../../commands/keytool",
        help="location of keytool",
    )
    parser.add_argument(
        "--ndautool",
        action="store",
        default="../../commands/ndau",
        help="location of ndautool",
    )
    parser.add_argument(
        "--input",
        action="store",
        default="",
        help="input csv file",
    )

    args = parser.parse_args()

    txmap = dict(
        SetValidation=SetValidation,
        ChangeValidation=ChangeValidation,
        ReleaseFromEndowment=ReleaseFromEndowment,
        Issue=Issue,
        Delegate=Delegate,
        CreditEAI=CreditEAI,
        RegisterNode=RegisterNode,
        Lock=Lock,
        SetRewardsDestination=SetRewardsDestination,
        NominateNodeReward=NominateNodeReward,
        ClaimNodeReward=ClaimNodeReward,
        Transfer=Transfer,
        TransferAndLock=TransferAndLock,
        RecordPrice=RecordPrice,
        SetSysvar=SetSysvar,
    )

    with open(args.input) as csvfile:
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
