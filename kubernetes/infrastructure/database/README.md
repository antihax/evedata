# Database Setup (required)

Run setup.sh as root on DB node and label with database=mysql
Apply sql.yaml to cluster
Create database users
Apply eve.sql and evedata.sql

# Database Backup (optional)

Edit the setup.sh file to include credentials for DB and B2 storage.
Run ./setup.sh and inspect logs to verify correct application.