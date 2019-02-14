# Genesis Key and Address Procedure and Logs

# YubiHSM Keypair Generation


1. Review the YubiHSM Quick Start Guide - https://developers.yubico.com/YubiHSM2/Usage_Guides/YubiHSM_quick_start_tutorial.html
2. Ensure the device has been reset to factory conditions by pressing the gold touch sensor down while inserting it into a USB port and holding it down for 10 seconds.
3. Install the YubiHSM software from https://developers.yubico.com/YubiHSM2/Releases/ and verify its signatures. Make that installation directory the current shell directory.
4. Start a command shell and launch the YubiHSM connector with the command
          ./yubihsm-connector -d
5. Check the connector status in a browser at http://127.0.0.1:12345/connector/status
6. Start the YubiHSM shell in another command shell, connect to the YubiHSM key, increase the keepalive timeout for convenience. Record the device information.
    ./yubihsm-shell
    yubihsm> connect
    yubihsm> keepalive 30
    yubihsm> get deviceinfo
7. Open a session using the device default authorization key and record device information.
    yubihsm> session open 1 password
8. Generate a replacement administrative AuthKey
    yubihsm> put authkey 0 2 "ndau Admin" 1 0x00003fffffffffff 0x00003fffffffffff <<new admin password>>
9. Close the current session, delete the factory default AuthKey, open a new session with the new administrative AuthKey
    yubihsm> delete 0 1 authkey
    yubihsm> session close 0
    yubihsm> session open 2 <<new admin password>>
10. Generate an operational AuthKey for key generation and signing
    yubihsm> put authkey 0 3 "ndau Key Generation" 1 asymmetric_gen,delete_asymmetric,asymmetric_sign_eddsa,asymmetric_sign_ecdsa asymmetric_sign_eddsa,asymmetric_sign_ecdsa <<new key generation password>>
11. Close the current session and open a new session with the ndau Key Generation AuthKey.
    yubihsm> session close 0
    yubihsm> session open 3 <<new key generation password>>
12. Generate a list of keypairs for initial special address creation.
    yubihsm> generate asymmetric 0 101 genesis_e_key1 1 asymmetric_sign_eddsa ed25519
    yubihsm> generate asymmetric 0 102 genesis_b_key1 1 asymmetric_sign_eddsa ed25519
    yubihsm> generate asymmetric 0 103 genesis_m_key1 1 asymmetric_sign_eddsa ed25519
    yubihsm> generate asymmetric 0 104 genesis_m_key2 1 asymmetric_sign_eddsa ed25519
    yubihsm> generate asymmetric 0 105 genesis_a_key1 1 asymmetric_sign_eddsa ed25519
    yubihsm> generate asymmetric 0 106 genesis_n_key1 1 asymmetric_sign_eddsa ed25519
    yubihsm> generate asymmetric 0 107 genesis_n_key2 1 asymmetric_sign_eddsa ed25519
13. Retrieve their public keys.
    yubihsm> get pubkey 0 101
    yubihsm> get pubkey 0 102
    yubihsm> get pubkey 0 103
    yubihsm> get pubkey 0 104
    yubihsm> get pubkey 0 105
    yubihsm> get pubkey 0 106
    yubihsm> get pubkey 0 107
14. Close the session and retrieve the device log.
    yubihsm> session close 0
    yubihsm> session open 2 cosimo
    yubihsm> audit get 0
    yubihsm> session close 0
    yubihsm> quit
15. Store the YubiHSM safely. The private keys generated in step 12 are the ownership keys used to sign the ClaimAccount transactions for these addresses. They should be treated with ordinary security but are immediately obsolete once those ClaimAccount transactions are submitted. Those transactions are prepared and signed prior to genesis and submitted as part of the genesis procedure. They are the first transactions on the blockchain and there is no risk of compromise.


# Initial Genesis Setup Results

Session Log:


    $ yubi/yubihsm-shell 
    Using default connector URL: http://127.0.0.1:12345
    yubihsm> connect
    keepalive 30
    get deviceinfo
    session open 1 password
    put authkey 0 2 "ndau Admin" 1 0x00003fffffffffff 0x00003fffffffffff cosimo
    delete 0 1 authkey
    session close 0
    session open 2 cosimo
    put authkey 0 3 "ndau Key Generation" 1 asymmetric_gen,delete_asymmetric,asymmetric_sign_eddsa,asymmetric_sign_ecdsa asymmetric_sign_eddsa,asymmetric_sign_ecdsa medici
    session close 0
    session open 3 medici
    generate asymmetric 0 101 genesis_e_key1 1 asymmetric_sign_eddsa ed25519
    generate asymmetric 0 102 genesis_b_key1 1 asymmetric_sign_eddsa ed25519
    generate asymmetric 0 103 genesis_m_key1 1 asymmetric_sign_eddsa ed25519
    generate asymmetric 0 104 genesis_m_key2 1 asymmetric_sign_eddsa ed25519
    generate asymmetric 0 105 genesis_a_key1 1 asymmetric_sign_eddsa ed25519
    generate asymmetric 0 106 genesis_n_key1 1 asymmetric_sign_eddsa ed25519
    generate asymmetric 0 107 genesis_n_key2 1 asymmetric_sign_eddsa ed25519
    get pubkey 0 101
    get pubkey 0 102
    get pubkey 0 103
    get pubkey 0 104
    get pubkey 0 105
    get pubkey 0 10Session keepalive set up to run every 15 seconds
    yubihsm> keepalive 30
    Session keepalive set up to run every 30 seconds
    yubihsm> get deviceinfo
    Version number:                2.0.0
    Serial number:                7550421
    Log used:                2/62
    Supported algorithms:        rsa-pkcs1-sha1, rsa-pkcs1-sha256, rsa-pkcs1-sha384, 
                            rsa-pkcs1-sha512, rsa-pss-sha1, rsa-pss-sha256, 
                            rsa-pss-sha384, rsa-pss-sha512, rsa2048, 
                            rsa3072, rsa4096, ecp256, 
                            ecp384, ecp521, eck256, 
                            ecbp256, ecbp384, ecbp512, 
                            hmac-sha1, hmac-sha256, hmac-sha384, 
                            hmac-sha512, ecdsa-sha1, ecdh, 
                            rsa-oaep-sha1, rsa-oaep-sha256, rsa-oaep-sha384, 
                            rsa-oaep-sha512, aes128-ccm-wrap, opaque, 
                            x509-cert, mgf1-sha1, mgf1-sha256, 
                            mgf1-sha384, mgf1-sha512, template-ssh, 
                            yubico-otp-aes128, yubico-aes-auth, yubico-otp-aes192, 
                            yubico-otp-aes256, aes192-ccm-wrap, aes256-ccm-wrap, 
                            ecdsa-sha256, ecdsa-sha384, ecdsa-sha512, 
                            ed25519, ecp224, 
    yubihsm> session open 1 password
    Created session 0
    yubihsm> put authkey 0 2 "ndau Admin" 1 0x00003fffffffffff 0x00003fffffffffff cosimo
    Stored Authentication key 0x0002
    yubihsm> delete 0 1 authkey
    yubihsm> session close 0
    yubihsm> session open 2 cosimo
    Created session 0
    yubihsm> put authkey 0 3 "ndau Key Generation" 1 asymmetric_gen,delete_asymmetric,asymmetric_sign_eddsa,asymmetric_sign_ecdsa asymmetric_sign_eddsa,asymmetric_sign_ecdsa medici
    Stored Authentication key 0x0003
    yubihsm> session close 0
    yubihsm> session open 3 medici
    Created session 0
    yubihsm> generate asymmetric 0 101 genesis_e_key1 1 asymmetric_sign_eddsa ed25519
    Generated Asymmetric key 0x0065
    yubihsm> generate asymmetric 0 102 genesis_b_key1 1 asymmetric_sign_eddsa ed25519
    Generated Asymmetric key 0x0066
    yubihsm> generate asymmetric 0 103 genesis_m_key1 1 asymmetric_sign_eddsa ed25519
    Generated Asymmetric key 0x0067
    yubihsm> generate asymmetric 0 104 genesis_m_key2 1 asymmetric_sign_eddsa ed25519
    Generated Asymmetric key 0x0068
    yubihsm> generate asymmetric 0 105 genesis_a_key1 1 asymmetric_sign_eddsa ed25519
    Generated Asymmetric key 0x0069
    yubihsm> generate asymmetric 0 106 genesis_n_key1 1 asymmetric_sign_eddsa ed25519
    Generated Asymmetric key 0x006a
    yubihsm> generate asymmetric 0 107 genesis_n_key2 1 asymmetric_sign_eddsa ed25519
    Generated Asymmetric key 0x006b
    yubihsm> get pubkey 0 101
    TCFyMCB6inwknLHre8T1itpFKnukr9tMTmWCdXS/1yc=
    yubihsm> get pubkey 0 102
    k+Wk8JwFpH6NrxDuPSrP+XTJy+W+CmQsi0kwtfNO2NQ=
    yubihsm> get pubkey 0 103
    aJsxh/RK4ZjF4vOGMyYKRCU2De8Ax5j112ZSTuUXVpU=
    yubihsm> get pubkey 0 104
    p2J9uY+EjivAQkDn9SnuvGjXYbyHBicrZd6ysjvsAXs=
    yubihsm> get pubkey 0 105
    flYi6sFqKbVAic0NBnMfjW5Cl3aii3Tv6elbogkDlGc=
    yubihsm> get pubkey 0 106
    QgHzLbiZKPMEw4W3i+1HYYJYfL2MHdd1IcG/+TtZjfs=
    yubihsm> get pubkey 0 107
    HjXUob90sUyd3gY0dvO+RLo18APmxjdlbly55fjE+bQ=
    yubihsm> session close 0
    yubihsm> session open 2 cosimo 
    Created session 0
    yubihsm> audit get 0
    0 unlogged boots found
    0 unlogged authentications found
    Found 30 items
    item:     1 -- cmd: 0xff -- length: 65535 -- session key: 0xffff -- target key: 0xffff -- second key: 0xffff -- result: 0xff -- tick: 4294967295 -- hash: 6830a758581287b99c4de5172126daf6
    item:     2 -- cmd: 0x00 -- length:    0 -- session key: 0xffff -- target key: 0x0000 -- second key: 0x0000 -- result: 0x00 -- tick: 0 -- hash: 838cc6336033da7208ef696638ba7dd9
    item:     3 -- cmd: 0x03 -- length:   10 -- session key: 0xffff -- target key: 0x0001 -- second key: 0xffff -- result: 0x83 -- tick: 547 -- hash: 7c81f977ea698b6d5b50d59332323308
    item:     4 -- cmd: 0x04 -- length:   17 -- session key: 0xffff -- target key: 0x0001 -- second key: 0xffff -- result: 0x84 -- tick: 547 -- hash: 997ad277645d08647f35f11b6f6c8838
    item:     5 -- cmd: 0x44 -- length:   93 -- session key: 0x0001 -- target key: 0x0002 -- second key: 0xffff -- result: 0xc4 -- tick: 548 -- hash: a8a022aa6fb2fd8d533217180c6b6437
    item:     6 -- cmd: 0x58 -- length:    3 -- session key: 0x0001 -- target key: 0x0001 -- second key: 0xffff -- result: 0xd8 -- tick: 549 -- hash: 3bbad64e3eb4a396a7515db9c94bbb99
    item:     7 -- cmd: 0x40 -- length:    0 -- session key: 0x0001 -- target key: 0xffff -- second key: 0xffff -- result: 0xc0 -- tick: 551 -- hash: f1e654e9071d97cc0dd6b5769723b240
    item:     8 -- cmd: 0x03 -- length:   10 -- session key: 0xffff -- target key: 0x0002 -- second key: 0xffff -- result: 0x83 -- tick: 551 -- hash: 258daeb76f7826722e2df09b234f5d7e
    item:     9 -- cmd: 0x04 -- length:   17 -- session key: 0xffff -- target key: 0x0002 -- second key: 0xffff -- result: 0x84 -- tick: 552 -- hash: 22519af11e1477848246be01c49d116b
    item:    10 -- cmd: 0x44 -- length:   93 -- session key: 0x0002 -- target key: 0x0003 -- second key: 0xffff -- result: 0xc4 -- tick: 553 -- hash: e27358411e555bc456345718cc5dd5db
    item:    11 -- cmd: 0x40 -- length:    0 -- session key: 0x0002 -- target key: 0xffff -- second key: 0xffff -- result: 0xc0 -- tick: 554 -- hash: 56a85509ea2ee4d90e4a5ff7c8cc2183
    item:    12 -- cmd: 0x03 -- length:   10 -- session key: 0xffff -- target key: 0x0003 -- second key: 0xffff -- result: 0x83 -- tick: 554 -- hash: f5f484d0d2f0209ea23df72ba3f2e676
    item:    13 -- cmd: 0x04 -- length:   17 -- session key: 0xffff -- target key: 0x0003 -- second key: 0xffff -- result: 0x84 -- tick: 555 -- hash: b1ea306244999dbe2cdf89e0f86f7906
    item:    14 -- cmd: 0x46 -- length:   53 -- session key: 0x0003 -- target key: 0x0065 -- second key: 0xffff -- result: 0xc6 -- tick: 556 -- hash: 4c644538e64a397ca5e41b4201564070
    item:    15 -- cmd: 0x46 -- length:   53 -- session key: 0x0003 -- target key: 0x0066 -- second key: 0xffff -- result: 0xc6 -- tick: 560 -- hash: cba8f069c619cbc5f955bf3bbc406bf5
    item:    16 -- cmd: 0x46 -- length:   53 -- session key: 0x0003 -- target key: 0x0067 -- second key: 0xffff -- result: 0xc6 -- tick: 564 -- hash: 15f6326398143bf0f4b41fb89b701cec
    item:    17 -- cmd: 0x46 -- length:   53 -- session key: 0x0003 -- target key: 0x0068 -- second key: 0xffff -- result: 0xc6 -- tick: 569 -- hash: de405bf2c73030e94425b169ea9042e6
    item:    18 -- cmd: 0x46 -- length:   53 -- session key: 0x0003 -- target key: 0x0069 -- second key: 0xffff -- result: 0xc6 -- tick: 573 -- hash: 8876c0fd1b088e24371c5c479408f776
    item:    19 -- cmd: 0x46 -- length:   53 -- session key: 0x0003 -- target key: 0x006a -- second key: 0xffff -- result: 0xc6 -- tick: 578 -- hash: e3cb5553ba3d53c5024f92bea10bb5f2
    item:    20 -- cmd: 0x46 -- length:   53 -- session key: 0x0003 -- target key: 0x006b -- second key: 0xffff -- result: 0xc6 -- tick: 582 -- hash: 4b4af2b8d9def2088cf1dd08565559e5
    item:    21 -- cmd: 0x54 -- length:    2 -- session key: 0x0003 -- target key: 0x0065 -- second key: 0xffff -- result: 0xd4 -- tick: 586 -- hash: 29840766129c512eb373e3d2b9977456
    item:    22 -- cmd: 0x54 -- length:    2 -- session key: 0x0003 -- target key: 0x0066 -- second key: 0xffff -- result: 0xd4 -- tick: 587 -- hash: 12c0acebd7e06be02ae37f75c1ec0152
    item:    23 -- cmd: 0x54 -- length:    2 -- session key: 0x0003 -- target key: 0x0067 -- second key: 0xffff -- result: 0xd4 -- tick: 587 -- hash: f4010733f0dad7ca6dd4fd81d2c9e600
    item:    24 -- cmd: 0x54 -- length:    2 -- session key: 0x0003 -- target key: 0x0068 -- second key: 0xffff -- result: 0xd4 -- tick: 588 -- hash: b83a38c59e75ec90d91a020232e69954
    item:    25 -- cmd: 0x54 -- length:    2 -- session key: 0x0003 -- target key: 0x0069 -- second key: 0xffff -- result: 0xd4 -- tick: 589 -- hash: fb09813f3be519f38f85f4b30fc9a65b
    item:    26 -- cmd: 0x54 -- length:    2 -- session key: 0x0003 -- target key: 0x006a -- second key: 0xffff -- result: 0xd4 -- tick: 589 -- hash: 7e89a3000d1af911c3545d27f5da1310
    item:    27 -- cmd: 0x54 -- length:    2 -- session key: 0x0003 -- target key: 0x006b -- second key: 0xffff -- result: 0xd4 -- tick: 590 -- hash: bc794a1fb88a1ede8e95df5cd8bd05d5
    item:    28 -- cmd: 0x40 -- length:    0 -- session key: 0x0003 -- target key: 0xffff -- second key: 0xffff -- result: 0xc0 -- tick: 590 -- hash: 08eb97d242067be5b8308c462c931085
    item:    29 -- cmd: 0x03 -- length:   10 -- session key: 0xffff -- target key: 0x0002 -- second key: 0xffff -- result: 0x83 -- tick: 591 -- hash: 25307b11826ec5ba0bc831384646cc30
    item:    30 -- cmd: 0x04 -- length:   17 -- session key: 0xffff -- target key: 0x0002 -- second key: 0xffff -- result: 0x84 -- tick: 592 -- hash: b077b25f6e432ea79dd5dc2fc9a0c6c5
    yubihsm> session close 0
    yubihsm> quit

Convert the public keys retrieved above into ndau format:


    $ keytool ed raw public TCFyMCB6inwknLHre8T1itpFKnukr9tMTmWCdXS/1yc= -b
    npuba8jadtbbebgcc6tseb7iw9bevu28y88e8yfpwtjkrquk9y4nj3u2e7mwz9muqgcnmi4yuyhp
    $ keytool ed raw public k+Wk8JwFpH6NrxDuPSrP+XTJy+W+CmQsi0kwtfNO2NQ= -b
    npuba8jadtbbecj8mjhsvsc4i9wpx6iq6rjk396zjuqm6y9aw3bntpevbprvj5npjjcn53b3uncz
    $ keytool ed raw public aJsxh/RK4ZjF4vOGMyYKRCU2De8Ax5j112ZSTuUXVpU= -b
    npuba8jadtbbebwjynnh8tfqdggf6m32nn3gbjcckpsp76anrghx47vfevzfc7mjk3f8ftw5arph
    $ keytool ed raw public p2J9uY+EjivAQkDn9SnuvGjXYbyHBicrZd6ysjvsAXs= -b
    npuba8jadtbbecvye9p3t8ci6k8aijaqr7jj748gtx5bzudsnj3mnzrmfnt57sazyhxryvwkwgs3
    $ keytool ed raw public flYi6sFqKbVAic0NBnMfjW5Cl3aii3Tv6elbogkDlGc= -b
    npuba8jadtbbeb9fnizk2fxcvpkathgs4bvvd8gy6swzq4tiy7hr7hwxzisjaqkgqi9zgpzc2iuw
    $ keytool ed raw public QgHzLbiZKPMEw4W3i+1HYYJYfL2MHdd1IcG/+TtZjfs= -b
    npuba8jadtbbebbad63pzcnut62e2qc5rc9pi7s2eyd6zygb5x5xeha598j5mgg9yeb3mx6zrwhj
    $ keytool ed raw public HjXUob90sUyd3gY0dvO+RLo18APmxjdlbly55fjE+bQ= -b
    npuba8jadtbbeardmxfbz74mcve752ddi7zvz3cmwprsarvnnp5fp3qmv3r22v65isbn39cem4av

Generate ndau addresses of the appropriate types from those public keys:


    $ keytool addr npuba8jadtbbebgcc6tseb7iw9bevu28y88e8yfpwtjkrquk9y4nj3u2e7mwz9muqgcnmi4yuyhp -e
    ndec3xgijh2khcjywtqdhp35mt77vwpucwg24zpxw6ekvhhx
    $ keytool addr npuba8jadtbbecj8mjhsvsc4i9wpx6iq6rjk396zjuqm6y9aw3bntpevbprvj5npjjcn53b3uncz -b
    ndbfj2q5udrti57x9gfix56puphk53jqcvu3uib8n8d8qy7a
    $ keytool addr npuba8jadtbbebwjynnh8tfqdggf6m32nn3gbjcckpsp76anrghx47vfevzfc7mjk3f8ftw5arph -m
    ndmp57vwxn7hqzg2xffactfx5q4kh9u4mc8med5ta268w7z4
    $ keytool addr npuba8jadtbbecvye9p3t8ci6k8aijaqr7jj748gtx5bzudsnj3mnzrmfnt57sazyhxryvwkwgs3 -m
    ndmqqjg2dmr7nv9gqqykxd3g4evsyfssfx5ivj6isthsnrtt
    $ keytool addr npuba8jadtbbeb9fnizk2fxcvpkathgs4bvvd8gy6swzq4tiy7hr7hwxzisjaqkgqi9zgpzc2iuw -a
    ndam32av92ekrqvjjzwkct4peri5c9wxyg4k9p34a3xtef3x
    $ keytool addr npuba8jadtbbebbad63pzcnut62e2qc5rc9pi7s2eyd6zygb5x5xeha598j5mgg9yeb3mx6zrwhj -n
    ndnhhr5kg3ku9vtaph7pwkz3mxgvizc565ec37prbkzja7ex
    $ keytool addr npuba8jadtbbeardmxfbz74mcve752ddi7zvz3cmwprsarvnnp5fp3qmv3r22v65isbn39cem4av -n
    ndngmp2f34wibp7w8y8rgfhtaeycfaz58g74djpzvs76z5s7

