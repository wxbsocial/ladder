# ladder
登高看世界
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ladder . 

./ladder -m client -s 82.157.32.37:1086 -socks 127.0.0.1:1080 -csi 82.157.32.37:1084
./ladder -m client -s 82.157.32.37:1086 -socks 82.157.32.37:10808 -csi 82.157.32.37:1084
./ladder -m server -o 0.0.0.0:1086



dV6#R5K6JUTA_wo{

scp -r ladder root@45.32.38.247:/home/root


scp -r ladder root@82.157.32.37:/home/qishi


配置shadowsocks服务端：
在 /etc目录下创建  shadowsocks.json 文件，将下面的内容放进去：
{
"server":"0.0.0.0",
"server_port":1080,
"password":"Zhf@135246",
"timeout":300,
"method":"aes-256-cfb",
"fast_open": false,
"workers": 1
}

vi /etc/systemd/system/shadowsocks.service
[Unit]
Description=Shadowsocks
[Service]
TimeoutStartSec=0
ExecStart=/usr/bin/ssserver -c /etc/shadowsocks.json
[Install]
WantedBy=multi-user.target

systemctl enable shadowsocks
systemctl start shadowsocks
systemctl status shadowsocks -l

curl --proxy "socks5://127.0.0.1:1080" --proxy-user test1:12345 \
https://job.toutiao.com/s/JxLbWby

curl --proxy "socks5://82.157.32.37:1080" --connect-timeout 1000 -m 2000 --proxy-user test1:12345 \
https://job.toutiao.com/s/JxLbWby

curl --proxy "socks5://82.157.32.37:10808" --proxy-user test1:12345 \
https://job.toutiao.com/s/JxLbWby


curl --proxy "socks5://127.0.0.1:1080" \
https://job.toutiao.com/s/JxLbWby