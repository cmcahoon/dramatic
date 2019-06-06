package actor

import "go.uber.org/zap"

// Zap is, allegedly, very performant. The SugaredLogger isn't even the fastest, so if a little extra juice needs to be
// found and the system is logging a lot, switch to the regular logger.
var logger *zap.SugaredLogger

func init() {
	plainLogger, _ := zap.NewProduction()
	logger = plainLogger.Sugar()
}
