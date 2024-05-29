package githandler

import (
	"fmt"

	// "sync"

	"log"
	"time"

	"os/exec"

	ftp_context "github.com/ftp_system_client/main_thread/context"
)

type GitEngine struct {
	ctx ftp_context.Context
}

func (gte *GitEngine) Init(ctx ftp_context.Context) {

	gte.ctx = ctx
}

func first_dir_init(path string) (err error) {

	cmd := exec.Command("git init " + path)
	err = cmd.Run()

	if err != nil {
		err = ftp_context.NewLogItem("first_dir_init", true).SetMessagef("exec.Command(\"git init %s) error:\n%s", path, err.Error())
	}
	return
}

func (gte *GitEngine) dir_commit(directory string) (err error) {
	loc := "dir_commit"
	ctx := gte.ctx.NewChild()
	var stderr string
	var output string
	for _, command := range [][]string{
		strings_split("git add .", " "),
		[]string{
			"git", "commit", fmt.Sprintf("-m \"%s\"", time.Now().Format(time.RFC1123)),
		},
	} {
		output, stderr, err = execute_commit_step(ctx, directory, command)
		if err != nil {
			set_stderr(ctx, loc, stderr, err)
			return
		}
		log.Println(output)
	}

	return
}

func (gte *GitEngine) Commit(path string) error {

	return gte.dir_commit(path)
}

func generate_commit(directory string, added []string) (commit_msg string) {

	return
}
