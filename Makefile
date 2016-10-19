all: 
	@./build.sh
clean:
	@rm -f jsoned
install: all
	@cp jsoned /usr/local/bin
uninstall: 
	@rm -f /usr/local/bin/jsoned

