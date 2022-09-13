

PROTO_BINARIES := $(GOPATH)/bin/protoc-gen-gogo $(GOPATH)/bin/protoc-gen-gogofast $(GOPATH)/bin/goimports

$(GOPATH)/bin/go-to-protobuf:
	go install k8s.io/code-generator/cmd/go-to-protobuf@v0.21.5
$(GOPATH)/src/github.com/gogo/protobuf:
	[ -e $(GOPATH)/src/github.com/gogo/protobuf ] || git clone --depth 1 https://github.com/gogo/protobuf.git -b v1.3.2 $(GOPATH)/src/github.com/gogo/protobuf
$(GOPATH)/bin/protoc-gen-gogo:
	go install github.com/gogo/protobuf/protoc-gen-gogo@v1.3.2
$(GOPATH)/bin/protoc-gen-gogofast:
	go install github.com/gogo/protobuf/protoc-gen-gogofast@v1.3.2
$(GOPATH)/bin/goimports:
	go install golang.org/x/tools/cmd/goimports@v0.1.7


pkg/apis/common/v1/generated.proto: $(GOPATH)/bin/go-to-protobuf $(PROTO_BINARIES) $(TYPES) $(GOPATH)/src/github.com/gogo/protobuf
	# These files are generated on a v3/ folder by the tool. Link them to the root folder
	# Format proto files. Formatting changes generated code, so we do it here, rather that at lint time.
	# Why clang-format? Google uses it.
	# [ -e ./vendor ] || go mod vendor

	$(GOPATH)/bin/go-to-protobuf \
		--go-header-file=./hack/boilerplate/boilerplate.go.txt  \
		--packages=github.com/kubeflow/common/pkg/apis/common/v1 \
		--apimachinery-packages=+k8s.io/apimachinery/pkg/util/intstr,+k8s.io/apimachinery/pkg/api/resource,k8s.io/apimachinery/pkg/runtime/schema,+k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/api/core/v1,k8s.io/api/policy/v1beta1 \
		--proto-import $(GOPATH)/src
	# Delete the link
	touch pkg/apis/common/v1/generated.proto