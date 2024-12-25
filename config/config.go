package config

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/chord-dht/chord-backend/aes"
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
		StorageDir:           "",
		BackupDir:            "",
		AESBool:              true,
		AESKeyPath:           "",
		TLSBool:              true,
		CaCert:               "",
		ServerCert:           "",
		ServerKey:            "",
	}
}

func JsonToConfig(json map[string]interface{}) (*Config, error) {
	cfg := NewConfig()

	if err := parseConfig(json, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseConfig(json map[string]interface{}, cfg *Config) error {
	var err error

	if cfg.IdentifierLength, err = getIntFromJson(json, "IdentifierLength"); err != nil {
		return err
	}

	if cfg.SuccessorsLength, err = getIntFromJson(json, "SuccessorsLength"); err != nil {
		return err
	}

	if cfg.IpAddress, err = getStringFromJson(json, "IpAddress"); err != nil {
		return err
	}

	if cfg.Port, err = getStringFromJson(json, "Port"); err != nil {
		return err
	}

	if cfg.Mode, err = getStringFromJson(json, "Mode"); err != nil {
		return err
	}

	if cfg.Mode == "join" {
		if cfg.JoinAddress, err = getStringFromJson(json, "JoinAddress"); err != nil {
			return err
		}

		if cfg.JoinPort, err = getStringFromJson(json, "JoinPort"); err != nil {
			return err
		}
	}

	if cfg.StabilizeTime, err = getIntFromJson(json, "StabilizeTime"); err != nil {
		return err
	}

	if cfg.FixFingersTime, err = getIntFromJson(json, "FixFingersTime"); err != nil {
		return err
	}

	if cfg.CheckPredecessorTime, err = getIntFromJson(json, "CheckPredecessorTime"); err != nil {
		return err
	}

	if cfg.StorageDir, err = getStringFromJson(json, "StorageDir"); err != nil {
		return err
	}

	if cfg.BackupDir, err = getStringFromJson(json, "BackupDir"); err != nil {
		return err
	}

	if cfg.AESBool, err = getBoolFromJson(json, "AESBool"); err != nil {
		return err
	}

	if cfg.AESBool {
		if cfg.AESKeyPath, err = getStringFromJson(json, "AESKeyPath"); err != nil {
			return err
		}
	}

	if cfg.TLSBool, err = getBoolFromJson(json, "TLSBool"); err != nil {
		return err
	}

	if cfg.TLSBool {
		if cfg.CaCert, err = getStringFromJson(json, "CaCert"); err != nil {
			return err
		}

		if cfg.ServerCert, err = getStringFromJson(json, "ServerCert"); err != nil {
			return err
		}

		if cfg.ServerKey, err = getStringFromJson(json, "ServerKey"); err != nil {
			return err
		}
	}

	return nil
}

func getIntFromJson(json map[string]interface{}, key string) (int, error) {
	val, ok := json[key]
	if !ok {
		return 0, fmt.Errorf("%s must be specified", key)
	}
	if v, ok := val.(float64); ok {
		return int(v), nil
	}
	return 0, fmt.Errorf("%s must be a float64, got %T", key, val)
}

func getStringFromJson(json map[string]interface{}, key string) (string, error) {
	val, ok := json[key]
	if !ok {
		return "", fmt.Errorf("%s must be specified", key)
	}
	if v, ok := val.(string); ok {
		return v, nil
	}
	return "", fmt.Errorf("%s must be a string, got %T", key, val)
}

func getBoolFromJson(json map[string]interface{}, key string) (bool, error) {
	val, ok := json[key]
	if !ok {
		return false, fmt.Errorf("%s must be specified", key)
	}
	if v, ok := val.(bool); ok {
		return v, nil
	}
	return false, fmt.Errorf("%s must be a bool, got %T", key, val)
}

func ValidateAndSetConfig(cfg *Config) error {
	if err := validateConfig(cfg); err != nil {
		return fmt.Errorf("failed to validate config: %v", err)
	}

	if err := determineAES(cfg); err != nil {
		return fmt.Errorf("failed to determine AES: %v", err)
	}

	if err := determineTLS(cfg); err != nil {
		return fmt.Errorf("failed to determine TLS: %v", err)
	}
	return nil
}

func CheckPortAvailability(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

func CheckRemoteAddressAvailability(address string, port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", address, port), 5*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func validateConfig(cfg *Config) error {
	if cfg.IdentifierLength < 1 || cfg.IdentifierLength > 160 {
		return fmt.Errorf("identifier length must be in the range of [1,160]")
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

	if cfg.Mode != "create" && cfg.Mode != "join" {
		return fmt.Errorf("mode must be either 'create' or 'join'")
	}

	if cfg.Mode == "join" {
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
			return fmt.Errorf("AES key path must be specified if AESBool is true")
		}
		if _, err := os.Stat(cfg.AESKeyPath); os.IsNotExist(err) {
			return fmt.Errorf("AES key file does not exist at specified path")
		}
	}

	if cfg.TLSBool {
		if cfg.CaCert == "" {
			return fmt.Errorf("CA certificate path must be specified if TLSBool is true")
		} else if _, err := os.Stat(cfg.CaCert); os.IsNotExist(err) {
			return fmt.Errorf("CA certificate file does not exist at specified path")
		}
		if cfg.ServerCert == "" {
			return fmt.Errorf("server certificate path must be specified if TLSBool is true")
		} else if _, err := os.Stat(cfg.ServerCert); os.IsNotExist(err) {
			return fmt.Errorf("server certificate file does not exist at specified path")
		}
		if cfg.ServerKey == "" {
			return fmt.Errorf("server key path must be specified if TLSBool is true")
		} else if _, err := os.Stat(cfg.ServerKey); os.IsNotExist(err) {
			return fmt.Errorf("server key file does not exist at specified path")
		}
	}

	return nil
}

func determineAES(cfg *Config) error {
	var err error
	if cfg.AESBool {
		cfg.AESKey, err = aes.LoadKey(cfg.AESKeyPath)
		return err
	}
	return nil
}

func determineTLS(cfg *Config) error {
	if cfg.TLSBool {
		var err error
		cfg.ServerTLSConfig, cfg.ClientTLSConfig, err = SetupTLS(cfg.CaCert, cfg.ServerCert, cfg.ServerKey)
		return err
	}
	return nil
}
