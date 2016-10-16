git pull
git submodule update
while :
do
	echo "Press [CTRL+C] to stop.."
	go run evedata-server.go
done
