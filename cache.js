
module.exports = function(cache){

	cache = cache || {};

	this.get = function(question){
		var msg = cache[getKey(question)];
		if(msg){
			var aRecords = msg.answer.filter(v => v.type == 1);
			if(!aRecords.length){
				console.log("m=get, status=anwser-without-a-record, questionName=%s, questionType=%s",
											question.name, question.type)
				this.remove(question);
				return null;
			}
			if(new Date().getTime() - msg.creationDate.getTime() > aRecords[0].ttl){
				console.log("m=get, status=expired, questionName=%s, questionType=%s, ttl=%s",
									question.name, question.type, aRecords[0].ttl)
				this.remove(question);
				return null;
			}
			console.log("m=get, status=expired, questionName=%s, questionType=%s, ttl=%s",
									question.name, question.type, aRecords[0].ttl)
			return msg;
		}
		console.log("m=get, status=found, questionName=%s, questionType=%s",
						question.name, question.type)
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

	function getKey(question){
		return question.name + '-' + question.type;
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
}

