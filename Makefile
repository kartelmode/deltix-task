BIN_PATH:=./bin/main
SOURCE_PATH=./cmd/main.go

build:
	mkdir -p ./bin
	GOARCH=amd64 GOOS=linux go build -o ${BIN_PATH} $(SOURCE_PATH)

copy_small:
	cp small/market_data.csv market_data.csv
	cp small/user_data.csv user_data.csv

copy_medium:
	cp medium/market_data.csv market_data.csv
	cp medium/user_data.csv user_data.csv
	
copy_big:
	cp big/market_data.csv market_data.csv
	cp big/user_data.csv user_data.csv

run: build
	# time ${BIN_PATH} 1h
	# time ${BIN_PATH} 1d
	time ${BIN_PATH} 30d

clean:
	rm market_data.csv
	rm user_data.csv
	rm ./bin/main
	rmdir ./bin