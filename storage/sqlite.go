package storage

import (
	"database/sql"
	"notes-api/models"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (storage *Storage) CreateNote(title, content string) (int64, error) {
	result, err := storage.db.Exec("INSERT INTO notes (title,content) VALUES (?, ?)", title, content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (storage *Storage) GetAllNotes() ([]models.Note, error) {
	rows, err := storage.db.Query("SELECT id, title, content FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.Id, &n.Title, &n.Content); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func (s *Storage) DeleteNote(id int) error {
	_, err := s.db.Exec("DELETE FROM notes WHERE id = ?", id)
	return err
}

func (s *Storage) UpdateNote(id int, title, content string) error {
	_, err := s.db.Exec("UPDATE notes SET title = ?, content = ? WHERE id = ?", title, content, id)
	return err
}
