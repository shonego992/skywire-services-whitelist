docker pull postgres
docker run --name skywire-whitelist-db -e POSTGRES_PASSWORD=supersecretpass -p 5433:5432 -d postgres