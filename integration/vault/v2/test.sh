#!/bin/bash

set -e

export PATH=/tmp/vault/bin:$PATH
export HOSTNAME="localhost"

export VAULT_ADDR="http://127.0.0.1:8200/"
export ROOT_TOKEN="$(vault read -field id auth/token/lookup-self)"

vault secrets enable -version 2 -path kv-v2 kv

vault kv put kv-v2/exists key=foobar
vault kv put kv-v2/database host=127.0.0.1 port=3306 username=confd password=p@sSw0rd
vault kv put kv-v2/upstream app1=10.0.1.10:8080 app2=10.0.1.11:8080
vault kv put kv-v2/nested/east app1=10.0.1.10:8080
vault kv put kv-v2/nested/west app2=10.0.1.11:8080

# Run confd
confd --onetime --log-level debug \
      --confdir ./integration/confdir \
      --backend vault \
      --auth-type token \
      --auth-token $ROOT_TOKEN \
      --node http://127.0.0.1:8200

vault secrets disable kv-v2