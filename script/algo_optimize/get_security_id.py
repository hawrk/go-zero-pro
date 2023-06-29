#!/usr/bin/python3

import pymysql

# 打开数据库连接
db = pymysql.connect(host='192.168.1.80', user='root', password='Root_123', database='jlp_algo')
# 使用 cursor() 方法创建一个游标对象 cursor
cursor = db.cursor()

# 查询所有证券 使用 execute()  方法执行 SQL 查询 使用 fetchall() 方法获取所有条件数据.
get_secs = 'select security_id from tb_security_info'
cursor.execute(get_secs)
secs = cursor.fetchall()
securitys = []
for sec in secs:
    securitys.append(sec[0])

print(securitys)
db.close()
