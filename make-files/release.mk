PACKAGE_STAGE_ROOT ?= build/pxp
PACKAGE_VERSION_DIR ?= $(PACKAGE_STAGE_ROOT)/$(VERSION)
PACKAGE_META_DIR ?= $(PACKAGE_VERSION_DIR)/meta
PACKAGE_BACKEND_DIR ?= $(PACKAGE_VERSION_DIR)/backend
PACKAGE_FRONTEND_DIR ?= $(PACKAGE_VERSION_DIR)/web-admin
PACKAGE_HASH_FILE ?= $(PACKAGE_VERSION_DIR)/hashes.txt
PACKAGE_AUDIT_LOG ?= $(PACKAGE_VERSION_DIR)/audit.log
PACKAGE_SIGNATURE ?= $(PACKAGE_VERSION_DIR)/signature.json
ADVISORY_SOURCE_DIR ?= build/security/advisories
ADVISORY_DIST_ROOT ?= dist/security
ADVISORY_DIST_VERSION ?= $(ADVISORY_DIST_ROOT)/$(VERSION)

.PHONY: package-pxp
package-pxp: verify-manifest build frontend-build
	@echo "[package] Staging artefacts in $(PACKAGE_VERSION_DIR)"
	@rm -rf $(PACKAGE_VERSION_DIR)
	@mkdir -p $(PACKAGE_META_DIR) $(PACKAGE_BACKEND_DIR) $(PACKAGE_FRONTEND_DIR)
	@cp $(PLUGIN_FILE) $(PACKAGE_META_DIR)/plugin.yaml
	@if [ -f "$(MANIFEST_FILE)" ]; then \
		cp "$(MANIFEST_FILE)" $(PACKAGE_META_DIR)/manifest.yaml; \
	else \
		echo "⚠️  manifest file $(MANIFEST_FILE) not found; produced bundle without manifest"; \
	fi
	@if [ -f "$(BUILD_DIR)/plugin" ]; then cp $(BUILD_DIR)/plugin $(PACKAGE_BACKEND_DIR)/; fi
	@if [ -f "$(BUILD_DIR)/migrate" ]; then cp $(BUILD_DIR)/migrate $(PACKAGE_BACKEND_DIR)/; fi
	@if [ -d "$(FRONTEND_OUTPUT)" ] && [ -n "$$\(ls -A $(FRONTEND_OUTPUT) 2>/dev/null)" ]; then \
		cp -R $(FRONTEND_OUTPUT)/. $(PACKAGE_FRONTEND_DIR)/; \
	else \
		echo "⚠️  No frontend output found at $(FRONTEND_OUTPUT)"; \
	fi
	@python scripts/hash_package.py "$(PACKAGE_VERSION_DIR)" "$(PACKAGE_HASH_FILE)"
	@{ \
		echo "created_at=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)"; \
		echo "plugin_version=$(VERSION)"; \
		echo "source_commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)"; \
		echo "staged_dir=$(PACKAGE_VERSION_DIR)"; \
	} > $(PACKAGE_AUDIT_LOG)
	@printf '{\n  "status": "pending",\n  "signed_by": "",\n  "signed_at": "",\n  "note": "Upload package to signing service to finalize signatures"\n}\n' > $(PACKAGE_SIGNATURE)
	@echo "[package] Artefacts staged. Hashes recorded in $(PACKAGE_HASH_FILE)"
	@rm -rf $(ADVISORY_DIST_VERSION)
	@mkdir -p $(ADVISORY_DIST_VERSION)
	@if [ -d "$(ADVISORY_SOURCE_DIR)" ] && [ -n "$$(ls -A $(ADVISORY_SOURCE_DIR) 2>/dev/null)" ]; then \
		echo "[package] Bundling advisories from $(ADVISORY_SOURCE_DIR) into $(ADVISORY_DIST_VERSION)"; \
		cp -R $(ADVISORY_SOURCE_DIR)/. $(ADVISORY_DIST_VERSION)/; \
	else \
		echo "[package] No advisories found at $(ADVISORY_SOURCE_DIR); skipping bundle"; \
	fi
