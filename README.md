# started-apps

This is a Cloud Foundry CLI plugin to allow easier identification of started apps compared to the standard 'cf apps' command.  By default all apps will be displayed - the stopped apps will display in a different color so as to allow the started apps to be highlighted.


Here's what a typical 'cf apps' looks like.
![Screenshot](screenshots/classic.png?raw=true)


Now here's what 'cf sa' looks like.  You'll only get the started apps!
![](screenshots/new.png?raw=true)


If you want to see all apps, but still easily tell which are actually started use 'cf sa -a'.
![](screenshots/new_show_all.png?raw=true)


## Build
```
$ glide up
$ glide install
$ go build
```


## Run
```
$ cf install-plugin started-apps
```


## Usage
```
$ cf started-apps
```
 or
 ```
$ cf sa
```

To only show started apps
```
$ cf sa -x
```
