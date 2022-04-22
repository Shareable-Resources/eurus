/* eslint-disable @typescript-eslint/camelcase */
import Route from '../../../foundation/server/CommonRoute';
import HelloWorldUserRoute from './HelloWorldUsersRoute';
import DbRoute from './DbRoute';
import LogRoute from './LogsRoute';
const router: Array<Route> = [
  new HelloWorldUserRoute(),
  new DbRoute(),
  new LogRoute(),
];

export default router;
