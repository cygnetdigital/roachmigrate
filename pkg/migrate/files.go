package migrate

import (
	"os"
	"path/filepath"
	"sort"
)

// LoadFiles from a given directory
func LoadFiles(dir string) ([]*File, error) {
	var out []*File

	// Get all files in the directory.
	// `ioutil.ReadDir` ensures that the files are sorted by name
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return out, nil
		}
		return nil, err
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		// skip non sql files
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		// read contents of file
		contents, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		out = append(out, &File{
			Name:     file.Name(),
			Contents: contents,
		})
	}

	return out, nil
}

// sortFiles without mutating the input
func sortFiles(files []*File) []*File {
	sorted := make([]*File, len(files))
	copy(sorted, files)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})
	return sorted
}
