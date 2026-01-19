package service

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransServiceImpl struct {
	client snap.Client
}

func NewMidtransService() MidtransService {
	var c snap.Client

	// ⚠️ JANGAN pakai string kosong
	c.New(midtrans.ServerKey, midtrans.Environment)

	return &MidtransServiceImpl{
		client: c,
	}
}

func (s *MidtransServiceImpl) CreateSnapPayment(
	paymentID string,
	amount int64,
) (*SnapResponse, error) {

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  paymentID,
			GrossAmt: amount,
		},
	}

	resp, err := s.client.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	return &SnapResponse{
		Token:       resp.Token,
		RedirectURL: resp.RedirectURL,
	}, nil
}
