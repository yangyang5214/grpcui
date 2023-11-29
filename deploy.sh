#!/bin/bash

cd web && yarn build && cd -

rm -rf cmd/build

mv web/build cmd/