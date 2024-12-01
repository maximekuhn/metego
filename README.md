# metego

This project contains the code of my weather station. It is intended to run on a Raspberry Pi
(model 3B+), with an external display. It is not guaranteed that it will work and/or display in a correct
way on yours.

- [ ] make it work
- [ ] make it pretty
- [ ] make it work again

## Tech stack
- Golang
- templ
- HTMX
- SQLite

## Requirements
- [Golang](https://go.dev/)
- [golangci-lint](https://github.com/golangci/golangci-lint)
- [templ](https://templ.guide)
- [Docker](https://www.docker.com/)
    - mainly for cross compiling to Raspberry Pi
- [go-task](https://taskfile.dev/installation/)
    - build tool

## Installation on a Raspberry Pi
### Setup LCD screen 
Install LCD screen drivers:
```shell
git clone https://github.com/goodtft/LCD-show
cd LCD-show
chmod -R 755 LCD-show/
sudo ./LCD7b-show
```

### Cross-compile
On your computer, compile for raspberry Pi:
```shell
task build-rpi
```
> you need Docker up and running

The binary will be located at `./bin/rpi/web`.

### Transfer the binary to the raspberry Pi

Find a way to transfer it to your raspberry Pi.
If you have python installed, you can create a quick http server.
```shell
python3 -m http.server 8000
```

On your raspberry Pi, simply get the binary from your host
```shell
wget http://<YOUR_HOST_IP>:8000/bin/rpi/web
```

Don't forget to make the file executable if it's not already the case:
```shell
chmod +x web
```

### Create a .env file
Next to your binary, create a .env file and set `OPEN_WEATHER_API_KEY` to your API key.

If you don't have one, you can get one for free: https://openweathermap.org/api.

### Create a service
To ensure the app starts everytime with the Raspberry Pi, simply create a service using systemd.

First, create a new file in systemd directory:
```
sudo touch /etc/systemd/system/metego.service
```

Then, edit the file:
```text
[Unit]
Description=Metego (Weather station)
After=network.target

[Service]
ExecStart=/home/pi/Documents/metego/web
WorkingDirectory=/home/pi/Documents/metego
Restart=always
RestartSec=3
StandardOutput=syslog
StandardError=syslog

[Install]
WantedBy=multi-user.target
```
> you might need to change `ExecStart` and `WorkingDirectory` to point to the correct bin/dir

Enable the service:
```shell
sudo systemctl enable metego.service
```

Start the service:
```shell
sudo systemctl start metego.service
```

Optionally, check the status:
```shell
sudo systemctl status metego
```

### Show the web app at Pi's startup
Copy autostart configuration:
```shell
cp /etc/xdg/lxsession/LXDE-pi/autostart ~/.config/lxsession/LXDE-pi
```

Edit the config:
```text
@lxpanel --profile LXDE-pi
@pcmanfm --desktop --profile LXDE-pi
#@xscreensaver -no-splash
point-rpi
chromium --start http://localhost:9004/weather/<YOUR_CITY>
```

// TODO: update docs

## Create a release
Tag the commit you want to create a release from.

For example:
```shell
git tag -a v.0.0.5.rc-1 -m "release candidate 1 for v0.0.5"
git push --tags
```

If the CI pass, a new release should be created and available [here](https://github.com/maximekuhn/metego/releases/latest).
