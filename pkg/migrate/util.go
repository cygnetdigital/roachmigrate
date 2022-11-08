package migrate

import (
	"sort"

	"github.com/cygnetdigital/roachmigrate/pkg/store"
)

// zip files with migrations
func zip(files []*File, migrations []store.RoachMigration) (out []*Migration) {

	migrationSeen := make(map[string]struct{})

	for _, file := range files {
		status := &Migration{
			File:     file,
			Filename: file.Name,
		}

		for _, migration := range migrations {
			if migration.Filename == file.Name {
				migrationSeen[migration.Filename] = struct{}{}
				status.Row = &migration
				break
			}
		}

		out = append(out, status)
	}

	// now for any migrations that exist without files
	for _, migration := range migrations {
		if _, ok := migrationSeen[migration.Filename]; !ok {
			out = append(out, &Migration{
				Row:      &migration,
				Filename: migration.Filename,
			})
		}
	}

	// sort by filename
	sort.Slice(out, func(i, j int) bool {
		return out[i].Filename < out[j].Filename
	})

	return out
}
