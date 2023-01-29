PROXY=20860
PROCESS_NAME="./KingMatchServer_1"
PORT=26862
kill -9  $(ps -ef |grep $PROXY | awk '{if($8=="'$PROCESS_NAME'" && $10=="'$PORT'")  print $2}')

