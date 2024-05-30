package githandler

import (
	"errors"
	"fmt"
	"os"

	// "sync"

	"log"
	"time"

	"os/exec"

	ftp_context "github.com/ftp_system_client/main_thread/context"
	filehandler "github.com/ftp_system_client/main_thread/file_handler"
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
	is_repo := false
	has_gitignore := false

	m, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, entry := range m {
		if entry.IsDir() && entry.Name() == ".git" {
			is_repo = true

		}
		if !entry.IsDir() && entry.Name() == ".gitignore" {
			has_gitignore = true
		}
	}

	if !is_repo {
		if err_pre := gte.git_init(ctx, directory); err_pre != nil {
			err = ftp_context.NewLogItem(loc, true).
				Set("after", "git_init").AppendParentError(err_pre)

			return
		}
		log.Println("initialised git repo in:", directory)

	}

	if !has_gitignore {
		if err_pre := gte.git_add_gitignore(ctx, directory); err_pre != nil {
			err = ftp_context.NewLogItem(loc, true).
				Set("after", "git_add_gitignore").AppendParentError(err_pre)

			return
		}
		log.Println("created .gitignore in:", directory)
	}

	if err_pre := gte.git_add(ctx, directory); err_pre != nil {
		err = ftp_context.NewLogItem(loc, true).
			Set("after", "git_add").AppendParentError(err_pre)

		return
	}
	log.Println("added files to .git in:", directory)

	pre := gte.git_commit(ctx, directory)
	if pre != nil {
		err = ftp_context.NewLogItem(loc, true).
			Set("after", "git_commit").
			Set("path", directory).
			AppendParentError(pre)

		set_stderr(ctx, loc, stderr, err)
		return
	}
	return
}

func (gte *GitEngine) Commit(path string) error {

	return gte.dir_commit(path)
}

func generate_commit(directory string, added []string) (commit_msg string) {

	return
}

func (gte *GitEngine) git_init(ctx ftp_context.Context, directory string) (err error) {
	o, stderr, err := execute_commit_step(ctx, directory, []string{"init", "."})
	if err != nil {
		log.Println(stderr)
		return
	}
	log.Println(o)
	return
}

func (gte *GitEngine) git_add_gitignore(ctx ftp_context.Context, directory string) (err error) {
	fo := filehandler.NewFileBasic(directory + "/.gitignore")
	b, err := os.ReadFile("./data/templates/.gitignore")
	log.Println(b)
	if err != nil {
		return
	}
	_, err = fo.Create().Write(b)
	if err != nil {
		log.Println(err, "\n", "fo.Create().Write(b)")
		return
	}
	fo.Close()

	return
}

func (gte *GitEngine) git_add(ctx ftp_context.Context, directory string) (err error) {
	loc := "git_add"
	if err_post := os.Remove(directory + "/.git/index.lock"); err_post != nil && !errors.Is(err_post, os.ErrNotExist) {
		err = ftp_context.NewLogItem(loc, true).
			Set("after", "os.Remove").
			AppendParentError(err_post, err)
		return
	}

	o, stderr, err := execute_commit_step(ctx, directory, []string{"add", "."})
	if err != nil {
		err = ftp_context.NewLogItem(loc, true).
			Set("after", "execute_commit_step").
			Set("std_err", stderr).
			Set("std_out", o).
			AppendParentError(err)
		// log.Println(stderr)

		return
	}
	log.Println(o)
	return
}

func (gte *GitEngine) git_commit(ctx ftp_context.Context, directory string) (err error) {
	loc := "git_commit"
	m := fmt.Sprintf("-m \"%s\"", time.Now().Format(time.RFC1123))
	o, stderr, err := execute_commit_step(ctx, directory, []string{"commit", m})

	if err != nil {
		if len(stderr) < 1 && len(o) > 0 {
			log.Println(o)
			return nil
		}

		err = ftp_context.NewLogItem(loc, true).
			Set("std_err", stderr).
			Set("std_out", o).
			AppendParentError(err)

		return
	}
	log.Println(o)

	return
}
