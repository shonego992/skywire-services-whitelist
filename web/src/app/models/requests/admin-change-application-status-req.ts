export class ChangeApplicationStatusReq {
  applicationId: number;
  status: number;
  userComment: string;
  adminComment: string;

  constructor() {
    this.applicationId = null;
    this.status = null;
    this.adminComment = '';
    this.userComment = '';
  }
}
