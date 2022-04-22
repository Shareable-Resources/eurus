const proxy = require('express-http-proxy');
const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');
require('dotenv/config');

const app = express();
const PORT = process.env.PORT;
const API_PATH = process.env.API_PATH;
const POXY_TARGET = process.env.POXY_TARGET;

console.log('##### PORT:', PORT);
console.log('##### API_PATH:', API_PATH);
console.log('##### POXY_TARGET:', POXY_TARGET);

app.use(cors())
app.use(bodyParser.json())

app.use(API_PATH, proxy(POXY_TARGET));

app.get('/', (req, res) => {
    res.send('We are on HomePage')
});

app.listen(PORT, () => {
    console.log("Running RESTful API on port", PORT);
});
