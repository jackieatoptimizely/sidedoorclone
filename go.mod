module github.com/optimizely/agent

go 1.13

require (
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/go-kit/kit v0.9.0
	github.com/google/uuid v1.1.1
	github.com/optimizely/go-sdk v1.0.0-rc1.0.20200108174332-d1b332cb2875
	github.com/orcaman/concurrent-map v0.0.0-20190826125027-8c72a8bb44f6
	github.com/rs/zerolog v1.15.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20191105231009-c1f44814a5cd // indirect
	gopkg.in/yaml.v2 v2.2.2
)

exclude (
	github.com/coreos/etcd v3.3.10+incompatible
	github.com/coreos/etcd v3.3.11+incompatible
	github.com/coreos/etcd v3.3.12+incompatible
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/gorilla/websocket v1.4.0
)
