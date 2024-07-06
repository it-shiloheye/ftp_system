# Commit Notes

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