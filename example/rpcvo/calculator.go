package rpcvo

import "errors"

type CalculatorService struct {}

type CalculatorArgs struct {
	A, B int
}

type CalculatorReply struct {
	Result int
}


func (this *CalculatorService) Div(args *CalculatorArgs, reply *CalculatorReply) error {
	if args.B == 0{
		return errors.New("B is zero.")
	}
	reply.Result = args.A * args.B
	return nil
}

func (h *CalculatorService) Multiply(args *CalculatorArgs, reply *CalculatorReply) error {
	reply.Result = args.A * args.B
	return nil
}
