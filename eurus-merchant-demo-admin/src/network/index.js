
function getEurusRpc() {
    return process.env.VUE_APP_EURUS_RPC_URL + ':' + process.env.VUE_APP_EURUS_RPC_PORT;
}

function getEurusChainId() {
    return process.env.VUE_APP_EURUS_CHAIN_ID;
}

function getEurusChainName() {
    return process.env.VUE_APP_EURUS_SIDECHAIN_CHAIN_NAME;
}

export function getEurusApiUrl() {
    return process.env.VUE_APP_EURUS_API_URL;
}

export function getAdminApiUrl() {
    return process.env.VUE_APP_ADMIN_API_URL;
}

export function getSideChainNetworks() {
    let sidechain = {
        name: getEurusChainName(),
        url: getEurusRpc(),
        chainId: getEurusChainId(),
    }
    return sidechain;
}

export function getUsdtSmartContractAddress() {
    return process.env.VUE_APP_EURUS_USDT_ADDRESS;
}

export function getDappSmartContractAddress() {
    return process.env.VUE_APP_EURUS_DAPP_ADDRESS;
}