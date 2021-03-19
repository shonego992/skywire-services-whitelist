CREATE TABLE miner_transfers (
  id              serial primary key,
  created_at      timestamp not null,
  updated_at      timestamp not null,
  deleted_at      timestamp null,
  old_username    varchar(255) not null,
  new_username    varchar(255) not null,
  miner_id        integer,
  foreign key     (old_username) references users (username),
  foreign key     (new_username) references users (username),
  foreign key     (miner_id)     references miners (id)
);

-- dropping the miner_import_id
-- no revert action should be performed for this step
ALTER TABLE nodes
DROP COLUMN IF EXISTS miner_import_id;