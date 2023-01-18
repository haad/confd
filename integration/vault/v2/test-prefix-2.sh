#!/bin/bash

export PATH=/tmp/vault/bin:$PATH
export HOSTNAME="localhost"

export VAULT_ADDR="http://127.0.0.1:8200/"
export ROOT_TOKEN="$(vault read -field id auth/token/lookup-self)"

export PREFIX2="secret_v2/platform/nested-prefixed-2/project/"
set -e

vault kv put /${PREFIX2}multiple_secrets FIRST_KEY=FIRST_VALUE ANOTHER_KEY=ANOTHER_VALUE

export ENV_VARS_SECRET=multiple_secrets
confd --onetime --log-level debug \
      --confdir ./integration/confdir-prefixed \
      --prefix /${PREFIX2} \
      --backend vault \
      --auth-type token \
      --auth-token $ROOT_TOKEN \
      --node http://127.0.0.1:8200

