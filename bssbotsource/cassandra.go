package main

import (
	"github.com/gocql/gocql"
)

func CassandraConn() (*gocql.Session, error) {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1", "127.0.0.2", "127.0.0.3")
	cluster.Keyspace = "botapi"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
