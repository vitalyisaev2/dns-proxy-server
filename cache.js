
module.exports = function(cache){

	cache = cache || {};

	this.get = function(question, cacheTimeout){
		var msg = cache[getKey(question)];
		if(msg){
			console.log("m=get, status=found, questionName=%s, questionType=%s",
						question.name, question.type)

			// test cache from DNS records
			console.log("m=get, status=test-cache-from-dns, questionName=%s, questionType=%s, cacheTimeout=%s",
					question.name, question.type, cacheTimeout)
			var dnsRecords;
			if( (dnsRecords = msg.authority.filter(v => v.type == 2)).length > 0 ){
				if( isValidCache(msg.creationDate, dnsRecords[0].ttl) ){
					console.log("m=get, status=valid-from-dns, questionName=%s, questionType=%s, ttl=%s",
								question.name, question.type, aRecords[0].ttl)
					return msg;
				}else{
					console.log("m=get, status=dns-expired, questionName=%s, questionType=%s, ttl=%s",
							question.name, question.type, dnsRecords[0].ttl)
				}
			}else{
				console.log("m=get, status=anwser-without-dns-record, questionName=%s, questionType=%s",
						question.name, question.type)
			}

			// Test DNS cache from A records
			console.log("m=get, status=test-cache-from-A, questionName=%s, questionType=%s, cacheTimeout=%s",
					question.name, question.type, cacheTimeout)
			var aRecords = msg.answer.filter(v => v.type == 1);
			if(!aRecords.length){
				console.log("m=get, status=anwser-without-a-record, questionName=%s, questionType=%s",
											question.name, question.type)
			}
			if( isValidCache(msg.creationDate, aRecords[0]) ){
				console.log("m=get, status=valid-from-a, questionName=%s, questionType=%s, ttl=%s",
							question.name, question.type, aRecords[0].ttl)
				return msg;
			}else {
				console.log("m=get, status=a-expired, questionName=%s, questionType=%s, ttl=%s",
						question.name, question.type, aRecords[0].ttl)
			}
		}

		// Test global cache timeout
			console.log("m=get, status=test-cache-from-global-cache, questionName=%s, questionType=%s, cacheTimeout=%s",
					question.name, question.type, cacheTimeout)
		if( cacheTimeout > 0 && isValidCache(msg.creationDate, cacheTimeout) ){
			console.log("m=get, status=valid-from-global-timeout, questionName=%s, questionType=%s, cacheTimeout=%s",
				question.name, question.type, cacheTimeout)
				return msg;
		}else{
			console.log("m=get, status=global-timeout-invalid, questionName=%s, questionType=%s, cacheTimeout=%s",
				question.name, question.type, cacheTimeout)
		}
		this.remove(question);
		return null;
	}

	/**
	 * Value must be a question
	 */
	this.put = function(question, value){
		var key = getKey(question);
		var tmp = cache[key];
		if(tmp == null){
			console.log("m=put, status=puttingNewQuestion, questionName=%s, questionType=%s",
					question.name, question.type)
			value.creationDate = new Date();
			cache[key] = value;
		}else{
			if(!this.get(question)){
				console.log("m=put, status=puttingNewType, questionName=%s, questionType=%s",
					question.name, question.type)
				tmp.additional = tmp.additional.concat(value.additional);
				tmp.answer= tmp.answer.concat(value.answer);
				tmp.authority= tmp.authority.concat(value.authority);
				tmp.question= tmp.question.concat(value.question);
			}else{
				console.log("m=put, status=already put, questionName=%s, questionType=%s",
				 question.name, question.type)
			}
		}
	}

	this.remove = function(question){
		var key = getKey(question);
		if(cache.hasOwnProperty(key)){
			delete cache[key];
		}
		return false;
	}

	this.getCache = function(){
		return cache;
	}

	this.keyset = function(){
		return Object.keys(cache);
	}

	this.size = function(){
		return this.keyset().length;
	}

	function getKey(question){
		return question.name + '-' + question.type;
	}

	function getAsSeconds(date){
		return date.getTime() / 1000.0;
	}

	function isValidCache(creationDate, record){
		return getAsSeconds(new Date()) - getAsSeconds(creationDate) <= record.ttl;
	}
}

