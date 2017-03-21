# Acquia Site factory local import

A small utility to import databases for ACSF for local setup.

## About
Download all the backups from a particular day. You can either download them one by one or rsync all of the once filtering based the date. Acquia adds a suffix to backup like 2017-03-21.

Create a mapping file, lets call is config.json - with name and id tuples. The name is the name of local DB and ID is the ACSF DB name. You can see the this ID is part of backup you have downloaded from the cloud.

```
  {
    "name": "sitefactory_site1",
    "id": "testsub01ldb100001"
  },
  {
    "name": "sitefactory_site2",
    "id": "testsub01ldb100002"
  }
```

## Usage
```
  -config-file string
    Config file path. (default ./config.json)
  -mysql-pass string
    MySQL password. (default "drupal")
  -mysql-user string
    MySQL user name. (default "drupal")
  -source-dir string
    B dump source directory. (default "current folder")
```

## Import Databases
Import all databases. You will have to create target databases manually.
```
./acsfdbimport -source-dir=some/dir -config-file=some/file/config.json
```
## TODO
* Add better error handling and messages.
* Better command line options.

## Contribution
* File and issue, feature request.
* Send a PR.
