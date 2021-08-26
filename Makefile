build:
	mkdir -p dist
	cd backend &&	go build -ldflags="-w" -trimpath -o ../dist/ipcamera-loaderd
	cp -r public dist/

build-debug:
	mkdir -p dist
	go build -o dist/ipcamera-loaderd

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

run: build-debug
	cd dist && ./ipcamera-loaderd

clean:
	rm -rf dist
