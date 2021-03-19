CREATE TABLE users (
  username    varchar(255) not null primary key,
  status      smallint default 0,
  created_at  timestamp not null,
  updated_at  timestamp not null,
  deleted_at  timestamp null,
  skycoin_address varchar(100)
);

CREATE TABLE api_keys (
  id          serial primary key,
  key         varchar(40) not null UNIQUE,
  created_at  timestamp not null,
  updated_at  timestamp not null,
  deleted_at  timestamp null,
  username    varchar(255) not null,
  foreign key (username) references users (username)
);