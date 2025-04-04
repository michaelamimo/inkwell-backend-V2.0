package config

import (
	"encoding/xml"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"os"
	"sync"
)

var (
	cfg  *APIConfig
	once sync.Once
)

// APIConfig represents the root element.
type APIConfig struct {
	XMLName        xml.Name             `xml:"API"`
	RequestDump    bool                 `xml:"REQUEST_DUMP,attr"`
	Context        ContextConfig        `xml:"CONTEXT"`
	Authentication AuthenticationConfig `xml:"AUTHENTICATION"`
	Pagination     PaginationConfig     `xml:"PAGINATION"`
	DB             DBConfig             `xml:"DB"`
	ThirdParty     ThirdPartyConfig     `xml:"THIRD_PARTY"`
}

// ContextConfig holds basic server settings.
type ContextConfig struct {
	Port            int                  `xml:"PORT"`
	Host            string               `xml:"HOST"`
	Path            string               `xml:"PATH"`
	TimeZone        string               `xml:"TIME_ZONE"`
	EnableBasicAuth bool                 `xml:"ENABLE_BASIC_AUTH"`
	Mode            string               `xml:"MODE"` // "release" or "debug"
	TrustedProxies  TrustedProxiesConfig `xml:"TRUSTED_PROXIES"`
}

// TrustedProxiesConfig holds a list of trusted proxy IP addresses.
type TrustedProxiesConfig struct {
	Proxies []string `xml:"PROXY"`
}

type ThirdPartyConfig struct {
	HFToken    string `xml:"HF_TOKEN"`
	OllamaHost string `xml:"OLLAMA_HOST"`
}

// AuthenticationConfig holds authentication settings.
type AuthenticationConfig struct {
	MultipleSameUserSessions bool `xml:"MULTIPLE_SAME_USER_SESSIONS,attr"`
	EnableTokenAuth          bool `xml:"ENABLE_TOKEN_AUTH"`
	SessionTimeout           int  `xml:"SESSION_TIMEOUT"`
}

// PaginationConfig holds pagination settings.
type PaginationConfig struct {
	PageSize int `xml:"PAGE_SIZE"`
}

// DBConfig holds database connection settings.
type DBConfig struct {
	Initialize bool         `xml:"INITIALIZE"`
	Server     string       `xml:"SERVER"`
	Host       string       `xml:"HOST"`
	Port       int          `xml:"PORT"`
	Driver     string       `xml:"DRIVER"`
	SSLMode    string       `xml:"SSL_MODE"`
	Names      DBNames      `xml:"NAMES"`
	Username   string       `xml:"USERNAME"`
	Password   DBPassword   `xml:"PASSWORD"`
	Pool       DBPoolConfig `xml:"POOL"`
}

// DBNames holds the names defined in the DB section.
type DBNames struct {
	INKWELL string `xml:"INKWELL,attr"`
}

// DBPassword holds password details.
type DBPassword struct {
	Type  string `xml:"TYPE,attr"`
	Value string `xml:",chardata"`
}

// DBPoolConfig holds database connection pooling settings.
type DBPoolConfig struct {
	MaxOpenConns    int `xml:"MAX_OPEN_CONNS"`
	MaxIdleConns    int `xml:"MAX_IDLE_CONNS"`
	ConnMaxLifetime int `xml:"CONN_MAX_LIFETIME"`
}

// LoadConfig loads and parses the XML configuration from the given file.
func LoadConfig(xmlPath string) (*APIConfig, error) {
	once.Do(func() {
		f, err := os.Open(xmlPath)
		if err == nil {
			defer f.Close()

			data, err := io.ReadAll(f)
			if err == nil {
				var newCfg APIConfig
				if err := xml.Unmarshal(data, &newCfg); err == nil {
					cfg = &newCfg
					return
				}
			}
		}

		// If XML file is not found, try loading from .env
		fmt.Println("Config file not found, attempting to load from environment...")

		_ = godotenv.Load() // Load .env file if present
		xmlConfig := os.Getenv("CONFIG_XML")

		if xmlConfig == "" {
			fmt.Println("No XML configuration found in environment variables")
			return
		}

		var newCfg APIConfig
		if err := xml.Unmarshal([]byte(xmlConfig), &newCfg); err == nil {
			cfg = &newCfg
		}
	})

	if cfg == nil {
		return nil, os.ErrInvalid
	}
	return cfg, nil
}

// GetConfig returns the loaded configuration.
func GetConfig() *APIConfig {
	return cfg
}
