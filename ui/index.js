'use strict';

let fs = require('fs');
let express = require('express');
let bodyParser = require('body-parser');

let entries = require('./records.json');
let password = 'cat';

let app = express();

app.use(bodyParser.json());
app.use(express.static(__dirname));

app.get('/load', (req, res) => {
	res.send(entries);
});

app.post('/save', (req, res) => {
	if (req.query.password == password) {
		entries = req.body;
		fs.writeFileSync('records.json', JSON.stringify(entries));
		res.send('ok');
	} else {
		res.status(401).send('wrong');
	}
});

app.listen(5380);
module.exports = app