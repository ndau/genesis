#! /usr/bin/env python3

import sys
import base64
import getpass

from yubihsm import YubiHsm
from yubihsm.objects import AsymmetricKey

if len(sys.argv) != 4:
   print("Usage: sign.py <authorization key number> <datafile> <signing key number>")
   exit()

authkeyID = int(sys.argv[1])
keynum = int(sys.argv[3])

password = getpass.getpass()

infile = open(sys.argv[2], 'r')
encoded_data = infile.read()
infile.close()

encoded_bytes = base64.standard_b64decode(encoded_data)

hsm = YubiHsm.connect("http://localhost:12345/connector/api")
session = hsm.create_session_derived(authkeyID, password)
key = AsymmetricKey(session, keynum)
sig = key.sign_eddsa(encoded_bytes)

encoded_sig = base64.standard_b64encode(sig).decode()
print(encoded_sig)

session.close()
hsm.close()
exit(0)