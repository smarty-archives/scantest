#!/usr/bin/make -f

compile:
	gb build all
	cp bin/scantest bin/scantest-web ~/bin

clean:
	rm -rf bin pkg