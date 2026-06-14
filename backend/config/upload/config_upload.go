package upload

type Upload struct {
	Image Image `mapstructure:"image"`
	Video Video `mapstructure:"video"`
}
