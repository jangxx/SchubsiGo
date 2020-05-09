# ![Icon](./icon/icon_32.png) SchubsiGo
An unofficial Pushover Client for Linux written in Go.

## Features

- Uses native desktop notifications
- Login and register device via simple webinterface
- Supports opening URLs directly from the notification
- Supports 2 Factor Authentication logins

## Missing features from the Pushover Open Client specification

- Playing sounds with notifications (Intentionally left out since I just hate notification sounds)
- Working with or acknowledging Emergency-Priority Messages (might be implemented in the future)
- Displaying a list of all notifications sorted by Application. This is not really in the scope of this program; I only really wanted to display incoming messages as desktop notifications. It would be possible to integrate this into the webinterface, however.

# Installation

    go get github.com/jangxx/SchubsiGo

You can also download a binary from the releases page.

# Usage

After starting the program, an icon appears in the system tray area.
Clicking on it opens the webinterface in the default webbrowser.
You can then proceed to login and register your device there.
Afterwards, you can use the webinterface to see your login status and log out.

If your OS does not use a system tray, you can use the `--no-tray` flag, to start the app without a tray icon.
In this case, you need to navigate to _http://localhost:33322_ manually to access the webinterface.

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
5. (optional) Run `rice clean`, to serve the webinterface from the built-in webserver in case you ran `rice embed-go` before.

You are now able to change the files around and gulp automatically updates the _build/_ directory.

If you are done, simply run `gulp build --production` to build and minify all assets of the webinterface.

# Attributions

The icon is a modified version of an icon by [Freepik](https://www.freepik.com/) from [www.flaticon.com](https://www.flaticon.com/) licensed by [CC 3.0 BY](http://creativecommons.org/licenses/by/3.0/").
