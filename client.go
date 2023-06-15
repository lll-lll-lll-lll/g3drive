package g3drive

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/drive/v3"
)

const QPath = " or '%s' in parents"

type ClientDrive struct {
	service *drive.Service
}

func New(srv *drive.Service) *ClientDrive {
	return &ClientDrive{service: srv}
}

func (cd *ClientDrive) UploadFile(ctx context.Context, g3dFile *G3DFile) (*drive.File, error) {
	targetFile, err := os.Open(g3dFile.path)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer targetFile.Close()
	res, err := cd.service.Files.
		Create(&drive.File{Name: g3dFile.Name, Parents: []string{g3dFile.ParentDirID}}).
		Media(targetFile).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return res, nil
}

func (cd *ClientDrive) CreateDir(ctx context.Context, g3df *G3DFile) (*drive.File, error) {
	res, err := cd.service.Files.
		Create(&drive.File{Name: g3df.ParentDirName, Parents: []string{g3df.ParentDirID}, MimeType: DirMimeType.String()}).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return res, nil
}

// if targetDirIDs is number of zero, default is 'DriveFolderID'
func (cd *ClientDrive) Search(ctx context.Context, targetDirIDs []string) (*drive.FileList, error) {
	DriveFolderID := os.Getenv("DRIVE_FOLDER_ID")
	var targetDirIDsForQ = make([]string, len(targetDirIDs))
	if len(targetDirIDs) == 0 {
		targetDirIDs = []string{DriveFolderID}
	}
	for i, tdid := range targetDirIDs {
		if i == 0 {
			targetDirIDsForQ = append(targetDirIDsForQ, fmt.Sprintf("'%s' in parents", tdid))
			continue
		}
		targetDirIDsForQ = append(targetDirIDsForQ, fmt.Sprintf(QPath, tdid))
	}
	files, err := cd.service.Files.List().PageSize(1000).
		Fields("files(id, name, mimeType, parents)").
		Q(strings.Join(targetDirIDsForQ, " ")).
		Context(ctx).
		Do()
	if err != nil {
		return nil, err
	}
	return files, nil
}

func SetParentID(g3df *G3DFile, parentID string) *G3DFile {
	g3df.ParentDirID = parentID
	return g3df
}

func SearchDirInList(g3df *G3DFile, fileOrDir *drive.FileList) (string, bool) {
	var existed bool = true
	for _, driveDir := range fileOrDir.Files {
		if driveDir.MimeType == DirMimeType.String() && driveDir.Name == g3df.ParentDirName {
			return driveDir.Id, existed
		}
	}
	return "", !existed
}

func SearchFileIDInList(g3df *G3DFile, fileOrDir *drive.FileList) (string, bool) {
	var existed bool = true
	for _, file := range fileOrDir.Files {
		if file.MimeType == g3df.MimeType.String() && file.Name == g3df.Name {
			return file.Id, existed
		}
	}
	return "", !existed
}

func Upload(ctx context.Context, client *ClientDrive, g3f *G3DFile) error {
	files, err := client.Search(ctx, nil)
	if err != nil {
		return err
	}
	dirID, ok := SearchDirInList(g3f, files)
	if !ok {
		dirRes, err := client.CreateDir(ctx, g3f)
		if err != nil {
			return err
		}
		dirID = dirRes.Id
	}
	g3f = SetParentID(g3f, dirID)
	filesDir, err := client.Search(ctx, []string{dirID})
	if err != nil {
		return err
	}
	// ok is true. file exist
	if _, ok = SearchFileIDInList(g3f, filesDir); ok {
		log.Println("same file name existed")
		return err
	}

	if _, err := client.UploadFile(ctx, g3f); err != nil {
		return err
	}
	return nil
}
