module github.com/sunshibao/go-utils

replace golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0

replace golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e => github.com/golang/sys v0.0.0-20180909124046-d0be0721c37e

replace golang.org/x/sys v0.0.0-20190422165155-953cdadca894 => github.com/golang/sys v0.0.0-20190422165155-953cdadca894

replace golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f => github.com/golang/sync v0.0.0-20180314180146-1d60e4601c6f

replace golang.org/x/net v0.0.0-20180906233101-161cd47e91fd => github.com/golang/net v0.0.0-20180906233101-161cd47e91fd

replace golang.org/x/lint v0.0.0-20190409202823-959b441ac422 => github.com/golang/lint v0.0.0-20190409202823-959b441ac422

replace golang.org/x/tools v0.0.0-20190208222737-3744606dbb67 => github.com/golang/tools v0.0.0-20190208222737-3744606dbb67

replace golang.org/x/crypto v0.0.0-20190208162236-193df9c0f06f => github.com/golang/crypto v0.0.0-20190208162236-193df9c0f06f

require (
	github.com/astaxie/beego v1.10.1
	github.com/coreos/etcd v3.3.9+incompatible
	github.com/denisenkom/go-mssqldb v0.0.0-20181014144952-4e0d7dc8888f
	github.com/go-redsync/redsync v1.3.0
	github.com/go-sql-driver/mysql v1.4.0
	github.com/gogo/protobuf v1.1.1
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/gqcn/structs v1.1.1
	github.com/kinwyb/go v0.0.0-20210607003725-9cbcb7c895c3
	github.com/lestrrat-go/file-rotatelogs v2.3.0+incompatible
	github.com/micro/protobuf v0.0.0-20180321161605-ebd3be6d4fdb
	github.com/opentracing/opentracing-go v1.0.2
	github.com/pborman/uuid v1.2.0
	github.com/pquerna/ffjson v0.0.0-20180717144149-af8b230fcd20
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/robertkrimen/otto v0.0.0-20180617131154-15f95af6e78d
	github.com/shirou/gopsutil v2.18.12+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/smartystreets/goconvey v0.0.0-20180222194500-ef6db91d284a
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v0.0.3
	github.com/stretchr/testify v1.3.0
	github.com/vmihailenco/msgpack v4.0.0+incompatible
	github.com/vrischmann/go-metrics-influxdb v0.1.1
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2
	golang.org/x/tools v0.0.0-20190208222737-3744606dbb67
)

go 1.13
