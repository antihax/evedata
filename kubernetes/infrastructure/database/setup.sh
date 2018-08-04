echo never | tee /sys/kernel/mm/transparent_hugepage/enabled
echo never | tee /sys/kernel/mm/transparent_hugepage/defrag
mkdir -p /data/mysqlconf
mkdir -p /data/mysql
mkdir -p /data/mysqlplugins
cp lib/udf_infusion.so /data/mysqlplugins
chown 1001:root /data/mysql
