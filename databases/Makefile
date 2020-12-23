.PHONY: pull
pull:
	docker pull mcr.microsoft.com/mssql/server:2019-latest

.PHONY: start
start:
	docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=<123>" \
       -p 1433:1433 --name sql1 -h sql1 \
    mcr.microsoft.com/mssql/server:2019-latest
#	 mcr.microsoft.com/mssql/server:2017-CU21-ubuntu-16.04

.PHONY: find_addr
find_addr:
	 docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' sql1