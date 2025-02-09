package karakuri_mod

import (
	"karakuripkgs"
	"os"
	"os/exec"
)

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
	createRegistryContainer(mod_info)
	startRegistryContainer(mod_info)
}

func disableRegistryModule() {
	removeModuleContainer("registry")
}
