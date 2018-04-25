
supervisordConfPath=~/etc/supervisord.conf
appConfDir=~/etc/supervisord_conf/
if [ ! -f "$supervisordConfPath" ]
    then
    echo "Error: supervisor config file miss"
    exit 1
fi

count=`ps -ef|grep supervisord|grep -v grep |grep $USER |wc -l`
if [ $count -lt 1 ]
    then
    echo "Error: supervisord not running"
    exit 1
fi

if [ ! -d "$appConfDir" ]
    then
    echo "$appConfDir not exist, will be maked"
    mkdir "$appConfDir"
fi

cd .. && bin_name=`pwd` && cd -
sed -i "s%intely%${bin_name}%g" install.etc
cp install.etc ~/etc/supervisord_conf/server_name.conf
supervisorctl -c ~/etc/supervisord.conf reread
