# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w

# Package lists
TOPLEVEL_PKG := main
INT_LIST := myinterface another/impl	#<-- Interface directories
IMPL_LIST := myimplementation another/pkg/impl	#<-- Implementation directories
CMD_LIST := mycmd1 another/pkg/mycmd1	#<-- Command directories

# List building
ALL_LIST = $(INT_LIST) $(IMPL_LIST) $(CMD_LIST)

BUILD_LIST = $(foreach int, $(ALL_LIST), $(int)_build)
RUN_LIST = $(foreach int, $(ALL_LIST), $(int)_build)
CLEAN_LIST = $(foreach int, $(ALL_LIST), $(int)_clean)
INSTALL_LIST = $(foreach int, $(ALL_LIST), $(int)_install)
IREF_LIST = $(foreach int, $(ALL_LIST), $(int)_iref)
TEST_LIST = $(foreach int, $(ALL_LIST), $(int)_test)
FMT_TEST = $(foreach int, $(ALL_LIST), $(int)_fmt)

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST)

all: build
build: $(BUILD_LIST)
run: $(RUN_LIST)
clean: $(CLEAN_LIST)
install: $(INSTALL_LIST)
test: $(TEST_LIST)
iref: $(IREF_LIST)
fmt: $(FMT_TEST)

$(BUILD_LIST): %_build: %_fmt %_iref
        $(GORUN) $(TOPLEVEL_PKG)/$*
$(BUILD_LIST): %_build: %_fmt %_iref
	$(GOBUILD) $(TOPLEVEL_PKG)/$*
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) $(TOPLEVEL_PKG)/$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) $(TOPLEVEL_PKG)/$*
$(IREF_LIST): %_iref:
	$(GODEP) $(TOPLEVEL_PKG)/$*
$(TEST_LIST): %_test:
	$(GOTEST) $(TOPLEVEL_PKG)/$*
$(FMT_TEST): %_fmt:
	$(GOFMT) ./$*
