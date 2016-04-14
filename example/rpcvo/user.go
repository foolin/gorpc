package rpcvo

type UserService struct {}

func (this *UserService) Say(name *string, reply *string) error {
	*reply = "Hello, " + *name + "!"
	return nil
}
