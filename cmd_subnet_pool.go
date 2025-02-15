package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

var subnetPoolCmds = []Command{

	{
		Description:  "Lists subnets",
		Subject:      "subnet-pool",
		AltSubject:   "subnet",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list subnet pools", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":     c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"filter":     c.FlagSet.String("filter", "*", "Filter to restrict the results. Defaults to '*'"),
				"datacenter": c.FlagSet.String("datacenter", _nilDefaultStr, "Quick filter to restrict the results to show only the subnets of a datacenter."),
			}
		},
		ExecuteFunc: subnetPoolListCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Get a subnet pool.",
		Subject:      "subnet-pool",
		AltSubject:   "subnet",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("Get a subnet pool", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"subnet_pool_id": c.FlagSet.Int("id", _nilDefaultInt, red("(Required)") + " Subnetpool's id"),
				"format":         c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"raw":            c.FlagSet.Bool("raw", false, green("(Flag)") + " When set the return will be a full dump of the object. This is useful when copying configurations. Only works with json and yaml formats."),
			}
		},
		ExecuteFunc: subnetPoolGetCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Create a subnet pool.",
		Subject:      "subnet-pool",
		AltSubject:   "subnet",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create subnet pool", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("config", _nilDefaultStr, red("(Required)") + " Read configuration from file"),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, green("(Flag)") + " If set, read configuration from pipe instead of from a file. Either this flag or the -config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created Useful for automating tasks."),
			}
		},
		ExecuteFunc: subnetPoolCreateCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Delete a subnet pool.",
		Subject:      "subnet-pool",
		AltSubject:   "subnet",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete subnet pool", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"subnet_pool_id": c.FlagSet.Int("id", _nilDefaultInt, red("(Required)") + " Subnet's's id"),
				"autoconfirm":    c.FlagSet.Bool("autoconfirm", false, green("(Flag)") + " If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: subnetPoolDeleteCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func subnetPoolListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	filter := getStringParam(c.Arguments["filter"])
	if datacenter, ok := getStringParamOk(c.Arguments["datacenter"]); ok {
		filter = fmt.Sprintf("datacenter_name: %s %s", datacenter, filter)
	}

	list, err := client.SubnetPoolSearch(filter)

	if err != nil {
		return "", err
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "DATACENTER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "DEST.",
			FieldType: tableformatter.TypeString,
			FieldSize: 3,
		},

		{
			FieldName: "PREFIX",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "NETWORK_EQUIPMENT",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "USER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "MANUAL_ONLY",
			FieldType: tableformatter.TypeBool,
			FieldSize: 3,
		},
		{
			FieldName: "AVAILABLE_IPS",
			FieldType: tableformatter.TypeString,
			FieldSize: 3,
		},
	}

	data := [][]interface{}{}
	for _, s := range *list {

		prefixStr := fmt.Sprintf("%s/%d", s.SubnetPoolPrefixHumanReadable, s.SubnetPoolPrefixSize)

		userEmail := ""
		if s.UserID != 0 {
			u, err := client.UserGet(s.UserID)
			if err != nil {
				return "", err
			}
			userEmail = u.UserEmail
		}

		utilization, err := client.SubnetPoolPrefixSizesStats(s.SubnetPoolID)

		if err != nil {
			return "", err
		}

		utilizationStr := fmt.Sprintf("%s (%s", utilization.IPAddressesUsableCountFree, utilization.IPAddressesUsableFreePercentOptimistic)

		networkEquipmentIdentifier := ""
		if s.NetworkEquipmentID != 0 {
			sw, err := client.SwitchDeviceGet(s.NetworkEquipmentID, false)
			if err != nil {
				return "", err
			}

			networkEquipmentIdentifier = sw.NetworkEquipmentIdentifierString
		}

		data = append(data, []interface{}{

			s.SubnetPoolID,
			s.DatacenterName,
			s.SubnetPoolDestination,
			prefixStr,
			networkEquipmentIdentifier,
			userEmail,
			s.SubnetPoolIsOnlyForManualAllocation,
			utilizationStr + "%%)",
		})

	}

	tableformatter.TableSorter(schema).OrderBy(
		schema[0].FieldName,
		schema[1].FieldName,
		schema[2].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Subnet pools", "", getStringParam(c.Arguments["format"]))
}

func subnetPoolGetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	id, ok := getIntParamOk(c.Arguments["subnet_pool_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	s, err := client.SubnetPoolGet(id)
	if err != nil {
		return "", err
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeInt,
			FieldSize: 6,
		},
		{
			FieldName: "DATACENTER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "DEST.",
			FieldType: tableformatter.TypeString,
			FieldSize: 3,
		},

		{
			FieldName: "PREFIX",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "NETWORK_EQUIPMENT",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "USER",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "MANUAL_ONLY",
			FieldType: tableformatter.TypeBool,
			FieldSize: 3,
		},
		{
			FieldName: "AVAILABLE_IPS",
			FieldType: tableformatter.TypeString,
			FieldSize: 3,
		},
	}

	prefixStr := fmt.Sprintf("%s/%d", s.SubnetPoolPrefixHumanReadable, s.SubnetPoolPrefixSize)

	userEmail := ""
	if s.UserID != 0 {
		u, err := client.UserGet(s.UserID)
		if err != nil {
			return "", err
		}
		userEmail = u.UserEmail
	}

	utilization, err := client.SubnetPoolPrefixSizesStats(s.SubnetPoolID)

	if err != nil {
		return "", err
	}

	utilizationStr := fmt.Sprintf("%s (%s", utilization.IPAddressesUsableCountFree, utilization.IPAddressesUsableFreePercentOptimistic)

	networkEquipmentIdentifier := ""
	if s.NetworkEquipmentID != 0 {
		sw, err := client.SwitchDeviceGet(s.NetworkEquipmentID, false)
		if err != nil {
			return "", err
		}

		networkEquipmentIdentifier = sw.NetworkEquipmentIdentifierString
	}

	data := [][]interface{}{{

		s.SubnetPoolID,
		s.DatacenterName,
		s.SubnetPoolDestination,
		prefixStr,
		networkEquipmentIdentifier,
		userEmail,
		s.SubnetPoolIsOnlyForManualAllocation,
		utilizationStr + "%)",
	}}

	var sb strings.Builder

	format := getStringParam(c.Arguments["format"])

	if getBoolParam(c.Arguments["raw"]) {
		ret, err := tableformatter.RenderRawObject(*s, format, "SubnetPool")
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	} else {
		table := tableformatter.Table{
			Data:   data,
			Schema: schema,
		}
		ret, err := table.RenderTransposedTable("subnet pool", "", format)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	}

	return sb.String(), nil
}

func subnetPoolCreateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	var sn metalcloud.SubnetPool

	err := getRawObjectFromCommand(c, &sn)
	if err != nil {
		return "", err
	}

	ret, err := client.SubnetPoolCreate(sn)
	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", ret.SubnetPoolID), nil
	}

	return "", err
}

func subnetPoolDeleteCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	id, ok := getIntParamOk(c.Arguments["subnet_pool_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}
	confirm := false

	obj, err := client.SubnetPoolGet(id)
	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["autoconfirm"]) {
		confirm = true
	} else {

		confirmationMessage := fmt.Sprintf("Deleting subnet %s/%d (%d).  Are you sure? Type \"yes\" to continue:",
			obj.SubnetPoolPrefixHumanReadable,
			obj.SubnetPoolPrefixSize,
			obj.SubnetPoolID)

		//this is simply so that we don't output a text on the command line under go test
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		confirm, err = requestConfirmation(confirmationMessage)
		if err != nil {
			return "", err
		}

	}

	if !confirm {
		return "", fmt.Errorf("Operation not confirmed. Aborting")
	}

	err = client.SubnetPoolDelete(obj.SubnetPoolID)

	return "", err
}
