'use strict';

let ui = require('./ui/index.js');
let dns = require('native-dns');
let Cache = require('./cache.js');
let server = dns.createServer();
let async = require('async');
let qtypeToName = require('native-dns-packet').consts.qtypeToName;

var cache = new Cache(ui.data.cache = {}, ui.data.dnsCacheTTL);
ui.data.containerEntries = [];
require('./hostnameController.js')(ui);
server.on('listening', () => console.log('server listening on', server.address()));
server.on('close', () => console.log('server closed', server.address()));
server.on('error', (err, buff, req, res) => console.error(err.stack));
server.on('socketError', (err, socket) => console.error(err));

server.serve(ui.data.dnsServerPort || 53);

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
}

console.log('adicionando containers ja em pe');
docker.listContainers({all: false}, function(err, containers) {
	containers.forEach(containerInfo => {
		addContainer(containerInfo.Id);
	});
	console.log('!ALL: ' + containers.length);
});

//////

server.on('request',
	(req, res) => {
		process.nextTick(handleRequest.bind(null, req, res))
	}
);

function handleRequest(request, response) {

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
				// insert the proxy function at the array to be called by parallel
				questionsToProxy.push(cb => proxy(question, response, cb));
			}
		}
	});
	async.parallel(questionsToProxy, function(msg) {
		// when all questions be done (end event) we will close the connection 
		// sending the response
		console.log('m=parallel, status=questions done, action=sending answers, msg=%s', msg);
		response.send();
	});
};

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
			 	a.records ? a.records[0] : null);
			response.answer.push(a);
		});
		return true;
	}
	console.log("m=resolveDnsLocally, status=notFound, question=%s", question.name);
	return false;
}


function proxy(question, response, cb) {
	console.log('m=proxy, status=begin, questionName=', question.name, ', type=', question.type,
		', cache=' + ui.data.cacheEnabled);
	var msg;
	if(ui.data.cacheEnabled && (msg = cache.get(question))){
		console.log('m=proxy, status=resolvedFromCache, host=%s, cacheSize=%s, qtd=%s',
			question.name, cache.size(),
				msg.answer.length);
		doAnswer(response, msg, question.name);
		cb('success');
		return ;
	}
	proxyToServer(question, response, cb, 0);

}

function doAnswer(response, msg, host){
	console.log("m=doAnswer, host=%s, answer=%s, authority=%s, additional=%s", host,
	 msg.answer.length, msg.authority.length, msg.additional.length);

	msg.answer.forEach(a => {
//		console.log('m=answerFound, type=%s, ttl=%s, ip=%s, server=%s', a.type, a.ttl,
//			 a.address, server.address);
		response.answer.push(a);
	});
	msg.authority.forEach(a => {
		response.authority.push(a);
	});
	msg.additional.forEach(a => {
		response.additional.push(a);
	});
}

function proxyToServer(question, response, cb, index){

	if(index >= ui.data.remoteDns.length){
		console.log('m=proxyToServer, status=noMoreServers');
		cb('not-found')
		return;
	}
	let server = ui.data.remoteDns[index];
	console.log('m=proxyToServer, status=resolvingFromRemote, server=%s, index=%s, dnsQtd=%s',
		server.address, index, ui.data.remoteDns.length);
	console.log("m=proxyToServer, host=%s, question=", question.name, question);
	let request = dns.Request({
		question: question, // forwarding the question
		server: server,  // this is the DNS server we are asking
		timeout: 3000
	});

	request.on('timeout', function () {
		console.log('m=timeout, question=%s, server=%s', question.name, server.address);
		proxyToServer(question, response, cb, index+1)
	});

	// when we get answers, append them to the response
	request.on('message', (err, msg) => {

		// mouting cache
		if(ui.data.cacheEnabled){
			cache.put(question, msg);
		}

		console.log('m=answerFound, status=cached, msg=', msg);
		doAnswer(response, msg, question.name);
		cb('success');
	});

	request.on('end', function(){
		response.answer.forEach(msg => {
			console.log('m=remote-end, type=%s, name=%s, address=%s, server=%s', msg.type,
				msg.name, msg.address, server.address);
		})
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