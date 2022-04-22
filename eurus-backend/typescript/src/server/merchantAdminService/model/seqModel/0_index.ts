import { MerchantClientFactory } from './MerchantClient';
import { MerchantAdminFactory } from './MerchantAdmin';
import { DepositLedgerFactory } from './DepositLedger';
import { WithdrawLedgerFactory } from './WIthdrawLedger';
import { WithdrawReqFactory } from './WithdrawReq';

export enum name {
  MerchantClient = 'merchant_clients',
  MerchantAdmins = 'merchant_admins',
  DepositLedgers = 'deposit_ledgers',
  WithdrawLedgers = 'withdraw_ledgers',
  WithdrawReqs = 'withdraw_reqs',
}

export const factory = {
  MerchantAdminFactory,
  MerchantClientFactory,
  DepositLedgerFactory,
  WithdrawLedgerFactory,
  WithdrawReqFactory,
};

export default { name, factory };
