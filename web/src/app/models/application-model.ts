import {MinerModel} from './miner-model';

export class WhitelistApplication {
  id: number;
  createdAt: string;
  currentStatus: any;
  userId?: string;
  changeHistory?: any[];
  disabled?: Date;
  miner: MinerModel;

  constructor() {
    this.id = null;
    this.currentStatus = null;
    this.changeHistory = [];
    this.createdAt = '';
    this.userId = null;
    this.disabled = null;
    this.miner = null;
  }
}
