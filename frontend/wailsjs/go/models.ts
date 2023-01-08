export namespace main {
	
	export class ControlBar {
	    selectedUid: string;
	
	    static createFrom(source: any = {}) {
	        return new ControlBar(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.selectedUid = source["selectedUid"];
	    }
	}
	export class GachaPieTotals {
	    t301: models.GachaTotal[];
	    t302: models.GachaTotal[];
	    t200: models.GachaTotal[];
	    t100: models.GachaTotal[];
	
	    static createFrom(source: any = {}) {
	        return new GachaPieTotals(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.t301 = this.convertValues(source["t301"], models.GachaTotal);
	        this.t302 = this.convertValues(source["t302"], models.GachaTotal);
	        this.t200 = this.convertValues(source["t200"], models.GachaTotal);
	        this.t100 = this.convertValues(source["t100"], models.GachaTotal);
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
	export class GachaPieDate {
	    usedCosts: models.GachaUsedCost[];
	    totals: GachaPieTotals;
	
	    static createFrom(source: any = {}) {
	        return new GachaPieDate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.usedCosts = this.convertValues(source["usedCosts"], models.GachaUsedCost);
	        this.totals = this.convertValues(source["totals"], GachaPieTotals);
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
	    controlBar: ControlBar;
	
	    static createFrom(source: any = {}) {
	        return new Option(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.showGacha = this.convertValues(source["showGacha"], ShowGacha);
	        this.otherOption = this.convertValues(source["otherOption"], OtherOption);
	        this.controlBar = this.convertValues(source["controlBar"], ControlBar);
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

export namespace models {
	
	export class GachaLog {
	    uid: string;
	    gacha_type: string;
	    time: string;
	    name: string;
	    lang: string;
	    item_type: string;
	    rank_type: string;
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new GachaLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uid = source["uid"];
	        this.gacha_type = source["gacha_type"];
	        this.time = source["time"];
	        this.name = source["name"];
	        this.lang = source["lang"];
	        this.item_type = source["item_type"];
	        this.rank_type = source["rank_type"];
	        this.id = source["id"];
	    }
	}
	export class GachaTotal {
	    total: number;
	    itemType: string;
	    rankType: string;
	
	    static createFrom(source: any = {}) {
	        return new GachaTotal(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.itemType = source["itemType"];
	        this.rankType = source["rankType"];
	    }
	}
	export class GachaUsedCost {
	    gachaType: string;
	    cost: number;
	
	    static createFrom(source: any = {}) {
	        return new GachaUsedCost(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gachaType = source["gachaType"];
	        this.cost = source["cost"];
	    }
	}

}

