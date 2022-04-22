#!/bin/zsh

set -e

mkdir -p ./oapi

~/go/bin/oapi-codegen -generate types,chi-server,spec -package oapi doc/openapi/reference/payment-gateway.openapi.yaml > oapi/dapp_payment_api.gen.go
