import {MinerNodeModel} from './miner-node-model';

export class MinerModel {
  id: number;
  createdAt: string;
  type: number;
  username: string;
  nodes: MinerNodeModel[];
  

  constructor() {
    this.id = null;
    this.createdAt = '';
    this.type = null;
    this.username = '';
    this.nodes = [];
  }
}
