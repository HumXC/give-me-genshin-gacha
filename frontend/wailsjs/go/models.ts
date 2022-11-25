export namespace main {
	
	export class ItemsInfo {
	    name: string;
	    usec: number;
	    time: number;
	
	    static createFrom(source: any = {}) {
	        return new ItemsInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.usec = source["usec"];
	        this.time = source["time"];
	    }
	}
	export class GachaInfo {
	    name: string;
	    count: number;
	    s3: number;
	    s4: number;
	    s5: number;
	    s5Items: ItemsInfo[];
	
	    static createFrom(source: any = {}) {
	        return new GachaInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.count = source["count"];
	        this.s3 = source["s3"];
	        this.s4 = source["s4"];
	        this.s5 = source["s5"];
	        this.s5Items = this.convertValues(source["s5Items"], ItemsInfo);
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

