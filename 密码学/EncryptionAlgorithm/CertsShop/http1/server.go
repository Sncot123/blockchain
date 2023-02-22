package main

import (
	"fmt"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello go")
}
func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServeTLS(":8080",
		"D:\\JetBrains\\GoLand_workspace\\go-study\\blockchain\\EncryptionAlgorithm\\CertsShop\\ca\\server\\server.crt",
		"D:\\JetBrains\\GoLand_workspace\\go-study\\blockchain\\EncryptionAlgorithm\\CertsShop\\ca\\server\\server.key",
		nil)
}
