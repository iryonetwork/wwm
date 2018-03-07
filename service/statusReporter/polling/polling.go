package polling

//go:generate sh ../../../bin/mockgen.sh service/statusReporter/polling URLStatusEndpoint $GOFILE

import (
	"context"
	"fmt"
	"time"

	"github.com/iryonetwork/wwm/status"
	"github.com/rs/zerolog"
)

const (
	// how many failed status responsed there needs to be to mark service as failing
	defaultCountThreshold int = 3
	// default interval for status calls
	defaultInterval time.Duration = time.Duration(3 * time.Second)
	// validity of archived status responses used to determine component status
	defaultStatusValidity time.Duration = time.Duration(30 * time.Second)
)

// URLStatusEndpoint is an interface that needs to be fulfilled by URL used by URLPollingComponent
type URLStatusEndpoint interface {
	String() string
	FetchStatus() (*status.Response, error)
}

type statusEntry struct {
	timestamp time.Time
	response  *status.Response
}

// Cfg is a config struct for status reporter service
type Cfg struct {
	CountThreshold *int
	Interval       *time.Duration
	StatusValidity *time.Duration
}

// URLPollingComponent is a struct for URL polling status component
type URLPollingComponent struct {
	url            URLStatusEndpoint
	countThreshold int
	interval       time.Duration
	statusValidity time.Duration
	statusLog      []statusEntry
	logger         zerolog.Logger
}

// Status returns status response for the URL Polling component
func (c *URLPollingComponent) Status() *status.Response {
	if len(c.statusLog) == 0 {
		return nil
	}

	// initialize responsesMap
	responsesMap := map[status.Value][]*status.Response{
		status.OK:      []*status.Response{},
		status.Warning: []*status.Response{},
		status.Error:   []*status.Response{},
	}

	lastStatus := c.statusLog[0]
	doLastStatusConsecutiveRepeatCount := true
	lastStatusConsecutiveRepeatCount := 0

	// cut statuses at validity and count them
	for i, st := range c.statusLog {
		if time.Since(st.timestamp) < c.statusValidity {
			if doLastStatusConsecutiveRepeatCount && lastStatus.response.Status == st.response.Status {
				lastStatusConsecutiveRepeatCount++
			} else {
				doLastStatusConsecutiveRepeatCount = false
			}
			responsesMap[st.response.Status] = append(responsesMap[st.response.Status], st.response)
		} else {
			// cut c.status as the rest is not valid anymore and break the loop
			c.statusLog = c.statusLog[:i+1]
			break
		}
	}

	// if last status was recorded 'countThreshold' in a row then return it
	if lastStatusConsecutiveRepeatCount >= c.countThreshold {
		return lastStatus.response
	}

	// otherwise return most frequent status
	count := len(responsesMap[status.OK])
	var resp *status.Response
	if count > 0 {
		resp = responsesMap[status.OK][0]
	}
	if len(responsesMap[status.Warning]) > count {
		count = len(responsesMap[status.Warning])
		resp = responsesMap[status.Warning][0]
	}
	if len(responsesMap[status.Error]) > count {
		count = len(responsesMap[status.Error])
		resp = responsesMap[status.Error][0]
	}

	return resp
}

// Start starts URL polling
func (c *URLPollingComponent) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(c.interval):
				c.currentStatus()
			}
		}
	}()
}

// NewURLPollingComponent returns new URL Polling status component
func New(url URLStatusEndpoint, cfg *Cfg, logger zerolog.Logger) *URLPollingComponent {
	c := &URLPollingComponent{
		url:            url,
		interval:       defaultInterval,
		countThreshold: defaultCountThreshold,
		statusValidity: defaultStatusValidity,
		logger:         logger,
	}

	if cfg != nil {
		if cfg.Interval != nil {
			c.interval = *cfg.Interval
		}
		if cfg.CountThreshold != nil {
			c.countThreshold = *cfg.CountThreshold
		}
		if cfg.StatusValidity != nil {
			c.statusValidity = *cfg.StatusValidity
		}
	}

	return c
}

func (c *URLPollingComponent) currentStatus() {
	resp, err := c.url.FetchStatus()
	if err != nil {
		c.logger.Error().Err(err).Str("url", c.url.String()).Msg("failed to fetch status")
		resp = &status.Response{Status: status.Error, Msg: fmt.Sprintf("failed to fetch status from %s", c.url)}
	}

	latest := statusEntry{
		timestamp: time.Now(),
		response:  resp,
	}

	c.statusLog = append([]statusEntry{latest}, c.statusLog...)
}
