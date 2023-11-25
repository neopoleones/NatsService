package subscriber

type Subscriber interface {
	Subscribe() (Entry, error)
}

type Entry interface {
	PullMessage() (*Message, error)
}

type Message struct {
	Subject string

	Header map[string][]string
	Data   []byte
}
