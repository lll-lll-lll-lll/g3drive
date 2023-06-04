package g3drive_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/lll-lll-lll/g3drive"
	"github.com/stretchr/testify/assert"
)

func Test_ReadDir(t *testing.T) {
	t.Parallel()
	file, err := os.Open("a/file1.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}

func Test_Parse_PDF(t *testing.T) {
	t.Parallel()
	fileName := "a/b/file1.pdf"
	f, _ := g3drive.Parse(fileName)
	assert.Equal(t, "b", f.ParentDirName)
	assert.Equal(t, "file1.pdf", f.Name)
}
