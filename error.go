package config

type cfgErr struct {
	msg   string
	inner error
}

func (e cfgErr) Error() string {
	return e.msg
}

func (e cfgErr) Unwrap() error {
	return e.inner
}

func wrap(err error, msg string) error {
	return cfgErr{inner: err, msg: msg}
}
