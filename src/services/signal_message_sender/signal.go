package signal_message_sender

import (
	"context"

	"github.com/dattito/purrmannplus-backend/config"
	"github.com/dattito/purrmannplus-backend/services/signal_message_sender/proto"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"google.golang.org/grpc"
)

var SignalMessageSender *_SignalMessageSender

type _SignalMessageSender struct {
	SenderNumber            string
	SIGNAL_CLI_GRPC_API_URL string
	Client                  proto.SignalServiceClient
}

// Sends a message to a given phone number
func (sms *_SignalMessageSender) Send(message, recipientPhoneNumber string) error {
	_, err := sms.Client.SendV2(
		context.Background(),
		&proto.SendV2Request{
			Number:     sms.SenderNumber,
			Message:    message,
			Recipients: []string{recipientPhoneNumber},
		},
	)

	if err != nil {
		logging.Errorf("Error sending signal message. Error: %s | Message tried to send: %s", err.Error(), message)
	} else {
		logging.Debugf("Signal message sent successfully. Message: %s | Recipient: %s", message, recipientPhoneNumber)
	}

	return err
}

// Creates a new signal message sender object from the given parameters
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

// Initializes the signal message sender
func Init() error {
	var err error
	SignalMessageSender, err = newSignalMessageSender(config.SIGNAL_SENDER_PHONENUMBER, config.SIGNAL_CLI_GRPC_API_URL)

	return err
}
