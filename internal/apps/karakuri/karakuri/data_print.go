package karakuri

import (
	"fmt"
	"hitoha"
	"karakuri_mod"
	"karakuripkgs"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func printContainerList(container_list hitoha.ContainerList, namespace string) {
	cluster := getTargetCluster().Target
	fmt.Println("[Cluster] " + cluster)
	if namespace != "all" {
		fmt.Printf("CONTAINER ID\t| Name\t\t\t| IMAGE\t\t\t| STATUS\t| PORT\t\t\t\t| COMMAND\n")
		fmt.Printf("----------------+-----------------------+-----------------------+---------------+-------------------------------+--------------------------------------\n")
		for _, entry := range container_list.List {
			// port info
			var port_info string = ""
			for _, port := range entry.Port {
				port_info += "0.0.0.0:" + strconv.Itoa(port.HostPort) + "->" + strconv.Itoa(port.ContainerPort) + "/" + port.Protocol + ","
			}
			port_info = strings.TrimRight(port_info, ",")
			fmt.Printf("%s\t| %-20s\t| %-20s\t| %-12s\t| %-24s\t| %-32s\n", entry.Id, entry.Name, entry.Image, entry.Status, port_info, entry.Command)
		}
	} else {
		fmt.Printf("CONTAINER ID\t| Name\t\t\t| NAMESPACE\t| IMAGE\t\t\t| STATUS\t| PORT\t\t\t\t| COMMAND\n")
		fmt.Printf("----------------+-----------------------+---------------+-----------------------+---------------+-------------------------------+------------------------------\n")
		for _, entry := range container_list.List {
			// port info
			var port_info string = ""
			for _, port := range entry.Port {
				port_info += "0.0.0.0:" + strconv.Itoa(port.HostPort) + "->" + strconv.Itoa(port.ContainerPort) + "/" + port.Protocol + ","
			}
			port_info = strings.TrimRight(port_info, ",")
			fmt.Printf("%s\t| %-20s\t| %-12s\t| %-20s\t| %-12s\t| %-24s\t| %-32s\n", entry.Id, entry.Name, entry.Namespace, entry.Image, entry.Status, port_info, entry.Command)
		}
	}
}

func printFileSystemMount(spec karakuripkgs.ConfigSpec) {
	var count = 0
	for i, entry := range spec.Mounts {

		if entry.MountType != "" {
			count += 1
			if i != 0 {
				fmt.Printf("                  ")
			}
			fmt.Printf("[%d] ", count)
			fmt.Printf("TYPE        : %s\n", entry.MountType)
			fmt.Printf("                      SOURCE      : %s\n", entry.Source)
			fmt.Printf("                      DESTINATION : %s\n", entry.Destination)
			fmt.Printf("                      OPTIONS     : ")
			// options
			var options = ""
			for _, opt_entry := range entry.Options {
				options += opt_entry + ","
			}
			options = strings.TrimRight(options, ",")
			fmt.Printf("%s\n", options)
		}
	}
}

func printHostDirectoryMount(spec karakuripkgs.ConfigSpec) {
	var count = 0
	for _, entry := range spec.Mounts {
		if entry.MountType == "" {
			count += 1
			if count != 1 {
				fmt.Printf("               ")
			}
			fmt.Printf("[%d] ", count)
			fmt.Printf("TYPE        : bind\n")
			fmt.Printf("                      SOURCE      : %s\n", entry.Source)
			fmt.Printf("                      DESTINATION : %s\n", entry.Destination)
		}
	}
}

func printContainerSpec(spec karakuripkgs.ConfigSpec) {
	fmt.Println("[BASIC]")
	// Hostname
	fmt.Printf(" HOSTNAME        : %s\n", spec.Hostname)
	// container directory
	fmt.Printf(" CONTAINER LAYER : %s\n", spec.Root.Path)
	// image directory
	fmt.Printf(" IMAGE LAYER     : %s\n", spec.Image.Path)

	fmt.Println()

	fmt.Println("[PROCESS]")
	// process
	fmt.Printf(" PROCESS ID  : %d\n", spec.Process.Pid)
	// command
	var command = ""
	for _, entry := range spec.Process.Args {
		command += entry + " "
	}
	fmt.Printf(" COMMAND     : %s\n", command)
	// env
	fmt.Printf(" ENVIRONMENT : ")
	for i, entry := range spec.Process.Env {
		if i != 0 {
			fmt.Printf("               ")
		}
		fmt.Printf("[%-2d] ", i+1)
		fmt.Printf("%s=%s\n", entry.Key, entry.Value)
	}

	fmt.Println()

	fmt.Println("[RESOURCE]")
	// cgroup path
	fmt.Printf(" CGROUP       : %s\n", spec.Cgroup.Path)
	// cpu
	fmt.Printf(" CPU LIMIT    : %s\n", spec.Cgroup.Cpu.Max)
	// memory
	fmt.Printf(" MEMORY LIMIT : %s\n", spec.Cgroup.Memory.Max)

	fmt.Println()

	fmt.Println("[NETWORK]")
	// link device
	fmt.Printf(" LINK DEVICE  : %s\n", spec.Network.HostDevice)
	// address
	fmt.Printf(" ADDRESS      : %s\n", spec.Network.Address)
	// gateway
	fmt.Printf(" GATEWAY      : %s\n", spec.Network.Gateway)
	// nameserver
	fmt.Printf(" NAMESERVER   : %s\n", spec.Network.Nameserver)
	// port
	fmt.Printf(" PORT FORWARD : ")
	for i, entry := range spec.Network.Port {
		if i != 0 {
			fmt.Printf("               ")
		}
		fmt.Printf("[%d] ", i+1)
		// host port
		fmt.Printf("HOST PORT      : %d\n", entry.HostPort)
		fmt.Printf("                    CONTAINER PORT : %d\n", entry.TargetPort)
		fmt.Printf("                    PROTOCOL       : %s\n", entry.Protocol)
	}

	fmt.Println()
	fmt.Println()

	fmt.Println("[MOUNT]")
	// system file
	fmt.Printf(" SYSTEM FILE    : ")
	printFileSystemMount(spec)
	// host file
	fmt.Printf(" HOST DIRECTORY : ")
	printHostDirectoryMount(spec)

	fmt.Println()
}

func printImageList(image_list hitoha.ImageList) {
	cluster := getTargetCluster().Target
	fmt.Println("[Cluster] " + cluster)
	// print image list
	fmt.Printf("REPOSITORY\t\t\t\t| TAG\t\t\t| ID\n")
	fmt.Printf("----------------------------------------+-----------------------+-----------------\n")
	for _, entry := range image_list.List {
		fmt.Printf("%-35s\t| %-15s\t| %s\n", entry.Image, entry.Tag, entry.ImageId)
	}
}

func printNamespaceList(namespace_list hitoha.NamespaceList) {
	cluster := getTargetCluster().Target
	fmt.Println("[Cluster] " + cluster)
	fmt.Printf("NAMESPACE\t| NETWORK I/F\t| ADDRESS\n")
	fmt.Printf("----------------+---------------+----------------\n")
	for _, entry := range namespace_list.Namespaces {
		fmt.Printf("%-12s\t| %s\t| %s\n", entry.Name, entry.Network.Name, entry.Network.Address)
	}
}

func PrintKarakuriVersion() {
	fmt.Printf("karakuri version %s\n", karakuripkgs.KARAKURI_VERSION)
}

func ShowContainerLog(id string, name string) {
	// Check if it is an available option
	cluster := getTargetCluster().Target
	if cluster != "localhost:9806" {
		fmt.Println("'karakuri logs' command is not available for remote clusters")
		os.Exit(1)
	}
	container_id := id
	if name != "none" {
		// retrieve container id
		if res, resp_id, message := karakuripkgs.RequestContainerId(name); !res {
			fmt.Println(message)
			return
		} else {
			container_id = resp_id
		}
	}

	cmd := exec.Command("less", karakuripkgs.FUTABA_ROOT+"/"+container_id+"/container.log")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return
	}
	cmd.Wait()
}

// module
func printModuleList(mod_list karakuri_mod.ModList) {
	cluster := getTargetCluster().Target
	fmt.Println("[Cluster] " + cluster)
	fmt.Println("Enabled:")
	for _, entry := range mod_list.List {
		if entry.Status == "enable" {
			fmt.Printf("  %-15s\t# %s\n", entry.Name, entry.Description)
		}
	}

	fmt.Println("Disabled:")
	for _, entry := range mod_list.List {
		if entry.Status == "disable" {
			fmt.Printf("  %-15s\t# %s\n", entry.Name, entry.Description)
		}
	}
}

// registry controlleer
func printTargetRegistry(registry_info hitoha.RegistryInfo) {
	cluster := getTargetCluster().Target
	fmt.Println("[Cluster] " + cluster)
	fmt.Println("Registry : " + registry_info.Target)
	fmt.Println("Status   : " + registry_info.Status)
}

func printRepository(repository_list hitoha.RepogitryList) {
	cluster := getTargetCluster().Target
	fmt.Println("[Cluster] " + cluster)
	fmt.Println("REPOSITORY")
	fmt.Println("--------------------------")
	for _, entry := range repository_list.Repository {
		fmt.Println(entry)
	}
}

func printTag(repository string, tag_list hitoha.TagList) {
	cluster := getTargetCluster().Target
	fmt.Println("[Cluster] " + cluster)
	fmt.Println("REPOSITORY: " + repository)
	fmt.Println("TAG")
	fmt.Println("--------------")
	for _, entry := range tag_list.Tag {
		fmt.Println(entry)
	}
}

// cluster controller
func printTargetCluster(cluster_info ClusterInfo) {
	fmt.Println("Mode           : " + cluster_info.Mode)
	fmt.Println()
	fmt.Println("Target Cluster : " + cluster_info.Target)
	fmt.Println("Status         : " + cluster_info.Status)
}
