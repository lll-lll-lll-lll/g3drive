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
		return nil, fmt.Errorf("failed to open file. %w", err)
	}
	defer targetFile.Close()
	res, err := cd.service.Files.
		Create(&drive.File{Name: g3dFile.Name, Parents: []string{g3dFile.ParentDirID}}).
		Media(targetFile).
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create file. %w", err)
	}
	return res, nil
}

func (cd *ClientDrive) MakeDir(ctx context.Context, g3df *G3DFile) (*drive.File, error) {
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
func (cd *ClientDrive) Find(ctx context.Context, driveForlderID string, targetDirIDs []string) (*drive.FileList, error) {
	var targetDirIDsForQ = make([]string, len(targetDirIDs))
	if len(targetDirIDs) == 0 {
		targetDirIDs = []string{driveForlderID}
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

func setParentDirID(g3df *G3DFile, parentID string) *G3DFile {
	g3df.ParentDirID = parentID
	return g3df
}

func findDirID(g3df *G3DFile, fileOrDir *drive.FileList) (string, bool) {
	var existed bool = true
	for _, driveDir := range fileOrDir.Files {
		if driveDir.MimeType == DirMimeType.String() && driveDir.Name == g3df.ParentDirName {
			return driveDir.Id, existed
		}
	}
	return "", !existed
}

// findFileID if file existed in fileOrDir, return fileID, true
func findFileID(g3df *G3DFile, fileOrDir *drive.FileList) (string, bool) {
	for _, file := range fileOrDir.Files {
		if file.MimeType == g3df.MimeType.String() && file.Name == g3df.Name {
			return file.Id, true
		}
	}
	return "", false
}

func Upload(ctx context.Context, client *ClientDrive, g3f *G3DFile, driveFolderID string) error {
	files, err := client.Find(ctx, driveFolderID, nil)
	if err != nil {
		return err
	}
	var dirID string
	var ok bool
	dirID, ok = findDirID(g3f, files)
	if !ok {
		dirRes, err := client.MakeDir(ctx, g3f)
		if err != nil {
			return err
		}
		dirID = dirRes.Id
	}
	g3f = setParentDirID(g3f, dirID)
	filesDir, err := client.Find(ctx, driveFolderID, []string{dirID})
	if err != nil {
		return err
	}
	// if ok is true, file already existed
	if _, ok = findFileID(g3f, filesDir); ok {
		log.Println("same file name existed")
		return err
	}

	if _, err := client.UploadFile(ctx, g3f); err != nil {
		return err
	}
	return nil
}
