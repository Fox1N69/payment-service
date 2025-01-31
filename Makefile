local-run:
	go run cmd/payment/main.go

docker-build: 
	docker-compose up --build -d

