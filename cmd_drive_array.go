package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	metalcloud "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
	"github.com/metalsoft-io/tableformatter"
)

//driveArrayCmds commands affecting instance arrays
var driveArrayCmds = []Command{

	{
		Description:  "Creates a drive array.",
		Subject:      "drive-array",
		AltSubject:   "da",
		Predicate:    "create",
		AltPredicate: "new",
		FlagSet:      flag.NewFlagSet("drive-array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label":                c.FlagSet.String("infra", _nilDefaultStr, red("(Required)") + " Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations."),
				"instance_array_id_or_label":                c.FlagSet.String("ia", _nilDefaultStr, red("(Required)") + " The id of the instance array it is attached to. It can be zero for unattached Drive Arrays"),
				"drive_array_label":                         c.FlagSet.String("label", _nilDefaultStr, red("(Required)") + " The label of the drive array"),
				"drive_array_storage_type":                  c.FlagSet.String("type", _nilDefaultStr, "Possible values: iscsi_ssd, iscsi_hdd"),
				"drive_size_mbytes_default":                 c.FlagSet.Int("size", _nilDefaultInt, "(Optional, default = 40960) Drive arrays's size in MBytes"),
				"drive_array_count":                         c.FlagSet.Int("count", _nilDefaultInt, "DriveArrays's drive count. Use this only for unconnected DriveArrays."),
				"drive_array_no_expand_with_instance_array": c.FlagSet.Bool("no-expand-with-ia", false, green("(Flag)") + " If set, auto-expand when the connected instance array expands is disabled"),
				"volume_template_id_or_label":               c.FlagSet.String("template", _nilDefaultStr, "DriveArrays's volume template to clone when creating Drives"),
				"return_id":                                 c.FlagSet.Bool("return-id", false, "(Optional) Will print the ID of the created Drive Array. Useful for automating tasks."),
			}
		},
		ExecuteFunc: driveArrayCreateCmd,
	},
	{
		Description:  "Edit a drive array.",
		Subject:      "drive-array",
		AltSubject:   "da",
		Predicate:    "edit",
		AltPredicate: "alter",
		FlagSet:      flag.NewFlagSet("edit_drive_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_array_id_or_label":                c.FlagSet.String("id", _nilDefaultStr, red("(Required)") + " Drive Array's ID or label. Note that using the label can be ambiguous and is slower."),
				"instance_array_id_or_label":             c.FlagSet.String("ia", _nilDefaultStr, red("(Required)") + " The id of the instance array it is attached to. It can be zero for unattached Drive Arrays"),
				"drive_array_label":                      c.FlagSet.String("label", _nilDefaultStr, red("(Required)") + " The label of the drive array"),
				"drive_array_storage_type":               c.FlagSet.String("type", _nilDefaultStr, "Possible values: iscsi_ssd, iscsi_hdd"),
				"drive_size_mbytes_default":              c.FlagSet.Int("size", _nilDefaultInt, "(Optional, default = 40960) Drive arrays's size in MBytes"),
				"drive_array_count":                      c.FlagSet.Int("count", _nilDefaultInt, "DriveArrays's drive count. Use this only for unconnected DriveArrays."),
				"drive_array_expand_with_instance_array": c.FlagSet.Bool("expand-with-ia", true, "Auto-expand when the connected instance array expands"),
				"volume_template_id_or_label":            c.FlagSet.String("template", _nilDefaultStr, "DriveArrays's volume template to clone when creating Drives"),
			}
		},
		ExecuteFunc: driveArrayEditCmd,
	},
	{
		Description:  "Lists all drive arrays of an infrastructure.",
		Subject:      "drive-array",
		AltSubject:   "da",
		Predicate:    "list",
		AltPredicate: "ls",
		FlagSet:      flag.NewFlagSet("list drive_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"infrastructure_id_or_label": c.FlagSet.String("infra", _nilDefaultStr, red("(Required)") + " Infrastructure's id or label. Note that the 'label' this be ambiguous in certain situations."),
				"format":                     c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: driveArrayListCmd,
	},
	{
		Description:  "Delete a drive array.",
		Subject:      "drive-array",
		AltSubject:   "da",
		Predicate:    "delete",
		AltPredicate: "rm",
		FlagSet:      flag.NewFlagSet("delete drive_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_array_id_or_label": c.FlagSet.String("id", _nilDefaultStr, red("(Required)") + " Drive Array's ID or label. Note that using the label can be ambiguous and is slower."),
				"autoconfirm":             c.FlagSet.Bool("autoconfirm", false, green("(Flag)") + " If set it will assume action is confirmed"),
			}
		},
		ExecuteFunc: driveArrayDeleteCmd,
	},
	{
		Description:  "Gets a drive array.",
		Subject:      "drive-array",
		AltSubject:   "da",
		Predicate:    "get",
		AltPredicate: "show",
		FlagSet:      flag.NewFlagSet("show drive_array", flag.ExitOnError),
		InitFunc: func(c *Command) {
			c.Arguments = map[string]interface{}{
				"drive_array_id_or_label": c.FlagSet.String("id", _nilDefaultStr, red("(Required)") + " Drive Array's ID or label. Note that using the label can be ambiguous and is slower."),
				"show_iscsi_credentials":  c.FlagSet.Bool("show-iscsi-credentials", false, green("(Flag)") + " If set returns the drives' iscsi credentials"),
				"format":                  c.FlagSet.String("format", "", "The output format. Supported values are 'json','csv','yaml'. The default format is human readable."),
			}
		},
		ExecuteFunc: driveArrayGetCmd,
	},
}

func driveArrayCreateCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	da := argsToDriveArray(c.Arguments)

	infra, err := getInfrastructureFromCommand("infra", c, client)
	if err != nil {
		return "", err
	}

	if v, ok := getStringParamOk(c.Arguments["instance_array_id_or_label"]); ok {

		iaID, err := getIDOrDo(v, func(label string) (int, error) {
			ia, err := client.InstanceArrayGetByLabel(label)
			if err != nil {
				return 0, err
			}
			return ia.InstanceArrayID, nil
		})

		if err != nil {
			return "", err
		}
		da.InstanceArrayID = iaID
	}

	if v, ok := getStringParamOk(c.Arguments["volume_template_id_or_label"]); ok {
		vtID, err := getIDOrDo(v, func(label string) (int, error) {
			vt, err := client.VolumeTemplateGetByLabel(label)
			if err != nil {
				return 0, err
			}
			return vt.VolumeTemplateID, nil
		},
		)
		if err != nil {
			return "", err
		}
		da.VolumeTemplateID = vtID
	}

	if da.DriveArrayLabel == "" {
		return "", fmt.Errorf("-label is required")
	}

	retDA, err := client.DriveArrayCreate(infra.InfrastructureID, *da)
	if err != nil {
		return "", err
	}

	if getBoolParam(c.Arguments["return_id"]) {
		return fmt.Sprintf("%d", retDA.DriveArrayID), nil
	}

	return "", err
}

func driveArrayEditCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	retDA, err := getDriveArrayFromCommand(c, client)
	if err != nil {
		return "", err
	}

	dao := retDA.DriveArrayOperation

	if v, ok := getStringParamOk(c.Arguments["instance_array_id_or_label"]); ok {
		iaID, err := getIDOrDo(v,
			func(label string) (int, error) {
				ia, err := client.InstanceArrayGetByLabel(label)
				if err != nil {
					return 0, err
				}
				return ia.InstanceArrayID, nil
			},
		)
		if err != nil {
			return "", err
		}

		dao.InstanceArrayID = iaID
	}

	if v, ok := getStringParamOk(c.Arguments["volume_template_id_or_label"]); ok {
		vtID, err := getIDOrDo(v,
			func(label string) (int, error) {
				vt, err := client.VolumeTemplateGetByLabel(label)
				if err != nil {
					return 0, err
				}
				return vt.VolumeTemplateID, nil
			},
		)
		if err != nil {
			return "", err
		}

		dao.VolumeTemplateID = vtID
	}

	updateIfIntParamSet(c.Arguments["drive_array_id"], &dao.DriveArrayID)
	updateIfStringParamSet(c.Arguments["drive_array_label"], &dao.DriveArrayLabel)
	updateIfStringParamSet(c.Arguments["drive_array_storage_type"], &dao.DriveArrayStorageType)
	updateIfIntParamSet(c.Arguments["drive_array_count"], &dao.DriveArrayCount)
	updateIfIntParamSet(c.Arguments["drive_size_mbytes_default"], &dao.DriveSizeMBytesDefault)
	updateIfBoolParamSet(c.Arguments["drive_array_expand_with_instance_array"], &dao.DriveArrayExpandWithInstanceArray)

	_, err = client.DriveArrayEdit(retDA.DriveArrayID, *dao)

	return "", err
}

func driveArrayListCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	infraIDStr, err := getParam(c, "infrastructure_id_or_label", "infra")
	if err != nil {
		return "", err
	}

	infraID, err := getIDOrDo(*infraIDStr.(*string), func(label string) (int, error) {
		ia, err := client.InfrastructureGetByLabel(label)
		if err != nil {
			return 0, err
		}
		return ia.InfrastructureID, nil
	},
	)
	if err != nil {
		return "", err
	}

	daList, err := client.DriveArrays(infraID)
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
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "SIZE (MB)",
			FieldType: tableformatter.TypeInt,
			FieldSize: 10,
		},
		{
			FieldName: "TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "ATTACHED TO",
			FieldType: tableformatter.TypeString,
			FieldSize: 30,
		},
		{
			FieldName: "DRV_CNT",
			FieldType: tableformatter.TypeInt,
			FieldSize: 10,
		},
		{
			FieldName: "TEMPLATE",
			FieldType: tableformatter.TypeString,
			FieldSize: 25,
		},
	}

	data := [][]interface{}{}
	for _, da := range *daList {
		status := da.DriveArrayServiceStatus

		if da.DriveArrayServiceStatus != "ordered" && da.DriveArrayOperation.DriveArrayDeployType == "edit" && da.DriveArrayOperation.DriveArrayDeployStatus == "not_started" {
			status = "edited"
		}

		if da.DriveArrayServiceStatus != "ordered" && da.DriveArrayOperation.DriveArrayDeployType == "delete" && da.DriveArrayOperation.DriveArrayDeployStatus == "not_started" {
			status = "marked for delete"
		}

		volumeTemplateName := ""
		if da.VolumeTemplateID != 0 {
			vt, err := client.VolumeTemplateGet(da.DriveArrayOperation.VolumeTemplateID)
			if err != nil {
				return "", err
			}

			volumeTemplateName = fmt.Sprintf("%s (#%d)", vt.VolumeTemplateDisplayName, vt.VolumeTemplateID)
		}

		instanceArrayLabel := ""
		if da.DriveArrayOperation.InstanceArrayID != nil && da.DriveArrayOperation.InstanceArrayID != 0 {
			var instanceArrayID int
			
			switch da.DriveArrayOperation.InstanceArrayID.(type) {
			case int:
				instanceArrayID = da.DriveArrayOperation.InstanceArrayID.(int)
			case float64:
				instanceArrayID = int(da.DriveArrayOperation.InstanceArrayID.(float64))
			default:
				return "", fmt.Errorf("Instance array ID type invalid.")
			}

			ia, err := client.InstanceArrayGet(instanceArrayID)
			if err != nil {
				return "", err
			}
			instanceArrayLabel = fmt.Sprintf("%s (#%d)", ia.InstanceArrayLabel, ia.InstanceArrayID)
		}

		data = append(data, []interface{}{
			da.DriveArrayID,
			da.DriveArrayOperation.DriveArrayLabel,
			status,
			da.DriveArrayOperation.DriveSizeMBytesDefault,
			da.DriveArrayOperation.DriveArrayStorageType,
			instanceArrayLabel,
			da.DriveArrayOperation.DriveArrayCount,
			volumeTemplateName})
	}

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)

	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}

	return table.RenderTable("Drive Arrays", "", getStringParam(c.Arguments["format"]))
}

func driveArrayDeleteCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	retDA, err := getDriveArrayFromCommand(c, client)
	if err != nil {
		return "", err
	}

	var retIA *metalcloud.InstanceArray

	if retDA.InstanceArrayID != 0 {
		retIA, err = client.InstanceArrayGet(retDA.InstanceArrayID)
		if err != nil {
			return "", err
		}
	}

	retInfra, err2 := client.InfrastructureGet(retDA.InfrastructureID)
	if err2 != nil {
		return "", err2
	}

	confirm, err := confirmCommand(c, func() string {

		var confirmationMessage string

		if retIA != nil {
			confirmationMessage = fmt.Sprintf("Deleting drive array %s (%d), attached to instance array (%s, %d) - from infrastructure %s (%d).  Are you sure? Type \"yes\" to continue:",
				retDA.DriveArrayLabel, retDA.DriveArrayID,
				retIA.InstanceArrayLabel, retIA.InstanceArrayID,
				retInfra.InfrastructureLabel, retInfra.InfrastructureID)
		} else {
			confirmationMessage = fmt.Sprintf("Deleting drive array %s (%d), unattached - from infrastructure %s (%d).  Are you sure? Type \"yes\" to continue:",
				retDA.DriveArrayLabel, retDA.DriveArrayID,
				retInfra.InfrastructureLabel, retInfra.InfrastructureID)
		}

		//this is simply so that we don't output a text on the command line
		if strings.HasSuffix(os.Args[0], ".test") {
			confirmationMessage = ""
		}

		return confirmationMessage
	})
	if err != nil {
		return "", err
	}

	if confirm {
		return "", client.DriveArrayDelete(retDA.DriveArrayID)
	}

	return "", fmt.Errorf("Operation not confirmed. Aborting")
}

func driveArrayGetCmd(c *Command, client metalcloud.MetalCloudClient) (string, error) {

	retDA, err := getDriveArrayFromCommand(c, client)
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
			FieldName: "STATUS",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "SIZE (MB)",
			FieldType: tableformatter.TypeInt,
			FieldSize: 10,
		},
		{
			FieldName: "TYPE",
			FieldType: tableformatter.TypeString,
			FieldSize: 10,
		},
		{
			FieldName: "ATTACHED TO",
			FieldType: tableformatter.TypeString,
			FieldSize: 30,
		},
		{
			FieldName: "TEMPLATE",
			FieldType: tableformatter.TypeString,
			FieldSize: 25,
		},
		{
			FieldName: "DETAILS",
			FieldType: tableformatter.TypeString,
			FieldSize: 25,
		},
	}

	drives, err := client.DriveArrayDrives(retDA.DriveArrayID)
	if err != nil {
		return "", err
	}

	data := [][]interface{}{}
	for _, d := range *drives {

		template := ""
		if d.TemplateIDOrigin != 0 {
			vt, err := client.VolumeTemplateGet(d.TemplateIDOrigin)
			if err != nil {
				return "", err
			}
			template = fmt.Sprintf("%s(#%d)", vt.VolumeTemplateDisplayName, vt.VolumeTemplateID)

		}

		details := ""

		if d.DriveOperatingSystem != nil {
			details = fmt.Sprintf("%s ", d.DriveOperatingSystem.OperatingSystemType)
		}

		if d.DriveFileSystem != nil {
			details = fmt.Sprintf("%s %s", details, d.DriveFileSystem.DriveFilesystemType)
		}

		dataRow := []interface{}{
			d.DriveID,
			d.DriveLabel,
			d.DriveServiceStatus,
			d.DriveSizeMBytes,
			d.DriveStorageType,
			fmt.Sprintf("instance-%d", d.InstanceID),
			template,
			details,
		}

		if getBoolParam(c.Arguments["show_iscsi_credentials"]) {

			credentials := fmt.Sprintf("Target: %s Port:%d IQN:%s LUN ID:%d",
				d.DriveCredentials.ISCSI.StorageIPAddress,
				d.DriveCredentials.ISCSI.StoragePort,
				d.DriveCredentials.ISCSI.TargetIQN,
				d.DriveCredentials.ISCSI.LunID)

			dataRow = append(dataRow, credentials)
		}

		data = append(data, dataRow)

	}

	if getBoolParam(c.Arguments["show_iscsi_credentials"]) {
		schema = append(schema, tableformatter.SchemaField{
			FieldName: "CREDENTIALS",
			FieldType: tableformatter.TypeString,
			FieldSize: 5,
		})
	}

	subtitle := fmt.Sprintf("Drive Array #%d", retDA.DriveArrayID)
	if retDA.InstanceArrayID != 0 {
		subtitle = fmt.Sprintf("%s attached to instance array %d", subtitle, retDA.InstanceArrayID)
	}
	subtitle = fmt.Sprintf("%s has the following drives:", subtitle)

	tableformatter.TableSorter(schema).OrderBy(schema[0].FieldName).Sort(data)
	table := tableformatter.Table{
		Data:   data,
		Schema: schema,
	}
	return table.RenderTable("Drives", subtitle, getStringParam(c.Arguments["format"]))
}

func argsToDriveArray(m map[string]interface{}) *metalcloud.DriveArray {
	return &metalcloud.DriveArray{
		DriveArrayID:                      getIntParam(m["drive_array_id"]),
		DriveArrayLabel:                   getStringParam(m["drive_array_label"]),
		DriveArrayStorageType:             getStringParam(m["drive_array_storage_type"]),
		DriveArrayCount:                   getIntParam(m["drive_array_count"]),
		DriveSizeMBytesDefault:            getIntParam(m["drive_size_mbytes_default"]),
		DriveArrayExpandWithInstanceArray: getBoolParam(m["drive_array_no_expand_with_instance_array"]),
	}
}

func getDriveArrayFromCommand(c *Command, client metalcloud.MetalCloudClient) (*metalcloud.DriveArray, error) {

	m, err := getParam(c, "drive_array_id_or_label", "id")
	if err != nil {
		return nil, err
	}

	id, label, isID := idOrLabel(m)

	if isID {
		return client.DriveArrayGet(id)
	}

	return client.DriveArrayGetByLabel(label)
}
