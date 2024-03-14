export namespace logger {
	
	export class Event {
	    message: string;
        done: boolean;
        error: string;
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.message = source["message"];
	        this.done = source["done"];
	        this.error = source["error"];
	    }
	}

}