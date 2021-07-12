package initialize

type Initialize interface {
	Init() (needRun bool, err error)
	WaitFinish()
	Name() string
}
