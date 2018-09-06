# ghlabels

A tool to synchronize labels on GitHub repositories sanely.

It uses the same credentials as [hub](http://github.com/github/hub) to access github repositories. In fact, it uses the same library
as hub so if you have multiple credentials with hub and expect to be able to select between them,
ghlabels will offer the same prompts. If you do not have hub setup, it will prompt for the credentials
the same way hub does.

## Install

```
go get nhooyr.io/ghlabels
```

## Usage

ghlabels uses a JSON file that looks like this and represents the labels for a given repository:

```json
[
    {
        "name": "blocked",
        "description": "",
        "color": "032e70"
    },
    {
        "name": "good first issue",
        "description": "",
        "color": "e7dff5"
    },
    {
        "name": "p1",
        "description": "",
        "color": "df0000"
    },
    {
        "name": "p2",
        "description": "",
        "color": "ff7474"
    },
    {
        "name": "p3",
        "description": "",
        "color": "ffabab"
    },
    {
        "name": "p4",
        "description": "",
        "color": "ffe3e3"
    }
]
``` 

You may create and edit a file manually or have ghlabels create this file for you via

```
ghlabels pull <org>/<repo> > labels.json
```

The above command will create a `labels.json` file in the current directory for the labels on `<org>/<repo>`

You can push this to an entire organization's repos (or a single repo) via

```
ghlabels push <org>[/<repo>] < labels.json 
```

You can delete labels across entire organizations or repos via the `delete` subcommand.

```
ghlabels delete [--defaults] [<label>] <org>[<repo>]
```

The `--defaults` flag will delete all the default labels.

You can rename labels across entire organizations or repos via the `rename` subcommand.

```
ghlabels rename <owner>[/<repo>] <old_name> <new_name>
```
