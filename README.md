# 基于 Xray 的被动扫描管理

## Install

1. Create Database
2. Modify config.yaml
3. Download xray (and lic), genca
4. Running...

## Usage

### Create Database

```sql
create database `weblisten` default character set utf8mb4 collate utf8mb4_unicode_ci;
create user weblisten@'127.0.0.1' identified by 'gyuawbdvuyiabdu';
grant all on `weblisten`.* to weblisten;
flush privileges;

set global innodb_large_prefix=on;
set global innodb_file_format=Barracuda;
```

## LICENSE

![WTFPL](http://www.wtfpl.net/wp-content/uploads/2012/12/wtfpl-badge-4.png)

[WTFPL](LICENSE)