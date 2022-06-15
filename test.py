
import os
for x in range (1000):
    os.system('mysql -h 34.101.243.58 -u user -puser -P3309  -e "show databases  "')
    print(x)