#!/bin/bash

#small script to run replay.go with a pre-set database and different vehicle identifiers.. usefull for quickly simulating multiple vehicles sending packets to the backend

go run replay.go -database=backend.db.randwick -vid BUS1A >/dev/null &
sleep 10 
go run replay.go -database=backend.db.randwick -vid BUS2 >/dev/null &
sleep 10
go run replay.go -database=backend.db.randwick -vid BUS3 >/dev/null &
sleep 10
go run replay.go -database=backend.db.randwick -vid YOLO >/dev/null &

