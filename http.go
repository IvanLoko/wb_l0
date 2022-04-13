package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	t "wb/subs"
)

type Todo struct {
	Name string
	Done bool
}

func IsNotDone(todo Todo) bool {
	return !todo.Done
}

// Запускаем локальные сервер http на http://localhost:8081
func Start_http() {
	log.Println("Server started on http://localhost:8081")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			id := strings.TrimPrefix(r.URL.Path, "/")
			// проверка на случай пустого поля в запросе
			if id == "" {
				http.Error(w, "id must be string", http.StatusBadRequest)
			} else {
				w.Header().Set("Content-Type", "application/json")
				//переводим в байты соответствующий json-файл из кеша
				json_resp, err := json.Marshal(t.Cash_data[id])

				if err != nil {
					http.Error(w, err.Error(), http.StatusConflict)
					return
				}

				w.WriteHeader(http.StatusOK)
				//отправляем на сервер
				w.Write(json_resp)
			}

		}
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
