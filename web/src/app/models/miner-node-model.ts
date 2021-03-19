import {UptimeModel} from './uptime.model';

export class MinerNodeModel {
  id: number;
  key: string;
  uptime: UptimeModel[];


  constructor() {
    this.id = null;
    this.key = '';
    this.uptime = [];
  }
}
