# LibreUniversity platform-api geliştirme komutları.
# Go yerelde yoksa: make GO="docker run --rm -v $$PWD:/src -w /src golang:1.26 go"

GO      ?= go
LINT    ?= golangci-lint
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)
# go.work kökünde ./... çalışmaz; çalışma alanındaki tüm modüller modül yolu deseniyle seçilir.
PKGS    := github.com/Libre-University/platform-api/...
# golangci-lint modül yolu desenini desteklemez; dizinler açıkça verilir.
LINT_PKGS := ./core/... ./platform/... ./modules/obs/... ./modules/lms/... ./cmd/libre-university/...

.PHONY: build test lint fmt vet tidy run migrate modules clean

build: ## Tek binary üret (bin/libre-university)
	$(GO) build -trimpath -ldflags '$(LDFLAGS)' -o bin/libre-university ./cmd/libre-university

test: ## Tüm çalışma alanı modüllerinde yarış dedektörüyle test
	$(GO) test -race -count=1 $(PKGS)

vet:
	$(GO) vet $(PKGS)

fmt: ## Biçimlendirme kontrolü (CI'da başarısız olur)
	@out=$$($(GO) fmt $(PKGS) ); if [ -n "$$out" ]; then echo "gofmt düzeltti:"; echo "$$out"; exit 1; fi

lint:
	$(LINT) run $(LINT_PKGS)

tidy: ## Her modülde go mod tidy
	@for m in core platform modules/obs modules/lms cmd/libre-university; do \
		echo "tidy $$m"; (cd $$m && $(GO) mod tidy) || exit 1; done

run: ## Sunucuyu çalıştır
	$(GO) run ./cmd/libre-university serve

migrate: ## Migration'ları uygula (LU_DATABASE_URL gerekli)
	$(GO) run ./cmd/libre-university migrate

modules:
	$(GO) run ./cmd/libre-university modules

clean:
	rm -rf bin
