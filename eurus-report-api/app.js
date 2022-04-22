const express = require('express');
const cors = require('cors');
const bodyParser = require('body-parser');
require('dotenv/config');

const app = express();

app.use(cors())
app.use(bodyParser.json())

//Routes
app.get('/', (req, res) => {
	res.send('We are on HomePage')
});

app.listen(process.env.PORT, () => {
	console.log(`Running RESTful API on port ${process.env.PORT}`);
});