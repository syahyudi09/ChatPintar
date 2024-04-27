package repository

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/syahyudi09/ChatPintar/App/model"
)

type FileRepository interface {
	SaveFileInfo(file *model.FileModel) error
	SaveFile(file io.Reader, fileName string) (string, error)
}

type fileRepository struct {
	db         *sql.DB
	storageDir string
}

func NewFileRepository(db *sql.DB) FileRepository {
	return &fileRepository{
		db: db,
	}
}

func (fr *fileRepository) SaveFileInfo(file *model.FileModel) error {
	query := "INSERT INTO files (file_name, file_path, uploader, created_at) VALUES ($1, $2, $3, $4)"
	_, err := fr.db.Exec(query, file.FileName, file.FilePath, file.Uploader, file.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save file info: %v", err)
	}
	return nil
}

func (fr *fileRepository) SaveFile(file io.Reader, fileName string) (string, error) {
	filePath := filepath.Join(fr.storageDir, fileName)

	// Buat file di sistem file
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	// Salin konten file
	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return filePath, nil
}
