package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/watcher",watcher2)
	http.ListenAndServe(":8001", nil)
}


// 做一个Linux 监控系统
func watcher2(w http.ResponseWriter, r *http.Request) {
	for {
		select {
		default:
			w.WriteHeader(202)
			w.Write([]byte("64 bytes or fewer"));
		}
	}
}
