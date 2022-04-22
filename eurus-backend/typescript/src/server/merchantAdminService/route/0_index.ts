/* eslint-disable @typescript-eslint/camelcase */
import Route from '../../../foundation/server/CommonRoute';
import UsersRoute from './MerchantClientsRoute';
import MerchantAdminsRoute from './MerchantAdminsRoute';
import DepositLedgersRoute from './DepositLedgersRoute';
import WithdrawLedgersRoute from './WithdrawLedgersRoute';
import WithdrawReqsRoute from './WithdrawReqsRoute';
import DbRoute from './DbRoute';
import LogRoute from './LogsRoute';
const router: Array<Route> = [
  new UsersRoute(),
  new MerchantAdminsRoute(),
  new DepositLedgersRoute(),
  new WithdrawLedgersRoute(),
  new WithdrawReqsRoute(),
  new DbRoute(),
  new LogRoute(),
];

export default router;
