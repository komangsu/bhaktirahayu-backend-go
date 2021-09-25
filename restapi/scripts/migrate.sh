#!/usr/bin/env bash
docker run -v $(pwd)/config/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://bhaktirahayu:123@localhost:7557/bhaktirahayu?sslmode=disable" up
