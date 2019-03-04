# krontab

[![Build Status](https://travis-ci.com/jacobtomlinson/krontab.svg?branch=master)](https://travis-ci.com/jacobtomlinson/krontab)
[![Current Release](https://img.shields.io/github/release/jacobtomlinson/krontab.svg)](https://github.com/jacobtomlinson/krontab/releases/latest)

A crontab replacement for Kubernetes.

Create `CronJob` resources on your Kubernetes cluster in the same way you would on your *nix system.
Krontab works by constructing a virtual crontab file from your CronJob resources and communicating changes back to the Kubernetes API. You can create more complex and customised jobs with custom templates and trigger your jobs manually any time from the command line.

Example crontab:

```shell
$ krontab -l
# template: default
0 1 * * * echo hello  # name: test
```

## Installation

### Quick install

```shell
curl -sL https://git.io/krontab | bash
```

### Linux (x86 64-bit)
```shell
LATEST_VERSION=$(curl -sL -o /dev/null -w %{url_effective} "https://github.com/jacobtomlinson/krontab/releases/latest" | rev | cut -f1 -d'/'| rev)
curl -L https://github.com/jacobtomlinson/krontab/releases/download/${LATEST_VERSION}/krontab-linux-x86_64 -o /usr/local/bin/krontab
chmod +x /usr/local/bin/krontab
```

### Linux (arm 32-bit)
```shell
LATEST_VERSION=$(curl -sL -o /dev/null -w %{url_effective} "https://github.com/jacobtomlinson/krontab/releases/latest" | rev | cut -f1 -d'/'| rev)
curl -L https://github.com/jacobtomlinson/krontab/releases/download/${LATEST_VERSION}/krontab-linux-arm -o /usr/local/bin/krontab
chmod +x /usr/local/bin/krontab
```

### OS X (64-bit)
```shell
LATEST_VERSION=$(curl -sL -o /dev/null -w %{url_effective} "https://github.com/jacobtomlinson/krontab/releases/latest" | rev | cut -f1 -d'/'| rev)
curl -L https://github.com/jacobtomlinson/krontab/releases/download/${LATEST_VERSION}/krontab-darwin-x86_64 -o /usr/local/bin/krontab
chmod +x /usr/local/bin/krontab
```

## Configuration

Authentication with your Kubernetes cluster will use your `~/.kube/config` credentials or a service account if being used from inside a pod on the cluster. Advanced configuration options are specified with environment variables.

| Env var  | Default | Description |
| ------------- | ------------- | ------------- |
| `KRONTAB_NAMESPACE` | `default` | The kubernetes namespace to use. |
| `KRONTAB_OWNER` | N/A | The owner of a cronjob. This will be used when creating jobs and will be used as a filter when listing them.  |

## Usage

```console
$ krontab help
Krontab is a crontab replacement for kubernetes.

You can use it to create cron jobs on a kubernetes cluster in a familiar crontab format.
Krontab works by allowing you to create job templates which are used in kubernetes. Then create
specific cron jobs using the crontab. Example krontab:

# Crontab example

# template: default
0 1 * * * echo hello  # name: test

Usage:
  krontab [flags]
  krontab [command]

Available Commands:
  create      Create a krontab resource
  delete      Delete a krontab resource
  edit        Edit a krontab resource
  get         Get krontab resources
  help        Help about any command
  list        List krontab resources
  run         Run a krontab job
  version     Print the version number of krontab

Flags:
  -e, --edit-crontab   Edit the crontab
  -h, --help           help for krontab
  -l, --list-crontab   List the crontab

Use "krontab [command] --help" for more information about a command.
```

### Templates

Krontab uses templates to turn your crontab into Kubernetes CronJob resources.
You will get a default template which you can view with `krontab get template default`.
You can create, edit and delete your own templates too.

```
$ krontab get template default
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: test
spec:
  schedule: 0 */6 * * *
  jobTemplate:
    metadata:
      name: test
      krontabTemplate: default
    spec:
      template:
        metadata:
        spec:
          restartPolicy: Never
          containers:
            - name: test
              image: busybox
              command: ["echo", "hello"]
              resources:
                limits:
                  cpu: "1"
                  memory: 2G
                requests:
                  cpu: "0.25"
                  memory: 0.5G
```

### Crontab

Krontab will automatically generate a crontab from your existing Kubernetes CronJob resources. You can edit
them with `krontab -e` or `krontab edit crontab`. Adding new lines to your crontab will be mapped into new
CrobJob resources. Editing lines will update them and deleteing lines will remove them.

```console
$ krontab -e
# Welcome to krontab, a crontab like editor for Kubernetes cron jobs.
#
# This is a virtual file and was generated from the kubernetes API. Any edits you make here will be
# sent back to the kubernetes API to create/update/delete CronJob resources. Next time you open this crontab
# you may notice different formatting and comments to how you save it. Comments are meaningful and contain
# metadata about the job.
#
# Example job
# -----------
#
# # template: default
# 0 1 * * * echo hello world  # name: hello-world
#
# Templates
# ---------
#
# Krontab uses templates to turn your your crontab into kubernetes compliant CronJob resources. You will find
# a default template by running 'krontab get template default'. This is a minimal CrobJob resource with runs your
# command in a ubuntu container. When krontab creates your jobs it will replace the container command and schedule pattern
# in the template with the command and schedule pattern from the crontab.
#
# All cronjobs following a '# template: <name>' comment will be created using that template. If no template comments exist
# then the default template will be used.
#
# Names
# -----
# Kubernetes requires each job to have a name and is specified at the end of a job with a '# name: <name>' comment.
# If you do not specify a name one will be autogenerated.


# template: default
0 1 * * * echo hello world  # name: test
```

## Examples

### Create a new cron job

To create a new cron job file you need to open your krontab file for editing. This will automatically open in `vim` or
whatever text edit you have configured with `VISUAL` or `EDITOR`.

```
krontab -e
```

Once this file is open you can create a new job in the same way you would when using regular `crontab`. You must specify your timings as
a five column cron time (minute, hour, day of month, month of year, day of week) followed by your shell command. This command will be triggered
every time the current time matches the schedule you have specified.

Let's create a job which echo's "hello` every hour.

```
0 * * * * echo hello
```

When creating a `CronJob` resource on the kubernetes cluster `krontab` will need to give it a name and choose a template to use. If you omit this information
like we have in our echo hello example above `krontab` will use the default template (`krontab get template default` if you want to see it) and random uuid will
be created for the name.

You can optionally choose to specify this information as yaml style comments in your crontab.

```
# template: default
0 * * * * echo hello  # name: hello-world
```

Here `template` has been specified at the document level which will affect all jobs below it until it is redefined. But the `name` info has been specified in line with the job and will only affect that job.

### Edit a job

Once you have some jobs when you run `crontab -e` you will see a crontab file which contains your existing jobs. They may not look exactly like when you created them as this
crontab file is dynamically generated on the fly. You can make changes to any of the information in the crontab and they will be reflected in ths `CronJob` resources on the kubernetes cluster.

_Note: All changes except renaming will perform an update task on the cluster. Renaming will result in a deletion and creation of a new job._

### Delete a job

You can delete a job by running `crontab -e` and removing the line you wish to delete. This will result in the `CronJob` resource being deleted from the cluster.

### Test a job

You can manually take a job for a test drive any time, even if it is not going to be triggered any time soon.

```
krontab run job <job name>
```

**Protip:** You can run one shot jobs without having to create a cron job first. This is useful for submitting one off batch jobs.

```
krontab run job --template <template name> --command <the command>
```

### Get running jobs

You can list any currently running jobs.

```
krontab list running
```

### View a template

You can print out any template.

```
krontab get template <template name>
```

### Edit a template

You can edit templates.

```
krontab edit template <template name>
```
_Note: You cannot edit the `default` template. Also editing a template will not result in jobs being recreated with the modified template._

### Deleting a template

You can delete templates.

```
krontab delete template <template name>
```

_Note: You cannot delete the template._

### Creating new templates

You can create new templates. These must be valid Kubernetes [`CronJob`](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) resources. The `schedule`, `command` and `name` sections will be overridden.

```
krontab create template <template name>
```

## Contributing

### Environment

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

```console
$ make
$ ./bin/krontab
```

#### Testing

```console
$ make test
```

### Releasing

Decide the new release version. Check out the [current releases](https://github.com/jacobtomlinson/krontab/releases) and follow SemVer to work out what it will be. Then run `make release` which will ask for the tag in the form `v{major}.{minor}.{patch}`. This will create the tag and push it to GitHub, Travis CI will then build the binaries and push them to github along with some autogenerated release notes.

```console
$ make release
```
