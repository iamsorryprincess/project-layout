#!/usr/bin/bash

# set actual credentials here

mysqldump --host=stage.com --user=test --password --databases test > mysql-dump.sql
mysql --host=127.0.0.1 --port=3306 --user=test --password test < mysql-dump.sql