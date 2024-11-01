package config

type Config struct {
	Mqtt struct {
		Connection struct {
			Host string `json:"host"`
			Port int    `json:"port"`
			TLS  struct {
				Enabled bool   `json:"enabled"`
				Ca      string `json:"ca"`
			} `json:"tls"`
			Authentication struct {
				Credentials struct {
					Enabled  bool   `json:"enabled"`
					Username string `json:"username"`
					Password string `json:"password"`
				} `json:"credentials"`
			} `json:"authentication"`
		} `json:"connection"`
	} `json:"mqtt"`
	Devices []Device `json:"devices"`
}
type Device struct {
	Name       string `json:"name"`
	MacAddress string `json:"macAddress"`
}
