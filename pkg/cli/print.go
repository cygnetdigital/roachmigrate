package cli

import (
	"time"

	"github.com/cygnetdigital/roachmigrate/pkg/migrate"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func formatTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return "-"
	}
	return t.Format(time.RFC3339)
}

func printMigrations(migrations []*migrate.Migration) {
	headerFmt := color.New(color.FgGreen).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Filename", "Status", "Updated at", "Error")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt).WithPadding(10)

	for _, m := range migrations {
		if m.Row == nil {
			tbl.AddRow(m.Filename, "not applied", "-", "-")
			continue
		}

		status := "RUNNING"
		if m.Row.Completed {
			status = "COMPLETE"
		}
		if m.Row.Failed {
			status = "FAILED"
		}

		tbl.AddRow(
			m.Filename,
			status,
			formatTime(&m.Row.UpdatedAt),
			m.Row.FailReason.String,
		)
	}

	tbl.Print()
}
