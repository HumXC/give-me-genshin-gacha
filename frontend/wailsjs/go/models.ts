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
	export class OtherOption {
	    autoSync: boolean;
	    useProxy: boolean;
	    darkTheme: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OtherOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.autoSync = source["autoSync"];
	        this.useProxy = source["useProxy"];
	        this.darkTheme = source["darkTheme"];
	    }
	}
	export class ShowGacha {
	    roleUp: boolean;
	    armsUp: boolean;
	    permanent: boolean;
	    start: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ShowGacha(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.roleUp = source["roleUp"];
	        this.armsUp = source["armsUp"];
	        this.permanent = source["permanent"];
	        this.start = source["start"];
	    }
	}
	export class Option {
	    showGacha: ShowGacha;
	    otherOption: OtherOption;
	
	    static createFrom(source: any = {}) {
	        return new Option(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.showGacha = this.convertValues(source["showGacha"], ShowGacha);
	        this.otherOption = this.convertValues(source["otherOption"], OtherOption);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

