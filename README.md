# awscp

**awscp** is Go module providing functions that check EC2 info and transfer files between Local and EC2 with SSH protocal.

**awscp** is built using `github.com/bramvdbogaerde/go-scp`, `golang.org/x/crypto/ssh`, and `github.com/aws/aws-sdk-go-v2` module.

# Examples

**localtoec2**: helps to distribute files (or directory using `-recursive`) in Local machine to multiple EC2 machines. A binary file is available in **Releases**.

```
Usage: localtoec2 [flags] [File1] [File2] ...
       localtoec2 [flags] -recursive [Dir]

[File1] [File2] ...: File path to be copied
[Dir] ...: Directory path to be copied

Flags:
  -instance-ids string
        EC2 instance IDs: e.g. i-1234567890abcdef0,i-1234567890abcdef1
  -key-path string
        Path of key pair
  -name string
        Name of EC2 instances
  -permission string
        Permission of remote file: default - 0755 (default "0755")
  -platfrom string
        OS platform of EC2 instances: amazonlinux, ubuntu, centos, rhel, debian, suse
        if empty, the platform will be predicted
  -recursive
        Copy files recursively
  -remote-dir string
        Path of remote directory which files are copied to: default - home directory, e.g. /home/{username}/dir = dir
  -tag-key string
        Tag key of EC2 instances
```

**ec2tolocal**: helps to bring remote files (or directory) in EC2 machines to Local machine. Files to be copied must have same paths throughout machines. Files in each machine will be saved in directory named each instance id. A binary file is available in **Releases**.

```
Usage: ec2tolocal [flags] [local-dir]

[local-dir]: directory path which files is copied to (required)

Flags:
  -instance-ids string
        EC2 instance IDs: e.g. i-1234567890abcdef0,i-1234567890abcdef1
  -key-path string
        Path of key pair (required)
  -name string
        Name of EC2 instances
  -platfrom string
        OS platform of EC2 instances: amazonlinux, ubuntu, centos, rhel, debian, suse
        if empty, the platform will be predicted
  -remote-dir string
        Path of remote directory which files are copied from: relative path - home directory, e.g. /home/{username}/dir = dir
  -remote-file string
        Path of remote file: relative path - home directory, e.g. /home/{username}/file.txt = file.txt
  -tag-key string
        Tag key of EC2 instances
```
