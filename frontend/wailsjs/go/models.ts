export namespace deploy {
  export class Deployment {
    nodeEthereumUrl: string;
    ethereumUrl: string;
    ethereumNetwork: string;
    indexerAddress: string;

    static createFrom(source: any = {}) {
      return new Deployment(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.nodeEthereumUrl = source['nodeEthereumUrl'];
      this.ethereumUrl = source['ethereumUrl'];
      this.ethereumNetwork = source['ethereumNetwork'];
      this.indexerAddress = source['indexerAddress'];
    }
  }
  export class Stacking {
    networkUrl: string;
    address: string;
    agentAddress: string;
    pledgeAmount: number;

    static createFrom(source: any = {}) {
      return new Stacking(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.networkUrl = source['networkUrl'];
      this.address = source['address'];
      this.agentAddress = source['agentAddress'];
      this.pledgeAmount = source['pledgeAmount'];
    }
  }
  export class Initialization {
    leaseTerm: number;
    userPublicKey: string;
    accountMnemonic: string;

    static createFrom(source: any = {}) {
      return new Initialization(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.leaseTerm = source['leaseTerm'];
      this.userPublicKey = source['userPublicKey'];
      this.accountMnemonic = source['accountMnemonic'];
    }
  }
  export class ParameterInfo {
    // Go type: Initialization
    initialization: any;
    // Go type: Stacking
    stacking: any;
    // Go type: Deployment
    deployment: any;

    static createFrom(source: any = {}) {
      return new ParameterInfo(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.initialization = this.convertValues(source['initialization'], null);
      this.stacking = this.convertValues(source['stacking'], null);
      this.deployment = this.convertValues(source['deployment'], null);
    }

    convertValues(a: any, classs: any, asMap: boolean = false): any {
      if (!a) {
        return a;
      }
      if (a.slice) {
        return (a as any[]).map((elem) => this.convertValues(elem, classs));
      } else if ('object' === typeof a) {
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
  export class DeployParameter {
    id: number;
    // Go type: ParameterInfo
    data: any;

    static createFrom(source: any = {}) {
      return new DeployParameter(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.id = source['id'];
      this.data = this.convertValues(source['data'], null);
    }

    convertValues(a: any, classs: any, asMap: boolean = false): any {
      if (!a) {
        return a;
      }
      if (a.slice) {
        return (a as any[]).map((elem) => this.convertValues(elem, classs));
      } else if ('object' === typeof a) {
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

export namespace application {
  export class AddApplicationParam {
    name: string;
    describe: string;

    static createFrom(source: any = {}) {
      return new AddApplicationParam(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.name = source['name'];
      this.describe = source['describe'];
    }
  }
  export class Application {
    id: number;
    name: string;
    describe: string;
    status: number;
    // Go type: time.Time
    createdAt: any;
    // Go type: time.Time
    updatedAt: any;
    // Go type: gorm.DeletedAt
    deletedAt: any;

    static createFrom(source: any = {}) {
      return new Application(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.id = source['id'];
      this.name = source['name'];
      this.describe = source['describe'];
      this.status = source['status'];
      this.createdAt = this.convertValues(source['createdAt'], null);
      this.updatedAt = this.convertValues(source['updatedAt'], null);
      this.deletedAt = this.convertValues(source['deletedAt'], null);
    }

    convertValues(a: any, classs: any, asMap: boolean = false): any {
      if (!a) {
        return a;
      }
      if (a.slice) {
        return (a as any[]).map((elem) => this.convertValues(elem, classs));
      } else if ('object' === typeof a) {
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
  export class PageApplicationVo {
    items: Application[];
    total: number;

    static createFrom(source: any = {}) {
      return new PageApplicationVo(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.items = this.convertValues(source['items'], Application);
      this.total = source['total'];
    }

    convertValues(a: any, classs: any, asMap: boolean = false): any {
      if (!a) {
        return a;
      }
      if (a.slice) {
        return (a as any[]).map((elem) => this.convertValues(elem, classs));
      } else if ('object' === typeof a) {
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
  export class ApplyVo {
    id: number;
    // Go type: time.Time
    createdAt: any;
    // Go type: time.Time
    updatedAt: any;
    name: string;
    describe: string;
    status: number;

    static createFrom(source: any = {}) {
      return new ApplyVo(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.id = source['id'];
      this.createdAt = this.convertValues(source['createdAt'], null);
      this.updatedAt = this.convertValues(source['updatedAt'], null);
      this.name = source['name'];
      this.describe = source['describe'];
      this.status = source['status'];
    }

    convertValues(a: any, classs: any, asMap: boolean = false): any {
      if (!a) {
        return a;
      }
      if (a.slice) {
        return (a as any[]).map((elem) => this.convertValues(elem, classs));
      } else if ('object' === typeof a) {
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
  export class UpdateApplicationParam {
    id: number;
    name: string;
    describe: string;

    static createFrom(source: any = {}) {
      return new UpdateApplicationParam(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.id = source['id'];
      this.name = source['name'];
      this.describe = source['describe'];
    }
  }
}

export namespace graph {
  export class GraphParameterVo {
    nodeEthereumUrl: string;
    ethereumUrl: string;
    ethereumNetwork: string;
    indexerAddress: string;
    mnemonic: string;
    applicationId: number;
    name: string;
    describe: string;
    status: number;
    // Go type: time.Time
    createdAt: any;
    // Go type: time.Time
    updatedAt: any;

    static createFrom(source: any = {}) {
      return new GraphParameterVo(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.nodeEthereumUrl = source['nodeEthereumUrl'];
      this.ethereumUrl = source['ethereumUrl'];
      this.ethereumNetwork = source['ethereumNetwork'];
      this.indexerAddress = source['indexerAddress'];
      this.mnemonic = source['mnemonic'];
      this.applicationId = source['applicationId'];
      this.name = source['name'];
      this.describe = source['describe'];
      this.status = source['status'];
      this.createdAt = this.convertValues(source['createdAt'], null);
      this.updatedAt = this.convertValues(source['updatedAt'], null);
    }

    convertValues(a: any, classs: any, asMap: boolean = false): any {
      if (!a) {
        return a;
      }
      if (a.slice) {
        return (a as any[]).map((elem) => this.convertValues(elem, classs));
      } else if ('object' === typeof a) {
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
    wsUrl: string;
    orderIndex: number;
    peerId: string;

    static createFrom(source: any = {}) {
      return new Account(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.publicKey = source['publicKey'];
      this.wsUrl = source['wsUrl'];
      this.orderIndex = source['orderIndex'];
      this.peerId = source['peerId'];
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
      this.peerId = source['peerId'];
      this.cpu = source['cpu'];
      this.memory = source['memory'];
      this.systemImage = source['systemImage'];
      this.vmType = source['vmType'];
      this.creator = source['creator'];
      this.expireTime = this.convertValues(source['expireTime'], null);
      this.user = source['user'];
      this.status = source['status'];
    }

    convertValues(a: any, classs: any, asMap: boolean = false): any {
      if (!a) {
        return a;
      }
      if (a.slice) {
        return (a as any[]).map((elem) => this.convertValues(elem, classs));
      } else if ('object' === typeof a) {
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

export namespace app {
  export class Config {
    publicKey: string;
    port: number;
    peerId: string;
    wsUrl: string;

    static createFrom(source: any = {}) {
      return new Config(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.publicKey = source['publicKey'];
      this.port = source['port'];
      this.peerId = source['peerId'];
      this.wsUrl = source['wsUrl'];
    }
  }
}

export namespace wallet {
  export class WalletVo {
    address: string;
    addressJson: string;

    static createFrom(source: any = {}) {
      return new WalletVo(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.address = source['address'];
      this.addressJson = source['addressJson'];
    }
  }
}
