

PROTO_BINARIES := $(GOPATH)/bin/protoc-gen-gogo $(GOPATH)/bin/protoc-gen-gogofast $(GOPATH)/bin/goimports $(GOPATH)/bin/protoc-gen-grpc-gateway $(GOPATH)/bin/protoc-gen-swagger /usr/local/bin/clang-format

# protoc,my.proto
define protoc
	# protoc $(1)
    [ -e ./vendor ] || go mod vendor
    protoc \
      -I /usr/local/include \
      -I $(CURDIR) \
      -I $(CURDIR)/vendor \
      -I $(GOPATH)/src \
      -I $(GOPATH)/pkg/mod/github.com/gogo/protobuf@v1.3.1/gogoproto \
      -I $(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
      --gogofast_out=plugins=grpc:$(GOPATH)/src \
      --grpc-gateway_out=logtostderr=true:$(GOPATH)/src \
      --swagger_out=logtostderr=true,fqn_for_swagger_name=true:. \
      $(1)

endef


$(GOPATH)/bin/mockery:
	go install github.com/vektra/mockery/v2@v2.10.0
$(GOPATH)/bin/controller-gen:
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1
$(GOPATH)/bin/go-to-protobuf:
	go install k8s.io/code-generator/cmd/go-to-protobuf@v0.21.5
$(GOPATH)/src/github.com/gogo/protobuf:
	[ -e $(GOPATH)/src/github.com/gogo/protobuf ] || git clone --depth 1 https://github.com/gogo/protobuf.git -b v1.3.2 $(GOPATH)/src/github.com/gogo/protobuf
$(GOPATH)/bin/protoc-gen-gogo:
	go install github.com/gogo/protobuf/protoc-gen-gogo@v1.3.2
$(GOPATH)/bin/protoc-gen-gogofast:
	go install github.com/gogo/protobuf/protoc-gen-gogofast@v1.3.2
$(GOPATH)/bin/protoc-gen-grpc-gateway:
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
$(GOPATH)/bin/protoc-gen-swagger:
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.16.0
$(GOPATH)/bin/openapi-gen:
	go install k8s.io/kube-openapi/cmd/openapi-gen@v0.0.0-20220124234850-424119656bbf
$(GOPATH)/bin/swagger:
	go install github.com/go-swagger/go-swagger/cmd/swagger@v0.28.0
$(GOPATH)/bin/goimports:
	go install golang.org/x/tools/cmd/goimports@v0.1.7

/usr/local/bin/clang-format:
ifeq ($(shell uname),Darwin)
	brew install clang-format
else
	sudo apt-get install clang-format
endif


pkg/apis/common/v1/generated.proto: $(GOPATH)/bin/go-to-protobuf $(PROTO_BINARIES) $(TYPES) $(GOPATH)/src/github.com/gogo/protobuf
	# These files are generated on a v3/ folder by the tool. Link them to the root folder
	# Format proto files. Formatting changes generated code, so we do it here, rather that at lint time.
	# Why clang-format? Google uses it.
	# [ -e ./vendor ] || go mod vendor

	$(GOPATH)/bin/go-to-protobuf \
		--drop-embedded-fields=github.com/kubeflow/common/pkg/apis/common/v1.SchedulingPolicy \
		--go-header-file=./hack/boilerplate/boilerplate.go.txt  \
		--apimachinery-packages=+k8s.io/apimachinery/pkg/util/intstr,+k8s.io/apimachinery/pkg/api/resource,+k8s.io/apimachinery/pkg/runtime/schema,+k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/apis/meta/v1beta1,k8s.io/api/core/v1,k8s.io/api/policy/v1beta1 \
		--packages=github.com/kubeflow/common/pkg/apis/common/v1 \
		--proto-import $(GOPATH)/src
	# Delete the link
	touch pkg/apis/common/v1/generated.proto
