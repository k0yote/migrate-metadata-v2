package k0yote3web

type SDKOptions struct {
	ThirdpartyProvier ThirdpartyProvider
	PrivateKey        string

	ApiKey string
}

type DownloadMetaOptions struct {
	BaseURL      string
	StartTokenID int
	EndTokenID   int
}

type DownloadCh struct {
	Endpoint string
	Data     []byte
}

type Attribute struct {
	DisplayType string `json:"display_type,omitempty"`
	TraitType   string `json:"trait_type"`
	Value       any    `json:"value"`
}

type MetaData struct {
	Name        string      `json:"name"`
	Image       string      `json:"image"`
	Description string      `json:"description,omitempty"`
	Attributes  []Attribute `json:"attributes"`
}

type ThirdpartyProvider string

const (
	INFURA  ThirdpartyProvider = "infura"
	ALCHEMY ThirdpartyProvider = "alchemy"
)

type GasPriority float64

const (
	Low     GasPriority = 2.0
	Average GasPriority = 2.1
	High    GasPriority = 2.2
)

func (g GasPriority) Value() float64 {
	return float64(g)
}
