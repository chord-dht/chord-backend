package config

import (
	"chord-backend/aes"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const Unspecified = "Unspecified"

type Config struct {
	IdentifierLength int
	SuccessorsLength int

	IpAddress            string
	Port                 string
	JoinAddress          string
	JoinPort             string
	StabilizeTime        int
	FixFingersTime       int
	CheckPredecessorTime int

	StorageDir string
	BackupDir  string

	Mode string // "create" or "join"

	AESBool    bool   // turn on/off AES encryption when storing files
	AESKeyPath string // path to the AES key file
	AESKey     []byte // AES key

	TLSBool         bool // turn on/off TLS connection
	CaCert          string
	ServerCert      string
	ServerKey       string
	ServerTLSConfig *tls.Config
	ClientTLSConfig *tls.Config
}

var NodeConfig *Config

func NewConfig() *Config {
	return &Config{
		IdentifierLength:     -1,
		SuccessorsLength:     -1,
		IpAddress:            Unspecified,
		Port:                 Unspecified,
		JoinAddress:          Unspecified,
		JoinPort:             Unspecified,
		StabilizeTime:        -1,
		FixFingersTime:       -1,
		CheckPredecessorTime: -1,
		StorageDir:           "storage",
		BackupDir:            "backup",
		AESBool:              false,
		AESKeyPath:           "",
		TLSBool:              false,
		CaCert:               "",
		ServerCert:           "",
		ServerKey:            "",
	}
}

func JsontToConfig(json map[string]interface{}) *Config {
	cfg := NewConfig()

	if val, ok := json["IdentifierLength"].(float64); ok {
		cfg.IdentifierLength = int(val)
	}
	if val, ok := json["SuccessorsLength"].(float64); ok {
		cfg.SuccessorsLength = int(val)
	}
	if val, ok := json["IpAddress"].(string); ok {
		cfg.IpAddress = val
	}
	if val, ok := json["Port"].(float64); ok {
		cfg.Port = fmt.Sprintf("%v", val)
	}
	if val, ok := json["JoinAddress"].(string); ok {
		cfg.JoinAddress = val
	}
	if val, ok := json["JoinPort"].(float64); ok {
		cfg.JoinPort = fmt.Sprintf("%v", val)
	}
	if val, ok := json["StabilizeTime"].(float64); ok {
		cfg.StabilizeTime = int(val)
	}
	if val, ok := json["FixFingersTime"].(float64); ok {
		cfg.FixFingersTime = int(val)
	}
	if val, ok := json["CheckPredecessorTime"].(float64); ok {
		cfg.CheckPredecessorTime = int(val)
	}
	if val, ok := json["AESBool"].(bool); ok {
		cfg.AESBool = val
	}
	if val, ok := json["AESKeyPath"].(string); ok {
		cfg.AESKeyPath = val
	}
	if val, ok := json["TLSBool"].(bool); ok {
		cfg.TLSBool = val
	}
	if val, ok := json["CaCert"].(string); ok {
		cfg.CaCert = val
	}
	if val, ok := json["ServerCert"].(string); ok {
		cfg.ServerCert = val
	}
	if val, ok := json["ServerKey"].(string); ok {
		cfg.ServerKey = val
	}
	return cfg
}

func ValidateAndSetConfig(cfg *Config) error {
	if err := validateConfig(cfg); err != nil {
		fmt.Println("Failed to validate config:", err)
		return err
	}

	determineMode(cfg)

	if err := determineAES(cfg); err != nil {
		fmt.Println("Failed to determine AES:", err)
		return err
	}

	if err := determineTLS(cfg); err != nil {
		fmt.Println("Failed to determine TLS:", err)
		return err
	}
	return nil
}

// CheckPortAvailability
func CheckPortAvailability(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// CheckRemoteAddressAvailability
func CheckRemoteAddressAvailability(address string, port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", address, port), 5*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Validate the configuration, returning an error if the configuration is invalid.
func validateConfig(cfg *Config) error {
	if cfg.IdentifierLength < 1 || cfg.IdentifierLength > 160 {
		return fmt.Errorf("identifier length must be in the range of [1,160]")
	}

	if cfg.IpAddress == Unspecified {
		return fmt.Errorf("ip address must be specified")
	}

	if net.ParseIP(cfg.IpAddress) == nil {
		return fmt.Errorf("invalid ip address format")
	}

	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port <= 1024 || port > 65535 {
		return fmt.Errorf("port must be in the range of (1024,65535]")
	}
	if !CheckPortAvailability(port) {
		return fmt.Errorf("port %d is not available", port)
	}

	if cfg.StabilizeTime < 1 || cfg.StabilizeTime > 60000 {
		return fmt.Errorf("stabilize time must be in the range of [1,60000] milliseconds")
	}

	if cfg.FixFingersTime < 1 || cfg.FixFingersTime > 60000 {
		return fmt.Errorf("fix fingers time must be in the range of [1,60000] milliseconds")
	}

	if cfg.CheckPredecessorTime < 1 || cfg.CheckPredecessorTime > 60000 {
		return fmt.Errorf("check predecessor time must be in the range of [1,60000] milliseconds")
	}

	if cfg.SuccessorsLength < 1 || cfg.SuccessorsLength > 32 {
		return fmt.Errorf("number of successors must be in the range of [1,32]")
	}

	if (cfg.JoinAddress != Unspecified && cfg.JoinPort == Unspecified) || (cfg.JoinAddress == Unspecified && cfg.JoinPort != Unspecified) {
		return fmt.Errorf("both --ja and --jp must be specified together")
	}

	if cfg.JoinAddress != Unspecified && cfg.JoinPort != Unspecified {
		if net.ParseIP(cfg.JoinAddress) == nil {
			return fmt.Errorf("invalid join address format")
		}

		joinPort, err := strconv.Atoi(cfg.JoinPort)
		if err != nil || joinPort <= 1024 || joinPort > 65535 {
			return fmt.Errorf("join port must be in the range of (1024,65535]")
		}
		if !CheckRemoteAddressAvailability(cfg.JoinAddress, joinPort) {
			return fmt.Errorf("join address %s:%d is not reachable", cfg.JoinAddress, joinPort)
		}
	}

	if cfg.AESBool {
		if cfg.AESKeyPath == "" {
			return fmt.Errorf("AES key path must be specified if --aes is specified")
		}
		if _, err := os.Stat(cfg.AESKeyPath); os.IsNotExist(err) {
			return fmt.Errorf("AES key file does not exist at specified path")
		}
	}

	if cfg.TLSBool {
		if cfg.CaCert == "" {
			return fmt.Errorf("CA certificate path must be specified if --tls is specified")
		} else if _, err := os.Stat(cfg.CaCert); os.IsNotExist(err) {
			return fmt.Errorf("CA certificate file does not exist at specified path")
		}
		if cfg.ServerCert == "" {
			return fmt.Errorf("server certificate path must be specified if --tls is specified")
		} else if _, err := os.Stat(cfg.ServerCert); os.IsNotExist(err) {
			return fmt.Errorf("server certificate file does not exist at specified path")
		}
		if cfg.ServerKey == "" {
			return fmt.Errorf("server key path must be specified if --tls is specified")
		} else if _, err := os.Stat(cfg.ServerKey); os.IsNotExist(err) {
			return fmt.Errorf("server key file does not exist at specified path")
		}
	}

	return nil
}

// Determine the mode of the Chord client.
// If the join address and join port are both specified, the mode is "join".
// Otherwise, the mode is "create".
func determineMode(cfg *Config) {
	if cfg.JoinAddress != Unspecified && cfg.JoinPort != Unspecified {
		cfg.Mode = "join"
	} else {
		cfg.Mode = "create"
	}
}

func determineAES(cfg *Config) error {
	var err error = nil
	if cfg.AESBool {
		cfg.AESKey, err = aes.LoadKey(cfg.AESKeyPath)
		return err
	}
	return nil
}

func determineTLS(cfg *Config) error {
	if cfg.TLSBool {
		var err error = nil
		cfg.ServerTLSConfig, cfg.ClientTLSConfig, err = SetupTLS(cfg.CaCert, cfg.ServerCert, cfg.ServerKey)
		return err
	}
	return nil
}
