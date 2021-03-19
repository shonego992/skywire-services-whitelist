import {UptimeModel} from './uptime.model';

export class WhitelistNodeModel {
  id: number;
  key: string;
  uptime: UptimeModel;

  constructor() {
    this.id = null;
    this.key = '';
    this.uptime = null;
  }
}
