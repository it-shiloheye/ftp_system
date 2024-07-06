
CREATE TYPE file_status_type AS ENUM (
    'created',
    'modified',
    'deleted'
);


CREATE TYPE peer_role_type AS ENUM (
    -- push changes to database
    -- @ uploads file to database
    -- @ fetch files through database
    'client',
    -- create a local copy of all files
    -- @ store files via hash (deleting data from db)
    -- @ upload files to database on request
    'storage',
    -- respond to rpc requests
    -- @ 
    'server'
);


-- register basic shared info of users
CREATE TABLE IF NOT EXISTS peers_table (
    id serial PRIMARY KEY,
    peer_id uuid DEFAULT gen_random_uuid() UNIQUE,
    ip_address TEXT NOT NULL,
    peer_role peer_role_type,
    peer_name TEXT UNIQUE,
    creation_time timestamptz default NOW(),
    pem BYTEA,
    peer_config jsonb
);

-- track dir paths of each peer for storage
CREATE TABLE IF NOT EXISTS peer_dirs (
    id serial primary key,
    peer_id uuid references peers_table(peer_id),
    creation_time timestamptz default NOW(),
    dir_path TEXT NOT NULL 
);

-- track file data specifically
CREATE TABLE IF NOT EXISTS file_data (
    file_data_id serial primary key,
    file_hash   VARCHAR(256) GENERATED ALWAYS AS (encode(sha256(file_data::bytea), 'hex')) STORED UNIQUE,
    prev_file_hash integer references file_data(id),
    file_status file_status_type,
    modification_date TIMESTAMP NOT NULL,
    creation_time timestamptz default NOW(),
    file_data BYTEA
);

-- track individual metadata for files and filehash
CREATE TABLE IF NOT EXISTS file_metadata (
    id serial PRIMARY KEY,
    peer_id uuid references peers_table(peer_id),
    dir_id integer references peer_dirs(id),
    file_data_id references file_data(id),
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_type VARCHAR(7) NOT NULL,
    creation_day DATE GENERATED ALWAYS AS (EXTRACT(DATE  from creation_time::timestamptz)) STORED,
    creation_time timestamptz default NOW()
);


CREATE TABLE IF NOT EXISTS file_transf_msg_table (
    id serial primary key,
    peer_reqing uuid references peers_table(peer_id) NOT NULL,
    file_hash integer references file_data(id) NOT NULL,
    req_time timestamptz default NOW(),
    peer_responding uuid references peers_table(peer_id),
    res_time timestamptz
);


CREATE TABLE IF NOT EXISTS file_change_log (
    id serial primary key,
    peer_id uuid references peers_table(peer_id),
    prev_file_state file_status_type,
    curr_file_state file_status_type,
    file_hash_id integer references file_data(id),
    messagef text not null,
    json_log jsonb
);