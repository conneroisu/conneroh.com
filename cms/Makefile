all: tt
tt: 
	docker build -t hello:latest .
	docker run -p 8055:8055 --privileged hello:latest
