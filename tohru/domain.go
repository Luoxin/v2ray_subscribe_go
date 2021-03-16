package tohru

type CheckUsableReq struct {
	// 当前Tohru的版本
	Version string `yaml:"version" json:"version" validate:"required"`
	Hello   string `yaml:"hello" json:"hello" validate:"required"`
}

type CheckUsableRsp struct {
	// 当前Kobayashi-san的版本
	Version string `yaml:"version" json:"version" validate:"required"`
	Token   string `yaml:"token" json:"token" validate:"required"`
}
