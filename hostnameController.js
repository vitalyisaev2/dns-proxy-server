module.exports = function(app, ce){
	var http = require('http'), util = require('util');

	app.get('/containers', (req, res) => {

		ce.forEach(c => {
			res.write(util.format('container=%s, domain=%s', c.container, c.domain));
			res.write('\n');
		});
		res.end();

	})
}