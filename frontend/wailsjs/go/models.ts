export namespace wallet {
	
	export class Wallet {
	    address: string;
	    address_json: string;
	
	    static createFrom(source: any = {}) {
	        return new Wallet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.address_json = source["address_json"];
	    }
	}

}

export namespace deploy {
	
	export class DeployParams {
	    id: number;
	    nodeEthereumUrl: string;
	    ethereumUrl: string;
	    ethereumNetwork: string;
	    indexerAddress: string;
	    mnemonic: string;
	
	    static createFrom(source: any = {}) {
	        return new DeployParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.nodeEthereumUrl = source["nodeEthereumUrl"];
	        this.ethereumUrl = source["ethereumUrl"];
	        this.ethereumNetwork = source["ethereumNetwork"];
	        this.indexerAddress = source["indexerAddress"];
	        this.mnemonic = source["mnemonic"];
	    }
	}

}

export namespace application {
	
	export class Application {
	    name: string;
	    abbreviation: string;
	    describe: string;
	    status: number;
	
	    static createFrom(source: any = {}) {
	        return new Application(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.abbreviation = source["abbreviation"];
	        this.describe = source["describe"];
	        this.status = source["status"];
	    }
	}

}

export namespace graph {
	
	export class GraphParameter {
	    nodeEthereumUrl: string;
	    ethereumUrl: string;
	    ethereumNetwork: string;
	    indexerAddress: string;
	    mnemonic: string;
	    application: application.Application;
	    applicationId: number;
	
	    static createFrom(source: any = {}) {
	        return new GraphParameter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nodeEthereumUrl = source["nodeEthereumUrl"];
	        this.ethereumUrl = source["ethereumUrl"];
	        this.ethereumNetwork = source["ethereumNetwork"];
	        this.indexerAddress = source["indexerAddress"];
	        this.mnemonic = source["mnemonic"];
	        this.application = this.convertValues(source["application"], application.Application);
	        this.applicationId = source["applicationId"];
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

export namespace account {
	
	export class Account {
	    publicKey: string;
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.publicKey = source["publicKey"];
	    }
	}

}

export namespace resource {
	
	export class Resource {
	    peerId: string;
	    cpu: string;
	    memory: string;
	    systemImage: string;
	    vmType: string;
	    creator: string;
	    // Go type: time.Time
	    expireTime: any;
	    user: string;
	    status: number;
	
	    static createFrom(source: any = {}) {
	        return new Resource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.peerId = source["peerId"];
	        this.cpu = source["cpu"];
	        this.memory = source["memory"];
	        this.systemImage = source["systemImage"];
	        this.vmType = source["vmType"];
	        this.creator = source["creator"];
	        this.expireTime = this.convertValues(source["expireTime"], null);
	        this.user = source["user"];
	        this.status = source["status"];
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
