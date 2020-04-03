package webdav

import (
	"os"
	"testing"

	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/rnd"
	"github.com/stretchr/testify/assert"
)

const (
	testUrl  = "http://webdav-dummy/"
	testUser = "admin"
	testPass = "photoprism"
)

func TestConnect(t *testing.T) {
	c := New(testUrl, testUser, testPass)

	assert.IsType(t, Client{}, c)
}

func TestClient_Files(t *testing.T) {
	c := New(testUrl, testUser, testPass)

	assert.IsType(t, Client{}, c)

	files, err := c.Files("Photos")

	if err != nil {
		t.Fatal(err)
	}

	if len(files) == 0 {
		t.Fatal("no files found")
	}
}

func TestClient_Directories(t *testing.T) {
	c := New(testUrl, testUser, testPass)

	assert.IsType(t, Client{}, c)

	t.Run("non-recursive", func(t *testing.T) {
		dirs, err := c.Directories("", false)

		if err != nil {
			t.Fatal(err)
		}

		if len(dirs) == 0 {
			t.Fatal("no directories found")
		}

		assert.IsType(t, fs.FileInfo{}, dirs[0])
		assert.Equal(t, "Photos", dirs[0].Name)
		assert.Equal(t, "/Photos", dirs[0].Abs)
		assert.Equal(t, true, dirs[0].Dir)
		assert.Equal(t, int64(0), dirs[0].Size)
	})

	t.Run("recursive", func(t *testing.T) {
		dirs, err := c.Directories("", true)

		if err != nil {
			t.Fatal(err)
		}

		if len(dirs) < 2 {
			t.Fatal("at least 2 directories expected")
		}
	})
}

func TestClient_Download(t *testing.T) {
	c := New(testUrl, testUser, testPass)

	assert.IsType(t, Client{}, c)

	files, err := c.Files("Photos")

	if err != nil {
		t.Fatal(err)
	}

	tempDir := os.TempDir() + rnd.UUID()
	tempFile := tempDir + "/foo.jpg"

	if len(files) == 0 {
		t.Fatal("no files to download")
	}

	if err := c.Download(files[0].Abs, tempFile); err != nil {
		t.Fatal(err)
	}

	if !fs.FileExists(tempFile) {
		t.Fatalf("%s does not exist", tempFile)
	}

	if err := os.RemoveAll(tempDir); err != nil {
		t.Fatal(err)
	}
}

func TestClient_DownloadDir(t *testing.T) {
	c := New(testUrl, testUser, testPass)

	assert.IsType(t, Client{}, c)

	t.Run("non-recursive", func(t *testing.T) {
		tempDir := os.TempDir() + rnd.UUID()

		if errs := c.DownloadDir("Photos", tempDir, false); len(errs) > 0 {
			t.Fatal(errs)
		}

		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("recursive", func(t *testing.T) {
		tempDir := os.TempDir() + rnd.UUID()

		if errs := c.DownloadDir("Photos", tempDir, true); len(errs) > 0 {
			t.Fatal(errs)
		}

		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatal(err)
		}
	})
}

func TestClient_UploadAndDelete(t *testing.T) {
	c := New(testUrl, testUser, testPass)

	assert.IsType(t, Client{}, c)

	tempName := rnd.UUID() + ".jpg"

	if err := c.Upload("testdata/example.jpg", tempName); err != nil {
		t.Fatal(err)
	}

	if err := c.Delete(tempName); err != nil {
		t.Fatal(err)
	}
}
