package crawler

type Crawler interface {
	Download() error
	Parse() error
}

