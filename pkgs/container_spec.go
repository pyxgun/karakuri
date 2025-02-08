package karakuripkgs

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type SpecFlag struct {
	RootPath    string
	ImagePath   string
	Command     string
	OutputPath  string
	HostDevice  string
	Address     string
	Gateway     string
	Nameserver  string
	Mount       string
	PortForward string
	EnvVars     string
}

type SpecEnv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SpecProcess struct {
	Args []string  `json:"args"`
	Env  []SpecEnv `json:"env"`
	Pid  int       `json:"pid"`
}

type SpecCpu struct {
	Max string `json:"max"`
}

type SpecMemory struct {
	Max string `json:"max"`
}

type SpecCgroup struct {
	Path   string     `json:"path"`
	Cpu    SpecCpu    `json:"cpu"`
	Memory SpecMemory `json:"memory"`
}

type SpecRoot struct {
	Path string `json:"path"`
}

type SpecImage struct {
	Path string `json:"path"`
}

type SpecMount struct {
	Destination string   `json:"destination"`
	MountType   string   `json:"type"`
	Source      string   `json:"source"`
	Options     []string `json:"options"`
}

type SpecPortForward struct {
	HostPort   int    `json:"host"`
	TargetPort int    `json:"target"`
	Protocol   string `json:"protocol"`
}

type SpecNetwork struct {
	HostDevice string            `json:"hostdevice"`
	Address    string            `json:"address"`
	Gateway    string            `json:"gateway"`
	Nameserver string            `json:"nameserver"`
	Port       []SpecPortForward `json:"port"`
}

type ConfigSpec struct {
	Version  string      `json:"version"`
	Process  SpecProcess `json:"process"`
	Cgroup   SpecCgroup  `json:"cgroup"`
	Root     SpecRoot    `json:"root"`
	Image    SpecImage   `json:"image"`
	Hostname string      `json:"hostname"`
	Fifo     string      `json:"fifo"`
	Mounts   []SpecMount `json:"mounts"`
	Network  SpecNetwork `json:"network"`
}

func CreateSpecFile(spec_flag SpecFlag) {
	var (
		root_path  string
		image_path string
		args       []string
	)

	// hostname
	hostname := (uuid.NewString())[24:]

	// set rootpath
	if spec_flag.RootPath == "" {
		root_path = FUTABA_ROOT + "/" + hostname
	} else {
		root_path = spec_flag.RootPath
	}

	// set imagepath
	if spec_flag.ImagePath == "" {
		image_path = root_path + "/rootfs"
	} else {
		image_path = spec_flag.ImagePath
	}

	// set command
	if spec_flag.Command == "" {
		args = []string{"sh"}
	} else {
		args = strings.Split(spec_flag.Command, ",")
	}

	// cgroup path
	cgroup_path := "/sys/fs/cgroup/" + hostname

	// set fifo
	fifo_path := root_path + "/fifo"

	// create json object
	config_spec := ConfigSpec{
		Version: KARAKURI_VERSION,
		Process: SpecProcess{
			Args: args,
			Pid:  0,
		},
		Cgroup: SpecCgroup{
			Path: cgroup_path,
			Cpu: SpecCpu{
				Max: "80%",
			},
			Memory: SpecMemory{
				Max: "1024M",
			},
		},
		Root: SpecRoot{
			Path: root_path,
		},
		Image: SpecImage{
			Path: image_path,
		},
		Hostname: hostname,
		Fifo:     fifo_path,
		Mounts: []SpecMount{
			// default mount
			// /proc
			{
				Destination: "/proc",
				MountType:   "proc",
				Source:      "proc",
				Options: []string{
					"nosuid",
					"noexec",
					"nodev",
				},
			},
			// /sys
			{
				Destination: "/sys",
				MountType:   "sysfs",
				Source:      "sysfs",
				Options: []string{
					"nosuid",
					"noexec",
					"nodev",
					"ro",
				},
			},
			// cgroup
			{
				Destination: "/sys/fs/cgroup",
				MountType:   "cgroup2",
				Source:      "cgroup",
			},
			// /dev
			{
				Destination: "/dev",
				MountType:   "tmpfs",
				Source:      "tmpfs",
				Options: []string{
					"nosuid",
					"mode=755",
					"size=65536k",
				},
			},
			// /dev/pts
			{
				Destination: "/dev/pts",
				MountType:   "devpts",
				Source:      "devpts",
				Options: []string{
					"rw",
					"nosuid",
					"noexec",
					"newinstance",
					"mode=620",
					"gid=5",
					"ptmxmode=0666",
				},
			},
			// /dev/mqueue
			{
				Destination: "/dev/mqueue",
				MountType:   "mqueue",
				Source:      "mqueue",
				Options: []string{
					"rw",
					"nosuid",
					"nodev",
					"noexec",
				},
			},
			// /dev/shm
			{
				Destination: "/dev/shm",
				MountType:   "tmpfs",
				Source:      "shm",
				Options: []string{
					"rw",
					"nosuid",
					"nodev",
					"noexec",
					"mode=1777",
					"size=65536k",
				},
			},
		},
		Network: SpecNetwork{
			HostDevice: spec_flag.HostDevice,
			Address:    spec_flag.Address,
			Gateway:    spec_flag.Gateway,
			Nameserver: spec_flag.Nameserver,
		},
	}

	// host mount
	if spec_flag.Mount != "" {
		mounts := strings.Split(spec_flag.Mount, ",")
		for _, entry := range mounts {
			mount := strings.Split(entry, ":")
			source := mount[0]
			destination := mount[1]
			// set
			config_spec.Mounts = append(
				config_spec.Mounts,
				SpecMount{
					Destination: destination,
					MountType:   "",
					Source:      source,
					Options: []string{
						"bind",
					},
				})
		}
	}

	// port forward
	if spec_flag.PortForward != "" {
		ports := strings.Split(spec_flag.PortForward, ",")
		for _, entry := range ports {
			port := strings.Split(entry, ":")
			host_port, _ := strconv.Atoi(port[0])
			target_port, _ := strconv.Atoi(port[1])
			protocol := port[2]
			// set
			config_spec.Network.Port = append(
				config_spec.Network.Port,
				SpecPortForward{
					HostPort:   host_port,
					TargetPort: target_port,
					Protocol:   protocol,
				})
		}
	}

	// environment variables
	if spec_flag.EnvVars != "" {
		envs := strings.Split(spec_flag.EnvVars, "!")
		for _, entry := range envs {
			env := strings.Split(entry, "=")
			key := env[0]
			value := env[1]
			value = strings.Replace(value, "*", " ", -1)
			// set
			config_spec.Process.Env = append(
				config_spec.Process.Env,
				SpecEnv{
					Key:   key,
					Value: value,
				},
			)
		}
	}

	// write file
	file, _ := json.MarshalIndent(config_spec, "", "  ")
	if err := os.WriteFile(spec_flag.OutputPath+"/config.json", file, 0644); err != nil {
		panic(err)
	}
}

func ReadSpecFile(spec string) ConfigSpec {
	bytes, err := os.ReadFile(spec + "/config.json")
	if err != nil {
		panic(err)
	}

	var config_spec ConfigSpec
	if err := json.Unmarshal(bytes, &config_spec); err != nil {
		panic(err)
	}

	return config_spec
}

func UpdateSpecCmd(spec string, cmd_args []string) {
	bytes, err := os.ReadFile(spec + "/config.json")
	if err != nil {
		panic(err)
	}

	var config_spec ConfigSpec
	if err := json.Unmarshal(bytes, &config_spec); err != nil {
		panic(err)
	}

	config_spec.Process.Args = cmd_args

	data, _ := json.MarshalIndent(config_spec, "", "  ")
	if err := os.WriteFile(spec+"/config.json", data, os.ModePerm); err != nil {
		panic(err)
	}
}

func UpdateSpecPid(spec string, pid int) {
	bytes, err := os.ReadFile(spec + "/config.json")
	if err != nil {
		panic(err)
	}

	var config_spec ConfigSpec
	if err := json.Unmarshal(bytes, &config_spec); err != nil {
		panic(err)
	}

	config_spec.Process.Pid = pid

	data, _ := json.MarshalIndent(config_spec, "", "  ")
	if err := os.WriteFile(spec+"/config.json", data, os.ModePerm); err != nil {
		panic(err)
	}
}
