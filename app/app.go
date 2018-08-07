package app

import (
	"fmt"
	"net/http"
)

func Start(addr string) error {
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	return http.ListenAndServe(addr, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	sess, err := SessaoIni(w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintln(w, sess.Valores)
}
