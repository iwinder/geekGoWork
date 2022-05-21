package store

var client Factory

//
type Factory interface {
	Users() UserStore
}

func GetClient() Factory {
	return client
}
func SetClient(factory Factory) {
	client = factory
}
