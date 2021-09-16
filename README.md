m3u8 下载工具。

# Example

```shell
# 直接写入 m3u 文件的 URL 地址下载
m3ugo -u "m3u文件链接" -o filename

# 添加请求 Header（也可以在配置文件 config.yml 中添加）
m3ugo -H 'User-Agent:Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36' -H 'Host:www.baidu.com'

# 资源的直接地址
# 需要实现 Fetcher 接口，实现从地址到 m3u8 地址的解析
m3ugo -l "link"
```