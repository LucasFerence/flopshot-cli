package edit

const ProviderType = "provider"

type Provider struct {
	Id   string `json:"_id" mapstructure:"_id"`
	Name string `json:"name"`
}

func (provider Provider) Label() string {
	return provider.Name
}

func init() {
	RegisterType(ProviderType, Provider{})
}
