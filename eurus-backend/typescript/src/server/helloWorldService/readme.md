# hello_world_service (Express.js+Sequelize.js)

hello_world_service is a ts library for setting up express server with sequelize. Default using winston logger. It is example project.

The project is ts based, meaning you need compile it from ts to js first (using tsc). The compile script is already defined in package.json

## Installation
npm install

## Usage
1. Entry point : hello_world_service/index.ts
2. You need to create the schema specified in config/HelloWorldServerConfig.json first (sequelize.schema)
3. Table creation and dummy data (hello_world_service/script/CreateTable.ts)
4. Express Server : hello_world_service/server.ts
5. Sequelize      : hello_world_service/sequelize/index.ts
6. Logger         : hello_world_service/util/serviceLogger.ts
7. Main File      : hello_world_service/index.ts
8. Config         : hello_world_service/config
9. Router         : hello_wolrd_service/route
10. Model         : hello_wolrd_service/model
-----
11. Use VSCode debugger or npm run script in package.json for testing
-----
## Contributing

## License