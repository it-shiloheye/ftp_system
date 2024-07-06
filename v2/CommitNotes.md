# Commit Notes

### 27th June 2024, 23:12 PM, GMT +3
```sh
1. Restructuring StorageStruct to:
    - store more options
    - function and update within the loop rather than as a global object
2. Restructure funcs:
    - WalkDir : moved to it's own file
    - set_up_storage: can now update saved config 
    - main loop now reloads the config on every 
        cycle
    - upload_func: able to more accurately seggregate what file to 
        upload
    - peer roles are saved both in db and locally 
3. Next add:
    - download_func to download files and updates to client and save 
    - update_func that checks if file hash has changed and updates prev
        saves 
    - storage_func that activates if storage func to store all peer 
        files and mark stored 
    - delete func that:
        - deletes stale files
        - marks locally deleted files as deleted on db 
4. Next then build:
    - React-based ui to allow:
        - easy configuration and searching directly in the browser
        - search other PCs and download/upload/send files
        - delete files easily
        - view change logs and history

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   CommitNotes.md
#	modified:   lib/db_access/generated/client_queries.sql.go
#	modified:   lib/db_access/generated/models.go
#	modified:   lib/db_access/sql/queries/client_queries.sql
#	modified:   lib/db_access/sql/schema/client_id.sql
#	modified:   package.json
#	modified:   peer/config/config.go
#	modified:   peer/config/data_storage_struct.go
#	modified:   peer/main.go
#	modified:   peer/main_thread/db_access/db_helpers.go
#	modified:   peer/main_thread/main_thread.go
#	new file:   peer/main_thread/storage_struct.go
#	new file:   peer/main_thread/walk_dir.go
#	modified:   peer/server/server_type.go
#
```

### 26th June 2024, 09:29 AM, GMT +3
```sh
1. Able to walk dir
2. Able to get files info from database
3. Able to upload files into database and hash
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   lib/db_access/db_main.go
#	modified:   lib/db_access/generated/client_queries.sql.go
#	modified:   lib/db_access/generated/models.go
#	modified:   lib/db_access/sql/queries/client_queries.sql
#	modified:   lib/db_access/sql/schema/client_id.sql
#	modified:   peer/config/config.go
#	modified:   peer/config/data_storage_struct.go
#	modified:   peer/main.go
#	modified:   peer/main_thread/db_access/db_helpers.go
#	modified:   peer/main_thread/main_thread.go
#	modified:   peer/network-peer/network_peer.go
#	modified:   peer/server/server_loop.go
#
```

### 25th June 2024, 08:16 AM, GMT +3
```sh
1. Database permissions sorted
2. Able to create and store CA_pem in db 
3. Next step is:
    - to upload files into db
    - work on config struct:
        - store in db
        - access throw browser
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   .air.toml
#	modified:   db_access/db_main.go
#	modified:   db_access/generated/client_queries.sql.go
#	modified:   db_access/generated/db.go
#	modified:   db_access/generated/models.go
#	modified:   db_access/sql/queries/client_queries.sql
#	modified:   db_access/sql/schema/client_id.sql
#	modified:   sqlc.yaml
#	modified:   tls_handler/v2/cert_handler_2.go
#
# Changes not staged for commit:
#	modified:   ../CommitNotes.md
#	modified:   ../go.mod
#	modified:   ../go.sum
#	modified:   ../peer/.air.toml
#	modified:   ../peer/browser-server/browser_server.go
#	modified:   ../peer/config/config.go
#	modified:   ../peer/main.go
#	modified:   ../peer/network-peer/network_peer.go
#	modified:   ../peer/server/init_server.go
#	modified:   ../peer/server/server_loop.go
#	modified:   ../peer/server/server_type.go
#	deleted:    ../sqlc.yaml
#
# Untracked files:
#	../peer/config/data_storage_struct.go
#	../peer/main_thread/
#
```

### 23rd June 2024, 19:25 PM GMT +3
```sh
1. Peer Server created successfully
2. Browser Server created successfully
3. Three steps:
    => connect to db
    => set up logging
    => set up browser and peer servers
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is ahead of 'origin/main' by 1 commit.
#   (use "git push" to publish your local commits)
#
# Changes to be committed:
#	new file:   peer/browser-server/browser_server.go
#	modified:   peer/main.go
#	modified:   peer/network-peer/network_peer.go
#	modified:   peer/server/init_server.go
#	modified:   peer/server/server_loop.go
#	modified:   peer/server/server_type.go
#
```

### 23rd June 2024, 18:34 PM GMT +3
```sh
1. Rewriting server and logging 
2. Rewriting db queries, to have single main table
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   lib/base/ip_handling.go
#	modified:   lib/db_access/generated/client_queries.sql.go
#	modified:   lib/db_access/generated/models.go
#	modified:   lib/db_access/sql/schema/client_id.sql
#	modified:   lib/logging/logging_struct.go
#	modified:   lib/sqlc.yaml
#	modified:   package.json
#	renamed:    server-peer/.air.toml -> peer/.air.toml
#	new file:   peer/config/config.go
#	renamed:    server-peer/main.go -> peer/main.go
#	new file:   peer/network-peer/network_peer.go
#	renamed:    server-peer/remove-item.ps1 -> peer/remove-item.ps1
#	renamed:    server-peer/server/init_server.go -> peer/server/init_server.go
#	renamed:    server-peer/server/server_loop.go -> peer/server/server_loop.go
#	new file:   peer/server/server_type.go
#	deleted:    server-peer/config/config.go
#	deleted:    server-peer/server/register_routes.go
#

```

### 22nd June 2024, 23:59 PM GMT +3
```sh
1. Shifting from server-client to peer-db-peer 
2. Using React-SWC-TS as browser frontend
3. Finalising connection to db
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   .gitignore
#	deleted:    _lib/.gitignore
#	deleted:    _lib/db_access/db_main.go
#	deleted:    _lib/db_access/sql/init.sql
#	deleted:    _lib/exports.go
#	new file:   client/.eslintrc.cjs
#	new file:   client/.gitignore
#	new file:   client/README.md
#	new file:   client/index.html
#	deleted:    client/main.go
#	new file:   client/public/vite.svg
#	new file:   client/src/App.css
#	new file:   client/src/App.tsx
#	new file:   client/src/assets/react.svg
#	new file:   client/src/index.css
#	new file:   client/src/main.tsx
#	new file:   client/src/vite-env.d.ts
#	new file:   client/vite.config.ts
#	modified:   dev_kill_script.ps1
#	modified:   dev_script.ps1
#	modified:   go.mod
#	modified:   go.sum
#	new file:   init.example.sql
#	renamed:    _lib/.air.toml -> lib/.air.toml
#	new file:   lib/.gitignore
#	renamed:    _lib/base/atomic_json.go -> lib/base/atomic_json.go
#	renamed:    _lib/base/base.go -> lib/base/base.go
#	renamed:    _lib/base/init.go -> lib/base/init.go
#	new file:   lib/base/ip_handling.go
#	renamed:    _lib/base/mutexed_map.go -> lib/base/mutexed_map.go
#	renamed:    _lib/base/mutexed_queue.go -> lib/base/mutexed_queue.go
#	renamed:    _lib/context/context.go -> lib/context/context.go
#	new file:   lib/db_access/db_main.go
#	new file:   lib/db_access/generated/client_queries.sql.go
#	new file:   lib/db_access/generated/db.go
#	new file:   lib/db_access/generated/models.go
#	new file:   lib/db_access/sql/queries/client_queries.sql
#	new file:   lib/db_access/sql/schema/client_id.sql
#	renamed:    _lib/file_handler/v2/bytes_store.go -> lib/file_handler/v2/bytes_store.go
#	renamed:    _lib/file_handler/v2/file_basic.go -> lib/file_handler/v2/file_basic.go
#	renamed:    _lib/file_handler/v2/file_hash.go -> lib/file_handler/v2/file_hash.go
#	renamed:    _lib/file_handler/v2/lock_file.go -> lib/file_handler/v2/lock_file.go
#	renamed:    _lib/logging/fake_logger.go -> lib/logging/fake_logger.go
#	renamed:    _lib/logging/log_item/error_type.go -> lib/logging/log_item/error_type.go
#	renamed:    _lib/logging/logging_struct.go -> lib/logging/logging_struct.go
#	renamed:    client/network_client/network_client.go -> lib/network_client/network_client.go
#	renamed:    client/network_client/network_engine.go -> lib/network_client/network_engine.go
#	new file:   lib/sqlc.yaml
#	renamed:    _lib/tls_handler/v2/cert_data.go -> lib/tls_handler/v2/cert_data.go
#	renamed:    _lib/tls_handler/v2/cert_handler_2.go -> lib/tls_handler/v2/cert_handler_2.go
#	modified:   package.json
#	renamed:    server/.air.toml -> server-peer/.air.toml
#	renamed:    server/config/config.go -> server-peer/config/config.go
#	new file:   server-peer/main.go
#	renamed:    client/remove-item.ps1 -> server-peer/remove-item.ps1
#	renamed:    server/server/init_server.go -> server-peer/server/init_server.go
#	renamed:    server/server/register_routes.go -> server-peer/server/register_routes.go
#	renamed:    server/server/server_loop.go -> server-peer/server/server_loop.go
#	deleted:    server/main.go
#	deleted:    server/remove-item.ps1
#
```