package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

type todo struct {
	ID      string
	Content string
}

var allTodos = new(sync.Map)
var idCounter int32

func returnOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func returnJson(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
func marshalMapToJson(m *sync.Map) ([]byte, error) {
	all := make(map[string]todo)
	m.Range(func(k, v interface{}) bool {
		all[k.(string)] = v.(todo)
		return true
	})
	return json.Marshal(all)
}

func server(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		jsonRes, err := marshalMapToJson(allTodos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(string(jsonRes))
		returnJson(w, jsonRes)
		break
	case http.MethodPost:
		bt, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		td := todo{}
		json.Unmarshal(bt, &td)
		newID := fmt.Sprintf("%d", atomic.AddInt32(&idCounter, 1))
		log.Printf("saving %s", td)
		allTodos.Store(newID, todo{Content: td.Content, ID: newID})
		log.Println(allTodos)
		returnOk(w)
		break
	case http.MethodDelete:
		parts := strings.Split(req.URL.Path, "/")
		ID := parts[len(parts)-1]
		if _, ok := allTodos.Load(ID); !ok {
			http.Error(w, fmt.Sprintf("todo with id %s was not found", ID), http.StatusBadRequest)
			return
		}
		allTodos.Delete(ID)
		returnOk(w)
		break
	default:
		http.Error(w, fmt.Sprintf("not valid http method"), http.StatusBadRequest)
	}

}

func main() {
	argPort := flag.String("port", "8090", "the port")
	flag.Parse()
	http.HandleFunc("/api/todo", server)
	port := fmt.Sprintf(":%s", *argPort)
	log.Printf("running server localhost%s", port)
	err := http.ListenAndServeTLS(port, "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal(err)
	}
}
