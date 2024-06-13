# Commit Notes

### 13th June 2024, 10:05 AM GMT +3
```txt
1. Able to read all files in a directory
2. Able to list loaded files
3. Saves loaded extensions
4. Regularly saves progress to file-tree.json
5. Server successfull receives files from client
Pending:
6. Reduce clientside memory use
    - only read when hashing and uploading (no dangling filehandlers)
7. Reduce serverside memory use
    - save directly to disk on upload, retaining only file address and info
8. Subscribe and download
    - clientside: 
        1. send subscriptions to server (to notify on update/changes to directory/files)
    - serverside: 
        1. track which clients have which files
        2. track which clients need which files
        3. push to clients on change ("/download/changes" route)
9. Load balancing
    - simple round robin queue
```

### 8th June 2024, 21:27 PM GMT+3
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

### 6th June 2024, 21:13 PM GMT+3
    1. Fixed all logging errors
    2. Moved std out and file logging to single thread
    3. Need to work on hashing

### 6th June 2024, 19:24 PM GMT+3
    1. Adding Logger struct
    2. Adding os.PathSeparator to config and excluded
    3. Improved FileBasic in lib to expose filehandler
    4. able to store and load file-tree.json (persistent state)
    5. Able to read files in directory