package tohru

type UserInfo struct {
	Hello         string `yaml:"hello" json:"hello" validate:"required"`
	TohruKey      string `yaml:"tohru_key" json:"tohru_key" validate:"required"`
	TohruPassword string `yaml:"tohru_password" json:"tohru_password" validate:"required"`
}

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

type SyncNodeReq struct {
	NodeList string `yaml:"node_list" json:"node_list" validate:"required"`
}

type SyncNodeRsp struct {
	NodeList string `yaml:"node_list" json:"node_list"`
}

type RegistrationReq struct {
	TohruKey      string `yaml:"tohru_key" json:"tohru_key" validate:"required"`
	TohruPassword string `yaml:"tohru_password" json:"tohru_password" validate:"required"`
}

type RegistrationRsp struct {
}

type ChangePasswordReq struct {
	// 当前Tohru的版本
	TohruKey         string `yaml:"tohru_key" json:"tohru_key" validate:"required"`
	OldTohruPassword string `yaml:"old_tohru_password" json:"old_tohru_password" validate:"required"`
	NewTohruPassword string `yaml:"old_tohru_password" json:"vtohru_password" validate:"required"`
}

type ChangePasswordRsp struct {
}
