package kvm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/digitalocean/go-qemu/qmp"
	"github.com/golang/protobuf/ptypes"
	"github.com/shirou/gopsutil/process"

	"code.cloudfoundry.org/bytefmt"
	"github.com/satori/go.uuid"

	n0stack "github.com/n0stack/proto"
	"github.com/n0stack/proto/device/vm"
	"github.com/n0stack/proto/resource/cpu"
)

type (
	Agent struct {
		// DB *gorm.DB
	}

	kvm struct {
		id      uuid.UUID
		workDir string

		args []string
		pid  int
		qmp  *qmp.SocketMonitor
	}
)

const (
	modelType = "device/vm/kvm"
)

// MakeNotification return Notification message with time.
// In the future, this method will be hooked some functions such as storing notifications for database.
// Hard-code string without some logic for making searching line easily, when you set arguments for operation. (ゴミ英語、grep検索を容易にするために関数呼び出しをするときにはtypoがあってもいいからstringをハードコードしろってこと、別にtypoがあってもバグにはならないし検索するときには問題ないため)
func MakeNotification(operation string, success bool, description string) *n0stack.Notification {
	return &n0stack.Notification{
		Operation:   operation,
		Success:     success,
		Description: description,
		NotifiedAt:  ptypes.TimestampNow(),
	}
}

func (k kvm) getInstanceName(n string) string {
	return fmt.Sprintf("n0core-%s", n)
}

func getVM(model *n0stack.Model) (*kvm, *n0stack.Notification) {
	k := &kvm{}

	var err error
	k.id, err = uuid.FromBytes(model.Id)
	if err != nil {
		return nil, MakeNotification("getVM.validateUUID", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	const basedir = "/var/lib/n0core"
	k.workDir = filepath.Join(basedir, modelType)
	k.workDir = filepath.Join(k.workDir, k.id.String())
	if err := os.MkdirAll(k.workDir, os.ModePerm); err != nil {
		return nil, MakeNotification("getVM.prepareWorkDir", false, fmt.Sprintf("error message '%s', when creating work directory, '%s'", k.workDir, err.Error()))
	}

	ps, err := process.Processes()
	if err != nil {
		return nil, MakeNotification("getVM.getProcessList", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	for _, p := range ps {
		c, _ := p.Cmdline() // エラーが発生する場合が考えられない
		// println(c)
		if strings.Contains(c, k.id.String()) {
			k.args, _ = p.CmdlineSlice()

			k.pid = int(p.Pid)
			return k, MakeNotification("getVM", true, fmt.Sprintf("Already running: pid=%d", k.pid))
		}
	}

	return k, MakeNotification("getVM", true, "Not running QEMU process")
}

func (k *kvm) runVM(spec *vm.Spec) *n0stack.Notification {
	switch spec.Cpu.Architecture {
	case cpu.Architecture_x86_64:
		k.args = []string{"qemu-system-x86_64"}
	}

	// -- QEMU metadata --
	k.args = append(k.args, "-uuid")
	k.args = append(k.args, k.id.String())
	k.args = append(k.args, "-name")
	k.args = append(k.args, fmt.Sprintf("guest=%s,debug-threads=on", k.getInstanceName(spec.Device.Model.Name)))
	k.args = append(k.args, "-msg")
	k.args = append(k.args, "timestamp=on")

	k.args = append(k.args, "-nodefaults")     // Don't create default devices
	k.args = append(k.args, "-no-user-config") // The "-no-user-config" option makes QEMU not load any of the user-provided config files on sysconfdir
	k.args = append(k.args, "-S")              // Do not start CPU at startup
	k.args = append(k.args, "-no-shutdown")    // Don't exit QEMU on guest shutdown

	// QMP
	const monitorFile = "monitor.sock"
	qmpPath := filepath.Join(k.workDir, monitorFile)
	k.args = append(k.args, "-chardev")
	k.args = append(k.args, fmt.Sprintf("socket,id=charmonitor,path=%s,server,nowait", qmpPath))
	k.args = append(k.args, "-mon")
	k.args = append(k.args, "chardev=charmonitor,id=monitor,mode=control")

	// -- BIOS --
	// boot priority
	k.args = append(k.args, "-boot")
	k.args = append(k.args, "menu=on,strict=on")

	// keyboard
	k.args = append(k.args, "-k")
	k.args = append(k.args, "en-us") // vm.Spec.Keymapみたいなので取得できるようにする

	// VNC
	k.args = append(k.args, "-vnc")
	k.args = append(k.args, ":0") // ぶつからないようにポートを設定する必要がある, unix socketでも可 unix:$workdir/vnc.sock,websocket

	// clock
	k.args = append(k.args, "-rtc")
	k.args = append(k.args, "base=utc,driftfix=slew")
	k.args = append(k.args, "-global")
	k.args = append(k.args, "kvm-pit.lost_tick_policy=delay")
	k.args = append(k.args, "-no-hpet")

	// CPU
	// TODO: 必要があればmonitorを操作してhotaddできるようにする
	// TODO: スケジューリングが可能かどうか調べる
	k.args = append(k.args, "-cpu")
	k.args = append(k.args, "host")
	k.args = append(k.args, "-smp")
	k.args = append(k.args, fmt.Sprintf("%d,sockets=1,cores=%d,threads=1", spec.Cpu.Vcpus, spec.Cpu.Vcpus))
	k.args = append(k.args, "-enable-kvm")
	// return true, "Succeeded to check cpu configurations"

	// Memory
	// TODO: スケジューリングが可能かどうか調べる
	k.args = append(k.args, "-m")
	k.args = append(k.args, fmt.Sprintf("%s", bytefmt.ByteSize(spec.Memory.Bytes)))
	k.args = append(k.args, "-device")
	k.args = append(k.args, "virtio-balloon-pci,id=balloon0,bus=pci.0") // dynamic configurations
	k.args = append(k.args, "-realtime")
	k.args = append(k.args, "mlock=off")

	// VGA controller
	k.args = append(k.args, "-device")
	k.args = append(k.args, "VGA,id=video0,bus=pci.0")

	// SCSI controller
	k.args = append(k.args, "-device")
	k.args = append(k.args, "lsi53c895a,bus=pci.0,id=scsi0")

	cmd := exec.Command(k.args[0], k.args[1:]...)
	if err := cmd.Start(); err != nil {
		return MakeNotification("startQEMUProcess.startProcess", false, fmt.Sprintf("error message '%s', args '%s'", err.Error(), k.args))
	}
	k.pid = cmd.Process.Pid

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(3 * time.Second):
		return MakeNotification("startQEMUProcess", true, "")
	case err := <-done:
		return MakeNotification("startQEMUProcess.waitError", true, fmt.Sprintf("error message '%s', args '%s'", err.Error(), k.args)) // stderrを表示できるようにする必要がある
	}
}

// func (k *kvm) connectQMP() *n0stack.Notification {
// 	qmpPath := ""

// 	var err error
// 	k.qmp, err = qmp.NewSocketMonitor("unix", qmpPath, 2*time.Second)
// 	if err != nil {
// 		return MakeNotification("connectQMP", true, "")
// 	}

// 	return MakeNotification("connectQMP", true, "")
// }

// Apply スペックを元にステートレスに適用する
func (a *Agent) Apply(ctx context.Context, spec *vm.Spec) (*n0stack.Notification, error) {
	// ps auxfww | grep $uuid
	k, n := getVM(spec.Device.Model)
	if !n.Success {
		return n, nil
	}

	// if vm is not running
	if k.args == nil {
		// check CPU usage
		// check Memory usage

		// qemu-system...
		n = k.runVM(spec)
		if !n.Success {
			return n, nil
		}
		return n, nil
	}

	// qmp-shell .../monitor.sock
	// k.connectQMP()

	// (QEMU) ...
	// conn :=
	// vcl := volume.NewRepositoryClient(conn)
	// vcl.
	// k.attachVolume()
	// k.attachNIC()

	return MakeNotification("Apply", true, ""), nil
}

func (k kvm) Kill() *n0stack.Notification {
	p, _ := os.FindProcess(k.pid)
	if err := p.Kill(); err != nil {
		return MakeNotification("Kill", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	return MakeNotification("Kill", true, "")
}

func (a *Agent) Delete(ctx context.Context, model *n0stack.Model) (*n0stack.Notification, error) {
	// ps auxfww | grep $uuid
	k, n := getVM(model)
	if !n.Success {
		return n, nil
	}

	// if vm is not running
	if k.args == nil {
		return MakeNotification("Delete", true, "Process is not existing"), nil
	}

	// kill $qemu
	n = k.Kill()
	if !n.Success {
		return n, nil
	}

	return MakeNotification("Delete", true, ""), nil
}
