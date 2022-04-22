
import axios from 'axios'
import { computeSHA256 } from "../utils";
import { getAdminApiUrl } from '../network';

export async function merchantLogin(username, password) {
    let result = null;

    try {
        var request = {
            username: username,
            merchantId: "1",
            passwordHash: computeSHA256(password)
        };
        let url = getAdminApiUrl() + "/merchant-admins/login";
        const headers = {
            "Content-Type": "application/json",
        };
        let response = await axios.post(url, request, { headers: headers });
        if (response && response.status === 200) {
            if (response.data) {
                let data = response.data;
                result = data;
            }
        }
    } catch (error) {
        console.error(error);
    }

    return result;
}

export async function getRefundRequestList() {
    let result = null;

    try {
        let url = getAdminApiUrl() + "/withdraw-reqs";
        const headers = {
            "Content-Type": "application/json",
        };
        let response = await axios.get(url, { headers: headers });
        if (response && response.status === 200) {
            if (response.data) {
                let data = response.data;
                result = data;
            }
        }
    } catch (error) {
        console.error(error);
    }

    return result;
}

export async function getRefundRequestLedgers() {
    let result = null;

    try {
        let url = getAdminApiUrl() + "/withdraw-ledgers";
        const headers = {
            "Content-Type": "application/json",
        };
        let response = await axios.get(url, { headers: headers });
        if (response && response.status === 200) {
            if (response.data) {
                let data = response.data;
                result = data;
            }
        }
    } catch (error) {
        console.error(error);
    }

    return result;
}

export async function getMerchantClients() {
    let result = null;

    try {
        let url = getAdminApiUrl() + "/merchant-clients";
        const headers = {
            "Content-Type": "application/x-www-form-urlencoded",
        };
        let response = await axios.get(url, { headers: headers });
        if (response && response.status === 200) {
            if (response.data) {
                let data = response.data;
                result = data;
            }
        }
    } catch (error) {
        console.error(error);
    }

    return result;
}

export async function getDepositLedgers() {
    let result = null;

    try {
        let url = getAdminApiUrl() + "/deposit-ledgers";
        const headers = {
            "Content-Type": "application/x-www-form-urlencoded",
        };
        let response = await axios.get(url, { headers: headers });
        if (response && response.status === 200) {
            if (response.data) {
                let data = response.data;
                result = data;
            }
        }
    } catch (error) {
        console.error(error);
    }

    return result;
}