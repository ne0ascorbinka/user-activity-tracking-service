package migrations

import (
	"io/fs"
	"testing"
)

func TestMigrationsFS(t *testing.T) {
	expectedFiles := []string{
		"000001_create_events_table.up.sql",
		"000001_create_events_table.down.sql",
		"000002_create_user_activity_stats_table.up.sql",
		"000002_create_user_activity_stats_table.down.sql",
	}

	for _, filename := range expectedFiles {
		content, err := fs.ReadFile(FS, filename)
		if err != nil {
			t.Errorf("failed to read embedded migration %s: %v", filename, err)
		}
		if len(content) == 0 {
			t.Errorf("embedded migration %s is empty", filename)
		}
	}
}
