# filter-apps

This is a Cloud Foundry CLI plugin to allow easier identification of started apps compared to the standard 'cf apps' command.


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
