package migrate

import (
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
)

// Validate a migration file
func Validate(file *File) error {
	if _, err := parser.Parse(string(file.Contents)); err != nil {
		return err
	}

	return nil
}
