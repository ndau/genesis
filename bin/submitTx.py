#  ----- ---- --- -- -
#  Copyright 2020 The Axiom Foundation. All Rights Reserved.
# 
#  Licensed under the Apache License 2.0 (the "License").  You may not use
#  this file except in compliance with the License.  You can obtain a copy
#  in the file LICENSE in the source distribution or at
#  https://www.apache.org/licenses/LICENSE-2.0.txt
#  - -- --- ---- -----


#! /usr/bin/env python3
import json
import time
import argparse
import requests


def load(filename, host, delay, action, startAt):
    f = open(filename)
    data = json.load(f)
    print(
        f"Will {action} {len(data)} objects to {host} with a delay of {delay} seconds."
    )

    counter = 0
    for obj in data:
        counter += 1
        txtype = obj["txtype"]
        tx = obj["tx"]
        comment = obj["comment"]

        if counter < startAt:
            print(f"{counter}) Skipping {txtype} ({comment})")
            continue

        print(action)
        
        if action == "prevalidate" or action == "both":
            print(f"{counter}) Prevalidating {txtype} ({comment})")
            presult = requests.post(f"{host}/tx/prevalidate/{txtype}", json=tx)
            if presult.status_code == 200:
                print(f"     Prevalidate OK on {txtype} ({comment})")
            else:
                print(
                    f"     Prevalidate for {txtype} ({comment}) got {presult.status_code} "
                    f"because {presult.reason}\n({presult.content})"
                )
                print(json.dumps(tx, indent=2))

        if action == "submit" or action == "both":
            print("Submitting.")
            time.sleep(delay)
            sresult = requests.post(f"{host}/tx/submit/{txtype}", json=tx)
            if sresult.status_code == 200:
                print("     Transaction succeeded.")
            else:
                print(
                    f"     Submit for {txtype} ({comment}) "
                    f"got {sresult.status_code} "
                    f"because {sresult.reason}\n({sresult.content})"
                )
                print(json.dumps(tx, indent=2))

        time.sleep(delay)

    print(f"finished")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--host",
        action="store",
        default="http://localhost:3030",
        dest="host",
        help="specify an arbitrary host location",
    )
    parser.add_argument(
        "--main",
        action="store_const",
        dest="host",
        const="https://mainnet-0.ndau.tech:3030",
        help="use mainnet as host",
    )
    parser.add_argument(
        "--staging",
        action="store_const",
        dest="host",
        const="https://api.ndau.tech:32300",
        help="use mainnet staging net as host",
    )
    parser.add_argument(
        "--test",
        action="store_const",
        dest="host",
        const="https://api.ndau.tech:31300",
        help="use testnet as host",
    )
    parser.add_argument(
        "--dev",
        action="store_const",
        dest="host",
        const="https://devnet-0.api.ndau.tech",
        help="use devnet as host",
    )
    parser.add_argument(
        "--local",
        action="store_const",
        dest="host",
        const="http://localhost:3030",
        help="use localhost:3030 as host",
    )
    parser.add_argument(
        "--input",
        action="store",
        default="genesis_tx_list.json",
        help="specify the json file to read transactions from",
    )
    parser.add_argument(
        "--delay",
        action="store",
        type=int,
        default=2,
        help="the amount of time (in seconds) to wait between submissions",
    )
    parser.add_argument(
        "--action",
        action="store",
        default="prevalidate",
        help="[submit|prevalidate|both]",
    )
    parser.add_argument(
        "--startAt",
        action="store",
        type=int,
        default=0,
        help="first (numbered) transaction to start with",
    )
    args = parser.parse_args()

    load(args.input, args.host, args.delay, args.action, args.startAt)
