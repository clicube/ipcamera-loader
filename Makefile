build:
	go build -o ipcamera-loaderd
	#go build -ldflags="-w" -trimpath -o ipcamera-loaderd

install: build
	sudo cp ipcamera-loaderd.service /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable ipcamera-loaderd
	sudo systemctl restart ipcamera-loaderd

pull:
	git fetch -p
	git clean -f
	git reset --hard
	git checkout raspi
	git reset --hard origin/raspi

update: pull install

build-arm6:
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-w" -trimpath -o ipcamera-loaderd

run: build
	./ipcamera-loaderd

clean:
	rm ipcamera-loaderd
