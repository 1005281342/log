package log

// Option ...
type Option struct {
	Address  string
	AppName  string
	FuncName string
}

// OptionFunc ...
type OptionFunc func(opt *Option)

func WithAddress(addr string) func(opt *Option) {
	return func(opt *Option) {
		opt.Address = addr
	}
}

func WithFuncName(name string) func(opt *Option) {
	return func(opt *Option) {
		opt.FuncName = name
	}
}

func WithAppName(name string) func(opt *Option) {
	return func(opt *Option) {
		opt.AppName = name
	}
}
