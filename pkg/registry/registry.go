package registry

type Registry interface {
	Storer
	Searcher
	Pager
}