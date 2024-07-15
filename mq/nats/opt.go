package nats

func WithSubjects(subjects ...string) SubOpt {
	return func(cfg *SubOption) {
		cfg.Subjects = subjects
	}
}
