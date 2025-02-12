package karakuri_mod

import (
	"karakuripkgs"
	"os"
	"os/exec"
	"strings"
)

func createRegistryModuleKarakurifile() {
	// create Karakurifile
	karakurifile := `FROM registry:latest

ENV REGISTRY_STORAGE_DELETE_ENABLED=true

CMD ["/entrypoint.sh", "/etc/docker/registry/config.yml"]
`

	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_REGISTRY + "/Karakurifile"); stat != nil {
		fd, err := os.Create(karakuripkgs.KARAKURI_MOD_REGISTRY + "/Karakurifile")
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

func setupRegistryModule() {
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_REGISTRY); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_REGISTRY, os.ModePerm); err != nil {
			panic(err)
		}
	}

	// volume mount
	if _, stat := os.Stat(karakuripkgs.KARAKURI_MOD_REGISTRY + "/data"); stat != nil {
		if err := os.MkdirAll(karakuripkgs.KARAKURI_MOD_REGISTRY+"/data", os.ModePerm); err != nil {
			panic(err)
		}
	}

	// Karakurifile
	createRegistryModuleKarakurifile()
}

func buildRegistryImage(mod_info ModInfo) {
	build_args := []string{"build", "--name", mod_info.ImageName, "--buildpath", mod_info.Path}
	build := exec.Command("karakuri", build_args...)
	if err := build.Run(); err != nil {
		panic(err)
	}
}

func createRegistryContainer(mod_info ModInfo) {
	create_args := []string{
		"create",
		"--name", mod_info.Name,
		"--image", mod_info.ImageName,
		"--mount", mod_info.Path + "/data:/var/lib/registry",
		"--port", "5000:5000:tcp",
		"--restart", "on-boot",
		"--ns", "system-mod",
	}
	create := exec.Command("karakuri", create_args...)
	if err := create.Run(); err != nil {
		panic(err)
	}
}

func startRegistryContainer(mod_info ModInfo) {
	start_args := []string{"start", "--name", mod_info.Name}
	start := exec.Command("karakuri", start_args...)
	if err := start.Run(); err != nil {
		panic(err)
	}
}

func enableRegistryModule(mod_info ModInfo) {
	image_info := strings.Split(mod_info.ImageName, ":")
	if !isImageExists(image_info[0], image_info[1]) {
		buildRegistryImage(mod_info)
	}
	createRegistryContainer(mod_info)
	startRegistryContainer(mod_info)
}

func disableRegistryModule() {
	removeModuleContainer("registry")
}
