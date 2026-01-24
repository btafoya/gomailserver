bench-perf:
	@echo "Building performance benchmark tool..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/bench-perf ./cmd/bench-perf
	@echo "Performance benchmark tool built: $(BUILD_DIR)/bench-perf"