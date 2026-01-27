.PHONY: bump

bump:
	@if [ -z "$(mod)" ]; then \
		echo "Error: mod is required. Usage: make bump mod=<module_name> ver=<version>"; \
		exit 1; \
	fi
	@if [ -z "$(ver)" ]; then \
		echo "Error: ver is required. Usage: make bump mod=<module_name> ver=<version>"; \
		exit 1; \
	fi
	@if [ ! -d "$(mod)" ]; then \
		echo "Error: module '$(mod)' does not exist"; \
		exit 1; \
	fi
	git tag $(mod)/$(ver)
	@echo "Successfully created tag $(mod)/$(ver)"
	@echo "Don't forget to push: git push origin $(mod)/$(ver)"
