package futaba

import (
	"errors"
	"fmt"
	"karakuripkgs"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

func setHostname(spec string) error {
	// read specfile
	config_spec := karakuripkgs.ReadSpecFile(spec)

	if err := syscall.Sethostname([]byte(config_spec.Hostname)); err != nil {
		return errors.New("failed to set hostname")
	}
	return nil
}

func setNameeserver(config_spec karakuripkgs.ConfigSpec) error {
	resolv_conf := config_spec.Root.Path + "/merged/etc/resolv.conf"
	if err := os.WriteFile(resolv_conf, []byte("nameserver "+config_spec.Network.Nameserver+"\nsearch karakuri.container\n"), 0644); err != nil {
		return errors.New("failed to create /etc/resolv.conf")
	}
	return nil
}

func setHostsfile(config_spec karakuripkgs.ConfigSpec) error {
	hosts_file := config_spec.Root.Path + "/merged/etc/hosts"
	if err := os.WriteFile(hosts_file, []byte("127.0.0.1 localhost\n127.0.1.1 "+config_spec.Hostname+"\n"), 0644); err != nil {
		return errors.New("failed to create /etc/hosts")
	}
	return nil
}

func setEnv(config_spec karakuripkgs.ConfigSpec) error {
	envs := config_spec.Process.Env
	for _, entry := range envs {
		if err := os.Setenv(entry.Key, entry.Value); err != nil {
			return errors.New("failed to set environmental variables: " + entry.Key + "=" + entry.Value)
		}
	}
	return nil
}

func mountOverlay(config_spec karakuripkgs.ConfigSpec) error {
	// overlay
	if err := syscall.Mount(
		"overlay",
		config_spec.Root.Path+"/merged",
		"overlay",
		0,
		"lowerdir="+config_spec.Image.Path+",upperdir="+config_spec.Root.Path+"/diff,workdir="+config_spec.Root.Path+"/work",
	); err != nil {
		return errors.New("failed to mount overlay filesystem")
	}
	return nil
}

func mountFs(config_spec karakuripkgs.ConfigSpec) error {
	// mount file systems
	for _, mount_info := range config_spec.Mounts {
		var (
			mount_flag   int
			mount_option string = ""
		)
		// retrieve options
		if mount_info.Options != nil {
			var option_tmp = ""
			for _, option := range mount_info.Options {
				switch option {
				case "nosuid":
					mount_flag |= syscall.MS_NOSUID
				case "noexec":
					mount_flag |= syscall.MS_NOEXEC
				case "nodev":
					mount_flag |= syscall.MS_NODEV
				case "ro":
					mount_flag |= syscall.MS_RDONLY
				case "rw":
					// mount_flag |= syscall.O_RDWR
				case "bind":
					mount_flag |= syscall.MS_BIND
				default:
					option_tmp += option + ","
				}
			}
			mount_option = strings.TrimRight(option_tmp, ",")
		} else {
			mount_flag = 0
		}
		// check directory
		if _, err := os.Stat(config_spec.Root.Path + "/merged" + mount_info.Destination); err != nil {
			if err := os.MkdirAll(config_spec.Root.Path+"/merged"+mount_info.Destination, os.ModePerm); err != nil {
				return errors.New("failed to create mount destination directory")
			}
		}
		// mount
		if err := syscall.Mount(mount_info.Source, config_spec.Root.Path+"/merged"+mount_info.Destination, mount_info.MountType, uintptr(mount_flag), mount_option); err != nil {
			return errors.New("failed to mount " + mount_info.Source + " to " + mount_info.Destination + ", type: " + mount_info.MountType)
		}
	}
	return nil
}

// mount standard device
func mountStdDeviceFile(config_spec karakuripkgs.ConfigSpec) error {
	files := []string{"random", "urandom", "null", "zero", "tty"}
	for _, file := range files {
		src := "/dev/" + file
		dst := config_spec.Root.Path + "/merged/dev/" + file
		if _, err := os.Stat(dst); os.IsNotExist(err) {
			if _, err := os.Create(dst); err != nil {
				return errors.New("failed to create " + dst)
			}
		}
		if err := syscall.Mount(src, dst, "", syscall.MS_BIND, ""); err != nil {
			return errors.New("failed to mount /dev file: " + file)
		}
	}
	return nil
}

// create symbolic link
func createSymbolicLink(config_spec karakuripkgs.ConfigSpec) error {
	devDir := config_spec.Root.Path + "/merged/dev"
	symlinks := []struct {
		link   string
		target string
	}{
		{devDir + "/fd", "/proc/self/fd"},
		{devDir + "/stdin", "/proc/self/fd/0"},
		{devDir + "/stdout", "/proc/self/fd/1"},
		{devDir + "/stderr", "/proc/self/fd/2"},
		{devDir + "/ptmx", "/dev/pts/ptmx"},
	}

	for _, s := range symlinks {
		if _, err := os.Lstat(s.link); err == nil {
			continue
		}
		if err := os.Symlink(s.target, s.link); err != nil {
			return fmt.Errorf("failed to create symlink %s -> %s: %w", s.link, s.target, err)
		}
	}

	return nil
}

func setCapabilities(pid int, config_spec karakuripkgs.ConfigSpec) error {
	// capability mapping
	mapping := map[string]int{
		"CHOWN":              unix.CAP_CHOWN,
		"DAC_OVERRIDE":       unix.CAP_DAC_OVERRIDE,
		"FSETID":             unix.CAP_FSETID,
		"FOWNER":             unix.CAP_FOWNER,
		"MKNOD":              unix.CAP_MKNOD,
		"NET_RAW":            unix.CAP_NET_RAW,
		"SETGID":             unix.CAP_SETGID,
		"SETUID":             unix.CAP_SETUID,
		"SETFCAP":            unix.CAP_SETFCAP,
		"SETPCAP":            unix.CAP_SETPCAP,
		"NET_BIND_SERVICE":   unix.CAP_NET_BIND_SERVICE,
		"KILL":               unix.CAP_KILL,
		"AUDIT_WRITE":        unix.CAP_AUDIT_WRITE,
		"SYS_CHROOT":         unix.CAP_SYS_CHROOT,
		"AUDIT_CONTROL":      unix.CAP_AUDIT_CONTROL,
		"AUDIT_READ":         unix.CAP_AUDIT_READ,
		"BLOCK_SUSPEND":      unix.CAP_BLOCK_SUSPEND,
		"BPF":                unix.CAP_BPF,
		"CHECKPOINT_RESTORE": unix.CAP_CHECKPOINT_RESTORE,
		"DAC_READ_SEARCH":    unix.CAP_DAC_READ_SEARCH,
		"IPC_LOCK":           unix.CAP_IPC_LOCK,
		"IPC_OWNER":          unix.CAP_IPC_OWNER,
		"LEASE":              unix.CAP_LEASE,
		"LINUX_IMMUTABLE":    unix.CAP_LINUX_IMMUTABLE,
		"MAC_ADMIN":          unix.CAP_MAC_ADMIN,
		"MAC_OVERRIDE":       unix.CAP_MAC_OVERRIDE,
		"NET_ADMIN":          unix.CAP_NET_ADMIN,
		"NET_BROADCAST":      unix.CAP_NET_BROADCAST,
		"PERFMON":            unix.CAP_PERFMON,
		"SYS_ADMIN":          unix.CAP_SYS_ADMIN,
		"SYS_BOOT":           unix.CAP_SYS_BOOT,
		"SYS_MODULE":         unix.CAP_SYS_MODULE,
		"SYS_NICE":           unix.CAP_SYS_NICE,
		"SYS_PACCT":          unix.CAP_SYS_PACCT,
		"SYS_PTRACE":         unix.CAP_SYS_PTRACE,
		"SYS_RAWIO":          unix.CAP_SYS_RAWIO,
		"SYS_RESOURCE":       unix.CAP_SYS_RESOURCE,
		"SYS_TIME":           unix.CAP_SYS_TIME,
		"SYS_TTY_CONFIG":     unix.CAP_SYS_TTY_CONFIG,
		"SYSLOG":             unix.CAP_SYSLOG,
		"WAKE_ALARM":         unix.CAP_WAKE_ALARM,
	}

	// mapping add capability
	addCaps := make([]int, len(config_spec.Capability.AddCapability))
	for i, s := range config_spec.Capability.AddCapability {
		if val, ok := mapping[s]; ok {
			addCaps[i] = val
		} else {
			return fmt.Errorf("failed to read capability")
		}
	}
	// mapping drop capability
	dropCaps := make([]int, len(config_spec.Capability.DropCapability))
	for i, s := range config_spec.Capability.DropCapability {
		if val, ok := mapping[s]; ok {
			dropCaps[i] = val
		} else {
			return fmt.Errorf("failed to read capability")
		}
	}

	// 1. Drop from Bounding set
	for _, c := range dropCaps {
		if err := unix.Prctl(unix.PR_CAPBSET_DROP, uintptr(c), 0, 0, 0); err != nil {
			return fmt.Errorf("failed to drop capability %d from bounding set: %w", c, err)
		}
	}

	// 2. prepare CapUserHeader
	var hdr unix.CapUserHeader
	hdr.Version = unix.LINUX_CAPABILITY_VERSION_3
	hdr.Pid = int32(pid)

	var data [2]unix.CapUserData
	var permittedLow, permittedHigh uint32

	for _, c := range addCaps {
		if c < 0 {
			continue
		}
		if c < 32 {
			permittedLow |= 1 << uint(c)
		} else {
			permittedHigh |= 1 << uint(c-32)
		}
	}

	data[0].Permitted = permittedLow
	data[0].Effective = permittedLow
	data[0].Inheritable = 0

	data[1].Permitted = permittedHigh
	data[1].Effective = permittedHigh
	data[1].Inheritable = 0

	if err := unix.Capset(&hdr, &data[0]); err != nil {
		return fmt.Errorf("capset failed: %w", err)
	}

	return nil
}

func pivotRoot(container_dir string) error {
	// change direrctory to container directory
	if err := os.Chdir(container_dir); err != nil {
		return errors.New("failed to enter container layer directory")
	}
	// mount merged
	if err := syscall.Mount("merged", container_dir+"/merged", "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return errors.New("failed to mount merged filesystem")
	}
	// create put_old directory
	if err := os.MkdirAll(container_dir+"/merged/put_old", 0700); err != nil {
		return errors.New("failed to create /put_old")
	}
	// pivot_root
	if err := syscall.PivotRoot("merged", container_dir+"/merged/put_old"); err != nil {
		return errors.New("failed to execute pivot_root")
	}
	// change directory to root after pivot_root
	if err := os.Chdir("/"); err != nil {
		return errors.New("failed to change directory to /")
	}
	// unmount put_old
	if err := syscall.Unmount("/put_old", syscall.MNT_DETACH); err != nil {
		return errors.New("failed to unmount /put_old")
	}
	// delete put_old
	if err := syscall.Rmdir("/put_old"); err != nil {
		return errors.New("failed to remove /put_old")
	}
	return nil
}

func createFifo(fifo_path string) error {
	if err := unix.Mkfifo(fifo_path+"/exec.fifo", 0o622); err != nil {
		return errors.New("failed to create named pipe: " + fifo_path + "/exec.fifo")
	}
	return nil
}

func waitParant(fifo_path string) error {
	_, err := unix.Open(fifo_path+"/exec.fifo", unix.O_WRONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return errors.New("failed to open named pipe: " + fifo_path + "/exec.fifo for waiting parant process")
	}
	return nil
}

func startChild(fifo_path string) error {
	_, err := os.OpenFile(fifo_path+"/exec.fifo", os.O_RDONLY, 0)
	if err != nil {
		return errors.New("failed to open named pipe: " + fifo_path + "/exec.fifo for starting child process")
	}
	return nil
}

func StartContainer(id string, terminal bool) {
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		fmt.Println(err)
		return
	}

	// config spec path
	container_spec_path := karakuripkgs.FUTABA_ROOT + "/" + id

	args := []string{"init", "--spec=" + container_spec_path}

	cmd := exec.Command(initCmd, args[0:]...)
	// set standard i/o/e
	if terminal {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		//devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
		logfile, err := os.OpenFile(container_spec_path+"/container.log", os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = logfile
		cmd.Stderr = logfile
	}

	// set clone flags
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWCGROUP |
			syscall.CLONE_NEWNET,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        65535,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        65535,
			},
		},
		GidMappingsEnableSetgroups: true,
	}

	// execute command
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return
	}
	// retrieve pid
	pid := cmd.Process.Pid
	// update pid
	karakuripkgs.UpdateSpecPid(container_spec_path, pid)

	// retrieve config spec
	config_spec := karakuripkgs.ReadSpecFile(container_spec_path)

	// setup network interface
	setupContainerNetwork(pid, config_spec.Network)

	// retrieve resource limit
	cpu_max := config_spec.Cgroup.Cpu.Max
	mem_max := config_spec.Cgroup.Memory.Max
	// create cgroup
	if err := createCgroup(id); err != nil {
		fmt.Println(err)
		return
	}
	// set cpu limit
	if err := setCpuLimit(id, cpu_max); err != nil {
		fmt.Println(err)
		return
	}
	// set memory limit
	if err := setMemoryLimit(id, mem_max); err != nil {
		fmt.Println(err)
		return
	}
	// set pid to cgroup
	if err := setCgourpPid(id, pid); err != nil {
		fmt.Println(err)
		return
	}

	// start child
	if err := startChild(config_spec.Fifo); err != nil {
		fmt.Println(err)
		return
	}

	if terminal {
		cmd.Wait()
	}
}

func execEntrypoint(config_spec karakuripkgs.ConfigSpec) error {
	args := config_spec.Process.Args
	cmd := exec.Command(args[0], args[1:]...)

	// set standard i/o/e
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return errors.New("failed to execute container entrypoint: " + strings.Join(args, " "))
	}
	cmd.Wait()

	return nil
}

func InitContainer(spec string) {
	// read specfile
	config_spec := karakuripkgs.ReadSpecFile(spec)

	// wait parent process
	if err := waitParant(config_spec.Fifo); err != nil {
		fmt.Println(err)
		return
	}

	// set hostname
	if err := setHostname(spec); err != nil {
		fmt.Println(err)
		return
	}
	// mount overlay
	if err := mountOverlay(config_spec); err != nil {
		fmt.Println(err)
		return
	}
	// mount file system
	if err := mountFs(config_spec); err != nil {
		fmt.Println(err)
		return
	}

	// mount standard device file
	if err := mountStdDeviceFile(config_spec); err != nil {
		fmt.Println(err)
		return
	}

	// create symbolic link
	if err := createSymbolicLink(config_spec); err != nil {
		fmt.Println(err)
		return
	}

	// env
	if err := setEnv(config_spec); err != nil {
		fmt.Println(err)
		return
	}

	// set nameserver
	setNameeserver(config_spec)

	// set hosts
	setHostsfile(config_spec)

	// pivot root
	if err := pivotRoot(config_spec.Root.Path); err != nil {
		fmt.Println(err)
		return
	}

	// set capability
	if err := setCapabilities(0, config_spec); err != nil {
		fmt.Println(err)
		return
	}

	// execute entry point
	if err := execEntrypoint(config_spec); err != nil {
		fmt.Println(err)
		return
	}
}

func ExecContainer(id string, terminal bool, command string) {
	// config spec path
	container_spec_path := karakuripkgs.FUTABA_ROOT + "/" + id
	// read config spec
	config_spec := karakuripkgs.ReadSpecFile(container_spec_path)
	// command
	cmd := strings.Split(command, ",")

	pid := config_spec.Process.Pid
	args := []string{"nsenter", "-t", strconv.Itoa(pid), "--all"}
	args = append(args, cmd...)
	nsenter := exec.Command("sudo", args...)

	nsenter.Stdin = os.Stdin
	nsenter.Stdout = os.Stdout
	nsenter.Stderr = os.Stderr

	if err := nsenter.Start(); err != nil {
		fmt.Println(err)
		return
	}
	if terminal {
		nsenter.Wait()
	}
}

func KillContainer(id string) {
	// config spec path
	container_spec_path := karakuripkgs.FUTABA_ROOT + "/" + id
	// read config spec
	config_spec := karakuripkgs.ReadSpecFile(container_spec_path)

	// kill process
	proc, err := os.FindProcess(config_spec.Process.Pid)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := proc.Kill(); err != nil {
		fmt.Println(err)
		return
	}
}
