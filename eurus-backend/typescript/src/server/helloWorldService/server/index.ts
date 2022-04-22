import express from 'express';
import config from '../config/HelloWorldServerConfig.json';
import {
  ContentTypeMiddleWare,
  timeMiddleware,
  accessControlAllowMiddleware,
  clientErrorHandler,
} from '../../../foundation/server/Middlewares';
import router from '../route/0_index';
import logger from '../util/ServiceLogger';
const app = express();
app.use(accessControlAllowMiddleware);
app.use(express.urlencoded({ extended: true }));
app.use(express.json());
app.use(timeMiddleware);
app.use(clientErrorHandler);

// We provide a root route just as an example
app.get('/', (req: any, res: any) => {
  res.send(
    `<h2>Hello, ${config.express.name} is listening!</h2><p>Admin Server Starts</p>`,
  );
});
for (const route of router) {
  app.use(route.getRouter());
}
//Set Content Type
app.use(ContentTypeMiddleWare);
logger.info(
  `Express server [${config.express.name}] started on port ${config.express.port}. Try some routes, such as '/api/users'.`,
);

export const printRoutes = (server: express.Express) => {
  const routes: any[] = [];
  server._router.stack.forEach((middleware) => {
    if (middleware.route) {
      // routes registered directly on the app
      routes.push(middleware.route);
    } else if (middleware.name === 'router') {
      // router middleware
      middleware.handle.stack.forEach((handler) => {
        const route = handler.route;
        if (route) {
          routes.push(route);
        }
      });
    }
  });

  routes.forEach((r) => {
    const methods: string[] = [];
    for (const key in r.methods) {
      methods.push(key.toUpperCase());
    }
    logger.info(
      `localhost:${config.express.port}${r.path}  (${methods.join(',')})`,
    );
  });
};

export default app;
