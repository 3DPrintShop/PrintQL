package printdb

type PrintDB interface {
	LoadSpoolInPrinter(id string, id2 string) error
}
