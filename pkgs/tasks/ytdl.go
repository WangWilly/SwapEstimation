package tasks

import (
	"context"
	"os/exec"
	"time"

	"github.com/WangWilly/swap-estimation/pkgs/cmd"
	"github.com/WangWilly/swap-estimation/pkgs/utils"
	"github.com/WangWilly/swap-estimation/pkgs/uuid"
)

////////////////////////////////////////////////////////////////////////////////

type DownloadTask struct {
	taskID    string
	targetUrl string
	filePath  string
	progress  int64

	retries      int
	retryDelay   time.Duration
	maxRetries   int
	retryChannel chan struct{}

	ctx        context.Context
	cancel     context.CancelFunc
	maxTimeout time.Duration
}

////////////////////////////////////////////////////////////////////////////////

func NewRetribleTaskWithCtx(
	ctx context.Context,
	uuidGen uuid.UUID,
	url string,
	filepath string,
	retryDelay time.Duration,
	maxRetries int,
) *DownloadTask {
	ctx, cancel := context.WithCancel(ctx)
	task := &DownloadTask{
		taskID:    uuidGen.New(),
		targetUrl: url,
		filePath:  filepath,
		progress:  0,

		retries:      0,
		retryDelay:   retryDelay,
		maxRetries:   maxRetries,
		retryChannel: make(chan struct{}, 1),

		ctx:        ctx,
		cancel:     cancel,
		maxTimeout: 0,
	}
	return task
}

func NewRetribleNamedTaskWithCtx(
	ctx context.Context,
	taskID string,
	url string,
	filepath string,
	retryDelay time.Duration,
	maxRetries int,
) *DownloadTask {
	ctx, cancel := context.WithCancel(ctx)
	task := &DownloadTask{
		taskID:    taskID,
		targetUrl: url,
		filePath:  filepath,
		progress:  0,

		retries:      0,
		retryDelay:   retryDelay,
		maxRetries:   maxRetries,
		retryChannel: make(chan struct{}, 1),

		ctx:        ctx,
		cancel:     cancel,
		maxTimeout: 0,
	}
	return task
}

////////////////////////////////////////////////////////////////////////////////

func (t *DownloadTask) WithMaxTimeout(timeout time.Duration) *DownloadTask {
	t.maxTimeout = timeout
	return t
}

////////////////////////////////////////////////////////////////////////////////

func (t *DownloadTask) GetID() string {
	return t.taskID
}

func (t *DownloadTask) GetProgress() int64 {
	return t.progress
}

func (t *DownloadTask) GetFilePath() string {
	return t.filePath
}

func (t *DownloadTask) GetTargetUrl() string {
	return t.targetUrl
}

func (t *DownloadTask) GetRetries() int {
	return t.retries
}

func (t *DownloadTask) GetMaxRetries() int {
	return t.maxRetries
}

func (t *DownloadTask) GetRetryDelay() time.Duration {
	return t.retryDelay
}

func (t *DownloadTask) GetMaxTimeout() time.Duration {
	return t.maxTimeout
}

////////////////////////////////////////////////////////////////////////////////

func (t *DownloadTask) execute(c cmd.Cmd) bool {
	logger := utils.GetDetailedLogger().With().Caller().Logger()

	// Execute
	if err := c.Run(); err != nil {
		t.progress = -1
		if t.ctx.Err() == context.Canceled {
			logger.Error().Err(err).Msg("Download canceled")
		} else {
			if t.retries < t.maxRetries {
				t.progress = -2
			}
			logger.Error().Err(err).Msg("Download failed")
		}
		return false
	}
	t.progress = 100

	// Cleanup
	logger.Info().Msgf("Download complete: %s", t.filePath)
	return true
}

func (t *DownloadTask) Execute() bool {
	// Setup
	t.progress = 30
	ctx := t.ctx
	if t.maxTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(t.ctx, t.maxTimeout)
		defer cancel()
	}

	// Command
	c := exec.CommandContext(
		ctx,
		"yt-dlp",
		"-o", t.filePath,
		"-f", "mp4",
		t.targetUrl,
	)

	// Execute
	return t.execute(c)
}

func (t *DownloadTask) SetRetrySignal() <-chan struct{} {
	logger := utils.GetDetailedLogger().With().Caller().Logger()

	go func() {
		if t.retries >= t.maxRetries {
			logger.Error().Msgf("Max retries reached for: %s", t.filePath)
			return
		}

		time.Sleep(t.retryDelay)
		t.retries++
		logger.Warn().Msgf("Retrying download: %s, attempt: %d", t.filePath, t.retries)
		t.retryChannel <- struct{}{}
	}()

	if t.retries >= t.maxRetries {
		logger.Error().Msgf("Max retries reached for: %s", t.filePath)
		return nil
	}
	return t.retryChannel
}

func (t *DownloadTask) Cancel() {
	logger := utils.GetDetailedLogger().With().Caller().Logger()
	logger.Error().Msgf("Canceling download: %s, ", t.filePath)
	logger.Error().Msgf("Canceling task: %s\n", t.taskID)
	t.cancel()
}
