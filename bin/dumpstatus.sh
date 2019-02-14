#! /bin/bash

CURL='curl --get '
LOGFILE='systemstate.log'

$CURL http://localhost:3030/system/all | jq . >> $LOGFILE
$CURL http://localhost:3030/order/current | jq . >> $LOGFILE
$CURL http://localhost:3030/account/list | jq . >> $LOGFILE

$CURL http://localhost:3030/order/current | jq . >> $LOGFILE

#   BPC operations

$CURL http://localhost:3030/account/account/ndbmgby86qw9bds9f8wrzut5zrbxuehum5kvgz9sns9hgknh | jq . >> $LOGFILE

#   Axiom Foundation

$CURL http://localhost:3030/account/account/ndeeh86uun6us9cenuck2uur679e37uczsmys33794gnvtfz | jq . >> $LOGFILE

#   ndau Network

$CURL http://localhost:3030/account/account/ndnf9ffbzhyf8mk7z5vvqc4quzz5i2exp5zgsmhyhc9cuwr4 | jq . >> $LOGFILE

#   ndev operations

$CURL http://localhost:3030/account/account/ndaea8w9gz84ncxrytepzxgkg9ymi4k7c9p427i6b57xw3r4 | jq . >> $LOGFILE

#   ntrd operations

$CURL http://localhost:3030/account/account/ndmmw2cwhhgcgk9edp5tiieqab3pq7uxdic2wabzx49twwxh | jq . >> $LOGFILE

# 5 ndau node operations accounts

$CURL http://localhost:3030/account/account/ndarw5i7rmqtqstw4mtnchmfvxnrq4k3e2ytsyvsc7nxt2y7 | jq . >> $LOGFILE
$CURL http://localhost:3030/account/account/ndaq3nqhez3vvxn8rx4m6s6n3kv7k9js8i3xw8hqnwvi2ete | jq . >> $LOGFILE
$CURL http://localhost:3030/account/account/ndahnsxr8zh7r6u685ka865wz77wb78xcn45rgskpeyiwuza | jq . >> $LOGFILE
$CURL http://localhost:3030/account/account/ndam75fnjn7cdues7ivi7ccfq8f534quieaccqibrvuzhqxa | jq . >> $LOGFILE
$CURL http://localhost:3030/account/account/ndaekyty73hd56gynsswuj5q9em68tp6ed5v7tpft872hvuc | jq . >> $LOGFILE
