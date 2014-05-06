#!/bin/bash

go run replay.go -database=backend.db.randwick -vid BUS1A >/dev/null &
sleep 60 
go run replay.go -database=backend.db.randwick -vid BUS2 >/dev/null &
sleep 60
go run replay.go -database=backend.db.randwick -vid BUS3 >/dev/null &


