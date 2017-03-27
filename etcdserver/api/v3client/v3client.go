// Copyright 2017 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v3client

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver"
	"github.com/coreos/etcd/etcdserver/api/v3rpc"
	"github.com/coreos/etcd/proxy/grpcproxy/adapter"
)

// New creates a clientv3 client that wraps an in-process EtcdServer. Instead
// of making gRPC calls through sockets, the client makes direct function calls
// to the etcd server through its api/v3rpc function interfaces.
func New(s *etcdserver.EtcdServer) *clientv3.Client {
	c := clientv3.NewCtxClient(context.Background())

	kvc := adapter.KvServerToKvClient(v3rpc.NewQuotaKVServer(s))
	c.KV = clientv3.NewKVFromKVClient(kvc)

	lc := adapter.LeaseServerToLeaseClient(v3rpc.NewQuotaLeaseServer(s))
	c.Lease = clientv3.NewLeaseFromLeaseClient(lc, time.Second)

	wc := adapter.WatchServerToWatchClient(v3rpc.NewWatchServer(s))
	c.Watcher = clientv3.NewWatchFromWatchClient(wc)

	mc := adapter.MaintenanceServerToMaintenanceClient(v3rpc.NewMaintenanceServer(s))
	c.Maintenance = clientv3.NewMaintenanceFromMaintenanceClient(mc)

	return c
}
