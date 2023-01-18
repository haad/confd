#!/bin/bash
set -ex

diff /tmp/confd-nested-prefixed-test.conf integration/expect/nested-prefixed.conf
