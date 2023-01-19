#!/bin/bash

export PATH=/tmp/vault/bin:$PATH
export HOSTNAME="localhost"

export VAULT_ADDR="http://127.0.0.1:8200/"
export ROOT_TOKEN="$(vault read -field id auth/token/lookup-self)"

export PREFIX="secret_v1/platform/nested-prefixed-1/project/"
set -e

vault write /${PREFIX}FIRST_KEY value=FIRST_VALUE
vault write /${PREFIX}ANOTHER_KEY value=ANOTHER_VALUE

# Run confd with prefix
confd --onetime --log-level debug \
      --confdir ./integration/confdir-prefixed \
      --prefix /${PREFIX} \
      --backend vault \
      --auth-type token \
      --auth-token $ROOT_TOKEN \
      --node http://127.0.0.1:8200
