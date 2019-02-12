#! /usr/bin/env python3
import json
import time
import argparse
import requests


def load(filename, host, delay, action, startAt):
    with open(filename) as f:
        data = json.load(f)
    print(
        f"Will {action} {len(data)} objects to {host} with a delay of {delay} seconds."
    )

    for counter, obj in enumerate(data):
        txtype = obj["txtype"]
        tx = obj["tx"]
        comment = obj["comment"]

        if counter < startAt:
            print(f"{counter}) Skipping {txtype} ({comment})")
            continue

        print(f"{counter}) Prevalidating {txtype} ({comment})")
        presult = requests.post(f"{host}/tx/prevalidate/{txtype}", json=tx)
        if presult.status_code == 200:
            print(f"     Prevalidate OK on {txtype} ({comment})")
            if action == "submit":
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
        else:
            try:
                contents = presult.json()
            except Exception:
                contents = presult.content
            print(
                f"     Prevalidate for {txtype} ({comment}) got {presult.status_code} "
                f"because {presult.reason}\n({contents})"
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
        const="https://node-0.main.ndau.tech",
        help="use mainnet as host",
    )
    parser.add_argument(
        "--test",
        action="store_const",
        dest="host",
        const="https://testnet-0.api.ndau.tech",
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
        "--submit",
        action="store_const",
        const="submit",
        default="prevalidate",
        help="submit each transaction after prevalidating it",
    )
    parser.add_argument(
        "--startAt",
        action="store",
        type=int,
        default=0,
        help="first (numbered) transaction to start with",
    )
    args = parser.parse_args()

    load(args.input, args.host, args.delay, args.submit, args.startAt)
