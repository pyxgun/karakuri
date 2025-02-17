package karakuri_mod

import (
	"encoding/json"
	"karakuripkgs"
	"os"
	"os/exec"
)

type ModInfo struct {
	Name        string `json:"module_name"`
	ImageName   string `json:"image_name"`
	Path        string `json:"path"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type ModList struct {
	List []ModInfo `json:"list"`
}

// image list
type ImageInfo struct {
	Image   string `json:"image"`
	Tag     string `json:"tag"`
	ImageId string `json:"image_id"`
	Rootfs  string `json:"rootfs"`
}

type ImageList struct {
	List []ImageInfo `json:"list"`
}

func SetupModules() {
	setupDnsModule()
	setupRegistryModule()
	setupRegistryBrowserModule()
	setupIngressModule()
}

func NewModList() {
	var mod_list_data ModList

	// module: dns
	mod_list_data.List = append(mod_list_data.List,
		ModInfo{
			Name:        "dns",
			ImageName:   "karakuri/dns:latest",
			Path:        karakuripkgs.KARAKURI_MOD_DNS,
			Status:      "disable",
			Description: "core DNS",
		},
	)

	// module: registry
	mod_list_data.List = append(mod_list_data.List,
		ModInfo{
			Name:        "registry",
			ImageName:   "karakuri/registry:latest",
			Path:        karakuripkgs.KARAKURI_MOD_REGISTRY,
			Status:      "disable",
			Description: "private registry listen on 5000",
		},
	)

	// module: registry-browser
	mod_list_data.List = append(mod_list_data.List,
		ModInfo{
			Name:        "registry-browser",
			ImageName:   "karakuri/registry-browser:latest",
			Path:        karakuripkgs.KARAKURI_MOD_REGISTRY_BROWSER,
			Status:      "disable",
			Description: "registry browser listen on 8081",
		},
	)

	// module: ingress
	mod_list_data.List = append(mod_list_data.List,
		ModInfo{
			Name:        "ingress",
			ImageName:   "karakuri/ingress:latest",
			Path:        karakuripkgs.KARAKURI_MOD_INGRESS,
			Status:      "disable",
			Description: "ingress controller for external access listen on 443",
		},
	)

	data, _ := json.MarshalIndent(mod_list_data, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_MOD_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}
}

func IsModuleEnabled(mod_name string) bool {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_MOD_LIST)
	if err != nil {
		panic(err)
	}

	var mod_list ModList
	if err := json.Unmarshal(bytes, &mod_list); err != nil {
		panic(err)
	}

	for _, entry := range mod_list.List {
		if entry.Name == mod_name {
			if entry.Status == "enable" {
				return true
			}
		}
	}
	return false
}

func isImageExists(image string, tag string) bool {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.HITOHA_IMAGE_LIST)
	if err != nil {
		panic(err)
	}

	var image_list_data ImageList
	if err := json.Unmarshal(bytes, &image_list_data); err != nil {
		panic(err)
	}

	for _, entry := range image_list_data.List {
		if entry.Image == image && entry.Tag == tag {
			return true
		}
	}
	return false
}

func EnableModule(mod_name string) ResponseEnableModule {
	// status check
	if IsModuleEnabled(mod_name) {
		return CreateResponseEnableModule("error", "Module: "+mod_name+" is already enabled.")
	}

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_MOD_LIST)
	if err != nil {
		panic(err)
	}

	var mod_list ModList
	if err := json.Unmarshal(bytes, &mod_list); err != nil {
		panic(err)
	}

	en_flag := false
	for i, entry := range mod_list.List {
		if entry.Name == mod_name {
			switch mod_name {
			case "dns":
				enableDnsModule(entry)
				en_flag = true
			case "registry":
				enableRegistryModule(entry)
				en_flag = true
			case "registry-browser":
				enableRegistryBrowserModule(entry)
				en_flag = true
			case "ingress":
				enableIngressModule(entry)
				en_flag = true
			}
			// update status
			mod_list.List[i].Status = "enable"
		}
	}

	if !en_flag {
		return CreateResponseEnableModule("error", "No such module, name: "+mod_name+".")
	}

	data, _ := json.MarshalIndent(mod_list, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_MOD_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}

	return CreateResponseEnableModule("success", "Enable module: "+mod_name+" success.")
}

func removeModuleContainer(mod_name string) {
	stop_args := []string{"stop", "--name", mod_name}
	stop := exec.Command("karakuri", stop_args...)
	if err := stop.Run(); err != nil {
		panic(err)
	}

	// remove container
	remove_args := []string{"rm", "--name", mod_name}
	remove := exec.Command("karakuri", remove_args...)
	if err := remove.Run(); err != nil {
		panic(err)
	}
}

func DisableModule(mod_name string) ResponseDisableModule {
	// status check
	if !IsModuleEnabled(mod_name) {
		return CreateResponseDisableModule("error", "Module: "+mod_name+" is still not enabled.")
	}

	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_MOD_LIST)
	if err != nil {
		panic(err)
	}

	var mod_list ModList
	if err := json.Unmarshal(bytes, &mod_list); err != nil {
		panic(err)
	}

	dis_flag := false
	for i, entry := range mod_list.List {
		if entry.Name == mod_name {
			switch mod_name {
			case "dns":
				disableDnsModule()
			case "registry":
				disableRegistryModule()
			case "registry-browser":
				disableRegistryBrowserModule()
			case "ingress":
				disableIngressModule()
			}
			// update status
			mod_list.List[i].Status = "disable"
			dis_flag = true
		}
	}

	if !dis_flag {
		return CreateResponseDisableModule("error", "No such module, name: "+mod_name+".")
	}

	data, _ := json.MarshalIndent(mod_list, "", "  ")
	if err := os.WriteFile(karakuripkgs.KARAKURI_MOD_LIST, data, os.ModePerm); err != nil {
		panic(err)
	}

	return CreateResponseDisableModule("success", "Disable module: "+mod_name+" success.")
}

func ShowModuleList() ResponseModuleList {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.KARAKURI_MOD_LIST)
	if err != nil {
		panic(err)
	}

	var mod_list ModList
	if err := json.Unmarshal(bytes, &mod_list); err != nil {
		panic(err)
	}

	return CreateResponseModuleList("success", mod_list)
}
