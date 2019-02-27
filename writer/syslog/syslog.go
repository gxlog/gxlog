package syslog

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type syslog struct {
	network string
	addr    string
	host    string
	conn    net.Conn
}

func syslogDial(network, addr string) (*syslog, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	log := &syslog{
		network: network,
		addr:    addr,
		host:    host,
	}
	if err := log.connect(); err != nil {
		return nil, err
	}
	return log, nil
}

func (log *syslog) Write(timestamp time.Time, priority int, tag string, msg []byte) error {
	if log.conn != nil {
		if err := log.write(timestamp, priority, tag, msg); err == nil {
			return nil
		} else {
			log.Close()
		}
	}
	if err := log.connect(); err != nil {
		return err
	}
	return log.write(timestamp, priority, tag, msg)
}

func (log *syslog) Close() error {
	if log.conn != nil {
		err := log.conn.Close()
		log.conn = nil
		return err
	}
	return nil
}

func (log *syslog) connect() error {
	var conn net.Conn
	var err error
	if log.network == "" {
		conn, err = dialLocal()
	} else {
		conn, err = net.Dial(log.network, log.addr)
	}
	if err != nil {
		return err
	}
	log.conn = conn
	return nil
}

func (log *syslog) write(timestamp time.Time, priority int, tag string, msg []byte) error {
	var err error
	if log.network == "" {
		_, err = fmt.Fprintf(log.conn, "<%d>%s %s[%d]: %s",
			priority, timestamp.Format(time.Stamp), tag, os.Getpid(), msg)
	} else {
		_, err = fmt.Fprintf(log.conn, "<%d>%s %s %s[%d]: %s",
			priority, timestamp.Format(time.RFC3339), log.host, tag, os.Getpid(), msg)
	}
	return err
}

func dialLocal() (net.Conn, error) {
	networks := []string{"unixgram", "unix"}
	paths := []string{"/dev/log", "/var/run/syslog", "/var/run/log"}
	for _, network := range networks {
		for _, path := range paths {
			conn, err := net.Dial(network, path)
			if err == nil {
				return conn, nil
			}
		}
	}
	return nil, errors.New("Unix syslog delivery error")
}
