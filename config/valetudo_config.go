package config

type ValetudoConfig struct {
	Embedded bool `json:"embedded"`
	Robot    struct {
		Implementation               string `json:"implementation"`
		ImplementationSpecificConfig struct {
			IP string `json:"ip"`
		} `json:"implementationSpecificConfig"`
	} `json:"robot"`
	Webserver struct {
		Port      int `json:"port"`
		BasicAuth struct {
			Enabled  bool   `json:"enabled"`
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"basicAuth"`
		BlockExternalAccess bool `json:"blockExternalAccess"`
	} `json:"webserver"`
	ZonePresets struct {
	} `json:"zonePresets"`
	GoToLocationPresets struct {
	} `json:"goToLocationPresets"`
	Mqtt struct {
		Enabled    bool `json:"enabled"`
		Connection struct {
			Host string `json:"host"`
			Port int    `json:"port"`
			TLS  struct {
				Enabled                 bool   `json:"enabled"`
				Ca                      string `json:"ca"`
				IgnoreCertificateErrors bool   `json:"ignoreCertificateErrors"`
			} `json:"tls"`
			Authentication struct {
				Credentials struct {
					Enabled  bool   `json:"enabled"`
					Username string `json:"username"`
					Password string `json:"password"`
				} `json:"credentials"`
				ClientCertificate struct {
					Enabled     bool   `json:"enabled"`
					Certificate string `json:"certificate"`
					Key         string `json:"key"`
				} `json:"clientCertificate"`
			} `json:"authentication"`
		} `json:"connection"`
		Identity struct {
			FriendlyName string `json:"friendlyName"`
			Identifier   string `json:"identifier"`
		} `json:"identity"`
		Interfaces struct {
			Homie struct {
				Enabled                   bool `json:"enabled"`
				AddICBINVMapProperty      bool `json:"addICBINVMapProperty"`
				CleanAttributesOnShutdown bool `json:"cleanAttributesOnShutdown"`
			} `json:"homie"`
			Homeassistant struct {
				Enabled                 bool `json:"enabled"`
				CleanAutoconfOnShutdown bool `json:"cleanAutoconfOnShutdown"`
			} `json:"homeassistant"`
		} `json:"interfaces"`
		Customizations struct {
			TopicPrefix    string `json:"topicPrefix"`
			ProvideMapData bool   `json:"provideMapData"`
		} `json:"customizations"`
	} `json:"mqtt"`
	NtpClient struct {
		Enabled  bool   `json:"enabled"`
		Server   string `json:"server"`
		Port     int    `json:"port"`
		Interval int    `json:"interval"`
		Timeout  int    `json:"timeout"`
	} `json:"ntpClient"`
	Timers struct {
	} `json:"timers"`
	LogLevel string `json:"logLevel"`
	Debug    struct {
		SystemStatInterval    bool `json:"systemStatInterval"`
		DebugHassAnchors      bool `json:"debugHassAnchors"`
		StoreRawUploadedMaps  bool `json:"storeRawUploadedMaps"`
		EnableDebugCapability bool `json:"enableDebugCapability"`
	} `json:"debug"`
	NetworkAdvertisement struct {
		Enabled bool `json:"enabled"`
	} `json:"networkAdvertisement"`
	Updater struct {
		Enabled        bool `json:"enabled"`
		UpdateProvider struct {
			Type                         string `json:"type"`
			ImplementationSpecificConfig struct {
			} `json:"implementationSpecificConfig"`
		} `json:"updateProvider"`
	} `json:"updater"`
}
