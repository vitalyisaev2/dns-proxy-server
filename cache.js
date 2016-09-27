module.exports = function(cache){

	cache = cache || {};

	this.get = function(question){
		var msg = cache[question.name];
		if(msg && msg.question.filter( q => q.type == question.type ).length > 0){
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

		var tmp = cache[question.name];
		if(tmp == null){
			console.log("m=put, status=puttingNewQuestion, questionName=%s, questionType=%s",
					question.name, question.type)
			cache[question.name] = value;
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

	this.getCache = function(){
		return cache;
	}
}

