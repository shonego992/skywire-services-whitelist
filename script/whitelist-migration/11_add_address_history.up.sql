CREATE TABLE addresses (
  id              serial primary key,
  skycoin_address varchar,
  created_at      timestamp not null,
  updated_at      timestamp not null,
  deleted_at      timestamp null,
  username        varchar(255) not null,
  foreign key     (username) references users (username)
); 