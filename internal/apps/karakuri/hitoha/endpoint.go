package hitoha

import (
	"encoding/json"
	"karakuri_mod"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// container
// GET
// show container list
func GetContainerList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// namespace
	namespace := params["namespace"]
	if namespace == "none" {
		namespace = "default"
	}

	resp := ShowContainerList(namespace)

	json.NewEncoder(w).Encode(resp)
}

// show container spec
func GetContainerSpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// id
	id := params["id"]

	resp := ShowContainerSpec(id)

	json.NewEncoder(w).Encode(resp)
}

func GetContainerId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// name
	name := params["name"]

	resp := retrieveContainerId(name)

	json.NewEncoder(w).Encode(resp)
}

// POST
// create container
func PostCreateContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// image
	image := strings.Replace(params["image"], "!", "/", -1)
	// port
	port := params["port"]
	// mount
	mount := strings.Replace(params["mount"], "-", "/", -1)
	// cmd
	cmd := strings.Replace(params["cmd"], "!", "/", -1)
	// registry
	registry := params["registry"]
	// name
	name := params["name"]
	if name == "none" {
		name = ""
	}
	// namespace
	namespace := params["namespace"]
	if namespace == "none" {
		namespace = "default"
	}
	// restart
	restart := params["restart"]

	//resp := CreateContainer(image, port, mount, cmd, registry)
	resp := CreateContainer(ParamsCreateContainer{
		ImageInfo: image,
		Name:      name,
		Namespace: namespace,
		Port:      port,
		Mount:     mount,
		Cmd:       cmd,
		Registry:  registry,
		Restart:   restart,
	})

	json.NewEncoder(w).Encode(resp)
}

func PostStartContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// id
	id := params["id"]
	terminal := params["terminal"]

	resp := StartContainer(id, terminal)

	json.NewEncoder(w).Encode(resp)
}

func PostRunContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// image
	image := strings.Replace(params["image"], "!", "/", -1)
	// port
	port := params["port"]
	// mount
	mount := strings.Replace(params["mount"], "-", "/", -1)
	// cmd
	cmd := strings.Replace(params["cmd"], "-", "/", -1)
	// registry
	registry := params["registry"]
	// name
	name := params["name"]
	// namespace
	namespace := params["namespace"]
	if namespace == "none" {
		namespace = "default"
	}
	// restart
	restart := params["restart"]
	// terminal
	terminal := params["terminal"]

	//resp := RunContainer(image, port, mount, cmd, registry)
	resp := RunContainer(ParamsRunContainer{
		ImageInfo: image,
		Name:      name,
		Namespace: namespace,
		Port:      port,
		Mount:     mount,
		Cmd:       cmd,
		Registry:  registry,
		Restart:   restart,
		Terminal:  terminal,
	})

	json.NewEncoder(w).Encode(resp)
}

func PostExecContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// id
	id := params["id"]
	cmd := strings.Replace(params["cmd"], "-", "/", -1)
	terminal := params["terminal"]

	resp := ExecContainer(id, terminal, cmd)

	json.NewEncoder(w).Encode(resp)
}

func PostKillContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// id
	id := params["id"]

	resp := KillContainer(id)

	json.NewEncoder(w).Encode(resp)
}

// DELETE
func DeleteDeleteContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// id
	id := params["id"]

	resp := DeleteContainer(id)

	json.NewEncoder(w).Encode(resp)
}

// image
// GET
// show images
func GetShowImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := ShowImageList()

	json.NewEncoder(w).Encode(resp)
}

// pull image
func GetPullImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// image/tag
	image_tag := strings.Replace(params["image-tag"], "!", "/", -1)
	// os/arch
	os_arch := params["os-arch"]
	// registry
	registry := params["registry"]

	resp := PullImage(image_tag, os_arch, registry)

	json.NewEncoder(w).Encode(resp)
}

// pull image
func PostPushImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// image/tag
	image_tag := strings.Replace(params["image-tag"], "!", "/", -1)
	// registry
	registry := params["registry"]

	resp := PushImage(registry, image_tag)

	json.NewEncoder(w).Encode(resp)
}

// DELET
// delete image
func DeleteDeleteImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// id
	id := params["id"]

	resp := DeleteImage(id)

	json.NewEncoder(w).Encode(resp)
}

// namespace
// GET
func GetNamespaceList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := showNamespaceList()

	json.NewEncoder(w).Encode(resp)
}

// POST
func PostNamespace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// namespace
	namespace := params["namespace"]

	resp := createNewNamespace(namespace)

	json.NewEncoder(w).Encode(resp)
}

// DELETE
func DeleteNamespace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// namespace
	namespace := params["namespace"]

	resp := deleteNamespace(namespace)

	json.NewEncoder(w).Encode(resp)
}

// module
// GET
func GetModuleList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := karakuri_mod.ShowModuleList()

	json.NewEncoder(w).Encode(resp)
}

// POST
func PostEnableModule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// namespace
	mod_name := params["mod_name"]

	resp := karakuri_mod.EnableModule(mod_name)

	json.NewEncoder(w).Encode(resp)
}

// DELETE
func DeleteDisableModule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// namespace
	mod_name := params["mod_name"]

	resp := karakuri_mod.DisableModule(mod_name)

	json.NewEncoder(w).Encode(resp)
}

// registry controller
// GET
// target registry
func GetTargetRegistry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := ShowTargetRegistry()

	json.NewEncoder(w).Encode(resp)
}

// get repositories
func GetShowRepositories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := ShowRepositories()

	json.NewEncoder(w).Encode(resp)
}

// get tags
func GetShowTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// namespace
	repository := strings.Replace(params["repository"], "!", "/", -1)

	resp := ShowTags(repository)

	json.NewEncoder(w).Encode(resp)
}

// POST
// connect registry
func PostConnectRegistry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// namespace
	registry := params["registry"]

	resp := ConnectRegistry(registry)

	json.NewEncoder(w).Encode(resp)
}

// DELETE
// disconnect registry
func DeleteDisconnectRegistry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := DisconnectRegistry()

	json.NewEncoder(w).Encode(resp)
}

// delete manifest
func DeleteDeleteManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// image/tag
	image_tag := strings.Replace(params["image-tag"], "!", "/", -1)

	resp := DeleteImageManifest(image_tag)

	json.NewEncoder(w).Encode(resp)
}
