package order

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/deadshvt/nats-streaming-service/internal/entity"

	"github.com/google/uuid"
)

const (
	Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = Charset[rand.Intn(len(Charset))]
	}
	return string(b)
}

func generateRandomIntInRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func generateRandomPhone() string {
	return fmt.Sprintf("+%d%d%d%d%d%d%d%d%d%d%d",
		generateRandomIntInRange(1, 9), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10))
}

func generateRandomLocale() string {
	locales := []string{"en", "es", "fr", "de", "ru", "zh", "ja"}
	return locales[rand.Intn(len(locales))]
}

func generateRandomCurrency() string {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "RUB", "CNY"}
	return currencies[rand.Intn(len(currencies))]
}

func GenerateOrder() *entity.Order {
	order := entity.Order{
		OrderUid:          uuid.New().String(),
		TrackNumber:       generateRandomString(12),
		Entry:             generateRandomString(4),
		Delivery:          *generateDelivery(),
		Locale:            generateRandomLocale(),
		InternalSignature: generateRandomString(10),
		CustomerId:        generateRandomString(10),
		DeliveryService:   generateRandomString(10),
		Shardkey:          generateRandomString(4),
		SmId:              rand.Intn(100),
		DateCreated:       time.Now(),
		OofShard:          generateRandomString(4),
	}

	order.Payment = *generatePayment(order.OrderUid)
	order.Items = *generateItems(generateRandomIntInRange(1, 6), order.TrackNumber)

	return &order
}

func generateDelivery() *entity.Delivery {
	return &entity.Delivery{
		Name:    generateRandomString(10),
		Phone:   generateRandomPhone(),
		Zip:     generateRandomString(6),
		City:    generateRandomString(8),
		Address: generateRandomString(15),
		Region:  generateRandomString(8),
		Email:   fmt.Sprintf("%s@example.com", generateRandomString(10)),
	}
}

func generatePayment(transaction string) *entity.Payment {
	return &entity.Payment{
		Transaction:  transaction,
		RequestId:    generateRandomString(10),
		Currency:     generateRandomCurrency(),
		Provider:     generateRandomString(10),
		Amount:       rand.Intn(10000),
		PaymentDt:    int(time.Now().Unix()),
		Bank:         generateRandomString(10),
		DeliveryCost: rand.Intn(500),
		GoodsTotal:   rand.Intn(1000),
		CustomFee:    rand.Intn(100),
	}
}

func generateItem(trackNumber string) *entity.Item {
	return &entity.Item{
		ChrtId:      rand.Intn(1000000),
		TrackNumber: trackNumber,
		Price:       rand.Intn(1000),
		Rid:         uuid.New().String(),
		Name:        generateRandomString(10),
		Sale:        rand.Intn(50),
		Size:        generateRandomString(3),
		TotalPrice:  rand.Intn(1000),
		NmId:        rand.Intn(100000),
		Brand:       generateRandomString(10),
		Status:      rand.Intn(1000),
	}
}

func generateItems(count int, trackNumber string) *[]entity.Item {
	items := make([]entity.Item, count)
	for i := range items {
		items[i] = *generateItem(trackNumber)
	}
	return &items
}
