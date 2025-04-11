#!/bin/sh

echo "Network Up"
minifab netup -s couchdb -e true -i 2.4.8 -o manufacturer.med.com

sleep 5

echo "Create channel"
minifab create -c medchannel

sleep 2

echo "Join channel"
minifab join -c medchannel

sleep 2

echo "Anchor Peer update"
minifab anchorupdate


