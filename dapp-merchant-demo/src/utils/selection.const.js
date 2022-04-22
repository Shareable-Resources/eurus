const DELETED = 'deleted'
const NORMAL = 'normal'
const CLOSED = 'closed'
const FROZEN = 'frozen'
const INITIAL = 'initial'
const PENDING = 'pending'
const INACTIVE = 'inactive'
const ACTIVE = 'active'
const REJECTED = 'rejected'
const NAMEFILLED = 'namefilled'
const RESTAURANT = 'restaurant'
const MARKET = 'market'
const TRAVEL = 'travel'
const OTHERS = 'others'
const DECENTRALIZED = 'decentralized'
const CENTRALIZED = 'centralized'
const MERCHANT = 'merchant'
const ALLWALLET = 'allWallet'
const CREATEWALLET = 'createWallet'
const GROUPWALLET = 'groupWallet'
const FEEWALLET = 'feeWallet'
const WITHDRAWWALLET = 'withdrawWallet'
const DISABLE = 'disable'
const DEPOSIT = 'deposit'
const WITHDRAW = 'withdraw'
const PAYMENT = 'payment'
const TRANSFER = 'transfer'
const LISTED = 'listed'
const UNLISTED = 'unlisted'
const VISIBLE = 'visible'
const INVISIBLE = 'invisible'
const FROM = 'from'
const TO = 'to'
const DEFAULT = 'default'
const TYPE_ETH = 'TYPE_ETH'
const EURUS_ETH = 'EURUS_ETH'
const AUTOFEE = 'autoFee'
const MANUAL = 'manual'
const CONNECTED = 'connected'
const DISCONNECTED = 'disconnected'
const FIXED_AMOUNT = 'fixAmount'
const BTC_USDT = 'BTC/USDT'
const ETH_USDT = 'ETH/USDT'
const TRX_USDT = 'TRX/USDT'
const DGT_USDT = 'DGT/USDT'
const GCT_USDT = 'GCT/USDT'
const BINANCE = 'Binance'
const COINBASE = 'Coinbase'
const GAINBIT = 'Gainbit'
const KRAKEN = 'Kraken'
const HUOBI = 'Huobi'
const MAIN_CHAIN = 'mainChain'
const SIDE_CHAIN = 'sideChain'
const MERCHANT_MAIN_CHAIN_SUB_ADDR = 'merchantMainChainSubAddr'
const ALL_WALLET = 'allWallet'
const GEN_ADDR_WALLET = 'genAddrWallet'
const SWEEP_WALLET = 'sweepWallet'
const SWEEP_FEE_WALLET = 'sweepFeeWallet'
const HOT_WALLET = 'hotWallet'
// const BACKENDVAL = 'backendVal'
// const map = {
//   [BACKENDVAL]: i18nLabel
// }
// transfer
export const TX_TYPE = {
  [DEPOSIT]: 'deposit',
  [WITHDRAW]: 'withdraw',
  [PAYMENT]: 'payment',
  [TRANSFER]: 'transfer'
}
// user
export const USER_STATUS = {
  [DELETED]: 'deleted',
  [NORMAL]: 'normal',
  [CLOSED]: 'closed',
  [FROZEN]: 'frozen'
}
export const USER_KYCSTATUS = {
  [INITIAL]: 'initial',
  [PENDING]: 'pending',
  [INACTIVE]: 'inactive',
  [ACTIVE]: 'active',
  [REJECTED]: 'rejected',
  [NAMEFILLED]: 'namefilled'
}
// merchant
export const MERCHANT_STATUS = {
  [NORMAL]: 'normal',
  [DISABLE]: 'disable'
}
export const MERCHANT_BUSINESSTYPE = {
  [RESTAURANT]: 'restaurant',
  [MARKET]: 'market',
  [TRAVEL]: 'travel',
  [OTHERS]: 'others'
}
// user-bal
export const USER_GROUP = {
  [DECENTRALIZED]: 'decentralized',
  [CENTRALIZED]: 'centralized',
  [MERCHANT]: 'merchant'
}
// wallet-bal
export const WALLET_BAL_ADDRTYPE = {
  [ALLWALLET]: 'allWallet',
  [CREATEWALLET]: 'createWallet',
  [GROUPWALLET]: 'groupWallet',
  [FEEWALLET]: 'feeWallet',
  [WITHDRAWWALLET]: 'withdrawWallet'
}
// config-asset
export const CONFIG_ASSET_STATUS = {
  [LISTED]: 'listed',
  [UNLISTED]: 'unlisted'
}
export const CONFIG_ASSET_WITHDRAW_FEE = {
  [FIXED_AMOUNT]: 'fixAmount',
  [DEFAULT]: 'default'
}
export const CONFIG_ASSET_APP_VISIBLE = {
  [VISIBLE]: 'visible',
  [INVISIBLE]: 'invisible'
}
export const CONFIG_ASSET_ADDR_TYPE = {
  [TYPE_ETH]: 'TYPE_ETH',
  [EURUS_ETH]: 'EURUS_ETH'
}
// config-fee
export const CONFIG_FEE_PAID_BY = {
  [FROM]: 'from',
  [TO]: 'to',
  [DEFAULT]: 'default'
}
export const CONFIG_FEE_STATUS = {
  [NORMAL]: 'normal',
  [DISABLE]: 'disable'
}
export const CONFIG_FEE_TRANS_TYPE = {
  [PAYMENT]: 'payment',
  [TRANSFER]: 'transfer'
}
// config-market price
export const CONFIG_MKP_FEE_TYPE = {
  [AUTOFEE]: 'autoFee',
  [MANUAL]: 'manual'
}
export const CONFIG_MKP_SOURCE_STATUS = {
  [CONNECTED]: 'connected',
  [DISCONNECTED]: 'disconnected'
}
export const CONFIG_MKP_TRANS_PAIR = {
  [BTC_USDT]: 'BTC_USDT',
  [ETH_USDT]: 'ETH_USDT',
  [TRX_USDT]: 'TRX_USDT',
  [DGT_USDT]: 'DGT_USDT',
  [GCT_USDT]: 'GCT_USDT'
}
export const CONFIG_MKP_SOURCE = {
  [BINANCE]: 'binance',
  [COINBASE]: 'coinbase',
  [GAINBIT]: 'gainbit',
  [KRAKEN]: 'kraken',
  [HUOBI]: 'huobi'
}
// stats-user addr
export const STATS_USER_ADDR_TYPE = {
  [MAIN_CHAIN]: 'mainChain',
  [SIDE_CHAIN]: 'sideChain',
  [MERCHANT_MAIN_CHAIN_SUB_ADDR]: 'merchantMainChainSubAddr'
}
export const STATS_PC_WALLET_GROUP = {
  [ALL_WALLET]: 'allWallet',
  [GEN_ADDR_WALLET]: 'genAddrWallet',
  [SWEEP_WALLET]: 'sweepWallet',
  [SWEEP_FEE_WALLET]: 'sweepFeeWallet',
  [HOT_WALLET]: 'hotWallet'
}
