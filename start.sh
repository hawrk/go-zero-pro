#!/bin/bash
echo "启动assess-api-server"
cd assess-api-server
nohup ./assess-api-server &
echo "启动assess-api-server 完成！"

cd ..
echo "mornano-rpc-server"
cd mornano-rpc-server
nohup ./mornano-rpc-server &
echo "mornano-rpc-server 完成！"


cd ..
echo "启动assess-rpc-server"
cd assess-rpc-server
nohup ./assess-rpc-server &
echo "启动assess-rpc-server 完成！"

cd ..
echo "启动assess-mq-server"
cd assess-mq-server
nohup ./assess-mq-server &
echo "启动assess-mq-server 完成！"

cd ..
echo "启动market-mq-server"
cd market-mq-server
nohup ./market-mq-server &
echo "启动market-mq-server 完成！"

