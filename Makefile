all:
	@./build.sh

clean:
	@rm -f jj

install: all
	@cp jj /usr/local/bin

uninstall: 
	@rm -f /usr/local/bin/jj

package:
	@./build.sh package
