echo never | tee /sys/kernel/mm/transparent_hugepage/enabled
echo never | tee /sys/kernel/mm/transparent_hugepage/defrag
mkdir -p /data/mysqlconf
mkdir -p /data/mysql
chown 1001:root /data/mysql
