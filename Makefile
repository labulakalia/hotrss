VERSION=`git tag | tail -1`
COMMIT=`git rev-parse --short HEAD`
BUILDDATE=`date "+%Y-%m-%d"`

BUILD_DIR=build
APP_NAME=hotrss

sources=$(wildcard *.go)

build = CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -o ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2)$(3)  cmd/main.go 
md5 = md5sum ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2)$(3) > ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2)_checksum.txt
tar = tar -cvzf ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2).tar.gz  -C ${BUILD_DIR}  $(APP_NAME)-$(1)-$(2)$(3) $(APP_NAME)-$(1)-$(2)_checksum.txt
delete = rm -rf ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2)$(3) ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2)_checksum.txt

LINUX = linux-amd64

WINDOWS = windows-amd64-.exe

DARWIN = darwin-amd64

ALL = $(LINUX) \
	$(WINDOWS) \
	$(DARWIN)

build_linux: $(LINUX:%=build/%)

build_windows: $(WINDOWS:%=build/%)

build_darwin: $(DARWIN:%=build/%)

build_all: $(ALL:%=build/%)

build/%: 
	$(call build,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)),$(word 3, $(subst -, ,$*)))
	$(call md5,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)),$(word 3, $(subst -, ,$*)))
	$(call tar,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)),$(word 3, $(subst -, ,$*)))
	$(call delete,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)),$(word 3, $(subst -, ,$*)))

clean:
	rm -rf ${BUILD_DIR}

vet:
	go vet cmd/main.go

build:
	go build -o hotrss -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" cms/main.go
