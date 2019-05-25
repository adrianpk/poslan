package mailer

// NewServiceForTests returns a configured sertvice
// Mainly used for tests.
// func NewServiceForTests(cfg *config.Config) Service {
// 	// Context
// 	ctx, cancel := context.WithCancel(context.Background())
// 	go checkSigTerm(cancel)

// 	// Logger
// 	logger := makeLogger()

// 	// Config
// 	cfg, err := config.Load()
// 	checkError(err)

// 	// Service
// 	svc := makeService(ctx, cfg, logger)
// 	return svc
// }
