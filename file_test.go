package g3drive_test

import (
	"testing"

	"github.com/lll-lll-lll/g3drive"
	"github.com/stretchr/testify/assert"
)

func Test_Parse_PDF(t *testing.T) {
	t.Parallel()
	fileName := "Drive/aws/aws-cognito-security.pdf"
	f, _ := g3drive.Parse(fileName)
	assert.Equal(t, "aws", f.ParentDirName)
	assert.Equal(t, "aws-cognito-security.pdf", f.Name)
}
