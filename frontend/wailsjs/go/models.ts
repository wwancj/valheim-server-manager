export namespace backup {
	
	export class Backup {
	    id: string;
	    serverId: string;
	    name: string;
	    path: string;
	    size: number;
	    createdAt: string;
	    files: string[];
	
	    static createFrom(source: any = {}) {
	        return new Backup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.serverId = source["serverId"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.createdAt = source["createdAt"];
	        this.files = source["files"];
	    }
	}

}

export namespace config {
	
	export class ServerConfig {
	    serverName: string;
	    worldName: string;
	    password: string;
	    port: number;
	    public: boolean;
	    preset: string;
	    modifiers: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new ServerConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.serverName = source["serverName"];
	        this.worldName = source["worldName"];
	        this.password = source["password"];
	        this.port = source["port"];
	        this.public = source["public"];
	        this.preset = source["preset"];
	        this.modifiers = source["modifiers"];
	    }
	}

}

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

export namespace mods {
	
	export class InstalledMod {
	    name: string;
	    fullName: string;
	    version: string;
	    description: string;
	    enabled: boolean;
	    dependencies: string[];
	    installedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new InstalledMod(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.fullName = source["fullName"];
	        this.version = source["version"];
	        this.description = source["description"];
	        this.enabled = source["enabled"];
	        this.dependencies = source["dependencies"];
	        this.installedAt = source["installedAt"];
	    }
	}

}

export namespace services {
	
	export class ModProfile {
	    id: string;
	    name: string;
	    description: string;
	    mods: string[];
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new ModProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.mods = source["mods"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class SafetyCheck {
	    name: string;
	    status: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new SafetyCheck(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.status = source["status"];
	        this.message = source["message"];
	    }
	}

}

export namespace system {
	
	export class ProcessInfo {
	    pid: number;
	    cpuPercent: number;
	    memRSS: number;
	    memPercent: number;
	    uptime: number;
	    running: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.cpuPercent = source["cpuPercent"];
	        this.memRSS = source["memRSS"];
	        this.memPercent = source["memPercent"];
	        this.uptime = source["uptime"];
	        this.running = source["running"];
	    }
	}
	export class SystemInfo {
	    cpuPercent: number;
	    memUsed: number;
	    memTotal: number;
	    memPercent: number;
	    goVersion: string;
	    os: string;
	    arch: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpuPercent = source["cpuPercent"];
	        this.memUsed = source["memUsed"];
	        this.memTotal = source["memTotal"];
	        this.memPercent = source["memPercent"];
	        this.goVersion = source["goVersion"];
	        this.os = source["os"];
	        this.arch = source["arch"];
	    }
	}

}

export namespace thunderstore {
	
	export class Latest {
	    namespace: string;
	    name: string;
	    version_number: string;
	    full_name: string;
	    description: string;
	    icon: string;
	    dependencies: string[];
	    download_url: string;
	    downloads: number;
	    date_created: string;
	
	    static createFrom(source: any = {}) {
	        return new Latest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.namespace = source["namespace"];
	        this.name = source["name"];
	        this.version_number = source["version_number"];
	        this.full_name = source["full_name"];
	        this.description = source["description"];
	        this.icon = source["icon"];
	        this.dependencies = source["dependencies"];
	        this.download_url = source["download_url"];
	        this.downloads = source["downloads"];
	        this.date_created = source["date_created"];
	    }
	}
	export class TSPackage {
	    name: string;
	    full_name: string;
	    owner: string;
	    description: string;
	    version_number: string;
	    total_downloads: number;
	    rating_score: number;
	    package_url: string;
	    icon: string;
	    latest?: Latest;
	
	    static createFrom(source: any = {}) {
	        return new TSPackage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.full_name = source["full_name"];
	        this.owner = source["owner"];
	        this.description = source["description"];
	        this.version_number = source["version_number"];
	        this.total_downloads = source["total_downloads"];
	        this.rating_score = source["rating_score"];
	        this.package_url = source["package_url"];
	        this.icon = source["icon"];
	        this.latest = this.convertValues(source["latest"], Latest);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class Version {
	    version_number: string;
	    download_url: string;
	    dependencies: string[];
	    file_size: number;
	    date_created: string;
	
	    static createFrom(source: any = {}) {
	        return new Version(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version_number = source["version_number"];
	        this.download_url = source["download_url"];
	        this.dependencies = source["dependencies"];
	        this.file_size = source["file_size"];
	        this.date_created = source["date_created"];
	    }
	}

}

