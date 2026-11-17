package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const photoSizePx = 16

// CreatePhotoFiles создаёт во временной папке теста count PNG-файлов и возвращает пути к ним; папка удаляется после теста
func CreatePhotoFiles(t *testing.T, count int) []string {
	t.Helper()
	dir := t.TempDir()

	var buffer bytes.Buffer
	require.NoError(t, png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, photoSizePx, photoSizePx))),
		"Не удалось сгенерировать PNG-изображение")

	paths := make([]string, 0, count)
	for i := range count {
		photoPath := filepath.Join(dir, fmt.Sprintf("photo_%d.png", i+1))
		require.NoErrorf(t, os.WriteFile(photoPath, buffer.Bytes(), 0o600), "Не удалось сохранить фотографию %s", photoPath)
		paths = append(paths, photoPath)
	}

	return paths
}
