package log

// Option ...
type Option struct {
	Address  string
	AppName  string
	FuncName string
}

// OptionFunc ...
type OptionFunc func(opt *Option)

// WithAddress with address
func WithAddress(addr string) func(opt *Option) {
	return func(opt *Option) {
		opt.Address = addr
	}
}

// WithFuncName with func name
func WithFuncName(name string) func(opt *Option) {
	return func(opt *Option) {
		opt.FuncName = name
	}
}

// WithAppName with app name
func WithAppName(name string) func(opt *Option) {
	return func(opt *Option) {
		opt.AppName = name
	}
}
