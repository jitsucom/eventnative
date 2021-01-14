package storages

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/jitsucom/eventnative/adapters"
	"github.com/jitsucom/eventnative/caching"
	"github.com/jitsucom/eventnative/events"
	"github.com/jitsucom/eventnative/logging"
	"github.com/jitsucom/eventnative/schema"
)

type GoogleAnalytics struct {
	name            string
	gaAdapter       *adapters.GoogleAnalytics
	tableHelper     *TableHelper
	processor       *schema.Processor
	streamingWorker *StreamingWorker
	fallbackLogger  *logging.AsyncLogger
	eventsCache     *caching.EventsCache
	staged          bool
}

func NewGoogleAnalytics(config *Config) (Storage, error) {
	if !config.streamMode {
		return nil, fmt.Errorf("Google Analytics destination doesn't support %s mode", BatchMode)
	}

	gaConfig := config.destination.GoogleAnalytics
	if err := gaConfig.Validate(); err != nil {
		return nil, err
	}

	requestDebugLogger := config.loggerFactory.CreateSQLQueryLogger(config.name)
	gaAdapter := adapters.NewGoogleAnalytics(gaConfig, requestDebugLogger)

	tableHelper := NewTableHelper(gaAdapter, config.monitorKeeper, config.pkFields, adapters.SchemaToGoogleAnalytics)

	ga := &GoogleAnalytics{
		name:           config.name,
		gaAdapter:      gaAdapter,
		tableHelper:    tableHelper,
		processor:      config.processor,
		fallbackLogger: config.loggerFactory.CreateFailedLogger(config.name),
		eventsCache:    config.eventsCache,
		staged:         config.destination.Staged,
	}

	ga.streamingWorker = newStreamingWorker(config.eventQueue, config.processor, ga, config.eventsCache, config.loggerFactory.CreateStreamingArchiveLogger(config.name), tableHelper)
	ga.streamingWorker.start()

	return ga, nil
}

func (ga *GoogleAnalytics) DryRun(payload events.Event) ([]adapters.TableField, error) {
	return dryRun(payload, ga.processor, ga.tableHelper)
}

func (ga *GoogleAnalytics) Insert(table *adapters.Table, event events.Event) (err error) {
	return ga.gaAdapter.Send(event)
}

func (ga *GoogleAnalytics) Store(fileName string, payload []byte, alreadyUploadedTables map[string]bool) (map[string]*StoreResult, int, error) {
	return nil, 0, errors.New("GoogleAnalytics doesn't support Store() func")
}

func (ga *GoogleAnalytics) StoreWithParseFunc(fileName string, payload []byte, skipTables map[string]bool, parseFunc func([]byte) (map[string]interface{}, error)) (map[string]*StoreResult, int, error) {
	return nil, 0, errors.New("GoogleAnalytics doesn't support StoreWithParseFunc() func")
}

func (ga *GoogleAnalytics) SyncStore(collectionTable string, objects []map[string]interface{}, timeIntervalValue string) (int, error) {
	return 0, errors.New("GoogleAnalytics doesn't support SyncStore() func")
}

func (ga *GoogleAnalytics) GetUsersRecognition() *UserRecognitionConfiguration {
	return disabledRecognitionConfiguration
}

//Fallback log event with error to fallback logger
func (ga *GoogleAnalytics) Fallback(failedEvents ...*events.FailedEvent) {
	for _, failedEvent := range failedEvents {
		ga.fallbackLogger.ConsumeAny(failedEvent)
	}
}

func (ga *GoogleAnalytics) Name() string {
	return ga.name
}

func (ga *GoogleAnalytics) Type() string {
	return GoogleAnalyticsType
}

func (ga *GoogleAnalytics) IsStaging() bool {
	return ga.staged
}

func (ga *GoogleAnalytics) Close() (multiErr error) {
	if err := ga.gaAdapter.Close(); err != nil {
		multiErr = multierror.Append(multiErr, fmt.Errorf("[%s] Error closing GoogleAnalytics client: %v", ga.Name(), err))
	}

	if ga.streamingWorker != nil {
		if err := ga.streamingWorker.Close(); err != nil {
			multiErr = multierror.Append(multiErr, fmt.Errorf("[%s] Error closing streaming worker: %v", ga.Name(), err))
		}
	}

	if err := ga.fallbackLogger.Close(); err != nil {
		multiErr = multierror.Append(multiErr, fmt.Errorf("[%s] Error closing fallback logger: %v", ga.Name(), err))
	}

	return
}
