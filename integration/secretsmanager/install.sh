#!/bin/bash

# feed zookeeper
export SSM_PATH="`dirname \"$0\"`"
sh -c "cd $SSM_PATH; go run main.go &;"