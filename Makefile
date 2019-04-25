BUILD_TAGS = ledger
BUILD_FLAGS = -tags "${BUILD_TAGS}"
install:
	go install $(BUILD_FLAGS) ./example/ledger-keys
