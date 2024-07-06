

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
    peer_role peer_role_type default 'client',
    peer_name TEXT UNIQUE,
    creation_time timestamptz default NOW(),
    pem BYTEA,
    peer_config jsonb
);

-- track dir paths of each peer for storage
-- ie.: "uploaded_dirs"
CREATE TABLE IF NOT EXISTS peer_dirs (
    id serial primary key,
    peer_id uuid references peers_table(peer_id),
    creation_time timestamptz default NOW(),
    dir_path TEXT NOT NULL 
);

CREATE UNIQUE INDEX uniq_dirs
ON peer_dirs(peer_id, dir_path);



-- track metadata for files 
CREATE TABLE IF NOT EXISTS file_metadata (
    id serial PRIMARY KEY,
    dir_id integer references peer_dirs(id),
    file_path TEXT NOT NULL UNIQUE, -- relative
    file_type VARCHAR(7) NOT NULL,
    file_state text not null, 
    file_data_id integer references file_data(id) not null,
    file_mode integer not null,
    mod_time timestamptz not null,
    creation_time timestamptz default NOW()
);

CREATE TABLE IF NOT EXISTS file_data (
    id serial primary key,
    file_hash    VARCHAR(256) GENERATED ALWAYS AS (encode(sha256(file_data_b::bytea), 'hex')) STORED UNIQUE,
    mod_time timestamptz not null,
    file_size integer not null,
    file_data_b BYTEA,
    creation_time timestamptz default NOW()
);


CREATE TABLE IF NOT EXISTS file_tracker (
    id serial primary key,
    peer_id uuid references peers_table(peer_id) not null,
    file_meta_id integer references file_metadata(id) not null,
    current_hash_id integer references file_data(id) not null,
    file_state text not null,
    log_time timestamptz default NOW()
);


GRANT ALL PRIVILEGES ON DATABASE ftp_system_db TO ftp_system_server;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ftp_system_server;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO ftp_system_server;