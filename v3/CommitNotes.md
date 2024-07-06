# Commit Notes

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