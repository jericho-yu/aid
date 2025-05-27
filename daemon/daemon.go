package daemon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// Daemon 守护进程服务提供者
type Daemon struct{}

var App Daemon

// New 实例化：
func (*Daemon) New(logger *zap.Logger) *Daemon { return &Daemon{} }

// Launch 启动守护进程
func (my *Daemon) Launch(title, logDir string) {
	var err error

	if syscall.Getppid() == 1 {
		if err := os.Chdir("./"); err != nil {
			panic(err)
		}
		syscall.Umask(0) // TODO TEST
		return
	}

	logFilename := fmt.Sprintf("%s/runtime.log", logDir)
	fp, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("【启动失败】创建总日志失败：%s", err.Error())
	}
	defer func() {
		_ = fp.Close()
	}()
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // TODO TEST
	cmd.Stdout = fp
	cmd.Stderr = fp
	cmd.Stdin = nil
	if err = cmd.Start(); err != nil {
		log.Fatalf("【启动失败】%s", err.Error())
	}

	fp.WriteString(fmt.Sprintf("(>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\r\n%s 程序启动成功 [进程号->%d] 启动于：%s\r\n", title, cmd.Process.Pid, time.Now().Format(string(time.DateTime+".000"))))

	os.Exit(0)
}
