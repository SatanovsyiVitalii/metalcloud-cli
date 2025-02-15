# metalcloud-cli

![Build](https://github.com/metalsoft-io/metalcloud-cli/actions/workflows/build.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/bigstepinc/metalcloud-cli/badge.svg?branch=master)](https://coveralls.io/github/bigstepinc/metalcloud-cli?branch=master)

This tool allows the manipulation of all Bigstep Metal Cloud elements via the command line.

![metalcloud-cli](https://bigstep.com/assets/images/blog/2019/metalcloud-cli-animated.gif)


### Installation

To install on Mac OS X:
```
brew tap metalsoft-io/homebrew-repo
brew install metalcloud-cli
```

To install on any CentOS/Redhat Linux distribution:
```
$ sudo yum install https://github.com/metalsoft-io/metalcloud-cli/releases/download/v1.0.3/metalcloud-cli_1.0.3_linux_amd64.rpm
```

To install on any Debian/Ubuntu distributions:
```
curl -skL $(curl -s https://api.github.com/repos/metalsoft-io/metalcloud-cli/releases/latest | grep -i browser_download_url  | grep "$(dpkg --print-architecture)" | grep deb | head -n 1 | cut -d'"' -f4) -o metalcloud-cli.deb && sudo dpkg -i metalcloud-cli.deb
```

To install on Windows:
Binaries are available [here](https://github.com/metalsoft-io/metalcloud-cli/releases/latest):
```
https://github.com/metalsoft-io/metalcloud-cli/releases/latest
```


To install using `go get` (this should also work on Windows):
```bash
go get github.com/metalsoft-io/metalcloud-cli
```

### Getting the API key

In the Metalcloud's Infrastructure Editor go to the upper left corner and click on your email. Then go to **Settings** > **API & SDKs** > **API credentials**

Copy the api key. It should be of the form <number>:<letters>

Configure credentials as environment variables:
```bash
export METALCLOUD_API_KEY="<your key>"
export METALCLOUD_ENDPOINT="https://api.bigstep.com"
export METALCLOUD_USER_EMAIL="<your email>"
```

### Getting a list of supported commands

Use `metalcloud-cli help` for a list of supported commands.


### Getting started

To create an infrastructure:
```bash
metalcloud-cli infrastructure create --label test --datacenter test --return-id
```

```bash
metalcloud-cli infrastructure list
+-------+-----------------------------------------+-------------------------------+-----------+-----------+---------------------+---------------------+
| ID    | LABEL                                   | OWNER                         | REL.      | STATUS    | CREATED             | UPDATED             |
+-------+-----------------------------------------+-------------------------------+-----------+-----------+---------------------+---------------------+
| 12345 | complex-demo                            | d.d@sdd.com                   | OWNER     | active    | 2019-03-28T15:23:08Z| 2019-03-28T15:23:08Z|
+-------+-----------------------------------------+-------------------------------+-----------+-----------+---------------------+---------------------+
```

To create an instance array in that infrastructure, get the ID of the infrastructure from above (12345):

```bash
metalcloud-cli instance-array create --infra 12345 --label master --proc 1 --proc-core-count 8 --ram 16
```

To view the id of the previously created drive array:

```bash
metalcloud-cli instance-array list --infra 12345
+-------+---------------------+---------------------+-----------+
| ID    | LABEL               | STATUS              | INST_CNT  |
+-------+---------------------+---------------------+-----------+
| 54321 | master              | ordered             | 1         |
+-------+---------------------+---------------------+-----------+
Total: 1 Instance Arrays
```

To create a drive array and attach it to the previous instance array:

```bash
metalcloud-cli drive-array create --infra 12345 --label master-da --ia 54321
```

To view the current status of the infrastructure

```bash
metalcloud-cli infrastructure get --id 12345
Infrastructures I have access to (as test@test.com)
+-------+----------------+-------------------------------+-----------------------------------------------------------------------+-----------+
| ID    | OBJECT_TYPE    | LABEL                         | DETAILS                                                               | STATUS    |
+-------+----------------+-------------------------------+-----------------------------------------------------------------------+-----------+
| 36791 | InstanceArray  | master                        | 1 instances (16 RAM, 8 cores, 1 disks)                                | ordered   |
| 47398 | DriveArray     | master-da                     | 1 drives - 40.0 GB iscsi_ssd (volume_template:0) attached to: 36791   | ordered   |
+-------+----------------+-------------------------------+-----------------------------------------------------------------------+-----------+
Total: 2 elements
```


### Apply support

Apply creates or updates a resource from a file. The supported format is yaml.

```bash
metalcloud-cli apply -f resources.yaml
```

The type of the requested resource needs to be specified using the field *kind*.

```
cat resources.yaml

kind: InstanceArray
apiVersion: 1.0
label: my-instance-array

---

kind: Secret
apiVersion: 1.0
name: my-secret

```

The objects and their fields can be found in the [SDK documentation](https://godoc.org/github.com/metalsoft-io/metal-cloud-sdk-go). The fields will be in the format specified in the yaml tag. For example `SubnetPool` object has a field named `subnet_pool_prefix_human_readable` in JSON format. In the YAML file used as imput for this command, the field should be called `prefix`. 

### Condensed format

The CLI also provides a "condensed format" for most of it's commands:
* instance-array = ia
* drive-array = da
* infrastructure = infra
* list = ls
* delete = rm
...

This allows commands such as:
```bash
metalcloud-cli infra ls
```

### Using label instead of IDs

Most commands also take a label instead of an id as a parameter. For example:
```bash
metalcloud-cli infra show --id complex-demo
```


### Permissions

Some commands depend on various permissions. For instance you cannot access another user's infrastructure unless you are a delegate of it. 


### Admin commands

To enable admin commands set the following environment variable:
```bash
export METALCLOUD_ADMIN="true"
```

## Debugging information

To enable debugging information in the output set the following environment variable:
```bash
export METALCLOUD_LOGGING_ENABLED=true
```
### Building the CLI

The build process is automated by travis. Just push into the repository using the appropriate tag:

Use `git tag` to get the last tag:
```
git tag
v1.6.7
v1.6.8
v1.6.9
v1.7.0
v1.7.1
v1.7.2
v1.7.3
...
v1.7.4
v1.7.5
v1.7.6
v1.7.7
v1.7.8
```
Push new changes with new tag:
```
git add .
git commit -m "commit comment"
git tag v1.0.1
git push --tags
```

A coverage report is generated automatically at each build by [coverall](https://coveralls.io/github/metalsoft-io/metalcloud-cli?branch=master). There is a lower limit to the coverage currently set at 20%. 

It is a good idea to update the master branch as well (with no tag):
```
git push
```

### Updating the SDK

To update the SDK update `go.mod` file then regenerate the interfaces used for testing. 
Ifacemaker is needed
```
go get ifacemaker
```

```
go generate
```
If new objects are added in the SDK `helpers/fix_package.go` will need to be updated.
