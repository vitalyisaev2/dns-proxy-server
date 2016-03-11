'use strict';

let fs = require('fs');
let express = require('express');
let bodyParser = require('body-parser');

let app = express();
let theRecordFile = process.cwd() + '/records.json';
let entries = app.data = getTheJson(theRecordFile);
let password = 'cat';

app.use(bodyParser.json());
app.use(express.static(__dirname));

app.get('/load', (req, res) => {
	res.send(entries);
});

app.post('/save', (req, res) => {
	if (req.query.password == password) {
		app.data.entries = req.body;
		writeJson(app.data);
		res.send('ok');
	} else {
		res.status(401).send('wrong');
	}
});

function writeJson(json){
	fs.writeFileSync(theRecordFile, JSON.stringify(json, null, '  '), {
		flag: 'w'
	});
}
function getTheJson(name){
	try{
		return JSON.parse(fs.readFileSync(name).toString());
	}catch(e){
		writeJson({
			"remoteDns": [
				{ "address": "8.8.8.8", "port": 53, "type": "udp" }
			],
			"entries" : []
		});
		return getTheJson(name);
	}
}

app.listen(80);
module.exports = app