package app

import "sync"

type ProcessExithandler func()

var (
	g_process_exit []ProcessExithandler
	g_process_exit_lock sync.Mutex
)

func RegisterProcessExit(p ProcessExithandler) {
	g_process_exit_lock.Lock()
	defer g_process_exit_lock.Unlock()
	g_process_exit = append(g_process_exit, p)
}

func OnProcessExit() {
	g_process_exit_lock.Lock()
	defer g_process_exit_lock.Unlock()

	for _, v := range g_process_exit {
		v()
	}
}