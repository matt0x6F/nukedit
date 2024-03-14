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
	export class Schedule {
	    id: number;
	    username: string;
	    cronExpression: string;
	    posts: boolean;
	    comments: boolean;
	    maxAge: number;
	    useMaxAge: boolean;
	    minScore: number;
	    useMinScore: boolean;
	    replacementTextLength: number;
	
	    static createFrom(source: any = {}) {
	        return new Schedule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.username = source["username"];
	        this.cronExpression = source["cronExpression"];
	        this.posts = source["posts"];
	        this.comments = source["comments"];
	        this.maxAge = source["maxAge"];
	        this.useMaxAge = source["useMaxAge"];
	        this.minScore = source["minScore"];
	        this.useMinScore = source["useMinScore"];
	        this.replacementTextLength = source["replacementTextLength"];
	    }
	}

}

export namespace reddit {
	
	export class NukeRequest {
	    scheduled: boolean;
	    cronExpression: string;
	    posts: boolean;
	    comments: boolean;
	    maxAge: number;
	    useMaxAge: boolean;
	    minScore: number;
	    useMinScore: boolean;
	    replacementTextLength: number;
	
	    static createFrom(source: any = {}) {
	        return new NukeRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scheduled = source["scheduled"];
	        this.cronExpression = source["cronExpression"];
	        this.posts = source["posts"];
	        this.comments = source["comments"];
	        this.maxAge = source["maxAge"];
	        this.useMaxAge = source["useMaxAge"];
	        this.minScore = source["minScore"];
	        this.useMinScore = source["useMinScore"];
	        this.replacementTextLength = source["replacementTextLength"];
	    }
	}
	export class NukeResult {
	    commentsDeleted: number;
	    postsDeleted: number;
	
	    static createFrom(source: any = {}) {
	        return new NukeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.commentsDeleted = source["commentsDeleted"];
	        this.postsDeleted = source["postsDeleted"];
	    }
	}

}

