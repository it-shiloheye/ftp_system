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

### 09th July 2024, 10:10 AM GMT +3
```sh
1. Able to create a properly named build
2. Client can query peer for list of files
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Your branch is up to date with 'origin/main'.
#
# Changes to be committed:
#	modified:   .gitignore
#	modified:   CommitNotes.md
#	new file:   v3/build_script.js
#	modified:   v3/client/index.html
#	new file:   v3/client/postcss.config.js
#	deleted:    v3/client/src/App.tsx
#	new file:   v3/client/src/api/fetch_files_list.ts
#	new file:   v3/client/src/components/main_body_wrapper.tsx
#	modified:   v3/client/src/index.css
#	modified:   v3/client/src/main.tsx
#	modified:   v3/client/src/routeTree.gen.ts
#	new file:   v3/client/src/routes/__root.tsx
#	new file:   v3/client/src/routes/about.lazy.tsx
#	modified:   v3/client/src/routes/index.tsx
#	new file:   v3/client/src/styles/base.general.css
#	new file:   v3/client/src/styles/base.tailwind.css
#	new file:   v3/client/src/styles/styles.css
#	modified:   v3/client/tailwind.config.js
#	modified:   v3/client/vite.config.ts
#	new file:   v3/lib/base/buffer_pool.go
#	new file:   v3/lib/base/os_termination_signal.go
#	new file:   v3/lib/base/string_operations.go
#	new file:   v3/lib/base/utils.go
#	modified:   v3/lib/context/context.go
#	new file:   v3/lib/cors/cors_test.go
#	new file:   v3/lib/cors/hash_javascript.go
#	new file:   v3/lib/cors/hash_javascript_directory.go
#	new file:   v3/lib/cors/types.go
#	modified:   v3/lib/db_access/db_main.go
#	modified:   v3/lib/logging/fake_logger.go
#	modified:   v3/lib/logging/logging_struct.go
#	modified:   v3/peer/config/data_storage_struct.go
#	modified:   v3/peer/main.go
#	new file:   v3/peer/mainthread/db_helpers/utils.go
#	modified:   v3/peer/mainthread/download_file.go
#	modified:   v3/peer/mainthread/loop.go
#	modified:   v3/peer/mainthread/walk_directory.go
#	new file:   v3/peer/network-peer/api/list_files.go
#	new file:   v3/peer/network-peer/api/register_apis.go
#	modified:   v3/peer/network-peer/network_peer.go
#	renamed:    v3/peer/server/init_server.go -> v3/peer/network-peer/server/init_server.go
#	renamed:    v3/peer/server/server_loop.go -> v3/peer/network-peer/server/server_loop.go
#	renamed:    v3/peer/server/server_type.go -> v3/peer/network-peer/server/server_type.go
#	new file:   v3/scripts/build_script.ps1
#
# Changes not staged for commit:
#	modified:   v1/client (modified content)
#	modified:   v1/lib (modified content)
#
```

### 03rd June 2024, 10:35 AM GMT +3
```sh
1.  Upload Works:
    -> add file_data bytes get hash and hash_id
    -> add file_metadata get file_path and file_id
    -> update file_tracker

2. Database Schema v1 settled
3. Next:
    a. Upload Directory (copy up with no download)
    b. Download file func
    c. React UI (or PreactUI?)

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Changes to be committed:
#	modified:   .gitignore
#	deleted:    lib/db_access/__sql/_peer_queries.sql
#	deleted:    lib/db_access/generated/_peer_queries.sql.go
#	deleted:    lib/db_access/generated/batch.go
#	modified:   lib/db_access/generated/db.go
#	modified:   lib/db_access/generated/models.go
#	new file:   lib/db_access/generated/peer_queries.sql.go
#	new file:   lib/db_access/sql/peer_queries.sql
#	renamed:    lib/db_access/__sql/_peer_schema.sql -> lib/db_access/sql/peer_schema.sql
#	modified:   lib/logging/logging_struct.go
#	modified:   lib/sqlc.yaml
#	modified:   peer/.air.toml
#	modified:   peer/config/data_storage_struct.go
#	modified:   peer/logs/log_err_file.txt
#	modified:   peer/logs/log_file.txt
#	modified:   peer/main.go
#	new file:   peer/mainthread/db_helpers/connect_to_db.go
#	new file:   peer/mainthread/db_helpers/upload.go
#	modified:   peer/mainthread/loop.go
#	new file:   peer/mainthread/upload_file.go
#	new file:   peer/mainthread/walk_directory.go
#
# Changes not staged for commit:
#	modified:   peer/mainthread/loop.go
#
```