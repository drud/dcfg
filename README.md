# drudconfig
A cli for configuration management.

## Use

Get the binary. Either build the source or head to the releases section.

### Create a Task Set (drud.yaml).

Each of the outer most list items in the drud.yaml file are Task Sets. The file can container one Task Set or multiple.
Each Task Set can be run on its own or with other Tas Sets in the file.

```
dcfg run all
dcfg run install
dcfg run uninstall
dcfg run install uninstall
```


#### The Attributes available to the Task Set are:

  -	name

       A descriptor that will be used to call this Task Set

  -	env

       A list of key => value paris that can be templated into the values of tasks. env values that are uppercase
       and start with a `$` wo;; be replaced with values from env vars on teh host.

  -	user (not implemented yet)

       Will allow you to set the user the Task Set should be run as.

  -	workdir

       Allows you to run the entire Task Set from a specified directory.

  -	tasks

       A list of task objects taht define what this Task Set does when executed.


#### A task can contain these default attributes:

  -	name  

       name of the task

  -	dest  

       what this action will be performed on

  -	workdir

       where this action will be called from

  -	wait

       how long to wait before this action is called

  -	repeat 

       how many times to run this action

  -	ignore

       ignore failures or not


##### The type of a task is defined by the `action` field.

This field's value is defined by eash plugin which imlpementes its own action.
For instance, the four currently existing plugins implement these actions:

  - command
       Runs an arbitrary bash command. No piping allowed for now.

  - write
       Writes a string or text block to file. Can use env vars.

  - replace
       Replace a string with another in a file. Works with regex.

  - config
       Add or change values in a config file.
  
  - template
      Configure an application in the container. Currently provides plugins for Drupal and WordPress.
      Template provides its own set of attributes, see [Template Task Type](docs/template-task-type.md) for details.


These actions when defined in the drud.yaml look like this:

###### command:

```
  - name: generic bash command
    action: command
    cmd: ls -l
```

###### write:

```
  - name: make content
    action: write
    workdir: testing
    write: |
        what what what
        what what {{.git_token}}
        whathwat
    dest: newfile.txt
    mode: 0777
```

###### replace:

given a file named turtle.txt that contains `name: bob` and you want to replace `bob` with `james`

```
  - name: replace name
    action: replace
    find: "(name:) ([a-z]*)"
    replace: "$1 james"
    dest: turtle.txt
```

###### config

```
  - name: config file
    action: config
    delim: ": "
    items:
      name: "configuration"
      count: "10"
      debug: "true"
    dest: turtle.conf
```


##### Full drud.yaml example:

```
- name: install
  env:
    site_name: MangoTango
    git_token: $GITHUB_TOKEN
  tasks:
  - action: command
    cmd: echo "{{.site_name}}"
    wait: 1s
    repeat: 3
  - name: create directory
    action: command
    cmd: mkdir testing
    ignore: yes
  - name: make file in new directory
    action: command
    workdir: testing
    cmd: touch newfile.txt
  - name: make content
    action: write
    workdir: testing
    write: |
        what what what
        what what {{.git_token}}
        whathwat
    dest: newfile.txt
    mode: 0777
  - name: look
    action: command
    cmd: ls -l testing

- name: uninstall
  tasks:
  - cmd: echo "testing this thing"
    action: command
    name: say something cool
  - cmd: rm newfile.txt
    action: command
    name: remove new file
    workdir: testing
  - name: rm directory
    action: command
    cmd: rm -rf testing
```

##This project is pluggable

In order to add new types of task actions to the projects you must do two things.

1. Implement the Action interface found in plugins/task.go in a new file. Name the file after whatever the new action
is called. 

2. Update the `TypeMap` at the bottom of plugins/task.go to map your action name to a pointer to an instance of your new action.

