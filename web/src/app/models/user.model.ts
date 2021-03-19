import {WhitelistApplication} from './application-model';
import {SkycoinAddressModel} from './skycoin-address.model';

export interface User {
  id: number;
  status: number;
  username: string;
  skycoinAddress?: string;
  applications: WhitelistApplication[];
  rights?: string[];
  disabled?: Date;
  createdAt: string;
  addressHistory?: SkycoinAddressModel[];
}
