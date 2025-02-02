package futaba

import (
	"encoding/json"
	"fmt"
	"karakuripkgs"
	"os"
)

type ContainerInfo struct {
	Id     string `json:"id"`
	Bundle string `json:"bundle"`
	Image  string `json:"image"`
}

type ContainerList struct {
	List []ContainerInfo `json:"list"`
}

func newContainerDirectory(spec karakuripkgs.ConfigSpec) {
	root_path := spec.Root.Path
	diff := root_path + "/diff"
	work := root_path + "/work"
	merged := root_path + "/merged"
	pipe := spec.Fifo

	if err := os.MkdirAll(diff, os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(work, os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(merged, os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(pipe, os.ModePerm); err != nil {
		panic(err)
	}
}

func NewContainerList() {
	var container_list ContainerList
	file, _ := json.MarshalIndent(container_list, "", "  ")
	if err := os.WriteFile(karakuripkgs.FUTABA_CONTAINER_LIST, file, os.ModePerm); err != nil {
		panic(err)
	}
}

func addContainerList(spec karakuripkgs.ConfigSpec) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.FUTABA_CONTAINER_LIST)
	if err != nil {
		NewContainerList()
		bytes, _ = os.ReadFile(karakuripkgs.FUTABA_CONTAINER_LIST)
	}

	var container_list ContainerList
	if err := json.Unmarshal(bytes, &container_list); err != nil {
		panic(err)
	}

	// set container info
	var container_info = ContainerInfo{
		Id:     spec.Hostname,
		Bundle: spec.Root.Path,
		Image:  spec.Image.Path,
	}
	container_list.List = append(container_list.List, container_info)

	file, _ := json.MarshalIndent(container_list, "", "  ")
	if err := os.WriteFile(karakuripkgs.FUTABA_CONTAINER_LIST, file, os.ModePerm); err != nil {
		panic(err)
	}
}

func createNewContainer(spec string) string {
	config_spec := karakuripkgs.ReadSpecFile(spec)

	container_id := config_spec.Hostname

	newContainerDirectory(config_spec)

	addContainerList(config_spec)

	// move spec file
	if err := os.Rename(spec+"/config.json", config_spec.Root.Path+"/config.json"); err != nil {
		panic(err)
	}

	return container_id
}

func CreateContainer(spec string) string {
	// create container
	container_id := createNewContainer(spec)

	container_dir := karakuripkgs.FUTABA_ROOT + "/" + container_id
	// read spec file
	config_spec := karakuripkgs.ReadSpecFile(container_dir)

	// create fifo
	createFifo(config_spec.Fifo)

	return container_id
}

func DeleteContainer(id string) {
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.FUTABA_CONTAINER_LIST)
	if err != nil {
		return
	}

	var container_list ContainerList
	if err := json.Unmarshal(bytes, &container_list); err != nil {
		panic(err)
	}

	var new_container_list ContainerList
	for _, entry := range container_list.List {
		if entry.Id == id {
			if err := os.RemoveAll(entry.Bundle); err != nil {
				panic(err)
			}
		} else {
			new_container_list.List = append(new_container_list.List, entry)
		}
	}

	file, _ := json.MarshalIndent(new_container_list, "", "  ")
	if err := os.WriteFile(karakuripkgs.FUTABA_CONTAINER_LIST, file, os.ModePerm); err != nil {
		panic(err)
	}

	// delete cgroup
	deleteCgroup(id)
}

func ShowContainerList() {
	fmt.Printf("ID\t\tBUNDLE\t\t\tIMAGE\n")
	var bytes []byte
	bytes, err := os.ReadFile(karakuripkgs.FUTABA_CONTAINER_LIST)
	if err != nil {
		return
	}

	var container_list ContainerList
	if err := json.Unmarshal(bytes, &container_list); err != nil {
		panic(err)
	}

	for _, entry := range container_list.List {
		fmt.Printf("%s\t%s\t%s\t\n", entry.Id, entry.Bundle, entry.Image)
	}
}
