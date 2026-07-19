package etl

import (
	"testing"

	"github.com/datazip-inc/olake-ui/server/internal/models"
)

const cdcConfig = `{"host":"db1","port":5432,"database":"app","update_method":{"type":"CDC","replication_slot":"olake_slot","publication":"olake_pub"}}`

func TestParsePostgresCDCSlot(t *testing.T) {
	tests := []struct {
		name       string
		sourceType string
		config     string
		want       pgCDCSlot
		wantOK     bool
	}{
		{
			name:       "postgres CDC",
			sourceType: "postgres",
			config:     cdcConfig,
			want:       pgCDCSlot{Host: "db1", Port: 5432, Database: "app", Slot: "olake_slot"},
			wantOK:     true,
		},
		{
			name:       "non-postgres source",
			sourceType: "mysql",
			config:     cdcConfig,
			wantOK:     false,
		},
		{
			name:       "standalone update method",
			sourceType: "postgres",
			config:     `{"host":"db1","port":5432,"database":"app","update_method":{"type":"Standalone"}}`,
			wantOK:     false,
		},
		{
			name:       "CDC without slot name",
			sourceType: "postgres",
			config:     `{"host":"db1","port":5432,"database":"app","update_method":{"type":"CDC"}}`,
			wantOK:     false,
		},
		{
			name:       "invalid JSON",
			sourceType: "postgres",
			config:     `not-json`,
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parsePostgresCDCSlot(tt.sourceType, tt.config)
			if ok != tt.wantOK {
				t.Fatalf("ok = %t, want %t", ok, tt.wantOK)
			}
			if ok && got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSourcesSharingSlot(t *testing.T) {
	slot := pgCDCSlot{Host: "db1", Port: 5432, Database: "app", Slot: "olake_slot"}
	sources := []*models.Source{
		{ID: 1, Type: "postgres", Config: cdcConfig}, // the source being deleted
		{ID: 2, Type: "postgres", Config: cdcConfig}, // same tuple -> match
		{ID: 3, Type: "postgres", Config: `{"host":"db2","port":5432,"database":"app","update_method":{"type":"CDC","replication_slot":"olake_slot"}}`}, // different host
		{ID: 4, Type: "postgres", Config: `{"host":"db1","port":5432,"database":"app","update_method":{"type":"CDC","replication_slot":"other_slot"}}`}, // different slot
		{ID: 5, Type: "mysql", Config: cdcConfig}, // different source type
		{ID: 6, Type: "postgres", Config: `{"host":"db1","port":5432,"database":"app","update_method":{"type":"Standalone"}}`}, // not CDC
	}

	matches := sourcesSharingSlot(sources, slot, 1)
	if len(matches) != 1 || matches[0].ID != 2 {
		t.Fatalf("expected only source 2 to match, got %+v", matches)
	}

	// Without exclusion the deleted source's own row matches too.
	matches = sourcesSharingSlot(sources, slot, 0)
	if len(matches) != 2 {
		t.Fatalf("expected sources 1 and 2 to match, got %d matches", len(matches))
	}
}
