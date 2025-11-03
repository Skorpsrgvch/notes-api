package main

import (
	"log"
	"net/http"
	"notes-api/handlers"
	"notes-api/storage"
)

func main() {
	store, err := storage.New("notes.db")
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	defer store.Close()

	// Роутинг
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateNoteHandler(store)(w, r)
		case http.MethodGet:
			handlers.GetAllNotesHandler(store)(w, r)
		default:
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) <= len("/notes/") {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodPut:
			handlers.UpdateNoteHandler(store)(w, r)
		case http.MethodDelete:
			handlers.DeleteNoteHandler(store)(w, r)
		default:
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
