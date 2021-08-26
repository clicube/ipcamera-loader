build:
	mkdir -p dist
	cd backend &&	go build -ldflags="-w" -trimpath -o ../dist/ipcamera-loaderd
	cd frontend && npm ci && npm run build
	mkdir -p dist/public
	cp -r frontend/out/* dist/public/
	cp ipcamera-loaderd.service dist/

build-debug:
	mkdir -p dist
	cd backend &&	go build -o ../dist/ipcamera-loaderd

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
