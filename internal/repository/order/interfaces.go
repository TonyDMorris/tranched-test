package order

type idGenerator interface {
	New() string
}
