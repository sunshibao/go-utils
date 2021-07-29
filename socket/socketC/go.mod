module socketC

go 1.15

require (
	github.com/sunshibao/go-utils v0.0.0-20210729073106-bc51f0f9972f
	golang.org/x/text v0.3.6
	google.golang.org/grpc v1.39.0 // indirect
)

replace google.golang.org/grpc v1.39.0 => google.golang.org/grpc v1.26.0
