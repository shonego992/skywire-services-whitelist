import {NodeKey} from './node-keys';

export class ApplicationWhitelistReq {
  description?: string;
  location?: string;
  nodes: NodeKey[] | string;
  oldImages: any;

  constructor() {
    this.description = '';
    this.location = '';
    this.nodes = [];
    this.oldImages = [];
  }
}
