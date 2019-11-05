# Hermes - Messenger of the Version Control Gods

**Hermes allows you to manage your git repositories from the command line**

**NOTE:** This project is a blatant ripoff of this project, [bitcar](https://github.com/carsdotcom/bitcar), created by a friend and former colleague. I decided to create hermes because I wasn't a fan of the node dependency with bitcar and wanted additional features.

**OTHER NOTE:** Currently the project is still in the very early stages of development, so, although some features work, there's still plenty more to come and there is no guarantee of backwards compatibility until the project reaches the v1.0.0 mark. 

### Table of Contents
- [Installation](#installation)
- [Setup](#setup)
- [Configuration](#configuration)
    - [Example File](#example-config)
- [Usage](#usage)
    - [Root/Get Command](#root-get-command)
    - [Alias Command](#alias-command)
    - [Repository Commands](#repository-commands)
        - [Repository Rm Command](#repository-rm-command)
    - [Remote Commands](#remote-commands)
        - [Remote Add Command](#remote-add-command)
        - [Remote Refresh Command](#remote-refresh-command)
    - [Setup Command](#setup-command)
    - [Version Command](#version-command)
- [Contributing Guidlines](./CONTRIBUTING.md)

----

## Installation

As the project is still in beta stages of development, it has not been distributed, and is only available to install from source. To install hermes, you must then have go setup with your GOPATH bin in your shell PATH. Once you have this you can install by running

    go get -u github.com/TheHipbot/hermes
    hermes -h

## Setup

Once you have hermes installed, you must setup the hermes config directory and add an alias to you shell profile. 

To create the config directory, run the following:

    hermes setup

And you must add the hermes alias to your shell profile (e.g. $HOME/.bash_profile, .$HOME/.profile, etc):

    cat >> ~/.bash_profile
    if which hermes > /dev/null; then
        eval "$(hermes alias)"
    fi

The above command will expect you to enter in the shell script text then end the command by hitting <CTRL+D>.

## Configuration

Hermes supports yaml configuration files stored in either of these two locations on the system:

    $HOME/.hermes.yml
    /etc/hermes/.hermes.yml

Here is a list of the current supported config keys and values along with their use:

* `repo_path` (default: `$HOME/hermes-repos/`) - tells hermes where to clone repos to on your system. From this base path, repos will be stored similar to the `go get` tool. For example hermes will store itself in `${repo_path}/github.com/TheHipbot/hermes`
* `config_path` (default: `$HOME/.hermes/`) - the directory where hermes will store configuration files such as its internal cache and the hermes target file. **NOTE:** you will want to set this in your hermes configuration file **BEFORE** you run `hermes setup` since that command will create the config folder. 
* `target_file` (default: `.hermes_target`) -  after running the hermes command, if there is a valid target (e.g. repo that you have cloned or want to jump to), hermes writes out the full path into the hermes target file. From there, the alias in your shell profile will read this, jump to the directory and then remove the target file. **NOTE:** `target_file` only specifies the file name, the file will be created in the `config_path`. if you change the target file after setting the hermes alias, you would need to open a new terminal session or re-source your profile file for the alias to realize the change
* `cache_file` (default: `cache.json`) - hermes stores a cache of repos it is aware of to allow for tab completion and prompts. this will be in json format. **NOTE:** `cache_file` only specifies the file name, the file will be created in the `config_path`
* `alias_name` (default: `hermes`) - the name of the alias function which calls through to the hermes binary. this will be the command you run when using hermes.
* `credentials_type` (default: `none`) - the type of storage which hermes user to store user provided credentials, supported types described below
    * `none` - this will not store the credentials at all, any time a call is made that requires authentication credentials must be passed into hermes
    * `file` - this is the default type and will store provided credentials in yaml file in plaintext *NOTE: this is by no means a secure solution and its recommended not to use this in conjunction with usernames and passwords*
* `credentials_file` (default: `credentials.yml`) - when using the `file` credential type, this is the filename in the config directory in which the credentials will be stored

<a name="example-config"></a>
### Example .hermes.yml File

```yaml
repo_path: /Users/jeremychambers/test-repos/
config_path: /Users/jeremychambers/.hermes-config/
cache_file: cached-repos.json
target_file: .hermes_target_file
alias_name: hit
credentials_type: file
credentials_file: my_credentials.yml
```

## Usage

<a name="root-get-command"></a>
### Root / Get Command

`hermes [REPO]` or `hermes get [REPO]`

Running the hermes command without any subcommands or with the `get` subcommand are synonymous. The command is used to jump to a repo in the hermes cache or pull down a repo then jump to its new location.

args - the expected arguments are a full repo path in this format: [remote address]/[project or user]/[repo nam] (e.g. github.com/TheHipbot/hermes) to clone a new repo, or a string to conduct a contains search on the repos in the hermes cache

When run, the following actions happen in order:

1. The hermes cache file is read in and then all repos are searched to see if they contain the text given as `args`
2. Based on the results of the search, 1 of 3 things will happen
    * If the search turns up a single result from the cache, hermes will set the target to the path of that repo and exit so the alias can move you to the directory
    * If the search turns up no results, hermes assumes this is a new repo and will attempt to clone it. If the clone is successful, the repo is added to the cache and the target is set to the new repo
    * If there are multiple results, the user is prompted to select a repo from the results. Once a repo is selected, hermes will continue with that repo.
3. Assuming the command has executed successfully a target path should be written to the target file. Hermes will exit 0 and the alias (assuming it has been setup) will read the path from the file, move the current working directory to that target directory, remove the target file and exit.

### Alias Command

`hermes alias`

This command is meant only to provide the alias for a terminal session so it should be added to a shell profile, but not used otherwise. It writes to stdout a bash function which runs the hermes binary with the given args, then if a target file was written, it read the content as a directory to cd into. This is necessary because it is the only way which hermes can move the shell session's current working directory.

### Remote Commands

`hermes remote [SUBCOMMAND] [FLAGS] [ARGS]`

This group of commands are used for managing remote git servers which hermes should track against. Remotes that have been added will be queried for repoistories to which the user has access. This repositories will be added to the cache file along with information on the remote itself.

#### Flags (Global to all remote commands)

**-a, --all**

When adding or refreshing a remote, this will index all repos available to the user (as opposed to starred or user owned repos which is the default depending on remote) if that option is available for the remote.

**-p, --protocol**

When adding a remote, you can add the protocol flag followed by a valid protocol (`http`, `https`, or `ssh`) inline with the remote command. If no protocol is present, the user will be prompted for a protocol. 

#### Remote Add Command

`hermes remote add [FLAGS] [REMOTE URL]`

This command is used to add new git remotes to your hermes cache to use for searching and managing repositories. Once you run the command with a valid remote url, you will be prompted for which type of remote it is, your credentials to that remote, and which repository access protocol you would prefer. Adding a remote that already exists will act as a refresh to the remote, in this case any new values provided for protocol will be update in the cache and 

Currently supported types of remotes:

- GitHub
- Gitlab

Currently supported protocols:

- http(s)
- ssh

##### Flags

**--token**

When adding a remote, you can add the token flag to provide an auth token inline to the command, which will then replace any existing entry for the token in the credentials storer. If authentication is not successful, user will be prompted for a new token.

**-t, --type**

When adding a remote, you can specify the type of remote inline to the command, valid options are `github` and `gitlab`.

#### Remote Refresh Command

`hermes remote refresh [FLAGS]`

The refresh command will attempt to "refresh" all currently tracked remotes by re-adding all repos of the remote. This will not affect existing repositories, but will add any new ones that have been created on the remote since the remote was first added or last refreshed. 

### Repository Commands

`hermes repo [SUBCOMMAND] [FLAGS] [ARGS]`

aliases: `repository`

This group of commands are used to manage repositories which hermes should track. Repositories will typically be added wholesale by remote, but can be added (soon) and removed individually from the cache and optionally from disk.

#### Repository Remove Command

`hermes repo rm [FLAGS] [REPOSITORY NAME]`

aliases: `remove`

The remove command will remove repositories from the hermes cache and optionally from the filesystem as well. Removing from the cache will keep the repository folder in the `repo_dir` but will remove it from the list of repos in the cache. If the repository still exists in the remote, refreshing the remote will add it back. 

##### Flags

**--hard**

When removing a repository from the repository, adding the hard remove flag will also delete the directory for the repo that was cloned into the `repo_dir` if it exists. It will also recursively remove any empty parent directories up to the `repo_dir`.

### Setup Command

`hermes setup`

Running hermes setup creates the `config_path` directory if specified in the .hermes.yml file or `$HOME/.hermes` by default. This ensures the directory is available for subsequent commands and should only be run once.

### Version Command

`hermes version`

This command will output version information for the hermes binary you are executing.