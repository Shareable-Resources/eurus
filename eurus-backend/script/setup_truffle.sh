#!/bin/sh

DIR=$(dirname "$0")
cd ${DIR}/../smartcontract
echo "Init npm at smartcontract folder"
npm init -y
echo "Installing truffle"
npm install truffle

npm install @truffle/hdwallet-provider

cd ${DIR}/../mainnet_smart_contract
echo "Init npm at mainnet smartcontract folder"
npm init -y
echo "Installing truffle"
npm install truffle

npm install @truffle/hdwallet-provider

echo "Install truffle completed"
