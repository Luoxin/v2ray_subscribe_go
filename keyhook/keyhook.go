package keyhook

type Init struct {
}

func (p *Init) Init() (needRun bool, err error) {
	return true, InitKeyHook()
}

func (p *Init) WaitFinish() {

}

func (p *Init) Name() string {
	return "key hook"
}
