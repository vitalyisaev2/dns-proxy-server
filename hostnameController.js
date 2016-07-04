module.exports = function(app){
	var http = require('http'), util = require('util');

	app.get('/containers', (req, res) => {

		app.data.containerEntries.forEach(c => {
			res.write(util.format('container=%s, ip=%s, domain=%s', c.container, c.ip, c.domain));
			res.write('\n');
		});
		res.end();

	})
}