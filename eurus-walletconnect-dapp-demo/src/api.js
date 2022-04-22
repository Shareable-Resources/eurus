import WalletConnect from "walletconnect";
import QRCodeModal from "@walletconnect/qrcode-modal";
import WalletConnectProvider from "@walletconnect/web3-provider";
import Web3 from "web3";

export async function walletConnect() {
    console.log("Wallet Connect API called")

    const bridge = "https://bridge.walletconnect.org";
    const wc = new WalletConnect({ bridge, qrcodeModal: QRCodeModal });
    const connector = await wc.connect();
    console.log("##### connector", connector);

    if (!connector.connected) {
        await connector.createSession();
    }
    return connector;

}

export async function eurusConnect() {
    const walletConnectProvider = new WalletConnectProvider({
        rpc: {
            2021: "http://13.228.169.25:8545",
        },
        infuraId: "c1aeacaf2c504e6283f13a2233cdb0d6",
        qrcode: true,
    });

    await walletConnectProvider.enable();
    const web3 = new Web3(walletConnectProvider);
    const wallet = await web3
    console.log("#####provider", web3);
    console.log("#####wallet", wallet);
    // console.log("#####wallet", chainId);

    return wallet;
}