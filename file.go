package g3drive

import (
	"fmt"
	"os"
	"path/filepath"
)

type CustomMimeType int

// ファイルのタイプを定義
const (
	DirMimeType CustomMimeType = iota + 1
	DocMimeType
	TextPlainMimeType
	PDFMimeType
)

func (cmt CustomMimeType) String() string {
	switch cmt {
	case DirMimeType:
		return "application/vnd.google-apps.folder"
	case DocMimeType:
		return "application/vnd.google-apps.document"
	case TextPlainMimeType:
		return "text/plain"
	case PDFMimeType:
		return "application/pdf"
	}
	return ""
}

type G3DFile struct {
	Ext           string
	Name          string
	path          string
	ParentDirName string
	ParentDirID   string
	MimeType      CustomMimeType
}

func Parse(filePath string) (*G3DFile, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	_, fileName := filepath.Split(filePath)
	_, parentDirName := filepath.Split(filepath.Dir(filePath))
	ext := filepath.Ext(filePath)
	mimeType, err := parseMimeType(ext)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &G3DFile{Ext: ext, Name: fileName, path: filePath, ParentDirName: parentDirName, MimeType: mimeType}, nil
}

func parseMimeType(ext string) (CustomMimeType, error) {
	switch ext {
	case ".pdf":
		return PDFMimeType, nil
	case ".txt":
		return TextPlainMimeType, nil
	case ".doc":
		return DocMimeType, nil
	}
	return 0, nil
}
