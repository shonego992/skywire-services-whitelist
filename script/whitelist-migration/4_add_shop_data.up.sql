CREATE TABLE shop_data (
  id              numeric primary key,
  status          varchar(255) not null,
  created_at      timestamp not null,
  updated_at      timestamp not null,
  deleted_at      timestamp null
);