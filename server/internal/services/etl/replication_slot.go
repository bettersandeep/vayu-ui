package etl

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/datazip-inc/olake-ui/server/internal/models"
	"github.com/datazip-inc/olake-ui/server/internal/utils"
	"github.com/datazip-inc/olake-ui/server/internal/utils/logger"
)

// pgCDCSlot identifies a postgres logical replication slot by its
// host+port+database+slot_name tuple.
type pgCDCSlot struct {
	Host     string
	Port     int
	Database string
	Slot     string
}

// parsePostgresCDCSlot extracts the replication-slot tuple from a decrypted
// postgres source config. Returns false for non-postgres sources, non-CDC
// update methods, and unparseable configs.
func parsePostgresCDCSlot(sourceType, config string) (pgCDCSlot, bool) {
	if !strings.EqualFold(sourceType, "postgres") {
		return pgCDCSlot{}, false
	}
	var cfg struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		Database     string `json:"database"`
		UpdateMethod struct {
			Type            string `json:"type"`
			ReplicationSlot string `json:"replication_slot"`
		} `json:"update_method"`
	}
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return pgCDCSlot{}, false
	}
	if !strings.EqualFold(cfg.UpdateMethod.Type, "CDC") || cfg.UpdateMethod.ReplicationSlot == "" {
		return pgCDCSlot{}, false
	}
	return pgCDCSlot{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Database: cfg.Database,
		Slot:     cfg.UpdateMethod.ReplicationSlot,
	}, true
}

// sourcesSharingSlot returns the sources (with decrypted configs) whose
// postgres CDC config points at the same slot tuple. excludeSourceID is
// skipped so callers can leave out the entity being deleted.
func sourcesSharingSlot(sources []*models.Source, slot pgCDCSlot, excludeSourceID int) []*models.Source {
	var matches []*models.Source
	for _, src := range sources {
		if src.ID == excludeSourceID {
			continue
		}
		if other, ok := parsePostgresCDCSlot(src.Type, src.Config); ok && other == slot {
			matches = append(matches, src)
		}
	}
	return matches
}

// jobsSharingSlot returns names of jobs other than excludeJobID whose source
// points at the same slot tuple (including other jobs on the same source row).
func (s Service) jobsSharingSlot(slot pgCDCSlot, excludeJobID int) ([]string, error) {
	sources, err := s.db.ListSources()
	if err != nil {
		return nil, fmt.Errorf("failed to list sources for shared-slot check: %s", err)
	}
	sourceIDs := make([]int, 0)
	for _, src := range sourcesSharingSlot(sources, slot, 0) {
		sourceIDs = append(sourceIDs, src.ID)
	}
	jobs, err := s.db.GetJobsBySourceID(sourceIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs for shared-slot check: %s", err)
	}
	var names []string
	for _, job := range jobs {
		if job.ID != excludeJobID {
			names = append(names, job.Name)
		}
	}
	return names, nil
}

// dropReplicationSlot runs the driver's `check --drop-replication-slot` for
// the given source (config must be decrypted) and returns a warning string
// when the drop did not succeed, or "" on success. The drop is idempotent on
// the driver side, so a failed drop can be retried manually.
//
// Version gating: driver versions predating the --drop-replication-slot flag
// fail with an unknown-flag error, which surfaces through this warning path -
// fail-safe by design (the job/source deletion has already succeeded).
func (s Service) dropReplicationSlot(ctx context.Context, source *models.Source) string {
	encryptedConfig, err := utils.Encrypt(source.Config)
	if err != nil {
		return fmt.Sprintf("replication slot was not dropped: failed to encrypt source config: %s", err)
	}
	workflowID := fmt.Sprintf("drop-slot-%s-%d", source.Type, time.Now().Unix())
	result, err := s.temporal.DropReplicationSlot(ctx, workflowID, source.Type, source.Version, encryptedConfig)
	if err != nil {
		return fmt.Sprintf("replication slot was not dropped: %s", err)
	}
	if status, _ := result["status"].(string); status != "succeeded" {
		message, _ := result["message"].(string)
		return fmt.Sprintf("replication slot was not dropped: %s", message)
	}
	logger.Infof("replication slot dropped for source id[%d] name[%s]", source.ID, source.Name)
	return ""
}
