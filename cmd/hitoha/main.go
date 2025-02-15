package main

import (
	"fmt"
	"hitoha"
	"karakuripkgs"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("karakuri version " + karakuripkgs.KARAKURI_VERSION)
	// initial setup
	hitoha.SetupEnvironment()

	router := mux.NewRouter()

	// Container
	// GET
	router.HandleFunc("/container/ls/{namespace}", hitoha.GetContainerList).Methods("GET")
	router.HandleFunc("/container/spec/{id}", hitoha.GetContainerSpec).Methods("GET")
	router.HandleFunc("/container/getid/{name}", hitoha.GetContainerId).Methods("GET")
	// POST
	router.HandleFunc("/container/create/{image}/{port}/{mount}/{cmd}/{registry}/{name}/{namespace}/{restart}", hitoha.PostCreateContainer).Methods("POST")
	router.HandleFunc("/container/start/{id}/{terminal}", hitoha.PostStartContainer).Methods("POST")
	router.HandleFunc("/container/run/{image}/{port}/{mount}/{cmd}/{registry}/{name}/{namespace}/{restart}/{terminal}", hitoha.PostRunContainer).Methods("POST")
	router.HandleFunc("/container/exec/{id}/{cmd}/{terminal}", hitoha.PostExecContainer).Methods("POST")
	router.HandleFunc("/container/kill/{id}", hitoha.PostKillContainer).Methods("POST")
	// DELETE
	router.HandleFunc("/container/delete/{id}", hitoha.DeleteDeleteContainer).Methods("DELETE")

	// Image
	// GET
	router.HandleFunc("/image/ls", hitoha.GetShowImages).Methods("GET")
	router.HandleFunc("/image/pull/{image-tag}/{os-arch}/{registry}", hitoha.GetPullImage).Methods("GET")
	// POST
	router.HandleFunc("/image/push/{image-tag}/{registry}", hitoha.PostPushImage).Methods("POST")
	// DELETE
	router.HandleFunc("/image/delete/{id}", hitoha.DeleteDeleteImage).Methods("DELETE")

	// namespcae
	// GET
	router.HandleFunc("/namespace/ls", hitoha.GetNamespaceList).Methods("GET")
	// POST
	router.HandleFunc("/namespace/create/{namespace}", hitoha.PostNamespace).Methods("POST")
	// DELETE
	router.HandleFunc("/namespace/delete/{namespace}", hitoha.DeleteNamespace).Methods("DELETE")

	// module
	// GET
	router.HandleFunc("/mod/list", hitoha.GetModuleList).Methods("GET")
	// POST
	router.HandleFunc("/mod/enable/{mod_name}", hitoha.PostEnableModule).Methods("POST")
	// DELETE
	router.HandleFunc("/mod/disable/{mod_name}", hitoha.DeleteDisableModule).Methods("DELETE")

	// registry controller
	// GET
	router.HandleFunc("/reg/target", hitoha.GetTargetRegistry).Methods("GET")
	router.HandleFunc("/reg/repository", hitoha.GetShowRepositories).Methods("GET")
	router.HandleFunc("/reg/tag/{repository}", hitoha.GetShowTags).Methods("GET")
	// POST
	router.HandleFunc("/reg/connect/{registry}", hitoha.PostConnectRegistry).Methods("POST")
	// DELETE
	router.HandleFunc("/reg/disconnect", hitoha.DeleteDisconnectRegistry).Methods("DELETE")
	router.HandleFunc("/reg/delete/{image-tag}", hitoha.DeleteDeleteManifest).Methods("DELETE")

	// execute server
	// local
	go func() {
		fmt.Println("Listen on \"127.0.0.1:9806\" ...")
		http.ListenAndServe("127.0.0.1:9806", router)
	}()
	// listen for cluster controller
	go func() {
		fmt.Println("Listen on \"0.0.0.0:9816\" ...")
		http.ListenAndServe("0.0.0.0:9816", router)
	}()

	select {}
}
