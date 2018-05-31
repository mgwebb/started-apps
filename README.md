# started-apps

This is a Cloud Foundry CLI plugin to allow easier identification of started apps compared to the standard 'cf apps' command.  By default only started apps will be displayed.  Using the optional '-a' parameter you'll see all apps, but the stopped apps will display in a different color so as to allow the started apps to still be easily identified.


Here's what a typical 'cf apps' looks like.
![Screenshot](screenshots/classic.png?raw=true)


Now here's what 'cf sa' looks like.  You'll only get the started apps!
![](screenshots/new.png?raw=true)


If you want to see all apps, but still easily tell which are actually started use 'cf sa -a'.
![](screenshots/new_show_all.png?raw=true)


## Install

Download the appropriate binary from the latest release.

To install for MacOS:
```
$ cf install-plugin started-apps-macos
```


## Usage
```
$ cf started-apps
```
 or
 ```
$ cf sa
```

To show all apps but with stopped apps in a different color use
```
$ cf sa -a
```
