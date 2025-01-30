package main

import (
	"hitoha"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// initial setup
	hitoha.SetupEnvironment()

	router := mux.NewRouter()

	// Container
	// GET
	router.HandleFunc("/container/ls/{namespace}", hitoha.GetContainerList).Methods("GET")
	router.HandleFunc("/container/spec/{id}", hitoha.GetContainerSpec).Methods("GET")
	router.HandleFunc("/container/getid/{name}", hitoha.GetContainerId).Methods("GET")
	// POST
	router.HandleFunc("/container/create/{image}/{port}/{mount}/{cmd}/{repositry}/{name}/{namespace}", hitoha.PostCreateContainer).Methods("POST")
	router.HandleFunc("/container/start/{id}", hitoha.PostStartContainer).Methods("POST")
	router.HandleFunc("/container/run/{image}/{port}/{mount}/{cmd}/{repositry}/{name}/{namespace}", hitoha.PostRunContainer).Methods("POST")
	router.HandleFunc("/container/exec/{id}", hitoha.PostExecContainer).Methods("POST")
	router.HandleFunc("/container/kill/{id}", hitoha.PostKillContainer).Methods("POST")
	// DELETE
	router.HandleFunc("/container/delete/{id}", hitoha.DeleteDeleteContainer).Methods("DELETE")

	// Image
	// GET
	router.HandleFunc("/image/ls", hitoha.GetShowImages).Methods("GET")
	router.HandleFunc("/image/pull/{image-tag}/{os-arch}/{repositry}", hitoha.GetPullImage).Methods("GET")
	// DELETE
	router.HandleFunc("/image/delete/{id}", hitoha.DeleteDeleteImage).Methods("DELETE")

	// namespcae
	// GET
	router.HandleFunc("/namespace/ls", hitoha.GetNamespaceList).Methods("GET")
	// POST
	router.HandleFunc("/namespace/create/{namespace}", hitoha.PostNamespace).Methods("POST")
	// DELETE
	router.HandleFunc("/namespace/delete/{namespace}", hitoha.DeleteNamespace).Methods("DELETE")

	// execute server
	http.ListenAndServe("127.0.0.1:9876", router)
}
