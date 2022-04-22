CONTRACT_DIR = smartcontract/build/contracts
SRC = $(CONTRACT_DIR)/*.json
EURUS_JSON_SRC_FILE = $(wildcard $(SRC))

MAINNET_CONTRACT_DIR = mainnet_smart_contract/build/contracts
MAINNET_SRC = $(MAINNET_CONTRACT_DIR)/*.json
MAINNET_JSON_SRC_FILE = $(wildcard $(MAINNET_SRC))


#ARCH = env GOOS=linux GOARCH=amd64 CGO_ENABLED=1
OUTPUT_FOLDER = bin/release
RUN_FOLDER = $(HOME)/eurus/bin
TAGS = dev
GIT_VERSION = `git rev-list -1 HEAD`
BUILD_DATE = `date "+%Y-%m-%d %H:%M:%S %Z"` 
LD_FLAGS= -ldflags="-X eurus-backend/foundation.GitCommit=$(GIT_VERSION) -X 'eurus-backend/foundation.BuildDate=$(BUILD_DATE)'"
-include tag.mak

.PHONY: force deploy

build: configServer authServer userServer kycServer blockChainIndexer approvalObserver depositObserver withdrawObserver signServer background adminServer merchantAdminServer sweepServer userObserver walletBackgroundIndexer serverControlClient passwordServer blockCypherBackground publicDataServer

genabi: $(EURUS_JSON_SRC_FILE) 
	mkdir -p smartcontract/build/typescript
	typescript/node_modules/.bin/typechain --target web3-v1 --out-dir smartcontract/build/typescript  'smartcontract/build/contracts/*.json'

gen_mainnet_abi: $(MAINNET_JSON_SRC_FILE)
	mkdir -p mainnet_smart_contract/build/typescript
	typescript/node_modules/.bin/typechain --target web3-v1 --out-dir mainnet_smart_contract/build/typescript  'mainnet_smart_contract/build/contracts/*.json'

#Setup development environment
setup: 
	sh script/setup.sh

compile_sc:
	cd smartcontract; truffle compile

compile_mainnet_sc:
	cd mainnet_smart_contract; truffle compile

$(EURUS_JSON_SRC_FILE): force
	echo "Generate GO smart contract $@"
	sh -x script/generate_go_smart_contract.sh 0 "$@" smartcontract/build/golang/contract

$(MAINNET_JSON_SRC_FILE): force
	echo "Generate GO smart contract $@"
	sh -x script/generate_go_smart_contract.sh 1 "$@" mainnet_smart_contract/build/golang/mainnet_contract

deploy:
	-bash script/deploy.sh $(filter-out $@,$(MAKECMDGOALS))

run:
	-bash script/run.sh $(filter-out $@,$(MAKECMDGOALS))

stop:
	-bash script/stop_observer.sh $(filter-out $@,$(MAKECMDGOALS))

serverlist:
	tmux list-session



configServer:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/configServer ./config_service

authServer:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/authServer ./auth_service

userServer:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/userServer ./user_service/user_service

blockChainIndexer:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/blockChainIndexer ./report_service/block_chain_indexer

approvalObserver:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/approvalObserver ./asset_service/approval_observer

depositObserver:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/depositObserver ./asset_service/deposit_observer

withdrawObserver:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/withdrawObserver ./asset_service/withdraw_observer

background:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/background ./asset_service/background

walletBackgroundIndexer:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/walletBackground ./report_service/wallet_background_indexer

signServer:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/signServer ./sign_service

adminServer:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/adminServer ./admin_service/admin_service/admin

merchantAdminServer:
#	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/merchantAdminServer ./admin_service/merchant_admin

sweepServer:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/sweepServer ./asset_service/sweep_service

kycServer:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/kycServer ./user_service/kyc_service

passwordServer:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/passwordServer ./password_service

encryptTool:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/encryptTool ./tool/encrypt_tool

extractAbi:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/extractAbi ./tool/extract_abi

extractBin:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/extractBin ./tool/extract_bin

keyGen:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/keyGen ./tool/key_gen

userObserver:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/userObserver ./asset_service/user_observer

serverControlClient:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/serverControlClient ./tool/server_control_client

decryptTool:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/decryptTool ./tool/decrypt_tool

blockCypherBackground:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/blockCypherBackground ./block_cypher_service/background

publicDataServer:
	typescript/node_modules/.bin/tsc -p "typescript/src/server/publicDataServer"
	cp -f typescript/package.json typescript/dist/tsc/publicDataServer
	cp -f typescript/package-lock.json typescript/dist/tsc/publicDataServer
	tar -czvf $(OUTPUT_FOLDER)/publicDataServer.tar.gz -C "typescript/dist/tsc/" publicDataServer
	

dataPatch1:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/dataPatch1 ./data_patch/patch1

dataPatch2:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/dataPatch2 ./data_patch/patch2

dataPatch3:
	$(ARCH) go build -tags $(TAGS) $(LD_FLAGS) -o $(OUTPUT_FOLDER)/dataPatch3 ./data_patch/patch3

# Usage, in terminal, type make truffle HASH=0x8e136fa5d3eb61f3a1c67f50d25af916876a8dd356a15aa0f0f8109ff71b6302
truffle:
	cd smartcontract; truffle debug $(HASH) --network besu
%:
	@:
