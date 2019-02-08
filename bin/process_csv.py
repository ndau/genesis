#!/usr/bin/env python3

"""
Process csv for ETL script
"""

import sys # to print to stderr
import csv
import pymongo
from pymongo import MongoClient
import pdb

""" 
# if getting timeout from MongoDB, you might have to whitelist your IP for access:
- login to AWS, 
- go to EC2
- N Virginia east
- click on the 14 instances
- click on the one starting with Primary
- Click on the security groups below
- click anyone of them
- there will be a filter put in the above search
- take it out
- and put Mongo in there
- you should see 3 security groups
- click on each one of those and look at the Inbound tab below
- that is where you add the firewall rule
"""

def main():

    import argparse
    parser = argparse.ArgumentParser(description="Merges ndau spreadsheet with MongoDB data")
    parser.add_argument('-v', '--verbose', action='store_true',
                        help=('print verbose info for debugging'
                              f'Default: false'))
    parser.add_argument('-i', '--input', default="input.csv", 
                        help=('input .csv file, default: input.csv'))
    parser.add_argument('-o', '--output', default="output.csv",
                        help=('output .csv file, default output.csv'))

    args = parser.parse_args()

    # allow verbose printing
    global verboseFlag
    verboseFlag = args.verbose

    if verboseFlag:
        for p in sys.path:
            print(p)

    # list of node addresses to be assigned to accounts
    # these nodes come from the "data and procedures/genesis_nodes.txt" file in the genesis repo
    node_list = ['ndarw5i7rmqtqstw4mtnchmfvxnrq4k3e2ytsyvsc7nxt2y7',
        'ndaq3nqhez3vvxn8rx4m6s6n3kv7k9js8i3xw8hqnwvi2ete',
        'ndahnsxr8zh7r6u685ka865wz77wb78xcn45rgskpeyiwuza',
        'ndam75fnjn7cdues7ivi7ccfq8f534quieaccqibrvuzhqxa',
        'ndaekyty73hd56gynsswuj5q9em68tp6ed5v7tpft872hvuc']

    node_list_index = 0

    # connect to MongoDB
    client = MongoClient('mongodb://admin:0n13r0Nd3v@34.228.30.229:27017')
    if verboseFlag:
        print(f'client = {client}')
        print(f'db names = {client.list_database_names()}')
    db = client['ndau_dashboard']
    if verboseFlag:
        print(f'db = {db}')
        print(f'collection names = {db.list_collection_names()}')
    collection = db['accountaddresses']
    if verboseFlag:
        print(f'collection = {collection}')
    first = collection.find_one()
    if verboseFlag:
        print(f'item = {first}')
#    pdb.set_trace()
    r = csv.reader(open(args.input))
    lines = list(r)
    if verboseFlag:
        print(f"addresses = {first['addresses']}")
    # loop through MongoDB records
    for record in collection.find():
        if verboseFlag:
            print(f'record = {record}')
        addr_index = 0
        addrs = record['addresses']
        # loop through rows in spreadsheet
        for line in lines:
            # if userID in Mongo = userID in spreadsheet row, insert address into spreadsheet row
            if record['userId'] == line[8]:
                if addr_index == 0:
                    first_line = line
#                pdb.set_trace()
                # if we don't have enough addresses in Mongo for this userID, print a message
                if addr_index >= len(addrs):
                    print(f'addr mismatch, num in Mongo: {len(addrs)}, num in CSV: {addr_index}')
                else:
                    line[3] = addrs[addr_index]
                addr_index += 1
            # pdb.set_trace()
            # if this row has data in it, assign a node address to this account in round robin
            if line[13] != '':
                line[12] = node_list[node_list_index]
                node_list_index = (node_list_index + 1) % len(node_list)
        # if we have too many addresses in Mongo for this userID, print a message.  This happens
        # on occasion because early on we generated too many addresses for each userID
        if addr_index != len(addrs):
            print(f'addr mismatch, num in Mongo: {len(addrs)}, num in CSV: {addr_index}')
            first_line[11] = f'Mongo: {len(addrs)}'
    writer = csv.writer(open(args.output, 'w'))
    writer.writerows(lines)

    print('All done.')



# kick it off
if __name__ == '__main__':
    main()
