/* eslint-disable */
import { 
    getEurusChainId,
    getEurusRpcDomain,
    getEurusRpcPort,
} from "@/utils/auth";

function getEurusRpc() {
    return getEurusRpcDomain() + ':' + getEurusRpcPort();
}

export function getEurusApiUrl() {
    return process.env.VUE_APP_EURUS_RPC_URL;
}

export function getSideChainNetworks() {
    let sidechain = {
        name: 'Eurus Dev',
        url: getEurusRpc(),
        chainId: getEurusChainId(),
    }
    return sidechain;
}