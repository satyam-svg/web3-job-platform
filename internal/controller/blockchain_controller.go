package controller

import (
	"encoding/json"
	"net/http"

	"github.com/satyam-svg/resume-parser/internal/service"
)

func VerifyPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TxHash string `json:"txHash"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	isValid := service.VerifyPayment(body.TxHash)
	if isValid {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Payment verified"}`))
	} else {
		http.Error(w, "Invalid payment", http.StatusBadRequest)
	}
}
