package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger yeni logrus loggeri yaradır və konfiqurasiya edir
func NewLogger() *logrus.Logger {
	logger := logrus.New()

	// Log formatını təyin et
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Çıxış mənbəyini təyin et (stderr)
	logger.SetOutput(os.Stderr)

	// Log səviyyəsini təyin et
	logger.SetLevel(logrus.InfoLevel)

	return logger
}
