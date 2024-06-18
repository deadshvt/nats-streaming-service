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

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = Charset[rand.Intn(len(Charset))]
	}
	return string(b)
}

func randomIntInRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func randomPhone() string {
	return fmt.Sprintf("+%d%d%d%d%d%d%d%d%d%d%d",
		randomIntInRange(1, 9), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10), rand.Intn(10),
		rand.Intn(10), rand.Intn(10))
}

func randomLocale() string {
	locales := []string{"en", "es", "fr", "de", "ru", "zh", "ja"}
	return locales[rand.Intn(len(locales))]
}

func randomCurrency() string {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "RUB", "CNY"}
	return currencies[rand.Intn(len(currencies))]
}

func RandomOrder() *entity.Order {
	order := entity.Order{
		OrderUid:          uuid.New().String(),
		TrackNumber:       randomString(12),
		Entry:             randomString(4),
		Delivery:          *randomDelivery(),
		Locale:            randomLocale(),
		InternalSignature: randomString(10),
		CustomerId:        randomString(10),
		DeliveryService:   randomString(10),
		Shardkey:          randomString(4),
		SmId:              rand.Intn(100),
		DateCreated:       time.Now(),
		OofShard:          randomString(4),
	}

	order.Payment = *randomPayment(order.OrderUid)
	order.Items = *randomItems(randomIntInRange(1, 6), order.TrackNumber)

	return &order
}

func randomDelivery() *entity.Delivery {
	return &entity.Delivery{
		Name:    randomString(10),
		Phone:   randomPhone(),
		Zip:     randomString(6),
		City:    randomString(8),
		Address: randomString(15),
		Region:  randomString(8),
		Email:   fmt.Sprintf("%s@example.com", randomString(10)),
	}
}

func randomPayment(transaction string) *entity.Payment {
	return &entity.Payment{
		Transaction:  transaction,
		RequestId:    randomString(10),
		Currency:     randomCurrency(),
		Provider:     randomString(10),
		Amount:       rand.Intn(10000),
		PaymentDt:    int(time.Now().Unix()),
		Bank:         randomString(10),
		DeliveryCost: rand.Intn(500),
		GoodsTotal:   rand.Intn(1000),
		CustomFee:    rand.Intn(100),
	}
}

func randomItem(trackNumber string) *entity.Item {
	return &entity.Item{
		ChrtId:      rand.Intn(1000000),
		TrackNumber: trackNumber,
		Price:       rand.Intn(1000),
		Rid:         uuid.New().String(),
		Name:        randomString(10),
		Sale:        rand.Intn(50),
		Size:        randomString(3),
		TotalPrice:  rand.Intn(1000),
		NmId:        rand.Intn(100000),
		Brand:       randomString(10),
		Status:      rand.Intn(1000),
	}
}

func randomItems(count int, trackNumber string) *[]entity.Item {
	items := make([]entity.Item, count)
	for i := range items {
		items[i] = *randomItem(trackNumber)
	}
	return &items
}
