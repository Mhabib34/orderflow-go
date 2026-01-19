package usecase

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"payment_service/internal/broker"
	"payment_service/internal/dto"
	"payment_service/internal/exception"
	"payment_service/internal/model"
	"payment_service/internal/repository"
	"payment_service/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PaymentUsecaseImpl struct {
	PaymentRepository repository.PaymentRepository
	Validate       *validator.Validate
	MidtransService service.MidtransService
	Publisher      broker.Publisher
}

func NewPaymentUsecase(paymentRepository repository.PaymentRepository, validate *validator.Validate, midtransService service.MidtransService, publisher broker.Publisher) PaymentUsecase {
	return &PaymentUsecaseImpl{
		PaymentRepository: paymentRepository, 
		Validate: validate, 
		MidtransService: midtransService,
		Publisher: publisher,
	}
}

func generateSignature(
	orderID, statusCode, grossAmount, serverKey string,
) string {
	raw := orderID + statusCode + grossAmount + serverKey
	hash := sha512.Sum512([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func mapMidtransStatus(midtransStatus string) string {
	switch midtransStatus {
	case "capture", "settlement":
		return "SUCCESS"
	case "pending":
		return "PENDING"
	case "expire":
		return "EXPIRED"
	default:
		return "FAILED"
	}
}

func (p *PaymentUsecaseImpl) CreatePayment(
	ctx context.Context,
	request dto.CreatePaymentRequest,
) (dto.PaymentResponse, error) {

	// 1. validate
	if err := p.Validate.Struct(request); err != nil {
		return dto.PaymentResponse{}, err
	}

	// 2. generate payment id
	paymentID := uuid.NewString()

	// 3. create payment to Midtrans
	snapResp, err := p.MidtransService.CreateSnapPayment(
		paymentID,
		request.Amount,
	)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	// 4. simpan ke DB
	payment := &model.Payments{
		OrderID:       request.OrderID,
		Amount:        request.Amount,
		Method:        request.Method,
		ProviderRefID: paymentID,
		PaymentURL:    snapResp.RedirectURL,
		Status:        "PENDING",
		PaymentID:     paymentID,
	}

	_, err = p.PaymentRepository.Create(ctx, payment)
	if err != nil {
		return dto.PaymentResponse{}, err
	}
		

	// 5. response
	return dto.PaymentResponse{
		PaymentID:  paymentID,
		PaymentURL: snapResp.RedirectURL,
		Status:     "PENDING",
	}, nil
}

func (p *PaymentUsecaseImpl) HandleMidtransCallback(
	ctx context.Context,
	payload dto.MidtransCallback,
) error {

	log.Printf("📩 Midtrans Callback Received\n")
	log.Printf("   OrderID: %s\n", payload.OrderID)
	log.Printf("   Status: %s\n", payload.TransactionStatus)
	log.Printf("   StatusCode: %s\n", payload.StatusCode)
	log.Printf("   GrossAmount: %s\n", payload.GrossAmount)

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		log.Println("❌ MIDTRANS_SERVER_KEY is empty!")
		return fmt.Errorf("server key not configured")
	}

	expectedSignature := generateSignature(
		payload.OrderID,
		payload.StatusCode,
		payload.GrossAmount,
		serverKey,
	)

	log.Printf("🔐 Expected Signature: %s\n", expectedSignature)
	log.Printf("🔐 Received Signature: %s\n", payload.SignatureKey)

	if payload.SignatureKey != expectedSignature {
		log.Printf("❌ Signature mismatch!\n")
		return fmt.Errorf("invalid midtrans signature")
	}

	log.Println("✅ Signature valid")

	status := mapMidtransStatus(payload.TransactionStatus)
	log.Printf("🔄 Mapping '%s' to '%s'\n", payload.TransactionStatus, status)

	log.Printf("💾 Updating payment_id: %s to status: %s\n", payload.OrderID, status)
	
	err := p.PaymentRepository.UpdateStatusByPaymentID(ctx, payload.OrderID, status)
	if err != nil {
		log.Printf("❌ Failed to update: %v\n", err)
		return err
	}

	//publish to rabbitmq
	event := dto.PaymentStatusChangedEvent{
		PaymentStatus: status,
		PaymentID:     uuid.MustParse(payload.OrderID),
		OrderID:       uuid.MustParse(payload.OrderID),
		PaymentMethod: payload.PaymentType,
	}
	body, err := json.Marshal(event)
	exception.PanicIfError(err)
	
	err = p.Publisher.Publish(ctx, "payment.status.updated", body)
	exception.PanicIfError(err) 

	log.Printf("✅ Payment %s updated successfully to %s\n", payload.OrderID, status)
	return nil
}
