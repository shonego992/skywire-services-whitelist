export class AdminClaims {
  flag_vip: boolean;
  create_user: boolean;
  disable_user: boolean;
  review_whitelist: boolean;
  missing_confirmation: boolean;

  constructor() {
    this.flag_vip = false;
    this.create_user = false;
    this.disable_user = false;
    this.review_whitelist = false;
    this.missing_confirmation = false;
  }
}
