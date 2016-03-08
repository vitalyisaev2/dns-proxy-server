'use strict';

let ui = require('./ui/index.js')
let dns = require('native-dns');
let server = dns.createServer();
let async = require('async');
let qtypeToName = require('native-dns-packet').consts.qtypeToName;

server.on('listening', () => console.log('server listening on', server.address()));
server.on('close', () => console.log('server closed', server.address()));
server.on('error', (err, buff, req, res) => console.error(err.stack));
server.on('socketError', (err, socket) => console.error(err));

server.serve(53);

function proxy(question, response, cb) {
	console.log('proxying: ', question.name, ', type: ', question.type);

	let server = ui.data.remoteDns[0];
	if(!server){
		throw "You need at least one remote server";
	}
	let request = dns.Request({
		question: question, // forwarding the question
		server: server,  // this is the DNS server we are asking
		timeout: 1000
	});

	request.on('timeout', function () {
		console.log('Timeout in making request no forwarding', question.name);
	});

	// when we get answers, append them to the response
	request.on('message', (err, msg) => {

		msg.answer.forEach(a => {
			response.answer.push(a);
		});
		msg.authority.forEach(a => {
			response.answer.push(a);
		});
	});

	request.on('end', function(){
		response.answer.forEach(msg => {
			console.log('remote DNS answer: type: ', msg.type, ', name: ', msg.name, ', address: ', msg.address);
		})
		cb();
	});
	request.send();
}


server.on('request', function handleRequest(request, response) {

	let questionsToProxy = [];

	/**
	 * DNS server can receive many questions on a same request, 
	 * so we needle 'ask' all requests locally or remotelly then
	 * respond to the reponse variable with the found results calling response.send()
	 * 
	 */
	var nQuestions = request.question.length;
	request.question.forEach(question => {

		console.log('request from:', request.address.address, ' for:', question.name, ' type:', qtypeToName(question.type));

		// finding a entry on local base that matches with question
		let entry = ui.data.entries.filter(r => new RegExp(r.domain, 'i').exec(question.name));
		if (entry.length) {
			entry[0].records.forEach(record => {
				record.name = question.name;
				record.ttl = record.ttl || 1800;
				if (record.type == 'CNAME') {
					record.data = record.address;
					questionsToProxy.push(cb => { 
						proxy({
							name: record.data, type: dns.consts.NAME_TO_QTYPE.A, class: 1
						}, response, cb);
					});
				}
				response.answer.push(dns[record.type](record));
			});
		} else {
			// forwarding host
			questionsToProxy.push(cb => proxy(question, response, cb));
		}
	});
	async.parallel(questionsToProxy, function() {
		// when all questions be done (end event) we will close the connection 
		// sending the response
		console.log('all questions be done, send client answers');
		response.send();
	});
});