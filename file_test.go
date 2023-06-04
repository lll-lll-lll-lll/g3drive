package g3drive_test

import (
	"testing"

	"github.com/lll-lll-lll/g3drive"
	"github.com/stretchr/testify/assert"
)

func Test_Parse_PDF(t *testing.T) {
	t.Parallel()
	fileName := "a/b/file1.pdf"
	f, _ := g3drive.Parse(fileName)
	assert.Equal(t, "b", f.ParentDirName)
	assert.Equal(t, "file1.pdf", f.Name)
}
