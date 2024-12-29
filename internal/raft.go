package raft

import (
	"log"
	"os"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

func SetupRaft(nodeID, dataDir string) *raft.Raft {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	logStore, err := raftboltdb.NewBoltStore(dataDir + "/raft-log.db")
	if err != nil {
		log.Fatalf("Failed to create log store: %v", err)
	}

	snapshots, err := raft.NewFileSnapshotStore(dataDir, 1, os.Stderr)
	if err != nil {
		log.Fatalf("Failed to create snapshot store: %v", err)
	}

	addr := "127.0.0.1:9000"
	transport, err := raft.NewTCPTransport(addr, nil, 3, 10*time.Second, os.Stderr)
	if err != nil {
		log.Fatalf("Failed to create transport: %v", err)
	}

	r, err := raft.NewRaft(config, nil, logStore, logStore, snapshots, transport)
	if err != nil {
		log.Fatalf("Failed to initialize raft: %v", err)
	}

	return r
}
