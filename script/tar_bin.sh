#! /bin/sh

if [ $# != 1 ]; then
	echo "num of parameter error"
	exit
fi

paltform=$1
this_server_name="rap_proxy"

if [ ${paltform} != "dev" -a ${paltform} != "release" ]; then
	echo "parameter ${paltform} error"
	exit
fi

cd ../server && go build

if [ $? -ne 0 ]; then
    echo "go build failed, please check out the files"
	exit
fi

cd -

cd  .. && mydir=`pwd` && cd -
dir_name=${mydir##*/}
bin_name=${dir_name}
echo "dir_name: ${dir_name}"
echo "bin_name: ${bin_name}"

tmp_dir=dir_$$

mkdir -p ${tmp_dir}/${dir_name}

mkdir -p ${tmp_dir}/${dir_name}/bin
mkdir -p ${tmp_dir}/${dir_name}/conf
mkdir -p ${tmp_dir}/${dir_name}/admin
mkdir -p ${tmp_dir}/${dir_name}/log
#mkdir -p ${tmp_dir}/${dir_name}/photos
#mkdir -p ${tmp_dir}/${dir_name}/static
#mkdir -p ${tmp_dir}/${dir_name}/tmp
#mkdir -p ${tmp_dir}/${dir_name}/report_stat

cp ../server/server ${tmp_dir}/${dir_name}/bin/
cp -rf ../conf/* ${tmp_dir}/${dir_name}/conf/
cp -rf admin/* ${tmp_dir}/${dir_name}/admin
#cp -rf ../photos/* ${tmp_dir}/${dir_name}/photos
#cp -rf ../static/* ${tmp_dir}/${dir_name}/static
#cp -rf ../report_stat/* ${tmp_dir}/${dir_name}/report_stat


if [ ${paltform} == "dev" ]; then
	echo dev
	rm -f ${tmp_dir}/${dir_name}/conf/config.json.*
	cp -rf ../conf/config.json.dev ${tmp_dir}/${dir_name}/conf/config.json
elif [ ${paltform} == "release" ]; then
	echo release
	rm -f ${tmp_dir}/${dir_name}/conf/config.json.*
	cp -rf ../conf/config.json.release ${tmp_dir}/${dir_name}/conf/config.json
	sed -i "s%~\/etc%\/etc%g" ${tmp_dir}/${dir_name}/admin/*.sh
fi

sed -i "s%server_name%${this_server_name}%g" ${tmp_dir}/${dir_name}/admin/*.*

find ${tmp_dir} -name "*.svn"|xargs rm -rf

now_date=`date "+%Y%m%d%H%M%S"`
filename="${paltform}_${this_server_name}_${now_date}_server.tar.gz"
echo $filename
cd ${tmp_dir} && tar zcvf ${filename} * && cd -
mv ${tmp_dir}/${filename} ./ && rm -rf ${tmp_dir}
