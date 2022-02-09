#!/usr/bin/env bash

#file=$1
#keyword=$2
#cat $file | while read line
#do
#    isGet=$(echo $line| grep -o "GET")
#    isPost=$(echo $line| grep -o "POST")
#    role='^(\d+\.\d+\.\d+\.\d+)\s-\s([^\s]*|\-)\s(\[[^\[\]]+\])\s([^\s]*)\s(\"(?:[^"]|\")+|-\")\s([^\s]*)\s(\"\d{3}\")\s(\d+|-)\s(\"(?:[^"]|\")+|-\")\s(\"(?:[^"]|\")+|-\")\s(\"\d+\.\d+\.\d+\.\d+\,\s\d+\.\d+\.\d+\.\d+\"|\"\d+\.\d+\.\d+\.\d+\")$'
#    if [[ $isGet == "GET" ]]; then
#        echo $(echo $line | grep -Eo '${role}')
#    elif [[ $isPost != "" ]]; then
#        echo ""
#    fi
#done

#    log_format  main        '$remote_addr - $remote_user [$time_local] $host "$request" $request_body '
#                                '"$status" $body_bytes_sent "$http_referer" '
#                                '"$http_user_agent" "$http_x_forwarded_for" ';

#(\d+\.\d+\.\d+\.\d+)\s-\s([^\s]*|\-)\s(\[[^\[\]]+\])\s([^\s]*)\s(\"(?:[^"]|\")+|-\")\s([^\s]*)\s(\"\d{3}\")\s(\d+|-)\s(\"(?:[^"]|\")+|-\")\s(\"(?:[^"]|\")+|-\")\s(\"\d+\.\d+\.\d+\.\d+\,\s\d+\.\d+\.\d+\.\d+\,\s\d+\.\d+\.\d+\.\d+\"|\"\d+\.\d+\.\d+\.\d+\,\s\d+\.\d+\.\d+\.\d+\"|\"\d+\.\d+\.\d+\.\d+\")
#https://rubular.com/r/8kAzUCByxqRr3S

#https://www.linuxhub.org/?p=4219

#
#
#'$http_host $server_addr $remote_addr "$http_x_forwarded_for" [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $request_time $upstream_response_time';
#www.linuxhub.cn 192.168.60.74 192.168.60.59 "113.105.222.200, 192.168.62.184" [14/Mar/2017:10:27:00 +0800] "GET /hello.php HTTP/1.0" 200 2146 "http://www.ddd.cn/test9.php" "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:52.0) Gecko/20100101 Firefox/52.0" 0.001 0.000
#
#
#([^\s]*)              #匹配 $http_host
#([^\s]*)              #匹配 $host
#([^\s]*|\-)         #  $remote_user
#(\d+\.\d+\.\d+\.\d+)  #匹配 $server_addr,$remote_addr
#(\"\d+\.\d+\.\d+\.\d+\,\s\d+\.\d+\.\d+\.\d+\"|\"\d+\.\d+\.\d+\.\d+\") #匹配 "$http_x_forwarded_for"
#(\"\d+\.\d+\.\d+\.\d+\,\s\d+\.\d+\.\d+\.\d+\,\s\d+\.\d+\.\d+\.\d+\"|\"\d+\.\d+\.\d+\.\d+\,\s\d+\.\d+\.\d+\.\d+\"|\"\d+\.\d+\.\d+\.\d+\") #匹配 "$http_x_forwarded_for"  3 个IP时
#(\"(?:[^"]|\")+|-\") 匹配  "$http_x_forwarded_for"  "::ffff:112.24.49.114, 58.247.212.61, 47.93.76.16"
#(\[[^\[\]]+\])     #匹配 [$time_local]
#(\"(?:[^"]|\")+|-\")   #匹配"$request","$http_referer"，"$http_user_agent"
#(\d{3})            #匹配$status
#(\"\d{3}\")            #匹配 "$status"
#(\d+|-)            #匹配$body_bytes_sent
#(\d*\.\d*|\-)      #匹配$request_time,$upstream_response_time'
#^                  #匹配每行数据的开头
#$                  #匹配每行数据的结局
#
#
#^([^\s]*)\s(\d+\.\d+\.\d+\.\d+)\s(\d+\.\d+\.\d+\.\d+)\s(\"\d+\.\d+\.\d+\.\d+\,\s\d+\.\d+\.\d+\.\d+\"|\"\d+\.\d+\.\d+\.\d+\")\s(\[[^\[\]]+\])\s(\"(?:[^"]|\")+|-\")\s(\d{3})\s(\d+|-)\s(\"(?:[^"]|\")+|-\")\s(\"(?:[^"]|\")+|-\")\s(\d*\.\d*|\-)\s(\d*\.\d*|\-)$
#
#
#日志匹配效果
#查看: http://www.rubular.com/r/WxbGSkXWRi

#grep只显示匹配到的部分
#isPost=$(echo $line| grep -o "POST")

## nginx日志 urldecode 解码
#cat data.log | grep POST | grep data=% | awk '{print $10}' | perl -pe 's/\+/ /g; s/%([0-9a-f]{2})/chr(hex($1))/eig'
#
## nginx日志 urldecode 解码 + base64 解码  暂不生效
#cat data.log | grep POST | grep -v data=% | awk '{split($10,a,"[=&]");print a[2]}' | perl -pe 's/\+/ /g; s/%([0-9a-f]{2})/chr(hex($1))/eig'| base64 -d
#
