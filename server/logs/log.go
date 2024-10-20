package logs

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	mu sync.Mutex // Mutex to synchronize access to the log file
)

// Log logs a message to the specified file path.
func Log(message string, filePath string) error {
	mu.Lock() // Lock to ensure that writes to the log file are synchronized
	defer mu.Unlock()

	// Create the directory for the log file if it doesn't exist
	logDir := filepath.Dir(filePath)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return err
	}

	// Open the log file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new logger that writes to the file
	logger := log.New(file, "", log.LstdFlags)

	// Log the message
	logger.Println(message)

	return nil
}
