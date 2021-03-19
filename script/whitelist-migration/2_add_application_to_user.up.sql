CREATE TABLE applications (
  id              serial primary key,
  current_status  smallint not null,
  created_at      timestamp not null,
  updated_at      timestamp not null,
  deleted_at      timestamp null,
  username        varchar(255) not null,
  foreign key     (username) references users (username)
);

CREATE TABLE change_history (
  id              serial primary key,
  status          smallint not null,
  description     varchar,
  location        varchar,
  created_at      timestamp not null,
  updated_at      timestamp not null,
  deleted_at      timestamp null,
  comment         varchar,
  application_id  integer,
  foreign key     (application_id) references applications (id)
);

CREATE TABLE nodes (
  id                      serial primary key,
  key                     varchar,
  created_at              timestamp not null,
  updated_at              timestamp not null,
  deleted_at              timestamp null
);

CREATE TABLE images (
  id                      serial primary key,
  path                    varchar,
  created_at              timestamp not null,
  updated_at              timestamp not null,
  deleted_at              timestamp null
);

CREATE TABLE change_nodes (
  change_history_id integer,
  node_id integer
);

CREATE TABLE change_images (
  change_history_id integer,
  image_id integer
);
