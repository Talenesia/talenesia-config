install:
	@go build -o talenesia
	@mv talenesia /usr/bin
	@ln talenesia.yaml /opt/config/talenesia.yaml