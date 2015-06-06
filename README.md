# Project template

Golang, bootstrap, startbootstrap-sb-admin-2 and MySQL

[SB Admin 2](http://startbootstrap.com/template-overviews/sb-admin-2/) is an open source, admin dashboard template for [Bootstrap](http://getbootstrap.com/) created by [Start Bootstrap](http://startbootstrap.com/). Copyright Iron Summit Media Strategies, LLC., Apache Licensed.

## Preparing the environment

Prerequisites:

- Git
- rsync
- GNU Make
- [Go](http://golang.org) 1.0.3 or newer

First, you should make a copy of this directory, and prepare the new project:

	cp -r simple gong
	cd gong
	./bootstrap.sh

Your project is now called **gong** and is ready to use.

Make sure the Go compiler is installed and `$GOPATH` is set.

Install dependencies, and compile:

	make deps
	make clean
	make all

Generate a self-signed SSL certificate (optional):

	cd ssl
	make

Set up MySQL (optional):

	sudo mysql < assets/files/database.sql

Edit the config file and run the server (check MySQL settings):

	vi gong.conf
	./gong

Install, uninstall. Edit Makefile and set PREFIX to the target directory:

	sudo make install
	sudo make uninstall

Allow non-root process to listen on low ports:

	/sbin/setcap 'cap_net_bind_service=+ep' /opt/gong/server

License:

MIT 
