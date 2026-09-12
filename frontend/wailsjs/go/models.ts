export namespace models {
	
	export class Server {
	    id: string;
	    name: string;
	    installPath: string;
	    worldName: string;
	    port: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Server(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.installPath = source["installPath"];
	        this.worldName = source["worldName"];
	        this.port = source["port"];
	        this.status = source["status"];
	    }
	}

}

