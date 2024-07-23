package main

import (
	"log"
	"os/exec"
)

type hookEvent string

const (
	evtWarnJSON   hookEvent = "WJSON"
	evtErrGitOps  hookEvent = "EGITOPS"
	evtErrUserAdd hookEvent = "EUSERADD"
	evtErrUserDel hookEvent = "EUSERDEL"
	evtErrUserMod hookEvent = "EUSERMOD"
	evtUserAdd    hookEvent = "IUSERADD"
	evtUserDel    hookEvent = "IUSERDEL"
	evtUserMod    hookEvent = "IUSERMOD"
)

type hookExecutor interface {
	Exec(event hookEvent, msg ...string)
}

type noopHookExecutor struct{}

func (noopHookExecutor) Exec(hookEvent, ...string) {
	// noop
}

type hookCommandExecutor struct {
	path string
}

func (h hookCommandExecutor) Exec(evt hookEvent, msgs ...string) {
	args := make([]string, 0, len(msgs)+1)
	args = append(args, string(evt))
	args = append(args, msgs...)

	cmd := exec.Command(h.path, args...)
	if err := cmd.Run(); err != nil {
		log.Printf("Warn: command %s execution failure: %s", h.path, err)
	}
}

func getHookExecutor(cmd string) (hookExecutor, error) {
	if cmd != "" {
		path, err := exec.LookPath(cmd)
		if err != nil {
			return nil, err
		}
		return &hookCommandExecutor{path}, nil
	}
	return &noopHookExecutor{}, nil
}
