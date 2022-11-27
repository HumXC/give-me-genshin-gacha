export namespace main {
	
	export class Message {
	    type: string;
	    msg: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.msg = source["msg"];
	    }
	}

}

