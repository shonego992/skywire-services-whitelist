CREATE TABLE export_records (
  id              serial primary key,
  payout_address  varchar(255) not null,
  miner_type      smallint not null,
  official_tx     varchar(255) not null,
  diy_tx          varchar(255) not null,
  correction_tx   varchar(255) not null,
  time_of_export  timestamp not null,
  created_at      timestamp not null,
  updated_at      timestamp not null,
  deleted_at      timestamp null,
  username        varchar(255) not null,
  number_of_nodes integer,
  foreign key     (username) references users (username)
);