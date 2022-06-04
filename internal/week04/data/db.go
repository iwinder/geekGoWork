package data

var client Factory

type Factory interface {
	UserRepo()
}

func GetClient() Factory {
	return client
}
func SetClient(factory Factory) {
	client = factory
}
