'use strict';

let ui = require('./ui/index.js')
let dns = require('native-dns');
let server = dns.createServer();
let async = require('async');
let qtypeToName = require('native-dns-packet').consts.qtypeToName;

ui.data.containerEntries = [];
server.on('listening', () => console.log('server listening on', server.address()));
server.on('close', () => console.log('server closed', server.address()));
server.on('error', (err, buff, req, res) => console.error(err.stack));
server.on('socketError', (err, socket) => console.error(err));

server.serve(53);

///////

// registra os containers do docker 
var Docker = require('dockerode');
var DockerEvents = require('docker-events');
var docker = new Docker({socketPath: '/var/run/docker.sock'});
var emitter = new DockerEvents({
	docker: docker
});
emitter.start(function(){
	console.log('adicionando containers ja em pe');
	docker.listContainers({all: false}, function(err, containers) {
	containers.forEach(containerInfo => {
		addContainer(containerInfo.Id);
	});
	console.log('!ALL: ' + containers.length);
});
emitter.on("start", function(message) {
	console.log("container started: %j", message);
	addContainer(message.id);
});
emitter.on("stop", function(message) {
	console.log("container stopped: %j", message);
	removeContainer(message.id);
});

//////

server.on('request', function handleRequest(request, response) {

	let questionsToProxy = [];

	/**
	 * DNS server can receive many questions on a same request, 
	 * so we needle 'ask' all requests locally or remotelly then
	 * respond to the reponse variable with the found results calling response.send()
	 * 
	 */
	var nQuestions = request.question.length;
	console.log('request with ', nQuestions, ' questions');
	request.question.forEach(question => {

		console.log('request from:', request.address.address, ' for:', question.name, ' type:', qtypeToName(question.type));

		// finding a entry on local base that matches with question
		let entry = ui.data.entries.filter(r => new RegExp(r.domain, 'i').exec(question.name));
		if (entry.length) {
			// primeiro vemos se nao esta no nosso registro
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
		}else{
			// procurando nos hostnames do container
			let entry = ui.data.containerEntries.filter(r => new RegExp(r.domain, 'i').exec(question.name));
			if (entry.length) {
				// primeiro vemos se nao esta no nosso registro
				entry[0].records.forEach(record => {
					console.log('resolvendo pelo hostname');
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
			}
				// se nenhum satisfazer vamos encaminhar pro remoto
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


function removeContainer(id){
	var index = false;
	ui.data.containerEntries.every((entry, i) => {
		if(entry._id == id){
			index = i;
			return false;
		}
		return true;
	});
	if(index !== false){
		console.log('removido: ', ui.data.containerEntries[index].domain);
		delete ui.data.containerEntries[index];
	}
}
function addContainer(id){
	var container = docker.getContainer(id);
	container.inspect(function (err, data) {
		var hostname = getHostname(data);
		var host = {
			"_id": id,
			"records": [
				{
					"type": "A",
					"address": getHostAddress(data),
					"ttl": 300,
					"name": hostname
				}
			],
			"domain": hostname
		};
		console.log('container adicionando:', host);
		ui.data.containerEntries.push(host);
	});
}
function getHostname(data){
	return data.Config.Hostname + '.' + data.Config.Domainname;
}
function getHostAddress(data){
	return data.NetworkSettings.Networks.bridge.IPAddress;
}