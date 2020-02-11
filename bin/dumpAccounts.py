#  ----- ---- --- -- -
#  Copyright 2020 The Axiom Foundation. All Rights Reserved.
# 
#  Licensed under the Apache License 2.0 (the "License").  You may not use
#  this file except in compliance with the License.  You can obtain a copy
#  in the file LICENSE in the source distribution or at
#  https://www.apache.org/licenses/LICENSE-2.0.txt
#  - -- --- ---- -----


#! /usr/bin/env python3

import requests
import time
import sys


def getData(base, path, parms=None):
    u = base + path
    try:
        r = requests.get(u, timeout=3, params=parms)
    except requests.Timeout:
        print(f"{time.asctime()}: Timeout in {u}")
        return {}
    except Exception as e:
        print(f"{time.asctime()}: Error {e} in {u}")
        return {}
    if r.status_code == requests.codes.ok:
        return r.json()
    print(f"{time.asctime()}: Error in {u}: ({r.status_code}) {r}")
    return {}


names = {
    "local": "http://localhost:3030",
    "main": "https://node-0.main.ndau.tech",
    "mainnet": "https://node-0.main.ndau.tech",
    "dev": "https://devnet-0.api.ndau.tech",
    "devnet": "https://devnet-0.api.ndau.tech",
    "test": "https://testnet-0.api.ndau.tech",
    "testnet": "https://testnet-0.api.ndau.tech",
}

if __name__ == "__main__":
    name = "dev"
    if len(sys.argv) > 1:
        name = sys.argv[1]

    node = names[name]

    page = 0
    pgsz = 100
    balances = []
    while True:
        qp = dict(pagesize=pgsz, pageindex=page)
        result = getData(node, "/account/list", parms=qp)
        if not result["Accounts"]:
            break
        page += 1

        accts = result["Accounts"]
        resp = requests.post(f"{node}/account/accounts", json=result["Accounts"])

        data = resp.json()
        for k in data:
            print (data[k])
