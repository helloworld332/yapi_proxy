1、按yapi官方安装yapi在本地机器上

2、在安装了yapi的机器上运行本目录下的脚本

3、运行脚本的依赖
   a. 安装python2.7
   b. 安装pymongo, 可通过pip install pymongo命令安装

4、脚本执行
  copy_project.py: 复制project, 需提供被复制的project_id, 复制后的工程默认命名为“XX_copy”(可通过项目的“配置”菜单进行修改)
  remove_project.py: 删除project, 主要是复制出来的库不需要时，可以删除掉

  project_id怎么找到? 可通过在yapi网站上点击进入某个项目，看项目的url看出来
      比如我进入一个项目的url是http://yapi.intely.cn/project/143/, 那么该项目的id是143
