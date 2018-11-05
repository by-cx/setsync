# SetSync tool

This tool solves problem with synchronization config files between multiple computers. It saves the config files into GitHub Gist and allows to download them back whenever you want. I use this tools to sync following files:

* .ssh/config
* .vimrc
* .config/Slack/storage/slack-teams
* .config/fish/config.fish
* .config/fish/functions/elxvpn.fish
* .config/terminator/config
* .gitconfig

So it covers my SSH config, settings of vim, Slack teams, script to connect our VPN, configuration of terminator or gitconfig.

This tool is not even close to modern synchronization features like Google Chrome or OS X has. But in combination with other tools it can help you to sync some files and make easier to set your computers up. I use this tool in combination with Ansible. I have my desktop defined in Ansible playbook so if I get a new notebook or desktop or I am upgrading my Elementary OS installations setsync and Ansible are huge help when the installation is done.

## Build

To build setsync you need:

* [Go 1.11](https://golang.org/) (lower version will be fine probably)
* [dep](https://github.com/golang/dep)

Then call:

    make build

At the end the binary file called *setsync* appears in your PWD.

## Usage

When you call *setsync help*, you will see this:

    Usage:
        setsync [FLAG|SET]...

    Flags:
        -h, --help    show help
        -v, --verbose show verbose help

    Sets:
        upload       uploads config files into the gist
        download     downloads config files from the gist
        add          adds file into the gist
        remove       removes file from the gist
        list         list files saved in the gist
        help         show help

Before you start *setsync* will ask you about GistID, your GitHub Token and the encryption key. The encryption key needs to be 32 characters long. If you do something wrong, you can always update the configuration in *.setsync* file. To get a gist ID just create a new gist and use ID from the URL. GitHub Token menu is your [GitHub profile settings](https://github.com/settings/tokens).

The upload and download has to be done manually so have this in mind when you update any "tracked" file in your local computer. To add new files use *add* commands. To pull configuration to the new computer or update the existing one use *download*. And if you change something call *upload*. Usually you have to turn off and on all affected application to load the new settings.
