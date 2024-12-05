package splash

import (
	"embed"
	"errors"
	"os"
	"path"

	"github.com/onflow/flowkit/v2"
	"github.com/spf13/afero"
)

type (
	EmbedLoader struct {
		embedFS *embed.FS
	}

	FileSystemLoader struct {
		baseDir  string
		fsLoader *afero.Afero
	}
)

var _ flowkit.ReaderWriter = (*EmbedLoader)(nil)
var _ flowkit.ReaderWriter = (*FileSystemLoader)(nil)

func NewEmbedLoader(embedFS *embed.FS) *EmbedLoader {
	return &EmbedLoader{
		embedFS: embedFS,
	}
}

func (el *EmbedLoader) ReadFile(source string) ([]byte, error) {
	return el.embedFS.ReadFile(source)
}

func (el *EmbedLoader) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return errors.New("operation WriteFile not supported by EmbedLoader")
}

func (el *EmbedLoader) MkdirAll(path string, perm os.FileMode) error {
	return errors.New("operation MkdirAll not supported by EmbedLoader")
}

func (el *EmbedLoader) Stat(path string) (os.FileInfo, error) {
	return nil, errors.New("operation Stat not supported by EmbedLoader")
}

func NewFileSystemLoader(baseDir string) *FileSystemLoader {
	return &FileSystemLoader{
		baseDir:  baseDir,
		fsLoader: &afero.Afero{Fs: afero.NewOsFs()},
	}
}

func (f *FileSystemLoader) ReadFile(source string) ([]byte, error) {
	source = path.Join(f.baseDir, source)
	return f.fsLoader.ReadFile(source)
}

func (f *FileSystemLoader) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return errors.New("operation WriteFile not supported by FileSystemLoader")
}

func (f *FileSystemLoader) MkdirAll(path string, perm os.FileMode) error {
	return errors.New("operation MkdirAll not supported by FileSystemLoader")
}

func (f *FileSystemLoader) Stat(path string) (os.FileInfo, error) {
	return nil, errors.New("operation Stat not supported by FileSystemLoader")
}
