package geerpc

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}
