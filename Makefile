SERVICE_NAME = email-services
WORKER_MAIN_FILE = server_app
BUILD_DIR = build
GO= go

#linux
# clean build file
cleanl:
	echo "remove bin exe"
	rm -rf $(BUILD_DIR)

# build binary
buildl:
	echo "build binary execute file"
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(WORKER_MAIN_FILE)_linux main.go

runl:
	make buildl

	echo "Run service application"
	cd $(BUILD_DIR) && ./$(WORKER_MAIN_FILE)_linux


#windows
cleanw:
	echo "remove bin exe"
	rd /s build

buildw:
	echo "build binary execute file"
	go build -o $(BUILD_DIR)/$(WORKER_MAIN_FILE)_win.exe main.go

runw:
	.\$(BUILD_DIR)\$(WORKER_MAIN_FILE)_win.exe

startw:
	make buildw
	.\$(BUILD_DIR)\$(WORKER_MAIN_FILE)_win.exe