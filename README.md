# filter-apps

This is a Cloud Foundry CLI plugin to allow easier identification of started apps compared to the standard 'cf apps' command.  By default all apps will be displayed - the stopped apps will display in a different color so as to allow the started apps to be highlighted.


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
