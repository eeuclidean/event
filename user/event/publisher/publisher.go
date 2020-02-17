package publisher





type EventPublisher interface {
	Publish() error
}
        