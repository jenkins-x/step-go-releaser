module github.com/jenkins-x/step-goreleaser

go 1.13

require (
	cloud.google.com/go v0.55.0 // indirect
	github.com/Azure/go-autorest/autorest v0.9.6 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/jenkins-x/jx-api v0.0.17
	github.com/jenkins-x/jx-helpers v1.0.44
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073 // indirect
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e // indirect
	k8s.io/api v0.17.6 // indirect
	k8s.io/apimachinery v0.17.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/utils v0.0.0-20200124190032-861946025e34 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.17.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.2
	k8s.io/client-go => k8s.io/client-go v0.16.5
)
