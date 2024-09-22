package v3

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	alarmNone    = "none"
	alarmNoSpace = "no space"
	alartCorrupt = "corrupt"
)

type Client struct {
	client    *clientv3.Client
	opTimeout time.Duration
}
