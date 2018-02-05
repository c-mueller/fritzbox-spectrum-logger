package repository

import (
	"github.com/Flaque/filet"
	"path/filepath"
	"testing"
)

func TestInitRepo(t *testing.T) {
	tmpdir := filet.TmpDir(t, "")
	defer filet.CleanUp(t)

	repo, err := NewRepository(filepath.Join(tmpdir, "test_db.db"))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	repo.Close()
}
