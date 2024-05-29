package githandler

import (
	"bytes"
	"fmt"
	"strings"
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
	var buf []byte
	var stderr string
	for _, command := range []string{"git add .", fmt.Sprintf("git commit -m \"%s\"", time.Now())} {
		c := strings.Split(command, " ")

		for retry := true; retry; {
			buf, stderr, err = ExecuteCommand(ctx, directory, c[0], c[1:]...)
			if err != nil {
				log.Println(stderr, "\n", err)
				return
			}
			log.Println(string(buf))
			retry, err = handle_common_git_errors(ctx, directory, stderr, err)
			if err != nil {
				err = ftp_context.NewLogItem(loc, true).
					ParentError(err).
					Set("after", "handle_common_git_errors").
					Set("error_msg", err.Error()).
					SetMessage("failed to retry")

				return
			}
		}

		clear_stderr(ctx)

	}

	return
}

func (gte *GitEngine) Commit(path string) error {

	return gte.dir_commit(path)
}

func ExecuteCommand(ctx ftp_context.Context, dir string, command string, arg ...string) (stdout []byte, stderr string, err error) {
	loc := "ExecuteCommand"
	cmd := exec.CommandContext(ctx, command, arg...)
	cmd.Dir = dir
	fmt.Println(cmd)
	var std_out bytes.Buffer
	var std_err bytes.Buffer
	cmd.Stdout = &std_out
	cmd.Stderr = &std_err
	if err = cmd.Start(); err != nil {
		msg := err.Error()
		err = ftp_context.NewLogItem("ExecuteCommand", true).
			ParentError(err).
			Set("after", "cmd.Start()").
			Set("error_msg", msg).
			SetMessage("")
		return
	}

	err = cmd.Wait()
	stdout = std_out.Bytes()
	stderr = std_err.String()

	if len(stdout) > 0 {
		log.Println("outout", stdout)
	}
	if err != nil {
		err = set_stderr(ctx, loc, stderr, err)
	}

	return
}

func handle_common_git_errors(ctx ftp_context.Context, directory string, stderr string, cmd_err error) (retry bool, err error) {
	loc := "handle_common_git_errors"
	var buf []byte
	child_count := child_count_f(ctx)

	if child_count > 5 {
		err = ftp_context.NewLogItem(loc, true).SetMessage("recursion too deep").ParentError(cmd_err)
		return
	}

	if strings.Contains(stderr, "not a git repository") {
		c := strings.Split("git init .", " ")
		buf, stderr, err = ExecuteCommand(ctx, directory, c[0], c[1:]...)

		if err != nil {
			log.Println(err)
			return handle_common_git_errors(ctx, directory, stderr, err)
		}
		log.Println(string(buf))

		return true, nil
	} else {
		retry = false
	}

	if string_contains_multiple(stderr, "fatal: Unable to create", ".git/index.lock: File exists.", "Another git process seems to be running in this repository") {
		c := strings.Split("taskkill -im git -f", " ")
		buf, stderr, err = ExecuteCommand(ctx, directory, c[0], c[1:]...)

		if err != nil {
			log.Fatalln(err)
			return
		}
		log.Println(string(buf))

		c = strings.Split("rm -Force .git/index.lock", " ")
		buf, stderr, err = ExecuteCommand(ctx, directory, c[0], c[1:]...)

		if err != nil {
			log.Fatalln(err)
			return
		}
		log.Println(string(buf))

		return true, nil
	} else {
		retry = false
	}

	return
}

func string_contains_multiple(str string, substrs ...string) bool {

	for _, substr := range substrs {
		if !strings.Contains(str, substr) {
			return false
		}
	}

	return true
}

func child_count_f(ctx ftp_context.Context) (n int) {
	cc := "child_count"
	_n, ok := ctx.Get(cc)
	if !ok {
		ctx.Set(cc, 1)
		return 1
	}
	n, _ = _n.(int)
	n += 1
	ctx.Set(cc, n)
	return
}

func set_stderr(ctx ftp_context.Context, loc string, stderr string, err error) (cmp_err error) {
	cc := "std_err"
	msg := err.Error()
	cmp_err = ftp_context.NewLogItem("ExecuteCommand", true).
		ParentError(err).
		Set("after", loc).
		Set("error_msg", msg).
		Set("stderr", string(stderr)).
		ParentError(err)

	ctx.Set(cc, cmp_err)
	return
}

func get_stderr(ctx ftp_context.Context) (stderr string, err ftp_context.LogErr, ok bool) {
	cc := "std_err"
	cmp_err, ok := ctx.Get(cc)
	if !ok {
		return
	}

	err, ok = cmp_err.(ftp_context.LogErr)
	if !ok {
		return
	}

	cmp_err, ok = err.Get("stderr")
	if !ok {
		return
	}

	stderr, ok = cmp_err.(string)
	return
}

func clear_stderr(ctx ftp_context.Context) (old_stderr ftp_context.LogErr) {
	cc := "std_err"
	_old_stderr, ok := ctx.Delete(cc)
	if !ok {
		return nil
	}

	old_stderr, ok = _old_stderr.(ftp_context.LogErr)
	return
}
