module.exports = function(app){
	var http = require('http'), util = require('util');

	app.get('/containers', (req, res) => {

		app.data.containerEntries.forEach(c => {
			res.write(util.format('container=%s, ip=%s, domain=%s', c.container, c.ip, c.domain));
			res.write('\n');
		});
		res.end();
	});

	app.get('/cache', (req, res) => {
		v2CacheAPI(req, res);
	});
	app.get('/v1/cache', (req, res) => {
			res.send(app.data.cache);
	});
	app.get('/v2/cache', (req, res) => {
		v2CacheAPI(req, res);
	});

	function v2CacheAPI(req, res){
		var keys = Object.keys(app.data.cache);
		res.send({
			keys: keys,
			size: keys.length,
			data: app.data.cache
		});
	}
}