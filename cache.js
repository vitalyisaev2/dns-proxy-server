module.exports = function(cache){

	cache = cache || {};

	this.get = function(question){
		var msg = cache[getKey(question)];
		if(msg){
			console.log("m=get, status=found, questionName=%s, questionType=%s",
						question.name, question.type)
			return msg;
		}
		console.log("m=get, status=not-found, questionName=%s, questionType=%s",
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

	function getKey(question){
		return question.name + '-' + question.type;
	}

	this.getCache = function(){
		return cache;
	}
}

