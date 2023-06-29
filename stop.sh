#!/bin/bash

echo "停止assess-api-server"
PID=`ps -ef|grep "assess-api-server"|grep -v grep |awk '{print $2}'`
if [ -n "${PID}" ]
then
  kill -9 ${PID}
fi
echo "停止assess-api-server 完成！"

echo "停止assess-rpc-server"
PID=`ps -ef|grep "assess-rpc-server"|grep -v grep |awk '{print $2}'`
if [ -n "${PID}" ]
then
  kill -9 ${PID}
fi
echo "停止assess-rpc-server 完成！"

echo "停止assess-mq-server"
PID=`ps -ef|grep "assess-mq-server"|grep -v grep |awk '{print $2}'`
if [ -n "${PID}" ]
then
  kill -9 ${PID}
fi
echo "停止assess-mq-server 完成！"


echo "停止market-mq-server"
PID=`ps -ef|grep "market-mq-server"|grep -v grep |awk '{print $2}'`
if [ -n "${PID}" ]
then
  kill -9 ${PID}
fi
echo "停止market-mq-server 完成！"

echo "停止mornano-rpc-server"
PID=`ps -ef|grep "mornano-rpc-server"|grep -v grep |awk '{print $2}'`
if [ -n "${PID}" ]
then
  kill -9 ${PID}
fi
echo "停止mornano-rpc-server 完成！"

