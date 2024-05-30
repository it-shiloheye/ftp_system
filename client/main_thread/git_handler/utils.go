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

func set_stderr(ctx ftp_context.Context, loc string, stderr string, err error) (cmp_err ftp_context.LogErr) {
	cc := "std_err"
	cmp_err = ftp_context.NewLogItem("ExecuteCommand", true).
		AppendParentError(err).
		Set("after", loc).
		Set("stderr", strings_split(string(stderr), "\n")).
		AppendParentError(err)

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

func ExecuteCommand(ctx ftp_context.Context, dir string, command string, arg ...string) (stdout []byte, stderr string, err ftp_context.LogErr) {
	loc := "ExecuteCommand"
	cmd := exec.CommandContext(ctx, command, arg...)
	cmd.Dir = dir
	log.Println(cmd, "\npwd:", dir)
	var std_out bytes.Buffer
	var std_err bytes.Buffer
	cmd.Stdout = &std_out
	cmd.Stderr = &std_err
	if err_ := cmd.Start(); err_ != nil {
		msg := err.Error()
		err = ftp_context.NewLogItem(loc, true).
			AppendParentError(err_).
			Set("after", "cmd.Start()").
			Set("error_msg", msg).
			SetMessage("")
		cmd.Cancel()
		return
	}

	err_ := cmd.Wait()
	stdout = std_out.Bytes()
	stderr = std_err.String()

	if err_ != nil {
		a := append([]string{command}, arg...)
		err = set_stderr(ctx, strings.Join(a, " "), stderr, err_)

	}

	return
}

// internally handles retrying commit command in case of error
func execute_commit_step(ctx ftp_context.Context, directory string, command []string) (output string, stderr string, err error) {
	var buf []byte
	buf, stderr, err = ExecuteCommand(ctx, directory, "git", command...)
	if len(stderr) < 1 {
		err = nil
	}
	output = string(buf)
	return
}
