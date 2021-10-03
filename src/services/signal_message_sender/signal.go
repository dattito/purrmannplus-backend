package signal_message_sender

import (
	"context"

	"github.com/dattito/purrmannplus-backend/config"
	"github.com/dattito/purrmannplus-backend/services/signal_message_sender/proto"
	"google.golang.org/grpc"
)

var SignalMessageSender *_SignalMessageSender

type _SignalMessageSender struct {
	SenderNumber            string
	SIGNAL_CLI_GRPC_API_URL string
	Client                  proto.SignalServiceClient
}

func (sms *_SignalMessageSender) Send(message, recipientPhoneNumber string) error {
	_, err := sms.Client.SendV2(
		context.Background(),
		&proto.SendV2Request{
			Number:     sms.SenderNumber,
			Message:    message,
			Recipients: []string{recipientPhoneNumber},
		},
	)

	return err
}

func newSignalMessageSender(senderNumber, signalCliGrpcApiUrl string) (*_SignalMessageSender, error) {
	conn, err := grpc.Dial(signalCliGrpcApiUrl, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &_SignalMessageSender{
		SenderNumber:            senderNumber,
		SIGNAL_CLI_GRPC_API_URL: signalCliGrpcApiUrl,
		Client:                  proto.NewSignalServiceClient(conn),
	}, nil
}

func Init() error {
	var err error
	SignalMessageSender, err = newSignalMessageSender(config.SIGNAL_SENDER_PHONENUMBER, config.SIGNAL_CLI_GRPC_API_URL)

	return err
}
