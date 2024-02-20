export namespace models {
	
	export class Account {
	    id: string;
	    client_id: string;
	    client_secret: string;
	    username: string;
	    password: string;
	    require_password: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.client_id = source["client_id"];
	        this.client_secret = source["client_secret"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.require_password = source["require_password"];
	    }
	}

}

