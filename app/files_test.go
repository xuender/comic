// nolint: testpackage
package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiles_Start(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	files := &Files{index: 0, ends: []int{3, 6, 9}}

	for index, start := range []int{0, 0, 0, 0, 4, 4, 4, 7, 7, 7} {
		files.index = index
		files.Start()
		ass.Equal(start, files.index)
	}
}

func TestFiles_End(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	files := &Files{index: 0, ends: []int{3, 6, 9}}

	for index, end := range []int{3, 3, 3, 3, 6, 6, 6, 9, 9, 9} {
		files.index = index
		files.End()
		ass.Equal(end, files.index)
	}
}
