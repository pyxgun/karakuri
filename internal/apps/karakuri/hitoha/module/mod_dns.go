package karakuri_mod

import (
	"bufio"
	"fmt"
	"karakuripkgs"
	"os"
	"os/exec"
	"strings"
)

func createDnsModuleKarakurifile() {
	// create Karakurifile
	karakurifile := `FROM alpine

RUN apk update
RUN apk add bind --no-cache

CMD ["/usr/sbin/named", "-c", "/conf/named.conf", "-g"]
`

	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_DNS + "/Karakurifile"); stat != nil {
		fd, err := os.Create(karakuripkgs.KARAKURI_MOD_DNS + "/Karakurifile")
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		bytes := []byte(karakurifile)
		if _, err := fd.Write(bytes); err != nil {
			panic(err)
		}
	}
}

func createDnsModuleConfig() {
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_DNS + "/conf"); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_DNS+"/conf", os.ModePerm); err != nil {
			panic(err)
		}
	}

	named_conf := `acl localnet {
    127.0.0.1;
    10.157.0.0/16;
};

options {
    version "unknown";
    directory "/var/bind";
    pid-file "/var/run/named/named.pid";
    recursion yes;
    notify no;
 
    listen-on { any; };
    listen-on-v6 { none; };

    allow-query { localnet; };
    allow-query-cache { localnet; };
    allow-recursion { localnet; };
    allow-transfer { none; };

    forwarders { 8.8.8.8; };
};

zone "karakuri.container" IN {
    type master;
    file "/conf/karakuri.container.zone";
};
`

	karakuri_container_zone := `$TTL 1h
@ IN SOA ns.karakuri.container root.karakuri.container. (
    2025012401 ; serial
    1h         ; refresh
    15m        ; retry
    1d         ; expire
    1h         ; minimum
);
 
@          IN NS ns.karakuri.container.
`

	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_DNS + "/conf/named.conf"); stat != nil {
		fd_1, err := os.Create(karakuripkgs.KARAKURI_MOD_DNS + "/conf/named.conf")
		if err != nil {
			panic(err)
		}
		defer fd_1.Close()
		bytes_1 := []byte(named_conf)
		if _, err := fd_1.Write(bytes_1); err != nil {
			panic(err)
		}
	}

	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_DNS + "/conf/karakuri.container.zone"); stat != nil {
		fd_2, err := os.Create(karakuripkgs.KARAKURI_MOD_DNS + "/conf/karakuri.container.zone")
		if err != nil {
			panic(err)
		}
		defer fd_2.Close()
		bytes_2 := []byte(karakuri_container_zone)
		if _, err := fd_2.Write(bytes_2); err != nil {
			panic(err)
		}
	}
}

func setupDnsModule() {
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_DNS); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_DNS, os.ModePerm); err != nil {
			panic(err)
		}
	}

	// Karakurifile
	createDnsModuleKarakurifile()
	// dns conf
	createDnsModuleConfig()
}

func buildDnsImage(mod_info ModInfo) {
	build_args := []string{"build", "--name", mod_info.ImageName, "--buildpath", mod_info.Path}
	build := exec.Command("karakuri", build_args...)
	if err := build.Run(); err != nil {
		panic(err)
	}
}

func createDnsContainer(mod_info ModInfo) {
	create_args := []string{"create", "--name", mod_info.Name, "--image", mod_info.ImageName, "--mount", mod_info.Path + "/conf:/conf", "--ns", "system-mod"}
	create := exec.Command("karakuri", create_args...)
	if err := create.Run(); err != nil {
		panic(err)
	}
}

func addNsRecord(mod_info ModInfo) {
	container_id := karakuripkgs.RetrieveContainerId("", mod_info.Name)
	config_spec := karakuripkgs.ReadSpecFile(karakuripkgs.FUTABA_ROOT + "/" + container_id)

	fd, err := os.OpenFile(karakuripkgs.KARAKURI_MOD_DNS+"/conf/karakuri.container.zone", os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	address := (strings.Split(config_spec.Network.Address, "/"))[0]
	dns_record := "\nns IN A " + address

	fmt.Fprintln(fd, dns_record)
}

func startDnsContainer(mod_info ModInfo) {
	start_args := []string{"start", "--name", mod_info.Name}
	start := exec.Command("karakuri", start_args...)
	if err := start.Run(); err != nil {
		panic(err)
	}
}

func enableDnsModule(mod_info ModInfo) {
	image_info := strings.Split(mod_info.ImageName, ":")
	if !isImageExists(image_info[0], image_info[1]) {
		buildDnsImage(mod_info)
	}
	createDnsContainer(mod_info)
	addNsRecord(mod_info)
	startDnsContainer(mod_info)
}

func disableDnsModule() {
	file, err := os.Open(karakuripkgs.KARAKURI_MOD_DNS + "/conf/karakuri.container.zone")
	if err != nil {
		panic(err)
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line_str := scanner.Text()
		if line_str != "" {
			lines = append(lines, line_str)
		}
	}
	file.Close()

	var new_zone_line []string
	for _, entry := range lines {
		record_info := strings.Split(entry, " ")
		if record_info[0] != "ns" {
			new_zone_line = append(new_zone_line, entry)
		}
	}
	data := strings.Join(new_zone_line, "\n")
	bytes := []byte(data)

	if err := os.WriteFile(karakuripkgs.KARAKURI_MOD_DNS+"/conf/karakuri.container.zone", bytes, os.ModePerm); err != nil {
		panic(err)
	}

	// remove container
	removeModuleContainer("dns")
}

// add entry
func AddDnsRecord(container_id string, container_address string) {
	fd, err := os.OpenFile(karakuripkgs.KARAKURI_MOD_DNS+"/conf/karakuri.container.zone", os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	address := (strings.Split(container_address, "/"))[0]
	dns_record := container_id + " IN A " + address + "\n"

	fmt.Fprint(fd, dns_record)

	// restart dns module
	restart_args := []string{"restart", "--name", "dns"}
	restart := exec.Command("karakuri", restart_args...)
	if err := restart.Run(); err != nil {
		panic(err)
	}
}

func DeleteDnsRecord(container_id string) {
	file, err := os.Open(karakuripkgs.KARAKURI_MOD_DNS + "/conf/karakuri.container.zone")
	if err != nil {
		panic(err)
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line_str := scanner.Text()
		if line_str != "" {
			lines = append(lines, line_str)
		}
	}
	file.Close()

	var new_zone_line []string
	for _, entry := range lines {
		record_info := strings.Split(entry, " ")
		if record_info[0] != container_id {
			new_zone_line = append(new_zone_line, entry)
		}
	}
	data := strings.Join(new_zone_line, "\n")
	data = data + "\n"
	bytes := []byte(data)

	if err := os.WriteFile(karakuripkgs.KARAKURI_MOD_DNS+"/conf/karakuri.container.zone", bytes, os.ModePerm); err != nil {
		panic(err)
	}

	// restart dns module
	restart_args := []string{"restart", "--name", "dns"}
	restart := exec.Command("karakuri", restart_args...)
	if err := restart.Run(); err != nil {
		panic(err)
	}
}
