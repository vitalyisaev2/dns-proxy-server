'use strict';

let ui = require('./ui/index.js')
let dns = require('native-dns');
let server = dns.createServer();
let async = require('async');
let qtypeToName = require('native-dns-packet').consts.qtypeToName;

ui.data.containerEntries = [];
require('./hostnameController.js')(ui);
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
emitter.start();
emitter.on("start", function(message) {
	console.log("container started: %j", message);
	addContainer(message.id);
});

emitter.on("stop", removeEvent);
emitter.on("die", removeEvent);
emitter.on("destroy", removeEvent);

function removeEvent(message) {
	console.log("m=removeEvent, msg=%j", message);
	removeContainer(message.id);
};

console.log('adicionando containers ja em pe');
docker.listContainers({all: false}, function(err, containers) {
	containers.forEach(containerInfo => {
		addContainer(containerInfo.Id);
	});
	console.log('!ALL: ' + containers.length);
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
	console.log('m=request, questions=', nQuestions);
	request.question.forEach(question => {

		console.log('m=solve, requestFrom=', request.address.address, ', for=', question.name,
			', type=', qtypeToName(question.type));
		console.log('m=solve, from=containers')
		// finding container that matches hostname
		if(!resolveDnsLocally(ui.data.containerEntries, question, questionsToProxy, response)){
			console.log('m=solve, from=json')
			// finding a entry on local base that matches with question
			if(!resolveDnsLocally(ui.data.entries, question, questionsToProxy, response)){
				console.log('m=solve, from=remote')
				// if not found locally then will find on WEB DNS
				questionsToProxy.push(cb => proxy(question, response, cb));
			}
		}
	});
	async.parallel(questionsToProxy, function() {
		// when all questions be done (end event) we will close the connection 
		// sending the response
		console.log('m=parallel, status=questions done, aciton=seding answers');
		response.send();
	});
});

function resolveDnsLocally(entries, question, questionsToProxy, response){
	console.log("m=resolveDnsLocally, status=begin, question=%s", question.name);
	let entry = entries.filter(r => new RegExp('^' + r.domain + '$', 'i').test(question.name));
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
			var a = dns[record.type](record);
			console.log("m=resolveDnsLocally, status=success, question=%s, answer=%s", question.name,
			 	a.records[0]);
			response.answer.push(a);
		});
		return true;
	}
	console.log("m=resolveDnsLocally, status=notFound, question=%s", question.name);
	return false;
}

function proxy(question, response, cb) {
	console.log('m=proxy, questionName=', question.name, ', type=', question.type);

	let server = ui.data.remoteDns[0];
	if(!server){
		throw "You need at least one remote server";
	}
	let request = dns.Request({
		question: question, // forwarding the question
		server: server,  // this is the DNS server we are asking
		timeout: 2000
	});

	request.on('timeout', function () {
		console.log('m=timeout, question=%s', question.name);
	});

	// when we get answers, append them to the response
	request.on('message', (err, msg) => {

		msg.answer.forEach(a => {
			console.log('m=answerFound, answer=', a)
			response.answer.push(a);
		});
		msg.authority.forEach(a => {
			response.answer.push(a);
		});
	});

	request.on('end', function(){
		response.answer.forEach(msg => {
			console.log('m=remote-end, type=', msg.type, ', name=', msg.name, ', address=', msg.address);
		})
		cb();
	});
	request.send();
}


function removeContainer(id){
	ui.data.containerEntries = ui.data.containerEntries.filter((entry, i) => {
		if(entry._id == id){
			console.log('M=removeContainer, container=%s, domain=%s', entry.container, entry.domain);
			return false
		}
		return true;
	});
}

function addContainer(id){
	var container = docker.getContainer(id);
	container.inspect(function (err, data) {
		console.info('m=addContainer, status=process-hostnames, container=', data.Name);
		getHostnames(data).forEach(hostname => {
			var ip = getHostAddress(data);
			var host = {
				"_id": id,
				"container": data.Name,
				"ip": ip,
				"records": [
					{
						"type": "A",
						"address": ip,
						"ttl": 300,
						"name": hostname
					}
				],
				"domain": hostname
			};
			console.info('m=addContainer, container=%s, hostname=%s', data.Name, hostname);
			ui.data.containerEntries.push(host);
		});
	});
}
function getHostnames(container){
	var hostnames = [getHostname(container)];
	if(Array.isArray(container.Config.Env)){
		container.Config.Env.forEach(function(env){
			var key = 'HOSTNAMES=';
			if(env.startsWith(key)){
				console.info('m=getHostnames, status=env fold');
				var strHosts = env.substring(key.length),
						arrHosts = strHosts.split(',');
				hostnames = hostnames.concat(arrHosts);
				console.log('m=getHostnames, container=%s, hosts=%s: ', container.Name, hostnames);
			}
		});
	}
	return hostnames;
}
function getHostname(data){
	var buff = data.Config.Hostname;
	if(data.Config.Domainname){
		buff += '.' + data.Config.Domainname;
	}
	return buff;
}
function getHostAddress(data){
	var networks = data.NetworkSettings.Networks;
	var network;
	Object.keys(networks).forEach(n => {
		if(!network){
			network = n
		}
	});
	return networks[network].IPAddress;
}