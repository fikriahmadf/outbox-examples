package service

import (
	"context"
	"encoding/json"

	"github.com/fikriahmadf/outbox-examples/external/domain/notif_publisher/model"
	"github.com/rs/zerolog/log"
)

func (s *ExternalNotifPublisherServiceImpl) SendMemoNotif(ctx context.Context, req model.SendMemoNotifRequest) (res model.SendMemoNotifResponse, err error) {

	resp, err := s.client.R().SetBody(req).Post(s.endpoints.SendMemoNotifPath)
	if err != nil {
		log.Warn().Err(err).Msg("[N8NService][SendMemoNotif] failed to send memo notif")
		return model.SendMemoNotifResponse{}, err
	}

	if resp.IsError() {
		log.Warn().Err(err).Msg("[N8NService][SendMemoNotif] failed to send memo notif")
		return res, err
	}

	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("[N8NService][Unmarshal] failed to unmarshal response")
		return model.SendMemoNotifResponse{}, err
	}

	return res, nil
}

type N8NService interface {
	SendMemoNotif(ctx context.Context, req model.SendMemoNotifRequest) (res model.SendMemoNotifResponse, err error)
}
