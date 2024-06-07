package sentinel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	ingestTimeout = time.Second * 30
)

func (s *Sentinel) IngestLog(ctx context.Context, endpoint, ruleID, streamName string, logs []map[string]string) error {
	logger := s.logger.WithField("module", "sentinel_ingest")

	ingest, err := azlogs.NewClient(endpoint, s.azCreds, nil)
	if err != nil {
		return fmt.Errorf("could not create azure ingest client: %v", err)
	}

	logPayload, err := json.Marshal(&logs)
	if err != nil {
		return fmt.Errorf("could not json encode log message: %v", err)
	}

	if s.logger.IsLevelEnabled(logrus.TraceLevel) {
		logger.Tracef("%s", string(logPayload))
	}

	logger.WithField("total", len(logs)).Debug("uploading logs")

	ctx, _ = context.WithTimeout(ctx, ingestTimeout)

	_, err = ingest.Upload(ctx, ruleID, streamName, logPayload, nil)
	if err != nil {
		return fmt.Errorf("could not upload logs: %v", err)
	}

	logger.WithField("total_logs", len(logs)).Debug("successfully uploaded gong logs")

	return nil
}
