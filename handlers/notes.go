package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"notes-api/models"
	"notes-api/storage"
	"strconv"
)

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func CreateNoteHandler(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		var req CreateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}

		if req.Title == "" || req.Content == "" {
			http.Error(w, "Заголовок и контент обязательны", http.StatusBadRequest)
			return
		}

		id, err := storage.CreateNote(req.Title, req.Content)
		if err != nil {
			log.Printf("Ошибка при создании заметки: %v", err)
			http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
			return
		}

		note := models.Note{
			Id:      int(id),
			Title:   req.Title,
			Content: req.Content,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

func GetNoteHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Path[len("/notes/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		note, err := s.GetNoteByID(id)
		if err != nil {
			log.Printf("Ошибка при получении заметки %d: %v", id, err)
			http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
			return
		}

		if note == nil {
			http.Error(w, "Заметка не найдена", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(note)
	}
}

func GetAllNotesHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		notes, err := s.GetAllNotes()
		if err != nil {
			log.Printf("Ошибка при получении заметок: %v", err)
			http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
	}
}

func DeleteNoteHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Path[len("/notes/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		err = s.DeleteNote(id)
		if err != nil {
			log.Printf("Ошибка при удалении заметки %d: %v", id, err)
			http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateNoteHandler(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Path[len("/notes/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		var req UpdateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}

		if req.Title == "" || req.Content == "" {
			http.Error(w, "Заголовок и контент обязательны", http.StatusBadRequest)
			return
		}

		err = s.UpdateNote(id, req.Title, req.Content)
		if err != nil {
			log.Printf("Ошибка при обновлении заметки %d: %v", id, err)
			http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
			return
		}

		// Возвращаем обновлённую заметку
		note := models.Note{
			Id:      id,
			Title:   req.Title,
			Content: req.Content,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(note)
	}
}
