package response

import (
	"encoding/json"
	"github.com/EugeneNail/actum/internal/service/log"
	"net/http"
)

type Sender struct {
	writer  http.ResponseWriter
	encoder json.Encoder
}

func NewSender(writer http.ResponseWriter) Sender {
	return Sender{writer, *json.NewEncoder(writer)}
}

func (sender *Sender) Send(data any, status int) {
	switch data.(type) {
	case error:
		err := data.(error)
		http.Error(sender.writer, err.Error(), status)
		log.Error(err)
	default:
		sender.writer.WriteHeader(status)
		if err := sender.encoder.Encode(data); err != nil {
			http.Error(sender.writer, err.Error(), status)
			log.Error(err)
		}
	}
}
