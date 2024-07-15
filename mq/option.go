package mq

type SubOption struct {
	Queue   string
	Consume string
}

type SubOpt func(*SubOption)

func Queue(name string) SubOpt {
	return func(o *SubOption) {
		o.Queue = name
	}
}

func Consume(name string) SubOpt {
	return func(o *SubOption) {
		o.Consume = name
	}
}

type PubOption struct {
	ReplyTo string
	Queue   string
}

type PubOpt func(*PubOption)

func ReplyTo(name string) PubOpt {
	return func(o *PubOption) {
		o.ReplyTo = name
	}
}

func QueueTo(name string) PubOpt {
	return func(o *PubOption) {
		o.Queue = name
	}
}
