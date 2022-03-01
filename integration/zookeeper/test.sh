#!/bin/bash

export HOSTNAME="localhost"
export PORT=2181

set -e

# feed zookeeper
export ZK_PATH="`dirname \"$0\"`"
sh -c "cd $ZK_PATH; go run main.go"

# Run confd
confd --onetime --log-level debug --confdir ./integration/confdir --interval 5 --backend zookeeper --node 127.0.0.1:2181 -watch
