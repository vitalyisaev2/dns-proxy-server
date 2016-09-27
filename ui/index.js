'use strict';

let fs = require('fs');
let express = require('express');
let bodyParser = require('body-parser');

let app = express();
let theRecordFile = process.cwd() + '/conf/records.json';
app.data = getTheJson(theRecordFile);
let password = 'cat';

app.use(bodyParser.json());
app.use(express.static(__dirname));
app.set('json spaces', 2);

app.get('/load', (req, res) => {
	res.send(app.data.entries);
});

app.post('/save', (req, res) => {
	if (req.query.password == password) {
		console.log('salvando: ', req.body);
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