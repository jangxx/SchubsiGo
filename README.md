# SchubsiGo
An unofficial Pushover Client for Linux written in Go.

## Features

- Uses native desktop notifications
- Login and register device via simple webinterface
- Supports opening URLs directly from the notification
- Supports 2 Factor Authentication logins

## Missing features from the Pushover Open Client specification

- Playing sounds with notifications (Intentionally left out since I just hate notification sounds)
- Working with or acknowledging Emergency-Priority Messages (might be implemented in the future)

# Installation

    go get github.com/jangxx/SchubsiGo

You can also download a binary from the releases page.

# Building the binary yourself

Install go.rice by running

    go get github.com/GeertJohan/go.rice
    go get github.com/GeertJohan/go.rice/rice

Install libgtk3. The exact way to do this differs by distro, so here is an example for Ubuntu (and it's derivatives):

    sudo apt install libgtk-3-dev

Clone this repository and run

    rice embed-go
    go build

# Modifying the web interface

1. Install node.js
2. Install gulp
```
npm i -g gulp-cli
```
3. Install dependencies:
```
cd webinterface
npm i
```
4. Run gulp
```
gulp
```

You are now able to change the files around and gulp automatically updates the _build/_ directory.

If you are done, simply run `gulp build --production` to build and minify all assets of the webinterface.