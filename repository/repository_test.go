package repository

import (
    "github.com/Flaque/filet"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/c-mueller/fritzbox-spectrum-logger/fritz"
    "os"
    "io/ioutil"
    "encoding/json"
    "time"
)

func TestInitRepo(t *testing.T) {
    tmpdir := filet.TmpDir(t, "")
    defer filet.CleanUp(t)

    repo, err := NewRepository(filepath.Join(tmpdir, "test_db.db"))
    assert.NoErrorf(t, err, "Initialization Failed")
    repo.Close()
}

func TestRepository_Insert(t *testing.T) {
    tmpdir := filet.TmpDir(t, "")
    defer filet.CleanUp(t)

    repo, err := NewRepository(filepath.Join(tmpdir, "test_db.db"))
    assert.NoErrorf(t, err, "Initialization Failed")

    spectrum := loadTestSpectrum(t)

    for i := 0; i < 1000; i++ {
        spectrum.Timestamp = time.Now().Unix() + int64(100*i)
        err = repo.Insert(spectrum)
        assert.NoErrorf(t, err, "Inserting element %d failed", i)
    }

    repo.Close()
}

func BenchmarkRepository_Insert(b *testing.B) {
    dir, err := ioutil.TempDir("", "")
    assert.NoError(b, err, "Creating tempdir failed")

    repo, err := NewRepository(filepath.Join(dir, "test_db.db"))
    assert.NoError(b, err, "Opening Repo Failed")

    spectrum := loadTestSpectrum(b)

    b.ResetTimer()

    count := 0
    start := time.Now()

    for i := 0; i < b.N; i++ {
        err = repo.Insert(spectrum)
        assert.NoError(b, err, "Insertion failed")
        spectrum.Timestamp = time.Now().Unix() + int64(100*i)
        count++
    }

    b.Log("Performed", count, "Operations During the benchmark")
    b.Log("The Benchmark ran", time.Since(start))

    repo.Close()
    os.Remove(dir)
}

func loadTestSpectrum(t testing.TB) *fritz.Spectrum {
    file, err := os.Open("testdata/example_spectrum.json")
    assert.NoError(t, err, "Loading Dummy Spectrum failed")
    var result *fritz.Spectrum
    data, err := ioutil.ReadAll(file)
    file.Close()
    err = json.Unmarshal(data, &result)
    return result
}
