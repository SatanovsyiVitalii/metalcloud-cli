package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

var networkProfileCmds = []Command{
	{
		Description:  "Lists all network profiles.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list network_profile", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"datacenter": c.FlagSet.String("datacenter", _nilDefaultStr, red("(Required)")+" Network profile datacenter"),
				"format":     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: networkProfileListCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Lists vlans of network profile.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "vlan-list",
		AltPredicate: "vlans",
		FlagSet:      flag.NewFlagSet("vlan-list network_profile", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", _nilDefaultInt, red("(Required)")+" Network profile's id."),
				"format":             c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: networkProfileVlansListCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Get network profile details.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("Get network profile details.", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", _nilDefaultInt, red("(Required)")+" Network profile's id."),
				"format":             c.FlagSet.String("format", _nilDefaultStr, "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
				"raw":                c.FlagSet.Bool("raw", false, green("(Flag)")+" If set returns the raw object serialized using specified format"),
			}
		},
		ExecuteFunc: networkProfileGetCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Create network profile.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("Create network profile", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"datacenter":            c.FlagSet.String("datacenter", _nilDefaultStr, red("(Required)")+" Label of the datacenter. Also used as an ID."),
				"format":                c.FlagSet.String("format", "json", "The input format. Supported values are 'json','yaml'. The default format is json."),
				"read_config_from_file": c.FlagSet.String("raw-config", _nilDefaultStr, red("(Required)")+" Read  configuration from file in the format specified with --format."),
				"read_config_from_pipe": c.FlagSet.Bool("pipe", false, green("(Flag)")+" If set, read  configuration from pipe instead of from a file. Either this flag or the --raw-config option must be used."),
				"return_id":             c.FlagSet.Bool("return-id", false, "Will print the ID of the created object. Useful for automating tasks."),
			}
		},
		ExecuteFunc: networkProfileCreateCmd,
		Endpoint:    DeveloperEndpoint,
		Example: `
#create file network-profile.yaml:
label: internet01
dc: us02-chi-qts01-dc
networkType: wan
vlans:
- vlanID: null
  portMode: native
  provisionSubnetGateways: false
  extConnectionIDs:
   - 10
  subnetPools: 
  - subnetPoolID: 13
	subnetPoolType: ipv4
- vlanID: 3205
  portMode: trunk
  provisionSubnetGateways: false
  extConnectionIDs: []

#create the actual profile from the file: 
metalcloud-cli network-profile create -datacenter us02-chi-qts01-dc -format yaml -raw-config ./network-profile.yaml

More details available https://docs.metalsoft.io/en/latest/guides/adding_a_network_profile.html
`,
	},
	{
		Description:  "Delete a network profile.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete network profile", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", _nilDefaultInt, red("(Required)")+" Network profile's id "),
				"autoconfirm":        c.FlagSet.Bool("autoconfirm", false, green("(Flag)")+" If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: networkProfileDeleteCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Add a network profile to an instance array.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "associate",
		AltPredicate: "assign",
		FlagSet:      flag.NewFlagSet("assign network profile to an instance array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"network_profile_id": c.FlagSet.Int("id", _nilDefaultInt, red("(Required)")+" Network profile's id"),
				"network_id":         c.FlagSet.Int("net", _nilDefaultInt, red("(Required)")+" Network's id"),
				"instance_array_id":  c.FlagSet.Int("ia", _nilDefaultInt, red("(Required)")+" Instance array's id"),
			}
		},
		ExecuteFunc: networkProfileAssociateToInstanceArrayCmd,
		Endpoint:    DeveloperEndpoint,
	},
	{
		Description:  "Remove network profile from an instance array.",
		Subject:      "network-profile",
		AltSubject:   "np",
		Predicate:    "remove",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("remove network profile of an instance array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"instance_array_id": c.FlagSet.String("ia", _nilDefaultStr, red("(Required)")+" Instance array's id"),
				"network_id":        c.FlagSet.String("net", _nilDefaultStr, red("(Required)")+" Network's id"),
			}
		},
		ExecuteFunc: networkProfileRemoveFromInstanceArrayCmd,
		Endpoint:    DeveloperEndpoint,
	},
}

func networkProfileListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	datacenter, ok := getStringParamOk(c.Arguments["datacenter"])
	if !ok {
		return "", fmt.Errorf("-datacenter is required")
	}

	npList, err := client.NetworkProfiles(datacenter)
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
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 30,
		},
		{
			FieldName: "NETWORK TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 30,
		},
		{
			FieldName: "VLANs",
			FieldType: tableformatter.TypeInterface,
			FieldSize: 30,
		},
		{
			FieldName: "CREATED",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "UPDATED",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
	}

	data := [][]interface{}{}
	for _, np := range *npList {
		vlans := ""

		for _, vlan := range np.NetworkProfileVLANs {
			if vlan.VlanID != nil {
				if vlans == "" {
					vlans = strconv.Itoa(*vlan.VlanID)

				} else {
					vlans = vlans + "," + strconv.Itoa(*vlan.VlanID)
				}
			}
		}

		data = append(data, []interface{}{
			np.NetworkProfileID,
			blue(np.NetworkProfileLabel),
			np.NetworkType,
			vlans,
			np.NetworkProfileCreatedTimestamp,
			np.NetworkProfileUpdatedTimestamp,
		})
	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	return table.RenderTable("Network Profiles", "", getStringParam(c.Arguments["format"]))
}

func networkProfileVlansListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	id, ok := getIntParamOk(c.Arguments["network_profile_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	retNP, err := client.NetworkProfileGet(id)
	if err != nil {
		return "", err
	}

	schemaConfiguration := []tableformatter.SchemaField{
		{
			FieldName: "VLAN",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "Port mode",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "External connections",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "Provision subnet gateways",
			FieldType: tableformatter.TypeBool,
			FieldSize: 6,
		},
	}

	dataConfiguration := [][]interface{}{}
	networkProfileVlans := retNP.NetworkProfileVLANs

	for _, vlan := range networkProfileVlans {

		externalConnectionIDs := vlan.ExternalConnectionIDs
		ecIds := ""
		for index, ecId := range externalConnectionIDs {

			retEC, err := client.ExternalConnectionGet(ecId)
			if err != nil {
				return "", err
			}

			if index == 0 {
				ecIds = retEC.ExternalConnectionLabel + " (#" + strconv.Itoa(ecId) + ")"
			} else {
				ecIds = ecIds + ", " + retEC.ExternalConnectionLabel + " (#" + strconv.Itoa(ecId) + ")"
			}
		}

		vlanid := "auto"
		if vlan.VlanID != nil {
			vlanid = strconv.Itoa(*vlan.VlanID)
		}

		dataConfiguration = append(dataConfiguration, []interface{}{
			vlanid,
			vlan.PortMode,
			ecIds,
			vlan.ProvisionSubnetGateways,
		})
	}

	tableConfiguration := tableformatter.Table{
		Data:   dataConfiguration,
		Schema: schemaConfiguration,
	}

	retConfigTable, err := tableConfiguration.RenderTableFoldable("", "", getStringParam(c.Arguments["format"]), 0)
	if err != nil {
		return "", err
	}

	return retConfigTable, err
}

func networkProfileGetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	id, ok := getIntParamOk(c.Arguments["network_profile_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	retNP, err := client.NetworkProfileGet(id)
	if err != nil {
		return "", err
	}

	schema := []tableformatter.SchemaField{
		{
			FieldName: "ID",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "LABEL",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "DATACENTER",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
		{
			FieldName: "DETAILS",
			FieldType: tableformatter.TypeString,
			FieldSize: 6,
		},
	}

	networkProfileVlans := retNP.NetworkProfileVLANs

	vlanListDescriptions := []string{}

	for _, vlan := range networkProfileVlans {

		externalConnectionIDs := vlan.ExternalConnectionIDs
		ecDescriptions := []string{}
		for _, ecId := range externalConnectionIDs {

			retEC, err := client.ExternalConnectionGet(ecId)
			if err != nil {
				return "", err
			}

			ecDescriptions = append(ecDescriptions, fmt.Sprintf("%s (#%d)", blue(retEC.ExternalConnectionLabel), ecId))
		}

		subnetPoolsDescriptions := []string{}
		subnetPools := vlan.SubnetPools

		for _, subnet := range subnetPools {
			if subnet.SubnetPoolID == nil { //if nil means that the subnet is automatically allocated
				subnetPoolsDescriptions = append(subnetPoolsDescriptions, blue(fmt.Sprintf("auto %s", subnet.SubnetPoolType)))
				continue
			}
			retSubnet, err := client.SubnetPoolGet(*subnet.SubnetPoolID)
			if err != nil {
				return "", err
			}

			subnetPoolsDescriptions = append(subnetPoolsDescriptions, fmt.Sprintf("%s/%s (#%d)", blue(retSubnet.SubnetPoolPrefixHumanReadable), blue(retSubnet.SubnetPoolPrefixSize), retSubnet.SubnetPoolID))
		}

		vlanid := "auto"
		if vlan.VlanID != nil {
			vlanid = strconv.Itoa(*vlan.VlanID)
		}

		gatewayIsProvisioned := ""
		if !vlan.ProvisionSubnetGateways {
			gatewayIsProvisioned = "no GW"
		}

		vlanDetails := fmt.Sprintf("VLAN ID: %s (%s) %s",
			yellow(vlanid),
			vlan.PortMode,
			red(gatewayIsProvisioned),
		)

		if len(ecDescriptions) > 0 {
			vlanDetails = fmt.Sprintf("%s EC:[%s]", vlanDetails, strings.Join(ecDescriptions, ","))
		}

		if len(subnetPoolsDescriptions) > 0 {
			vlanDetails = fmt.Sprintf("%s Subnets:[%s]", vlanDetails, strings.Join(subnetPoolsDescriptions, ","))
		}

		vlanListDescriptions = append(vlanListDescriptions, vlanDetails)
	}

	data := [][]interface{}{
		{
			"#" + strconv.Itoa(retNP.NetworkProfileID),
			blue(retNP.NetworkProfileLabel),
			retNP.DatacenterName,
			strings.Join(vlanListDescriptions, "\n"),
		},
	}

	var sb strings.Builder

	format := getStringParam(c.Arguments["format"])

	if getBoolParam(c.Arguments["raw"]) {
		ret, err := tableformatter.RenderRawObject(*retNP, format, "Server interfaces")
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	} else {

		table := tableformatter.Table{
			Data:   data,
			Schema: schema,
		}

		ret, err := table.RenderTable("", "", format)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	}

	return sb.String(), nil
}

func networkProfileCreateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	datacenter, ok := getStringParamOk(c.Arguments["datacenter"])
	if !ok {
		return "", fmt.Errorf("-datacenter is required")
	}

	readContentfromPipe := getBoolParam((c.Arguments["read_config_from_pipe"]))

	var err error
	content := []byte{}

	if readContentfromPipe {
		content, err = readInputFromPipe()
	} else {

		if configFilePath, ok := getStringParamOk(c.Arguments["read_config_from_file"]); ok {

			content, err = readInputFromFile(configFilePath)
		} else {
			return "", fmt.Errorf("-raw-config <path_to_json_file> or -pipe is required")
		}
	}

	if err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", fmt.Errorf("Content cannot be empty")
	}

	format := getStringParam(c.Arguments["format"])

	var npConf metalcloud.NetworkProfile
	switch format {
	case "json":
		err := json.Unmarshal(content, &npConf)
		if err != nil {
			return "", err
		}
	case "yaml":
		err := yaml.Unmarshal(content, &npConf)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("input format \"%s\" not supported", format)
	}

	ret, err := client.NetworkProfileCreate(datacenter, npConf)
	if err != nil {
		return "", err
	}

	if c.Arguments["return_id"] != nil && *c.Arguments["return_id"].(*bool) {
		return fmt.Sprintf("%d", ret.NetworkProfileID), nil
	}

	return "", err

}

func networkProfileDeleteCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	networkProfileId, ok := getIntParamOk(c.Arguments["network_profile_id"])
	if !ok {
		return "", fmt.Errorf("-id is required")
	}

	confirm := getBoolParam(c.Arguments["autoconfirm"])

	networkProfile, err := client.NetworkProfileGet(networkProfileId)
	if err != nil {
		return "", err
	}

	if !confirm {

		confirmationMessage := fmt.Sprintf("Deleting network profile %s (%d).  Are you sure? Type \"yes\" to continue:",
			networkProfile.NetworkProfileLabel, networkProfile.NetworkProfileID)

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

	err = client.NetworkProfileDelete(networkProfileId)

	return "", err
}

func networkProfileAssociateToInstanceArrayCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {
	id, ok := getIntParamOk(c.Arguments["network_profile_id"])
	if !ok {
		return "", fmt.Errorf("-id required")
	}

	net, ok := getIntParamOk(c.Arguments["network_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	ia, ok := getIntParamOk(c.Arguments["instance_array_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	_, err := client.InstanceArrayNetworkProfileSet(ia, net, id)
	if err != nil {
		return "", err
	}

	return "", nil
}

func networkProfileRemoveFromInstanceArrayCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	instance_array_id, ok := getStringParamOk(c.Arguments["instance_array_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	ia, err := strconv.Atoi(instance_array_id)
	if err != nil {
		return "", err
	}

	network_id, ok := getStringParamOk(c.Arguments["network_id"])
	if !ok {
		return "", fmt.Errorf("-net required")
	}

	net, err := strconv.Atoi(network_id)
	if err != nil {
		return "", err
	}

	return "", client.InstanceArrayNetworkProfileClear(ia, net)
}
