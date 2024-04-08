# ![Icon](./icon/icon_32.png) SchubsiGo
An unofficial Pushover Client for Linux written in Go.

## Features

- Uses native desktop notifications
- Login and register device via simple web interface
- Supports opening URLs directly from the notification
- Supports 2-FA logins

## Missing features from the Pushover Open Client specification

- Playing sounds with notifications (Intentionally left out since I just hate notification sounds)
- Working with or acknowledging Emergency-Priority Messages (might be implemented in the future)
- Displaying a list of all notifications sorted by Application. This is not really in the scope of this program; I only really wanted to display incoming messages as desktop notifications. It would be possible to integrate this into the webinterface, however.

# Installation

    go install github.com/jangxx/SchubsiGo@latest

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

    go install github.com/GeertJohan/go.rice/rice@latest

Install libgtk3 and libayatana-appindicator3. The exact way to do this differs by distro, so here is an example for Ubuntu (and it's derivatives):

    sudo apt install libgtk-3-dev libayatana-appindicator3-dev

Clone this repository and run

    rice embed-go
    go build

# Modifying the web interface

1. Install node.js
2. Install dependencies:
```
cd webinterface
npm i
```
3. Run vite in watch mode:
```
npm run watch
```
4. (optional) Run `rice clean`, to serve the webinterface from the built-in webserver in case you ran `rice embed-go` before.

You are now able to change the files around and vite automatically updates the _/webinterface/dist/_ directory.

If you are done, simply run `npm run build` to build and minify all assets of the webinterface.

### Alternative

Run `npm run dev` in step 3 instead to spin up a hot-reloading development server instead of serving the files through the Go server.
This has the advantage of faster iteration and not needing the Go server to be running, but it's also further away from the way everything works in production.

# Attributions

The icon is a modified version of an icon by [Freepik](https://www.freepik.com/) from [www.flaticon.com](https://www.flaticon.com/) licensed by [CC 3.0 BY](http://creativecommons.org/licenses/by/3.0/").
