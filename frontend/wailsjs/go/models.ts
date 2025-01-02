export namespace backend {
	
	export class DocInfo {
	    Id: string;
	    ParentId: string;
	    IsFolder: boolean;
	    Name: string;
	    Bookmarked: boolean;
	    // Go type: time
	    LastModified?: any;
	    FileType?: string;
	
	    static createFrom(source: any = {}) {
	        return new DocInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.ParentId = source["ParentId"];
	        this.IsFolder = source["IsFolder"];
	        this.Name = source["Name"];
	        this.Bookmarked = source["Bookmarked"];
	        this.LastModified = this.convertValues(source["LastModified"], null);
	        this.FileType = source["FileType"];
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

}

