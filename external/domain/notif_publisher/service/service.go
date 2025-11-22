package service

import (
	"errors"
	"io"
	"net/http"

	"github.com/fikriahmadf/outbox-examples/configs"
	"github.com/go-resty/resty/v2"
)

type ExternalNotifPublisherService interface {
	N8NService
}

type ExternalNotifPublisherServiceImpl struct {
	client    *resty.Client
	cfg       *configs.Config
	endpoints Endpoints
}

type Endpoints struct {
	SendMemoNotifPath string
}

func ProvideNotifPublisherService(cfg *configs.Config) *ExternalNotifPublisherServiceImpl {
	r := resty.New()
	r.SetRetryCount(cfg.External.N8N.RetryCount)
	r.SetRetryWaitTime(cfg.External.N8N.RetryWaitTime)
	r.SetBaseURL(cfg.External.N8N.BaseURL)
	r.AddRetryCondition(func(r *resty.Response, err error) bool {
		return errors.Is(err, io.EOF) ||
			r.StatusCode() == http.StatusServiceUnavailable ||
			r.StatusCode() == http.StatusBadGateway
	})

	return &ExternalNotifPublisherServiceImpl{
		client:    r,
		cfg:       cfg,
		endpoints: Endpoints(cfg.External.N8N.Endpoints),
	}
}
