package g3drive_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/lll-lll-lll/g3drive/g3drive"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v3"
)

func Test_SearchDirInList(t *testing.T) {
	t.Parallel()
	wantDirName := "test_dir"
	g3df := g3drive.G3DFile{ParentDirName: wantDirName}
	list := &drive.FileList{Files: []*drive.File{{MimeType: g3drive.DirMimeType.String(), Name: wantDirName}}}
	_, ok := g3drive.SearchDirIDInList(&g3df, list)
	if !ok {
		t.Log("File not Exist in Drive")
	}
	assert.Equal(t, ok, true)
}

// attention: actual google drive api used
func Test_Search(t *testing.T) {
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./credentials.json")
	ctx := context.Background()
	srv, err := drive.NewService(ctx)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	client := g3drive.New(srv)
	fs, err := client.Search(ctx, []string{"1C6orkHG_KwAOGDy9BQJ6uJSHjZdUDujU"})
	if err != nil {
		t.Error(err)
	}
	for _, f := range fs.Files {
		if f.Name == "test.doc" {
			t.Log("f.Name == test.doc")
		}
		fmt.Println(f.Id, f.Name)
	}
}
