'use strict';

let fs = require('fs');
let express = require('express');
let bodyParser = require('body-parser');

let app = express();

let entries = app.entries = require(process.cwd() + '/records.json');
let password = 'cat';

app.use(bodyParser.json());
app.use(express.static(__dirname));

app.get('/load', (req, res) => {
	res.send(entries);
});

app.post('/save', (req, res) => {
	if (req.query.password == password) {
		entries = app.entries = req.body;
		fs.writeFileSync('records.json', JSON.stringify(entries, null, '  '));
		res.send('ok');
	} else {
		res.status(401).send('wrong');
	}
});

app.listen(5380);
module.exports = app