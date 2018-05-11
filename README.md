# Hermes
**Messenger of the Version Control Gods**

**Hermes allows you to manage your git repositories from the command line**

**NOTE:** This project is a blatant ripoff of this project, [bitcar](https://github.com/carsdotcom/bitcar), created by a friend and former colleague. I decided to create hermes because I wasn't a fan of the node dependency with bitcar and wanted additional features.

**OTHER NOTE:** Currently the project is still in the very early stages of development, so, although some features work, there's still plenty more to come and there is no guarantee of backwards compatibility until the project reaches the v1.0.0 mark. 

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

### Example .hermes.yml File

```yaml
repo_path: /Users/jeremychambers/test-repos/
config_path: /Users/jeremychambers/.hermes-config/
cache_file: cached-repos.json
target_file: .hermes_target_dir
```