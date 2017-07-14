# filter-apps

This is a Cloud Foundry CLI plugin to allow easier identification of started apps compared to the standard 'cf apps' command.

<B>Build</B>

$ glide up
$ glide install
$ go build


<B>Run</B>

cf install-plugin started-apps


<B>Usage</B>

$ cf started-apps
 or
$ cf sa

To only show started apps
$ cf sa -x
