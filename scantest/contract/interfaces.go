package contract

type Handler interface {
	Handle(*Context)
}
