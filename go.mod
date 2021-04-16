module github.com/banzaicloud/cluster-registry-controller

go 1.15

require (
	emperror.dev/errors v0.8.0
	github.com/banzaicloud/cluster-registry v0.0.0-20210408202748-0ca595389ef1
	github.com/banzaicloud/k8s-objectmatcher v1.5.0
	github.com/banzaicloud/operator-tools v0.21.1
	github.com/go-logr/logr v0.2.1
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/gomodule/redigo v1.8.4 // indirect
	github.com/onsi/ginkgo v1.15.1
	github.com/onsi/gomega v1.10.1
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.4.0
	github.com/throttled/throttled v2.2.5+incompatible
	go.uber.org/zap v1.13.0
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	sigs.k8s.io/controller-runtime v0.6.5
)

replace github.com/banzaicloud/cluster-registry-controller/static => ./static
