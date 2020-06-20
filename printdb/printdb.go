package printdb

// PrintDB represents an interface that something must implement to support the PrintDB standard
type PrintDB interface {
	LoadSpoolInPrinter(spoolID string, printerID string) error
}
