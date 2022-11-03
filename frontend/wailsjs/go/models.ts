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
  export class Staking {
    networkUrl: string;
    address: string;
    agentAddress: string;
    pledgeAmount: number;

    static createFrom(source: any = {}) {
      return new Staking(source);
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
    publicKey: string;
    accountMnemonic: string;

    static createFrom(source: any = {}) {
      return new Initialization(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.leaseTerm = source['leaseTerm'];
      this.publicKey = source['publicKey'];
      this.accountMnemonic = source['accountMnemonic'];
    }
  }
  export class DeployParameter {
    // Go type: Initialization
    initialization: any;
    // Go type: Staking
    staking: any;
    // Go type: Deployment
    deployment: any;

    static createFrom(source: any = {}) {
      return new DeployParameter(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.initialization = this.convertValues(source['initialization'], null);
      this.staking = this.convertValues(source['staking'], null);
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
}

export namespace v2 {
  export class AddParam {
    name: string;
    selectNodeType: string;
    leaseTerm: number;
    thegraphIndexer: string;
    stakingAmount: number;

    static createFrom(source: any = {}) {
      return new AddParam(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.name = source['name'];
      this.selectNodeType = source['selectNodeType'];
      this.leaseTerm = source['leaseTerm'];
      this.thegraphIndexer = source['thegraphIndexer'];
      this.stakingAmount = source['stakingAmount'];
    }
  }
  export class AddApplicationVo {
    id: number;
    result: boolean;

    static createFrom(source: any = {}) {
      return new AddApplicationVo(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.id = source['id'];
      this.result = source['result'];
    }
  }
}

export namespace application {
  export class ListVo {
    id: number;
    name: string;
    selectNodeType: string;
    status: number;
    grtIncome: number;

    static createFrom(source: any = {}) {
      return new ListVo(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.id = source['id'];
      this.name = source['name'];
      this.selectNodeType = source['selectNodeType'];
      this.status = source['status'];
      this.grtIncome = source['grtIncome'];
    }
  }
  export class PageApplicationVo {
    items: ListVo[];
    total: number;

    static createFrom(source: any = {}) {
      return new PageApplicationVo(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.items = this.convertValues(source['items'], ListVo);
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
    name: string;
    selectNodeType: string;
    status: number;
    leaseTerm: number;
    p2pForwardPort: number;
    cliForwardPort: number;
    peerId: string;
    orderIndex: number;
    thinkingTime: number;

    static createFrom(source: any = {}) {
      return new ApplyVo(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.id = source['id'];
      this.createdAt = this.convertValues(source['createdAt'], null);
      this.name = source['name'];
      this.selectNodeType = source['selectNodeType'];
      this.status = source['status'];
      this.leaseTerm = source['leaseTerm'];
      this.p2pForwardPort = source['p2pForwardPort'];
      this.cliForwardPort = source['cliForwardPort'];
      this.peerId = source['peerId'];
      this.orderIndex = source['orderIndex'];
      this.thinkingTime = source['thinkingTime'];
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
    selectNodeType: string;

    static createFrom(source: any = {}) {
      return new UpdateApplicationParam(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.id = source['id'];
      this.name = source['name'];
      this.selectNodeType = source['selectNodeType'];
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
    plugin: string;
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
      this.plugin = source['plugin'];
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

    static createFrom(source: any = {}) {
      return new Account(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.publicKey = source['publicKey'];
      this.wsUrl = source['wsUrl'];
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
  export class GraphRulesInfo {
    info: GraphRule[];

    static createFrom(source: any = {}) {
      return new GraphRulesInfo(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.info = this.convertValues(source['info'], GraphRule);
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
  export class QueueInfo {
    info: StatusInfo[];

    static createFrom(source: any = {}) {
      return new QueueInfo(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.info = this.convertValues(source['info'], StatusInfo);
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
