// Shared variables to be used in application

export class ApiRoutes {

  public static USER = {
        'Users' : '/users',
        'Keys': '/users/keys',
        'Address' : '/users/address'
  };

  public static WHITELIST = {
    'Application': '/whitelist/application',
    'UpdateApplication': '/whitelist/updateApplication',
    'ApplicationNoImages': '/whitelist/updateApplicationNoImageChange',
    'Whitelists': '/whitelist/whitelists',
    'Whitelist': '/whitelist/whitelist',
    'Miners': '/miners/miners',
    'AdminAllMiners': '/miners/allMiners',
    'Miner': '/miners/miner',
    'MinerForAdmin': '/miners/minerForAdmin',
    'UserList': '/miners/uploadUserList',
    'Import': '/miners/import',
    'ProcessImport': '/miners/import/process',
    'GetMinersForUser': '/miners/minersForUser',
    'ExportMiners': '/miners/exportMiners',
    'ExportMinersNoLimitations': '/miners/exportMinersNoLimitations',
    'TransferMiner': '/miners/transferMiner'
  };

  public static AUTH = {
    'Refresh': '/auth/refresh',
    'Info': '/info',
    'Login': '/auth/login'
  };

  public static INFO = {
    'UptimeInfo': '/info/getNodeInfo'
  };

  public static ADMIN = {
    'Users': '/admin/users',
    'Admins': '/admin/admins',
    'DisableUser': '/admin/disableUser',
    'EnableUser': '/admin/enableUser'
  };
}
