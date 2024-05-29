package githandler

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	// "sync"

	ftp_context "github.com/ftp_system_client/main_thread/context"
)

func string_contains_multiple(str string, substrs ...string) bool {

	for _, substr := range substrs {
		if !strings.Contains(str, substr) {
			return false
		}
	}

	return true
}

func inc_child_count_f(ctx ftp_context.Context) (n int) {
	fmt.Println(ctx)
	return
}
func dec_child_count_f(ctx ftp_context.Context) (n int) {
	fmt.Println(ctx)
	return
}

func set_stderr(ctx ftp_context.Context, loc string, stderr string, err error) (cmp_err error) {
	cc := "std_err"
	msg := err.Error()
	cmp_err = ftp_context.NewLogItem("ExecuteCommand", true).
		ParentError(err).
		Set("after", loc).
		Set("error_msg", msg).
		Set("stderr", strings_split(string(stderr), "\n")).
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

func strings_split(str string, substr string) (out []string) {
	a := strings.Split(str, substr)
	b := ""
	for _, s := range a {
		b = strings.Trim(s, "\t\n\r")
		if len(s) > 0 {
			out = append(out, b)
		}
	}

	return

}

func ExecuteCommand(ctx ftp_context.Context, dir string, command string, arg ...string) (stdout []byte, stderr string, err error) {
	loc := "ExecuteCommand"
	cmd := exec.CommandContext(ctx, command, arg...)
	cmd.Dir = dir
	fmt.Println(cmd, "\npwd:", dir)
	var std_out bytes.Buffer
	var std_err bytes.Buffer
	cmd.Stdout = &std_out
	cmd.Stderr = &std_err
	if err = cmd.Start(); err != nil {
		msg := err.Error()
		err = ftp_context.NewLogItem(loc, true).
			ParentError(err).
			Set("after", "cmd.Start()").
			Set("error_msg", msg).
			SetMessage("")
		cmd.Cancel()
		return
	}

	err = cmd.Wait()
	stdout = std_out.Bytes()
	stderr = std_err.String()

	if err != nil {
		a := append([]string{command}, arg...)
		err = set_stderr(ctx, strings.Join(a, " "), stderr, err)

	}

	return
}

// internally handles retrying commit command in case of error
func execute_commit_step(ctx ftp_context.Context, directory string, command []string) (output string, stderr string, err error) {
	loc := "execute_commit_step"
	var buf []byte

	for retry := true; retry; {
		buf, stderr, err = ExecuteCommand(ctx, directory, command[0], command[1:]...)
		if err != nil {
			log.Println("trying to handle error:\n", stderr)
			retry, err = handle_common_git_errors(ctx, directory, stderr, err)
			if err != nil {
				err = ftp_context.NewLogItem(loc, true).
					ParentError(err).
					Set("after", "handle_common_git_errors").
					Set("error_msg", err.Error()).
					SetMessage("failed to retry")

				return
			} else {
				log.Println("able to handle error")
			}
			if !retry {
				return
			}
		} else {
			retry = false
			output = string(buf)
			break
		}

	}

	clear_stderr(ctx)
	return
}
