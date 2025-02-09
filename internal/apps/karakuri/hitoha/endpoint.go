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
	image := params["image"]
	// port
	port := params["port"]
	// mount
	mount := strings.Replace(params["mount"], "-", "/", -1)
	// cmd
	cmd := strings.Replace(params["cmd"], "!", "/", -1)
	// repositry
	repositry := params["repositry"]
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

	//resp := CreateContainer(image, port, mount, cmd, repositry)
	resp := CreateContainer(ParamsCreateContainer{
		ImageInfo: image,
		Name:      name,
		Namespace: namespace,
		Port:      port,
		Mount:     mount,
		Cmd:       cmd,
		Repositry: repositry,
	})

	json.NewEncoder(w).Encode(resp)
}

func PostStartContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// id
	id := params["id"]

	resp := StartContainer(id)

	json.NewEncoder(w).Encode(resp)
}

func PostRunContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// image
	image := params["image"]
	// port
	port := params["port"]
	// mount
	mount := strings.Replace(params["mount"], "-", "/", -1)
	// cmd
	cmd := strings.Replace(params["cmd"], "-", "/", -1)
	// repositry
	repositry := params["repositry"]
	// name
	name := params["name"]
	// namespace
	namespace := params["namespace"]
	if namespace == "none" {
		namespace = "default"
	}

	//resp := RunContainer(image, port, mount, cmd, repositry)
	resp := RunContainer(ParamsRunContainer{
		ImageInfo: image,
		Name:      name,
		Namespace: namespace,
		Port:      port,
		Mount:     mount,
		Cmd:       cmd,
		Repositry: repositry,
	})

	json.NewEncoder(w).Encode(resp)
}

func PostExecContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// retrieve parameter
	params := mux.Vars(r)
	// id
	id := params["id"]

	resp := ExecContainer(id)

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
	// repositry
	repositry := params["repositry"]

	resp := PullImage(image_tag, os_arch, repositry)

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
