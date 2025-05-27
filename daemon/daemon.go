package daemon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jericho-yu/aid/zapProvider"
)

// Daemon 守护进程服务提供者
type Daemon struct {
	logger *zap.Logger
}

var App Daemon

// New 实例化：
func (*Daemon) New(logger *zap.Logger) *Daemon {
	var (
		err error
		d   *Daemon
	)

	d = &Daemon{}

	if logger == nil {
		if d.logger, err = zapProvider.ZapProviderApp.New(
			zapProvider.ZapProviderConfig.New(zapcore.ErrorLevel).
				SetPath(".").
				SetPathAbs(false).
				SetInConsole(false).
				SetEncoderType(zapProvider.EncoderTypeConsole).
				SetNeedCompress(true).
				SetMaxBackup(30).
				SetMaxSize(10).
				SetMaxDay(30),
		); err != nil {
			log.Fatal(err)
		}
	} else {
		d.logger = logger
	}

	return d
}

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
		my.logger.Error(title, zap.String("创建总日志失败", err.Error()))
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
		my.logger.Error(title, zap.String("启动失败", err.Error()))
		log.Fatalf("【启动失败】%s", err.Error())
	}

	my.logger.Info(title, zap.String("启动成功", ""), zap.Int("进程号", cmd.Process.Pid), zap.Time("启动时间", time.Now()))
	log.Printf("%s 程序启动成功 [进程号->%d] 启动于：%s\r\n", title, cmd.Process.Pid, time.Now().Format(string(time.DateTime+".000")))
	os.Exit(0)
}
