package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
	"code.cloudfoundry.org/cli/plugin/models"
	"github.com/cloudfoundry/cli/plugin"
)

//	"github.com/cloudfoundry/cli/plugin/models"

// AppInfo is a struct representing the output seen on screen.
type AppInfo struct {
	Name      string
	Instances string
	Memory    string
	Disk      string
	Urls      string
	State     string
}

// FilterApps represents this application.
type FilterApps struct {
	Connection plugin.CliConnection
	UI         terminal.UI
}

// Run is the main method.
func (c *FilterApps) Run(cliConnection plugin.CliConnection, args []string) {
	c.Connection = cliConnection
	var err error

	ui := terminal.NewUI(
		os.Stdin,
		os.Stdout,
		terminal.NewTeePrinter(os.Stdout),
		trace.NewLogger(os.Stdout, false, "false", ""),
	)
	c.UI = ui
	var apps []plugin_models.GetAppsModel
	all := true

	/* check for proper usage */
	if len(args) > 2 ||
		(len(args) == 2 && args[1] != "-a") {
		c.UI.Say(terminal.FailureColor("\nError - invalid usage.\n"))
		_, err := cliConnection.CliCommand("help", "started-apps")
		if err != nil {
			c.UI.Say(err.Error())
			return
		}
		return
	}

	/* check for excluded flag */
	if len(args) == 2 && args[1] == "-a" {
		all = true
	} else {
		all = false
	}

	username, _ := c.Connection.Username()
	var org, space string
	var currentOrg plugin_models.Organization
	var currentSpace plugin_models.Space
	currentOrg, err = c.Connection.GetCurrentOrg()
	if err == nil {
		org = currentOrg.OrganizationFields.Name
	}
	currentSpace, err = c.Connection.GetCurrentSpace()
	if err == nil {
		space = currentSpace.Name
	}
	c.UI.Say("\nGetting apps in org " + terminal.PromptColor(org) + " / space " + terminal.PromptColor(space) + " as " + terminal.PromptColor(username) + "...")

	apps, err = c.Connection.GetApps()

	if err != nil {
		c.UI.Say(terminal.FailureColor("FAILED"))
		c.UI.Say(err.Error())
		return
	}

	appInfo := RetrieveApps(apps, all)

	c.UI.Say("OK\n")

	if len(appInfo) > 0 {
		table := c.UI.Table([]string{"name", "requested state", "instances", "memory", "disk", "urls"})
		for _, app := range appInfo {
			if app.State == "started" {
				table.Add(app.Name, app.State, app.Instances, app.Memory, app.Disk, app.Urls)
			} else {
				table.Add(terminal.FailureColor(app.Name), terminal.FailureColor(app.State),
					terminal.FailureColor(app.Instances), terminal.FailureColor(app.Memory),
					terminal.FailureColor(app.Disk), terminal.FailureColor(app.Urls))
			}
		}
		table.Print()
		c.UI.Say("\n")
	} else {
		fmt.Printf("No %s apps found", terminal.PromptColor("started"))
		fmt.Println("\n")
	}
}

// RetrieveApps function iterates through the results of the GetApps() function and creates an AppInfo array.
func RetrieveApps(apps []plugin_models.GetAppsModel, all bool) (appInfo []AppInfo) {
	appInfo = make([]AppInfo, 0, 50)
	for _, app := range apps {
		if app.State == "started" || all == true {
			routes := app.Routes
			u := []string{}
			for _, r := range routes {
				u = append(u, r.Host+"."+r.Domain.Name)
			}
			inst := strconv.Itoa(app.RunningInstances) + "/" + strconv.Itoa(app.TotalInstances)
			a := AppInfo{Name: app.Name, Instances: inst, Memory: ConvertSize(app.Memory), Disk: ConvertSize(app.DiskQuota), Urls: strings.Join(u, ", "), State: app.State}
			appInfo = append(appInfo, a)
		}
	}
	return appInfo
}

// ConvertSize will properly format the size value in Megabytes or Gigabytes.
func ConvertSize(size int64) (formattedSize string) {
	var val string
	if size < 1024 {
		val = strconv.FormatInt(size, 10) + "M"
	} else {
		gig := size / 1024
		val = strconv.FormatInt(gig, 10) + "G"
	}
	return val
}

// GetMetadata function is required for cli plugin.
func (c *FilterApps) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "started-apps",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 1,
		},
		Commands: []plugin.Command{
			{
				Name:     "started-apps",
				Alias:    "sa",
				HelpText: "List started apps in the target space",
				UsageDetails: plugin.Usage{
					Usage: "cf started-apps",
					Options: map[string]string{
						"a": "Include stopped apps",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(FilterApps))
}
