Opinionated tool that helps with secrets / deployments.

## Installation

```console
foo@bar:~$ brew tap mantro/gops
foo@bar:~$ brew install mantro/gops/gops
```

## Development setup

```console
foo@bar:~$ brew install golang
foo@bar:~$ ./install.sh
```

## What does it do?

On first time invocation `gops` will create a `<GIT>/.gops.yaml` file with the default settings.

`gops` will also check `<GIT>/.gitignore` to see if `.gops.yaml` is on ignore (create/modify if needed).

## What repository layout is expected?

In the default settings `gops` expects a `ops/config` folder with yaml files, e.g.:

```bash
./ops/config/all/all.config.yaml
./ops/config/dev/dev.config.yaml
./ops/config/dev/dev.secrets.yaml
./ops/config/prod/prod.config.yaml
./ops/config/prod/prod.secrets.yaml
```

These files will be merged in-memory into one file by using directory names as keys and merging all files on the same level into it, e.g.:

```yaml
all:
  setting1: 1
  setting2: 2
dev:
  some_config: "config"
  some_secret: "secret"
prod:
  production: true
  a_secret: true
```

`gops` will try to find files names `*.template.*` and apply the Jet template engine to it to generate a new file named `*.generated.*`, e.g.:

```bash
scripts/connect-db.template.sh --> scripts/connect-db.generated.sh`
```

In order to see which variables you can use execute `gops dump`. Most of the data you will be interested in is in the `Data` key.

## Templates

The template engine being used is https://github.com/CloudyKit/jet.

The main use case looks like this:

```yaml
# file: configfile.template.yaml
ConnectionStrings:
  Default: "Host=#{{ .Data.dev.database.hostname }}#; Port=#{{ .Data.dev.database.port }}#"
```

```bash
# invoke
gops
```

```yaml
# generates: configfile.generated.yaml
ConnectionStrings:
  Default: "Host=your-database.azure.com; Port=5432"
```

## Usage

You should be somewhere within a git repository for it to work.

```bash
gops help           # Displays a command overview
gops config         # Shows the current gops configuration
gops dump           # Merges yaml config files and dumps it on stdout
gops target         # Target management, shows or sets current target

gops templates      # Generate the templates

gops                # Same as 'gops templates'

```
