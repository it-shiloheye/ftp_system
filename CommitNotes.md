# Commit Notes

### 10th June 2024, 10:00 AM GMT +3
```sh
1. Able to read directories and list files
2. Able to upload files to database
3. Successfully reconnects on loss of connection
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   v3/CommitNotes.md
#	modified:   v3/lib/db_access/db_main.go
#	modified:   v3/lib/db_access/generated/upload_queries.sql.go
#	modified:   v3/lib/db_access/sql/queries/upload_queries.sql
#	modified:   v3/lib/logging/fake_logger.go
#	modified:   v3/peer/main.go
#	modified:   v3/peer/mainthread/db_helpers/connect_to_db.go
#	modified:   v3/peer/mainthread/download_file.go
#	modified:   v3/peer/mainthread/loop.go
#	modified:   v3/peer/mainthread/upload_file.go
#	modified:   v3/peer/mainthread/walk_directory.go
#	modified:   v3/peer/network-peer/api/list_files.go
#	modified:   v3/scripts/remove_executables.ps1
#	modified:   v3/scripts/test_storage.ps1
#
# Changes not staged for commit:
#	modified:   v1/client (modified content)
#	modified:   v1/lib (modified content)
#
```

### 06th June 2024, 22:18 PM GMT +3 
```sh
1. Set up v2 and v3 as part of main repo
2. Able to automatically upload and download file_basic
3. Will need to improve UI
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is ahead of 'origin/main' by 1 commit.
#   (use "git push" to publish your local commits)
#
# Changes to be committed:
#	modified:   .gitignore
#	deleted:    v2
#	new file:   v2/.gitignore
#	new file:   v2/CommitNotes.md
#	new file:   v2/LICENSE
#	new file:   v2/README.md
#	new file:   v2/client/.eslintrc.cjs
#	new file:   v2/client/.gitignore
#	new file:   v2/client/README.md
#	new file:   v2/client/index.html
#	new file:   v2/client/public/vite.svg
#	new file:   v2/client/src/App.css
#	new file:   v2/client/src/App.tsx
#	new file:   v2/client/src/assets/react.svg
#	new file:   v2/client/src/index.css
#	new file:   v2/client/src/main.tsx
#	new file:   v2/client/src/vite-env.d.ts
#	new file:   v2/client/vite.config.ts
#	new file:   v2/dev_kill_script.ps1
#	new file:   v2/dev_script.ps1
#	new file:   v2/go.mod
#	new file:   v2/go.sum
#	new file:   v2/init.example.sql
#	new file:   v2/lib/.air.toml
#	new file:   v2/lib/.gitignore
#	new file:   v2/lib/base/atomic_json.go
#	new file:   v2/lib/base/base.go
#	new file:   v2/lib/base/init.go
#	new file:   v2/lib/base/ip_handling.go
#	new file:   v2/lib/base/mutexed_map.go
#	new file:   v2/lib/base/mutexed_queue.go
#	new file:   v2/lib/context/context.go
#	new file:   v2/lib/db_access/db_main.go
#	new file:   v2/lib/db_access/generated/batch.go
#	new file:   v2/lib/db_access/generated/client_queries.sql.go
#	new file:   v2/lib/db_access/generated/db.go
#	new file:   v2/lib/db_access/generated/models.go
#	new file:   v2/lib/db_access/generated/peer_queries.sql.go
#	new file:   v2/lib/db_access/sql/queries/peer_queries.sql
#	new file:   v2/lib/db_access/sql/schema/peer_tables.sql
#	new file:   v2/lib/file_handler/v2/bytes_store.go
#	new file:   v2/lib/file_handler/v2/file_basic.go
#	new file:   v2/lib/file_handler/v2/file_hash.go
#	new file:   v2/lib/file_handler/v2/lock_file.go
#	new file:   v2/lib/logging/fake_logger.go
#	new file:   v2/lib/logging/log_item/error_type.go
#	new file:   v2/lib/logging/logging_struct.go
#	new file:   v2/lib/network_client/network_client.go
#	new file:   v2/lib/network_client/network_engine.go
#	new file:   v2/lib/sqlc.yaml
#	new file:   v2/lib/tls_handler/v2/cert_data.go
#	new file:   v2/lib/tls_handler/v2/cert_handler_2.go
#	new file:   v2/peer/.air.toml
#	new file:   v2/peer/browser-server/browser_server.go
#	new file:   v2/peer/config/config.go
#	new file:   v2/peer/config/data_storage_struct.go
#	new file:   v2/peer/main.go
#	new file:   v2/peer/main_thread/db_access/db_helpers.go
#	new file:   v2/peer/main_thread/main_thread.go
#	new file:   v2/peer/main_thread/storage_struct.go
#	new file:   v2/peer/main_thread/walk_dir.go
#	new file:   v2/peer/network-peer/network_peer.go
#	new file:   v2/peer/remove-item.ps1
#	new file:   v2/peer/server/init_server.go
#	new file:   v2/peer/server/server_loop.go
#	new file:   v2/peer/server/server_type.go
#	new file:   v2/plan.md
#	new file:   v2/postgres.bat
#	deleted:    v3
#	new file:   v3/.gitignore
#	new file:   v3/CommitNotes.md
#	new file:   v3/LICENSE
#	new file:   v3/client/.gitignore
#	new file:   v3/client/README.md
#	new file:   v3/client/index.html
#	new file:   v3/client/public/vite.svg
#	new file:   v3/client/src/App.css
#	new file:   v3/client/src/App.tsx
#	new file:   v3/client/src/assets/react.svg
#	new file:   v3/client/src/index.css
#	new file:   v3/client/src/main.tsx
#	new file:   v3/client/src/routeTree.gen.ts
#	new file:   v3/client/src/routes/index.tsx
#	new file:   v3/client/src/vite-env.d.ts
#	new file:   v3/client/tailwind.config.js
#	new file:   v3/client/vite.config.ts
#	new file:   v3/go.mod
#	new file:   v3/go.sum
#	new file:   v3/lib/.air.toml
#	new file:   v3/lib/.gitignore
#	new file:   v3/lib/base/atomic_json.go
#	new file:   v3/lib/base/base.go
#	new file:   v3/lib/base/init.go
#	new file:   v3/lib/base/ip_handling.go
#	new file:   v3/lib/base/mutexed_map.go
#	new file:   v3/lib/context/context.go
#	new file:   v3/lib/db_access/db_main.go
#	new file:   v3/lib/db_access/generated/db.go
#	new file:   v3/lib/db_access/generated/download_queries.sql.go
#	new file:   v3/lib/db_access/generated/models.go
#	new file:   v3/lib/db_access/generated/peer_queries.sql.go
#	new file:   v3/lib/db_access/generated/store_queries.sql.go
#	new file:   v3/lib/db_access/generated/upload_queries.sql.go
#	new file:   v3/lib/db_access/sql/peer_schema.sql
#	new file:   v3/lib/db_access/sql/queries/download_queries.sql
#	new file:   v3/lib/db_access/sql/queries/peer_queries.sql
#	new file:   v3/lib/db_access/sql/queries/store_queries.sql
#	new file:   v3/lib/db_access/sql/queries/upload_queries.sql
#	new file:   v3/lib/file_handler/v2/bytes_store.go
#	new file:   v3/lib/file_handler/v2/file_basic.go
#	new file:   v3/lib/file_handler/v2/file_hash.go
#	new file:   v3/lib/file_handler/v2/lock_file.go
#	new file:   v3/lib/logging/fake_logger.go
#	new file:   v3/lib/logging/log_item/error_type.go
#	new file:   v3/lib/logging/logging_struct.go
#	new file:   v3/lib/network_client/network_client.go
#	new file:   v3/lib/network_client/network_engine.go
#	new file:   v3/lib/sqlc.yaml
#	new file:   v3/lib/tls_handler/v2/cert_data.go
#	new file:   v3/lib/tls_handler/v2/cert_handler_2.go
#	new file:   v3/peer/.air.toml
#	new file:   v3/peer/build_and_run.ps1
#	new file:   v3/peer/config/config.go
#	new file:   v3/peer/config/data_storage_struct.go
#	new file:   v3/peer/main.go
#	new file:   v3/peer/mainthread/db_helpers/connect_to_db.go
#	new file:   v3/peer/mainthread/db_helpers/upload.go
#	new file:   v3/peer/mainthread/download_file.go
#	new file:   v3/peer/mainthread/file_map_type.go
#	new file:   v3/peer/mainthread/loop.go
#	new file:   v3/peer/mainthread/upload_file.go
#	new file:   v3/peer/mainthread/walk_directory.go
#	new file:   v3/peer/network-peer/network_peer.go
#	new file:   v3/peer/remove-item.ps1
#	new file:   v3/peer/server/init_server.go
#	new file:   v3/peer/server/server_loop.go
#	new file:   v3/peer/server/server_type.go
#	new file:   v3/plan.md
#	new file:   v3/scripts/dev_kill_script.ps1
#	new file:   v3/scripts/dev_script.ps1
#	new file:   v3/scripts/remove_executables.ps1
#	new file:   v3/scripts/remove_filemap_json.ps1
#	new file:   v3/scripts/remove_locks.ps1
#	new file:   v3/scripts/test_storage.ps1
#	new file:   v3/test.ps1
#
# Changes not staged for commit:
#	modified:   v1/client (modified content)
#	modified:   v1/lib (modified content)
#
```

### 06th June 2024, 22:00 PM GMT +3
```sh
setting up parent repository
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   .gitignore
#	modified:   CommitNotes.md
#	deleted:    client/.air.toml
#	deleted:    client/.gitignore
#	deleted:    client/CommitNotes.md
#	deleted:    client/LICENSE
#	deleted:    client/README.md
#	deleted:    client/client_config.schema.json
#	deleted:    client/go.mod
#	deleted:    client/go.sum
#	deleted:    client/init_client/config.go
#	deleted:    client/init_client/struct_options.go
#	deleted:    client/main.go
#	deleted:    client/main_thread/dir_handler/file_tree_json.go
#	deleted:    client/main_thread/dir_handler/read_files_in_directory.go
#	deleted:    client/main_thread/init.go
#	deleted:    client/main_thread/main_thread.go
#	deleted:    client/main_thread/network_client/init.go
#	deleted:    client/main_thread/network_client/network_client.go
#	deleted:    client/main_thread/network_client/network_engine.go
#	deleted:    client/main_thread/utils.go
#	deleted:    lib
#	deleted:    server/main_thread/gin_server/file_upload_routes.go
#	deleted:    server/remove-item.ps1
#	deleted:    utils/.air.toml
#	deleted:    utils/.gitignore
#	deleted:    utils/main.go
#	new file:   v1/.gitignore
#	new file:   v1/CommitNotes.md
#	renamed:    LICENSE -> v1/LICENSE
#	new file:   v1/client
#	renamed:    dev_kill_script.ps1 -> v1/dev_kill_script.ps1
#	renamed:    dev_script.ps1 -> v1/dev_script.ps1
#	new file:   v1/lib
#	renamed:    server/.air.toml -> v1/server/.air.toml
#	renamed:    server/.gitignore -> v1/server/.gitignore
#	renamed:    server/config.schema -> v1/server/config.schema
#	renamed:    server/go.mod -> v1/server/go.mod
#	renamed:    server/go.sum -> v1/server/go.sum
#	renamed:    server/initialise_server/config.go -> v1/server/initialise_server/config.go
#	renamed:    server/main.go -> v1/server/main.go
#	renamed:    server/main_thread/actions/file_operations.go -> v1/server/main_thread/actions/file_operations.go
#	renamed:    server/main_thread/dir_handler/dir_handler_type.go -> v1/server/main_thread/dir_handler/dir_handler_type.go
#	renamed:    server/main_thread/dir_handler/utils.go -> v1/server/main_thread/dir_handler/utils.go
#	new file:   v1/server/main_thread/gin_server/file_upload_routes.go
#	renamed:    server/main_thread/gin_server/init_server.go -> v1/server/main_thread/gin_server/init_server.go
#	renamed:    server/main_thread/gin_server/register_routes.go -> v1/server/main_thread/gin_server/register_routes.go
#	renamed:    server/main_thread/gin_server/utils_func.go -> v1/server/main_thread/gin_server/utils_func.go
#	renamed:    server/main_thread/init.go -> v1/server/main_thread/init.go
#	renamed:    client/remove-item.ps1 -> v1/server/remove-item.ps1
#	new file:   v2
#	new file:   v3
#
# Changes not staged for commit:
#	modified:   v1/client (modified content)
#	modified:   v1/lib (modified content)
#
```

### 19th June 2024, 12:51 PM GMT +3
```sh
1. Uploading files to server
2. Storing files serverside with hash
3. bulk uploads
    - upload 10 files at once
4. bulk confirms
    - upload filetree to server, 
        server responds with missing files

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   client/.air.toml
#	modified:   client/CommitNotes.md
#	modified:   client/go.mod
#	modified:   client/go.sum
#	modified:   client/init_client/config.go
#	modified:   client/main.go
#	deleted:    client/main_thread/actions/file_operations.go
#	modified:   client/main_thread/dir_handler/file_tree_json.go
#	modified:   client/main_thread/dir_handler/read_files_in_directory.go
#	deleted:    client/main_thread/logging/logging_struct.go
#	modified:   client/main_thread/main_thread.go
#	deleted:    client/main_thread/network_client/client_engine.go
#	modified:   client/main_thread/network_client/init.go
#	deleted:    client/main_thread/network_client/make_get_request.go
#	deleted:    client/main_thread/network_client/make_post_request.go
#	modified:   client/main_thread/network_client/network_client.go
#	new file:   client/main_thread/network_client/network_engine.go
#	deleted:    client/main_thread/network_client/read_json_from_response.go
#	new file:   client/main_thread/utils.go
#	modified:   lib
#	modified:   server/.air.toml
#	modified:   server/go.mod
#	modified:   server/go.sum
#	modified:   server/initialise_server/config.go
#	modified:   server/main.go
#	modified:   server/main_thread/actions/file_operations.go
#	new file:   server/main_thread/dir_handler/utils.go
#	new file:   server/main_thread/gin_server/file_upload_routes.go
#	deleted:    server/main_thread/gin_server/handle_file_uploads.go
#	modified:   server/main_thread/gin_server/init_server.go
#	new file:   server/main_thread/gin_server/register_routes.go
#	new file:   server/main_thread/gin_server/utils_func.go
#	deleted:    server/main_thread/logging/logging_struct.go
#
# Changes not staged for commit:
#	modified:   lib (modified content)
#

```

### 13th June 2024, 19:17 PM GMT +3
```sh
1. 1st attempt to save files on server
2. 3 steps:
    - send file to server
    - save file server side
    - respond to confirm or deny presence of file
3. client-side:
    - read all files "included"
    - upload files where necessary
    - confirm presence of files
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is ahead of 'origin/main' by 3 commits.
#   (use "git push" to publish your local commits)
#
# Changes to be committed:
#	modified:   CommitNotes.md
#	modified:   client/.air.toml
#	modified:   client/CommitNotes.md
#	modified:   client/go.sum
#	modified:   client/init_client/struct_options.go
#	modified:   client/main.go
#	modified:   client/main_thread/dir_handler/file_tree_json.go
#	modified:   client/main_thread/logging/logging_struct.go
#	modified:   client/main_thread/main_thread.go
#	modified:   server/.air.toml
#	deleted:    server/gin_server/handle_file_uploads.go
#	deleted:    server/gin_server/tmp_file_hold.go
#	modified:   server/go.mod
#	modified:   server/main.go
#	modified:   server/main_thread/dir_handler/dir_handler_type.go
#	new file:   server/main_thread/gin_server/handle_file_uploads.go
#	renamed:    server/gin_server/init_server.go -> server/main_thread/gin_server/init_server.go
#	modified:   server/main_thread/logging/logging_struct.go
#	deleted:    server/main_thread/main_thread.go
#
```

### 13th June 2024, 10:05 AM GMT +3
```txt
1. Able to read all files in a directory
2. Able to list loaded files
3. Saves loaded extensions
4. Regularly saves progress to file-tree.json
5. Server successfull receives files from client
Pending:
6. Reduce clientside memory use
    - only read when hashing and uploading 
        (no dangling filehandlers)
7. Reduce serverside memory use
    - save directly to disk on upload, 
        retaining only file address and info
8. Subscribe and download
    - clientside: 
        1. send subscriptions to server 
        (to notify on update/changes to 
            directory/files)
    - serverside: 
        1. track which clients have which files
        2. track which clients need which files
        3. push to clients on change
            (eg.: "/download/changes" route)
9. Load balancing
    - simple round robin queue
```

### 08th June 2024, 21:27 PM GMT+3
```sh
1. Set up ftp_system/client as depenedency of 
    server
2. Able to transmit files from client to server
3. Need to work on:
    - storing files serverside
    - fetching files from server to client
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   client/.gitignore
#	new file:   client/CommitNotes.md
#	new file:   client/LICENSE
#	new file:   client/README.md
#	modified:   client/go.mod
#	modified:   client/main.go
#	modified:   client/main_thread/dir_handler/read_files_in_directory.go
#	modified:   client/main_thread/logging/logging_struct.go
#	modified:   client/main_thread/main_thread.go
#	modified:   client/main_thread/network_client/client_engine.go
#	modified:   client/main_thread/network_client/make_get_request.go
#	modified:   client/main_thread/network_client/make_post_request.go
#	modified:   client/main_thread/network_client/network_client.go
#	modified:   client/main_thread/network_client/read_json_from_response.go
#	modified:   server/gin_server/init_server.go
#	modified:   server/go.mod
#	modified:   server/go.sum
#	modified:   server/initialise_server/config.go
#	modified:   server/main_thread/logging/logging_struct.go
#

```

### 06th June 2024, 21:13 PM GMT+3
    1. Fixed all logging errors
    2. Moved std out and file logging to single thread
    3. Need to work on hashing

### 06th June 2024, 19:24 PM GMT+3
    1. Adding Logger struct
    2. Adding os.PathSeparator to config and excluded
    3. Improved FileBasic in lib to expose filehandler
    4. able to store and load file-tree.json (persistent state)
    5. Able to read files in directory