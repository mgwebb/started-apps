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
// added comment

// AppInfo is a struct representing the output seen on screen.
type AppInfo struct {
	Space     string
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
	all := true
	allorg := false

	//c.UI.Say(terminal.FailureColor("\nargs len = " + strconv.Itoa(len(args))))

	invalidargs := CheckProperUsage(args)

	if invalidargs == true {
		c.UI.Say(terminal.FailureColor("\nError - invalid usage.\n"))
		_, err := cliConnection.CliCommand("help", "started-apps")
		if err != nil {
			c.UI.Say(err.Error())
			return
		}
		return
	}
	/*
	  We've already checked for correct number of args and valid values,
	  so now just determine what user is asking for.
	*/
	if len(args) == 2 {
		if args[1] == "-a" {
			all = true
			allorg = false
		} else if args[1] == "-o" {
			all = false
			allorg = true
		}
	} else if len(args) == 1 {
		all = false
		allorg = false
	} else {
		all = true
		allorg = true
	}

	//TODO - remove these once testing is completed.
	//c.UI.Say(terminal.FailureColor("all = " + strconv.FormatBool(all)))
	//c.UI.Say(terminal.FailureColor("allorg = " + strconv.FormatBool(allorg)))

	var org string
	var currentOrg plugin_models.Organization
	var currentSpace plugin_models.Space
	currentOrg, err = c.Connection.GetCurrentOrg()
	username, _ := c.Connection.Username()
	if err == nil {
		org = currentOrg.OrganizationFields.Name
	}
	currentSpace, err = c.Connection.GetCurrentSpace()

	// Check for error, then whether we actually find a targeted space
	if err != nil {
		c.UI.Say(terminal.FailureColor("FAILED"))
		c.UI.Say(err.Error())
	} else if len(currentSpace.Guid) < 1 {
		c.UI.Say(terminal.FailureColor("FAILED"))
		c.UI.Say("\nNo space targeted, use " + terminal.PromptColor("'cf target -s'") + "to target a space.\n")
		return
	}

	if allorg == true {
		RetrieveOrgApps(all, c, org, username, currentSpace)
	} else {
		RetrieveSpaceApps(all, c, org, username, currentSpace)
	}

}

//PrintApps will iterate through a list of apps and print in same format as 'cf apps' command.
//NOTE - this will be common for printing apps for a single space or for all spaces in an org.  Keep track of the
//change in "space" to print the header.
func PrintApps(appInfo []AppInfo, c *FilterApps, org string, username string, spaceName string) {
	// always print the header
	if len(appInfo) > 0 {
		var currentSpaceName = ""
		table := c.UI.Table([]string{"name", "requested state", "instances", "memory", "disk", "urls"})
		for _, app := range appInfo {
			if currentSpaceName != app.Space {
				currentSpaceName = app.Space //reset the var that holds the current space name
				c.UI.Say("Getting apps in org " + terminal.PromptColor(org) + " / space " + terminal.PromptColor(currentSpaceName) + " as " + terminal.PromptColor(username) + "...")
				c.UI.Say(terminal.SuccessColor("OK\n"))
			}
			if app.State == "started" {
				if strings.HasPrefix(app.Instances, "0/") {
					table.Add(app.Name, terminal.FailureColor(app.State), terminal.FailureColor(app.Instances), app.Memory, app.Disk, app.Urls)
				} else {
					table.Add(app.Name, app.State, app.Instances, app.Memory, app.Disk, app.Urls)
				}
			} else {
				table.Add(terminal.FailureColor(app.Name), terminal.FailureColor(app.State),
					terminal.FailureColor(app.Instances), terminal.FailureColor(app.Memory),
					terminal.FailureColor(app.Disk), terminal.FailureColor(app.Urls))
			}
		}
		table.Print()
		c.UI.Say("\n")
	} else {
		c.UI.Say("Getting apps in org " + terminal.PromptColor(org) + " / space " + terminal.PromptColor(spaceName) + " as " + terminal.PromptColor(username) + "...")
		c.UI.Say(terminal.SuccessColor("OK\n"))
		fmt.Printf("No %s apps found", terminal.PromptColor("started"))
		fmt.Println("")
	}

}

// CheckProperUsage function ensures the proper syntax for invoking this plugin
func CheckProperUsage(args []string) (invalidargs bool) {
	if len(args) > 3 {
		invalidargs = true
	}
	if len(args) == 2 && (args[1] != "-a" && args[1] != "-o") {
		invalidargs = true
	}
	if len(args) == 3 {
		if (args[1] == "-a" && args[2] == "-o") ||
			(args[1] == "-o" && args[2] == "-a") {
			invalidargs = false
		} else {
			invalidargs = true
		}
	}
	return invalidargs
}

//RetrieveSpaceApps handles a request for showing apps in only the current space
func RetrieveSpaceApps(all bool, c *FilterApps,
	org string, username string, currentSpace plugin_models.Space) {

	var err error
	var apps []plugin_models.GetAppsModel

	apps, err = c.Connection.GetApps()
	if err != nil {
		c.UI.Say(terminal.FailureColor("FAILED"))
		c.UI.Say(err.Error())
		return
	}
	var appInfo = BuildAppInfo(apps, all, currentSpace.Name)

	//now that we have our array, let's print...
	PrintApps(appInfo, c, org, username, currentSpace.Name)

}

// RetrieveOrgApps function iterates throught the spaces in an org,
//   then gets the apps in each space,  then call new print PrintOrgApps method to display.
func RetrieveOrgApps(all bool, c *FilterApps, org string, username string,
	currentSpace plugin_models.Space) {

	var err error
	var apps []plugin_models.GetAppsModel
	var spaceList []plugin_models.GetSpaces_Model
	var appInfo []AppInfo

	//Get the list of spaces in the org
	spaceList, err = c.Connection.GetSpaces()
	if err != nil {
		c.UI.Say(terminal.FailureColor("FAILED"))
		c.UI.Say(err.Error())
		return
	}
	//Iterate through spaces and print apps for each
	for _, space := range spaceList {
		c.Connection.CliCommandWithoutTerminalOutput("target", "-s", space.Name) //switch to the next space
		apps, err = c.Connection.GetApps()                                       //get the list of apps in this space
		appInfo = BuildAppInfo(apps, all, space.Name)
		PrintApps(appInfo, c, org, username, space.Name) // pass the space.Name since the appInfo array may be empty if there are no apps
		c.UI.Say("\n")                                   //print an extra empty line to make it easier to read
	}

	//Now switch back to the original space they were in

	c.Connection.CliCommandWithoutTerminalOutput("target", "-s", currentSpace.Name)

}

// BuildAppInfo - takes a []plugin_models.GetAppsModel and builds an AppInfo array
func BuildAppInfo(apps []plugin_models.GetAppsModel, all bool, spaceName string) (appInfo []AppInfo) {
	appInfo = make([]AppInfo, 0, 50)
	for _, app := range apps {
		if app.State == "started" || all == true {
			routes := app.Routes
			u := []string{}
			for _, r := range routes {
				u = append(u, r.Host+"."+r.Domain.Name)
			}
			inst := strconv.Itoa(app.RunningInstances) + "/" + strconv.Itoa(app.TotalInstances)
			a := AppInfo{Space: spaceName, Name: app.Name, Instances: inst, Memory: ConvertSize(app.Memory), Disk: ConvertSize(app.DiskQuota), Urls: strings.Join(u, ", "), State: app.State}
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
			Minor: 2,
			Build: 0,
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
						"o": "Show apps for entire org",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(FilterApps))
}
