STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo
install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/smirl/cortex@latest/steampipe-plugin-cortex.plugin -tags "${BUILD_TAGS}" *.go
install-local:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/local/cortex/steampipe-plugin-cortex.plugin -tags "${BUILD_TAGS}" *.go
