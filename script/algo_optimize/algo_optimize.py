#!/usr/bin/python3

from datetime import datetime

import pymysql

# 打开数据库连接
db = pymysql.connect(host='192.168.1.80', user='root', password='Root_123', database='jlp_algo')
target_db = pymysql.connect(host='192.168.1.80', user='root', password='Root_123', database='assess')
# 使用 cursor() 方法创建一个游标对象 cursor
cursor = db.cursor()
target_cursor = target_db.cursor()

# 清空原有表的数据
target_cursor.execute('delete from assess.tb_algo_optimize_base')
target_cursor.execute('delete from assess.tb_algo_optimize')
target_db.commit()

# 查询所有证券 使用 execute()  方法执行 SQL 查询 使用 fetchall() 方法获取所有条件数据.
get_secs = 'select security_id ,security_name from tb_security_info where id <= 4142'
cursor.execute(get_secs)
secs = cursor.fetchall()

# 查询所有算法
get_algos = 'select id,uuser_id,provider_name,algo_name,algorithm_type from tb_algo_info where algorithm_type=1 and id <=17'
cursor.execute(get_algos)
algos = cursor.fetchall()

for algo in algos:
    for sec in secs:
        add_optimize_base = 'insert into assess.tb_algo_optimize_base (provider_id,provider_name,sec_id,sec_name,algo_id,algo_type,algo_name,open_rate,income_rate,basis_point,create_time,update_time) values (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)'
        optimize_bases = (
            algo[1], str(algo[2]).replace('\0', ''), sec[0], str(sec[1]).replace(' ', ''), algo[0], algo[4],
            str(algo[3]).replace('\0', ''), 0, 0, 0, datetime.now(), datetime.now())
        target_cursor.execute(add_optimize_base, optimize_bases)
        get_algo_optimize = 'select * from assess.tb_algo_optimize where sec_id=%s '
        target_cursor.execute(get_algo_optimize, (sec[0]))
        # if target_cursor.rowcount == 0:
        #     # 优选表中没有数据就插入
        #     add_optimize = 'insert into assess.tb_algo_optimize (provider_id,provider_name,sec_id, sec_name, algo_id, algo_type, algo_name, score, create_time, upate_time) values (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)'
        #     optimizes = (
        #         algo[1], str(algo[2]).replace('\0', ''), sec[0], sec[1], algo[0], algo[4],
        #         str(algo[3]).replace('\0', ''), 0,
        #         datetime.now(), datetime.now())
        #     target_cursor.execute(add_optimize, optimizes)

target_db.commit()
target_cursor.execute("select * from tb_algo_optimize_base")
optimize_bases = target_cursor.fetchall()
for optimize_base in optimize_bases:
    print(optimize_base)
# print("===================================================")
# target_cursor.execute("select * from tb_algo_optimize")
# optimizes = target_cursor.fetchall()
# for optimize in optimizes:
#     print(optimize)

# 关闭数据库连接
target_db.close()
db.close()
