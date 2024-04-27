package usecase

import (
	"fmt"
	"io"
	"time"

	"github.com/syahyudi09/ChatPintar/App/model"
	"github.com/syahyudi09/ChatPintar/App/repository"
)

type FileUsecase interface {
	UploadFile(file io.Reader, fileName, uploader string) error
}

type fileUsecase struct {
	repository repository.FileRepository
}

func NewFileUsecase(repository repository.FileRepository) FileUsecase {
	return &fileUsecase{
		repository: repository,
	}
}

func (uc *fileUsecase) UploadFile(file io.Reader, fileName, uploader string) error {
	if file == nil || fileName == "" || uploader == "" {
		return fmt.Errorf("invalid input")
	}

	// Simpan file secara local
	filePath, err := uc.repository.SaveFile(file, fileName)
	if err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}

	fileModel := &model.FileModel{
		FileName:  fileName,
		FilePath:  filePath,
		Uploader:  uploader,
		CreatedAt: time.Now(), // Waktu pembuatan
	}

	// Simpan file di database
	err = uc.repository.SaveFileInfo(fileModel)
	if err != nil {
		return fmt.Errorf("error saving file info: %v", err)
	}

	return nil
}
