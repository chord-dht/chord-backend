package handlers

import (
	"chord-backend/config"
	"fmt"
	"time"

	"github.com/chord-dht/chord-core/node"
	"github.com/chord-dht/chord-core/storage"
)

var LocalNode *node.Node

func NewNodeWithConfig(
	cfg *config.Config,
	storageFactory func(string) (storage.Storage, error),
) (*node.Node, error) {
	chordNode, err := node.NewNode(
		cfg.IdentifierLength,
		cfg.SuccessorsLength,
		cfg.IpAddress,
		cfg.Port,
		storageFactory,
		cfg.StorageDir,
		cfg.BackupDir,
		time.Duration(cfg.StabilizeTime),
		time.Duration(cfg.FixFingersTime),
		time.Duration(cfg.CheckPredecessorTime),
		cfg.TLSBool,
		cfg.ServerTLSConfig,
		cfg.ClientTLSConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating node: %w", err)
	}

	return chordNode, nil
}
