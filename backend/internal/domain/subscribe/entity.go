package subscribe

import (
	"errors"
	"time"
)

type SubscriberStatus string

const (
	SubscriberPending      SubscriberStatus = "pending"
	SubscriberActive       SubscriberStatus = "active"
	SubscriberUnsubscribed SubscriberStatus = "unsubscribed"
	SubscriberBounced      SubscriberStatus = "bounced"
)

var (
	ErrSubscriberNotFound = errors.New("subscriber not found")
	ErrSubscriberExists   = errors.New("subscriber already exists")
	ErrInvalidSubscriber  = errors.New("invalid subscriber")
)

type Subscriber struct {
	ID               string
	Email            string
	Name             string
	Status           SubscriberStatus
	VerifyToken      string
	UnsubscribeToken string
	VerifiedAt       *time.Time
	UnsubscribedAt   *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (s SubscriberStatus) Valid() bool {
	return s == SubscriberPending || s == SubscriberActive || s == SubscriberUnsubscribed || s == SubscriberBounced
}

type EmailDeliveryStatus string

const (
	EmailDeliveryPending EmailDeliveryStatus = "pending"
	EmailDeliverySent    EmailDeliveryStatus = "sent"
	EmailDeliveryFailed  EmailDeliveryStatus = "failed"
)

type EmailDelivery struct {
	ID           string
	SubscriberID *string
	Subject      string
	RefType      string
	RefID        string
	Status       EmailDeliveryStatus
	ErrorMessage string
	SentAt       *time.Time
	CreatedAt    time.Time
}
