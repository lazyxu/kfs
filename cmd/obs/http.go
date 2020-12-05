package main

import (
	"io"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/lazyxu/kfs/storage"
)

const httpPort = 9099

func serverHttp(s storage.Storage) {
	http.HandleFunc("/api/obs/write", func(w http.ResponseWriter, r *http.Request) {
		typ := r.Header.Get("kfs_type")
		typInt, err := strconv.Atoi(typ)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid type: typ"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = s.Write(typInt, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/api/obs/read", func(w http.ResponseWriter, r *http.Request) {
		typ := r.Header.Get("kfs_type")
		typInt, err := strconv.Atoi(typ)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid type: typ"))
			return
		}
		hash := r.Header.Get("kfs_hash")
		reader, err := s.Read(typInt, hash)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = io.Copy(w, reader)
		if err != nil {
			logrus.Error(err)
		}
	})
	http.HandleFunc("/api/obs/delete", func(w http.ResponseWriter, r *http.Request) {
		typ := r.Header.Get("kfs_type")
		typInt, err := strconv.Atoi(typ)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid type: typ"))
			return
		}
		hash := r.Header.Get("kfs_hash")
		err = s.Delete(typInt, hash)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	http.ListenAndServe(":9099", nil)
}
