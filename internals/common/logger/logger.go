package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// levelNames is an array of string representations of the logging levels.
// The order of the strings matches the order of the Level constants.
var levelNames = [...]string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

// Logger is a struct that provides logging functionality with configurable log levels and output destinations.
// It supports writing logs to one or more io.Writer instances, as well as to a file with configurable max size and backup count.
// The Logger struct is thread-safe.
type Logger struct {
	sync.Mutex
	level        Level
	writers      []io.Writer
	fileWriter   *os.File
	fileLogPath  string
	fileMaxSize  int64
	fileBackups  int
	fileLogLevel Level
}

// New creates a new Logger instance with the specified configuration.
//
// The level parameter sets the minimum log level that will be written.
// The fileLogPath parameter specifies the path to the log file.
// The fileMaxSize parameter sets the maximum size of the log file in bytes before it is rotated.
// The fileBackups parameter sets the number of backup log files to keep.
// The fileLogLevel parameter sets the minimum log level that will be written to the log file.
//
// The Logger instance returned will write logs to both the console (os.Stdout) and the specified log file.
func New(level Level, fileLogPath string, fileMaxSize int64, fileBackups int, fileLogLevel Level) (*Logger, error) {
	writers := []io.Writer{os.Stdout}
	var file *os.File

	// Check if the log file already exists, create if it doesn't, or open in append mode if it does
	fileExists, err := fileExists(fileLogPath)
	if err != nil {
		return nil, err
	}

	if fileExists {
		file, err = os.OpenFile(fileLogPath, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.Create(fileLogPath)
		if err != nil {
			return nil, err
		}
	}

	writers = append(writers, file)

	return &Logger{
		level:        level,
		writers:      writers,
		fileWriter:   file,
		fileLogPath:  fileLogPath,
		fileMaxSize:  fileMaxSize,
		fileBackups:  fileBackups,
		fileLogLevel: fileLogLevel,
	}, nil
}

// fileExists checks if the log file exists or not. It returns true if the file exists, false if it does not exist, and an error if there was a problem checking the file's existence.
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Print logs a message with the given level, module, struct name, function name, and arguments.
// The log message is written to both the console and the log file, if configured.
// If the log level is FatalLevel, the program will exit with a status of 1.
// If a log file is configured and the file size exceeds the maximum size, the log file will be rotated.
func (l *Logger) Print(level Level, module, structName, funcName string, args ...interface{}) {
	l.Lock()
	defer l.Unlock()

	if level < l.level && level < l.fileLogLevel {
		return
	}

	ts := time.Now().Format("2006-01-02 15:04:05.000")
	levelName := levelNames[level]
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	} else {
		file = filepath.Base(file)
	}

	logFormat := fmt.Sprintf("%s [%s] %s :: %s :: %s:%d :: %s", ts, levelName, module, structName, file, line, funcName)
	logMessage := fmt.Sprint(args...)

	for _, w := range l.writers {
		fmt.Fprintf(w, "%s - %s\n", logFormat, logMessage)
	}

	if level == FatalLevel {
		os.Exit(1)
	}

	if l.fileWriter != nil && l.fileMaxSize > 0 {
		l.rotateLogFile()
	}
}

// rotateLogFile rotates the log file if the current file size exceeds the maximum allowed size.
// It creates a backup of the current log file, opens a new log file, and updates the writers.
// It also spawns a goroutine to prune old log files.
func (l *Logger) rotateLogFile() {
	fi, err := l.fileWriter.Stat()
	if err != nil {
		log.Println("Failed to get file info:", err)
		return
	}

	if fi.Size() >= l.fileMaxSize {
		l.fileWriter.Close()

		backupPath := fmt.Sprintf("%s.%d", l.fileLogPath, time.Now().UnixNano())
		if err := os.Rename(l.fileLogPath, backupPath); err != nil {
			log.Printf("Failed to rotate log file: %v", err)
			return
		}

		file, err := os.OpenFile(l.fileLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Printf("Failed to open new log file: %v", err)
			return
		}

		l.fileWriter = file
		l.writers = append(l.writers, file)

		go l.pruneLogFiles()
	}
}

// pruneLogFiles removes old log files from the log directory, keeping only the most recent fileBackups number of files.
// It first reads all files in the log directory, finds the log files, sorts them by modification time, and then removes the oldest files
// until the number of log files is less than or equal to the fileBackups setting.
func (l *Logger) pruneLogFiles() {
	dir := filepath.Dir(l.fileLogPath)
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Failed to read directory: %v", err)
		return
	}

	var logFiles []os.FileInfo
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			log.Printf("Failed to get file info: %v", err)
			continue
		}
		if strings.HasPrefix(fileInfo.Name(), filepath.Base(l.fileLogPath)) && strings.HasSuffix(fileInfo.Name(), ".log") {
			logFiles = append(logFiles, fileInfo)
		}
	}

	if len(logFiles) <= l.fileBackups {
		return
	}

	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].ModTime().Before(logFiles[j].ModTime())
	})

	for i := 0; i < len(logFiles)-l.fileBackups; i++ {
		path := filepath.Join(dir, logFiles[i].Name())
		if err := os.Remove(path); err != nil {
			log.Printf("Failed to remove log file: %v", err)
		}
	}
}

// Debug logs a debug-level message with the specified module, struct name, function name, and arguments.
func Debug(module, structName, funcName string, args ...interface{}) {
	defaultLogger.Print(DebugLevel, module, structName, funcName, args...)
}

//  Info logs an info-level message with the specified module, struct name, function name, and arguments.
func Info(module, structName, funcName string, args ...interface{}) {
	defaultLogger.Print(InfoLevel, module, structName, funcName, args...)
}

func Warn(module, structName, funcName string, args ...interface{}) {
	defaultLogger.Print(WarnLevel, module, structName, funcName, args...)
}

//  Error logs an error-level message with the specified module, struct name, function name, and arguments.
func Error(module, structName, funcName string, args ...interface{}) {
	defaultLogger.Print(ErrorLevel, module, structName, funcName, args...)
}

// Fatal logs a fatal-level message with the specified module, struct name, function name, and arguments.
// This function will cause the program to exit after logging the message.
func Fatal(module, structName, funcName string, args ...interface{}) {
	defaultLogger.Print(FatalLevel, module, structName, funcName, args...)
}

var defaultLogger *Logger

// init initializes the default logger for the application. It sets up the log file path, max size, number of backups, and log level.
// The logs directory is created if it doesn't exist. If there are any errors during initialization, the program will exit with a fatal error.
func init() {
	level := InfoLevel
	currentDate := time.Now().Format("2006-01-02")
	fileLogPath := fmt.Sprintf("./logs/app_%s.log", currentDate)
	fileMaxSize := int64(10 * 1024 * 1024) // 10 MB
	fileBackups := 5
	fileLogLevel := DebugLevel

	// Create the logs directory if it doesn't exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	logger, err := New(level, fileLogPath, fileMaxSize, fileBackups, fileLogLevel)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	defaultLogger = logger
}
