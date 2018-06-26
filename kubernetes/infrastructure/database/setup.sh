echo never | tee /sys/kernel/mm/transparent_hugepage/enabled
echo never | tee /sys/kernel/mm/transparent_hugepage/defrag
mkdir -p /data/mysqlconf
mkdir -p /data/mysql
chown 1001:root /data/mysql

echo -e "
[mysqld_safe]
nice            = -19

[client]
default-character-set=utf8

[mysql]
default-character-set=utf8

[mysqld]
collation-server = utf8_unicode_ci
init-connect='SET NAMES utf8'
character-set-server = utf8
require_secure_transport = ON

# log to syslog
log_syslog              = 1

skip-log-bin
default_storage_engine=TokuDB

bind-address            = 0.0.0.0

# optimizations
skip-name-resolve
max_prepared_stmt_count = 104857
thread_pool_size        = 16
innodb_buffer_pool_size = 8M
innodb_buffer_pool_instances = 1
innodb_log_buffer_size  = 8M
innodb_log_file_size    = 8M
innodb_flush_log_at_trx_commit = 0
innodb_thread_concurrency = 0
wait_timeout            = 340
interactive_timeout     = 580
max_connections         = 2046
thread_cache_size       = 10M
query_cache_size        = 0
query_cache_type        = 0
sort_buffer_size        = 20M
read_rnd_buffer_size    = 20M
join_buffer_size        = 512K
" > /data/mysqlconf/my.cnf

docker run -d --name percona -e INIT_TOKUDB=yes --log-driver syslog -e MYSQL_RANDOM_ROOT_PASSWORD=yes \
-p 0.0.0.0:3306:3306 -v /data/mysql:/var/lib/mysql -v /data/mysqlconf:/etc/mysql/conf.d percona/percona-server:5.7