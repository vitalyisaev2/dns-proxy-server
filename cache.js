module.exports = function(cache){

	cache = cache || {};

	this.get = function(question){
		var msg = cache[question.name];
		if(msg && msg.question.filter( q => q.type == question.type ).length > 0){
			return msg;
		}
		return null;
	}

	/**
	 * Value must be a question
	 */
	this.put = function(key, value){
	var tmp = cache[key];
		if(tmp == null){
			cache[key] = value;
		}else{
			tmp.additional = tmp.additional.concat(value.additional);
			tmp.answer= tmp.answer.concat(value.answer);
			tmp.authority= tmp.authority.concat(value.authority);
			tmp.question= tmp.question.concat(value.question);
		}
	}

	this.getCache = function(){
		return cache;
	}
}

