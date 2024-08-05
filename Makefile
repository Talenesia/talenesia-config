install:
	@go build -o talenesia
	@mv talenesia /usr/bin
	@ln talenesia.yml /opt/config/talenesia.yml