CREATE TABLE miners (
  id              serial primary key,
  username        varchar(255) not null,
  created_at      timestamp not null,
  updated_at      timestamp not null,
  deleted_at      timestamp null,
  type            smallint not null,
  foreign key     (username) references users (username)
);

ALTER TABLE nodes
ADD COLUMN miner_id integer null;

-- ALTER TABLE nodes 
--    ADD CONSTRAINT fk_minerId
--    FOREIGN KEY (miner_id) 
--    REFERENCES miners (id);

CREATE TABLE miner_imports (
  id                  serial primary key,
  username            varchar(255) not null,
  created_at          timestamp not null,
  updated_at          timestamp not null,
  deleted_at          timestamp null,
  number_of_miners    integer
);
