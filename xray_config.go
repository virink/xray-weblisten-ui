package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// XrayConfig - xray config https://xray.cool/xray/#/configration/config
type XrayConfig struct {
	BasicCrawler struct {
		AllowVisitParentPath bool  `yaml:"allow_visit_parent_path"`
		MaxCountOfLinks      int64 `yaml:"max_count_of_links"`
		MaxDepth             int64 `yaml:"max_depth"`
		Restriction          struct {
			Excludes []string      `yaml:"excludes"`
			Includes []interface{} `yaml:"includes"`
		} `yaml:"restriction"`
	} `yaml:"basic_crawler"`
	HTTP struct {
		AllowMethods []string `yaml:"allow_methods"`
		Cookies      struct {
			Key string `yaml:"key"`
		} `yaml:"cookies"`
		DialTimeout int64 `yaml:"dial_timeout"`
		EnableHTTP2 bool  `yaml:"enable_http2"`
		FailRetries int64 `yaml:"fail_retries"`
		Headers     struct {
			UserAgent []string `yaml:"User-Agent"`
		} `yaml:"headers"`
		MaxConnsPerHost int64  `yaml:"max_conns_per_host"`
		MaxQPS          int64  `yaml:"max_qps"`
		MaxRedirect     int64  `yaml:"max_redirect"`
		MaxRespBodySize int64  `yaml:"max_resp_body_size"`
		Proxy           string `yaml:"proxy"`
		ReadTimeout     int64  `yaml:"read_timeout"`
		TLSSkipVerify   bool   `yaml:"tls_skip_verify"`
	} `yaml:"http"`
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	Mitm struct {
		AllowIPRange []interface{} `yaml:"allow_ip_range"`
		Auth         struct {
			Password string `yaml:"password"`
			Username string `yaml:"username"`
		} `yaml:"auth"`
		CaCert      string `yaml:"ca_cert"`
		CaKey       string `yaml:"ca_key"`
		ProxyHeader struct {
			Via        string `yaml:"via"`
			XForwarded bool   `yaml:"x_forwarded"`
		} `yaml:"proxy_header"`
		Queue struct {
			MaxLength int64 `yaml:"max_length"`
		} `yaml:"queue"`
		Restriction struct {
			Excludes []string `yaml:"excludes"`
			Includes []string `yaml:"includes"`
		} `yaml:"restriction"`
		UpstreamProxy string `yaml:"upstream_proxy"`
	} `yaml:"mitm"`
	Plugins struct {
		Baseline struct {
			DetectChinaIDCardNumber         bool `yaml:"detect_china_id_card_number"`
			DetectCookieHttponly            bool `yaml:"detect_cookie_httponly"`
			DetectCookiePasswordLeak        bool `yaml:"detect_cookie_password_leak"`
			DetectCorsHeaderConfig          bool `yaml:"detect_cors_header_config"`
			DetectHTTPHeaderConfig          bool `yaml:"detect_http_header_config"`
			DetectOutdatedSslVersion        bool `yaml:"detect_outdated_ssl_version"`
			DetectSerializationDataInParams bool `yaml:"detect_serialization_data_in_params"`
			DetectServerErrorPage           bool `yaml:"detect_server_error_page"`
			DetectUnsafeScheme              bool `yaml:"detect_unsafe_scheme"`
			Enabled                         bool `yaml:"enabled"`
		} `yaml:"baseline"`
		BruteForce struct {
			DetectDefaultPassword   bool   `yaml:"detect_default_password"`
			DetectUnsafeLoginMethod bool   `yaml:"detect_unsafe_login_method"`
			Enabled                 bool   `yaml:"enabled"`
			PasswordDictionary      string `yaml:"password_dictionary"`
			UsernameDictionary      string `yaml:"username_dictionary"`
		} `yaml:"brute_force"`
		CmdInjection struct {
			DetectBlindInjection bool `yaml:"detect_blind_injection"`
			Enabled              bool `yaml:"enabled"`
		} `yaml:"cmd_injection"`
		CrlfInjection struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"crlf_injection"`
		Dirscan struct {
			Depth      int64  `yaml:"depth"`
			Dictionary string `yaml:"dictionary"`
			Enabled    bool   `yaml:"enabled"`
		} `yaml:"dirscan"`
		Fastjson struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"fastjson"`
		Jsonp struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"jsonp"`
		MaxParallel   int64 `yaml:"max_parallel"`
		PathTraversal struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"path_traversal"`
		Phantasm struct {
			Depth   int64    `yaml:"depth"`
			Enabled bool     `yaml:"enabled"`
			Poc     []string `yaml:"poc"`
		} `yaml:"phantasm"`
		Redirect struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"redirect"`
		Sqldet struct {
			BooleanBasedDetection      bool `yaml:"boolean_based_detection"`
			DangerouslyUseCommentInSQL bool `yaml:"dangerously_use_comment_in_sql"`
			DangerouslyUseOrInSQL      bool `yaml:"dangerously_use_or_in_sql"`
			Enabled                    bool `yaml:"enabled"`
			ErrorBasedDetection        bool `yaml:"error_based_detection"`
			TimeBasedDetection         bool `yaml:"time_based_detection"`
		} `yaml:"sqldet"`
		Ssrf struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"ssrf"`
		Struts struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"struts"`
		Thinkphp struct {
			DetectThinkPhpSQLInjection bool `yaml:"detect_think_php_sql_injection"`
			Enabled                    bool `yaml:"enabled"`
		} `yaml:"thinkphp"`
		Upload struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"upload"`
		XSS struct {
			Enabled       bool `yaml:"enabled"`
			IeFeature     bool `yaml:"ie_feature"`
			IncludeCookie bool `yaml:"include_cookie"`
		} `yaml:"xss"`
		XXE struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"xxe"`
	} `yaml:"plugins"`
	Reverse struct {
		Client struct {
			DNSServerIP   string `yaml:"dns_server_ip"`
			HTTPBaseURL   string `yaml:"http_base_url"`
			RemoteServer  bool   `yaml:"remote_server"`
			RmiServerAddr string `yaml:"rmi_server_addr"`
		} `yaml:"client"`
		DBFilePath string `yaml:"db_file_path"`
		DNS        struct {
			Domain             string `yaml:"domain"`
			Enabled            bool   `yaml:"enabled"`
			IsDomainNameServer bool   `yaml:"is_domain_name_server"`
			ListenIP           string `yaml:"listen_ip"`
			Resolve            []struct {
				Record string `yaml:"record"`
				TTL    int64  `yaml:"ttl"`
				Type   string `yaml:"type"`
				Value  string `yaml:"value"`
			} `yaml:"resolve"`
		} `yaml:"dns"`
		HTTP struct {
			Enabled    bool   `yaml:"enabled"`
			ListenIP   string `yaml:"listen_ip"`
			ListenPort string `yaml:"listen_port"`
		} `yaml:"http"`
		Rmi struct {
			Enabled    bool   `yaml:"enabled"`
			ListenIP   string `yaml:"listen_ip"`
			ListenPort string `yaml:"listen_port"`
		} `yaml:"rmi"`
		Token string `yaml:"token"`
	} `yaml:"reverse"`
	Subdomain struct {
		AllowRecursive bool     `yaml:"allow_recursive"`
		DNSServers     []string `yaml:"dns_servers"`
		MainDictionary string   `yaml:"main_dictionary"`
		MaxDepth       int64    `yaml:"max_depth"`
		Modes          []string `yaml:"modes"`
		SubDictionary  string   `yaml:"sub_dictionary"`
		WorkerCount    int64    `yaml:"worker_count"`
	} `yaml:"subdomain"`
	Version float64 `yaml:"version"`
}

const webhook = "/vul_webhook"

var xrayConfigTemplate = `version: 2.3

# 配置解释见 https://chaitin.github.io/xray/#/configration/plugins
plugins:
  max_parallel: 10
  xss:
    enabled: true
    ie_feature: false
    include_cookie: false
  baseline:
    enabled: true
    detect_outdated_ssl_version: false
    detect_http_header_config: false
    detect_cors_header_config: true
    detect_server_error_page: false
    detect_china_id_card_number: false
    detect_serialization_data_in_params: true
    detect_cookie_password_leak: true
    detect_unsafe_scheme: false
    detect_cookie_httponly: false
  cmd_injection:
    enabled: true
    detect_blind_injection: false
  crlf_injection:
    enabled: true
  dirscan:
    enabled: true
    depth: 1
    dictionary: ""
  jsonp:
    enabled: true
  path_traversal:
    enabled: true
  redirect:
    enabled: true
  sqldet:
    enabled: true
    error_based_detection: true
    boolean_based_detection: true
    time_based_detection: true
    # 下面两个选项很危险，开启之后可以增加检测率，但是有破坏数据库数据的可能性，请务必了解工作原理之后再开启
    dangerously_use_comment_in_sql: false
    dangerously_use_or_in_sql: false
  ssrf:
    enabled: true
  xxe:
    enabled: true
  upload:
    enabled: true
  brute_force:
    enabled: true
    detect_default_password: true
    detect_unsafe_login_method: true
    username_dictionary: ""
    password_dictionary: ""
  struts:
    enabled: true
  thinkphp:
    enabled: true
    detect_think_php_sql_injection: true
  fastjson:
    enabled: true
  phantasm:
    enabled: true
    depth: 1
    poc:
      - poc-yaml-74cms-sqli-1
      - poc-yaml-74cms-sqli-2
      - poc-yaml-74cms-sqli
      - poc-yaml-activemq-cve-2016-3088
      - poc-yaml-bash-cve-2014-6271
      - poc-yaml-cacti-weathermap-file-write
      - poc-yaml-citrix-cve-2019-19781-path-traversal
      - poc-yaml-coldfusion-cve-2010-2861-lfi
      - poc-yaml-confluence-cve-2015-8399
      - poc-yaml-confluence-cve-2019-3396-lfi
      - poc-yaml-coremail-cnvd-2019-16798
      - poc-yaml-couchcms-cve-2018-7662
      - poc-yaml-couchdb-cve-2017-12635
      - poc-yaml-couchdb-unauth
      - poc-yaml-dedecms-carbuyaction-fileinclude
      - poc-yaml-dedecms-cve-2018-6910
      - poc-yaml-dedecms-cve-2018-7700-rce
      - poc-yaml-dedecms-guestbook-sqli
      - poc-yaml-dedecms-membergroup-sqli
      - poc-yaml-dedecms-url-redirection
      - poc-yaml-discuz-ml3x-cnvd-2019-22239
      - poc-yaml-discuz-v72-sqli
      - poc-yaml-discuz-wechat-plugins-unauth
      - poc-yaml-discuz-wooyun-2010-080723
      - poc-yaml-dlink-850l-info-leak
      - poc-yaml-dlink-cve-2019-16920-rce
      - poc-yaml-dlink-cve-2019-17506
      - poc-yaml-docker-api-unauthorized-rce
      - poc-yaml-docker-registry-api-unauth
      - poc-yaml-druid-monitor-unauth
      - poc-yaml-drupal-cve-2019-6340
      - poc-yaml-drupal-drupalgeddon2-rce
      - poc-yaml-drupalgeddon-cve-2014-3704-sqli
      - poc-yaml-duomicms-sqli
      - poc-yaml-dvr-cve-2018-9995
      - poc-yaml-ecology-filedownload-directory-traversal
      - poc-yaml-ecology-javabeanshell-rce
      - poc-yaml-ecology-springframework-directory-traversal
      - poc-yaml-ecology-syncuserinfo-sqli
      - poc-yaml-ecology-validate-sqli
      - poc-yaml-ecology-workflowcentertreedata-sqli
      - poc-yaml-ecshop-360-rce
      - poc-yaml-elasticsearch-cve-2014-3120
      - poc-yaml-elasticsearch-cve-2015-1427
      - poc-yaml-elasticsearch-cve-2015-3337-lfi
      - poc-yaml-elasticsearch-unauth
      - poc-yaml-etcd-unauth
      - poc-yaml-etouch-v2-sqli
      - poc-yaml-feifeicms-lfr
      - poc-yaml-finereport-directory-traversal
      - poc-yaml-glassfish-cve-2017-1000028-lfi
      - poc-yaml-hadoop-yarn-unauth
      - poc-yaml-ifw8-router-cve-2019-16313
      - poc-yaml-influxdb-unauth
      - poc-yaml-jboss-cve-2010-1871
      - poc-yaml-jboss-unauth
      - poc-yaml-jenkins-cve-2018-1000600
      - poc-yaml-jenkins-cve-2018-1000861-rce
      - poc-yaml-jira-cve-2019-11581
      - poc-yaml-jira-ssrf-cve-2019-8451
      - poc-yaml-joomla-cnvd-2019-34135-rce
      - poc-yaml-joomla-cve-2015-7297-sqli
      - poc-yaml-joomla-cve-2017-8917-sqli
      - poc-yaml-joomla-ext-zhbaidumap-cve-2018-6605-sqli
      - poc-yaml-kibana-unauth
      - poc-yaml-laravel-debug-info-leak
      - poc-yaml-maccms-rce
      - poc-yaml-maccmsv10-backdoor
      - poc-yaml-metinfo-cve-2019-16996-sqli
      - poc-yaml-metinfo-cve-2019-16997-sqli
      - poc-yaml-metinfo-cve-2019-17418-sqli
      - poc-yaml-metinfo-lfi-cnvd-2018-13393
      - poc-yaml-mongo-express-cve-2019-10758
      - poc-yaml-msvod-sqli
      - poc-yaml-myucms-lfr
      - poc-yaml-nagio-cve-2018-10735
      - poc-yaml-nagio-cve-2018-10736
      - poc-yaml-nagio-cve-2018-10737
      - poc-yaml-nagio-cve-2018-10738
      - poc-yaml-netgear-cve-2017-5521
      - poc-yaml-nextjs-cve-2017-16877
      - poc-yaml-nexus-cve-2019-7238
      - poc-yaml-nhttpd-cve-2019-16278
      - poc-yaml-nuuo-file-inclusion
      - poc-yaml-pandorafms-cve-2019-20224-rce
      - poc-yaml-php-cgi-cve-2012-1823
      - poc-yaml-phpcms-cve-2018-19127
      - poc-yaml-phpmyadmin-cve-2018-12613-file-inclusion
      - poc-yaml-phpmyadmin-setup-deserialization
      - poc-yaml-phpok-sqli
      - poc-yaml-phpstudy-backdoor-rce
      - poc-yaml-phpunit-cve-2017-9841-rce
      - poc-yaml-pulse-cve-2019-11510
      - poc-yaml-pyspider-unauthorized-access
      - poc-yaml-qibocms-sqli
      - poc-yaml-rails-cve-2018-3760-rce
      - poc-yaml-razor-cve-2018-8770
      - poc-yaml-rconfig-cve-2019-16663
      - poc-yaml-resin-cnnvd-200705-315
      - poc-yaml-resin-inputfile-fileread-or-ssrf
      - poc-yaml-resin-viewfile-fileread
      - poc-yaml-satellian-cve-2020-7980-rce
      - poc-yaml-seacms-rce
      - poc-yaml-seacms-sqli
      - poc-yaml-seacms-v654-rce
      - poc-yaml-seeyon-wooyun-2015-0108235-sqli
      - poc-yaml-solr-cve-2017-12629-xxe
      - poc-yaml-solr-cve-2019-0193
      - poc-yaml-solr-velocity-template-rce
      - poc-yaml-spark-unauth
      - poc-yaml-spring-cve-2016-4977
      - poc-yaml-springcloud-cve-2019-3799
      - poc-yaml-supervisord-cve-2017-11610
      - poc-yaml-tensorboard-unauth
      - poc-yaml-thinkcmf-write-shell
      - poc-yaml-thinkphp-v6-file-write
      - poc-yaml-thinkphp5-controller-rce
      - poc-yaml-thinkphp5023-method-rce
      - poc-yaml-tomcat-cve-2017-12615-rce
      - poc-yaml-tomcat-cve-2018-11759
      - poc-yaml-tpshop-sqli
      - poc-yaml-typecho-rce
      - poc-yaml-uwsgi-cve-2018-7490
      - poc-yaml-vbulletin-cve-2019-16759
      - poc-yaml-weblogic-cve-2017-10271-reverse
      - poc-yaml-weblogic-cve-2019-2729-1
      - poc-yaml-weblogic-cve-2019-2729-2
      - poc-yaml-weblogic-ssrf
      - poc-yaml-weblogic-cve-2017-10271
      - poc-yaml-weblogic-cve-2019-2725
      - poc-yaml-webmin-cve-2019-15107-rce
      - poc-yaml-wordpress-ext-adaptive-images-lfi
      - poc-yaml-wordpress-ext-mailpress-rce
      - poc-yaml-wuzhicms-v410-sqli
      - poc-yaml-yccms-rce
      - poc-yaml-youphptube-encoder-cve-2019-5127
      - poc-yaml-youphptube-encoder-cve-2019-5128
      - poc-yaml-youphptube-encoder-cve-2019-5129
      - poc-yaml-yungoucms-sqli
      - poc-yaml-zabbix-authentication-bypass
      - poc-yaml-zabbix-cve-2016-10134-sqli
      - poc-yaml-zcms-v3-sqli
      - poc-yaml-zimbra-cve-2019-9670-xxe
      - poc-yaml-zzcms-zsmanage-sqli
      - poc-go-ecology-db-config-info-leak
      - poc-go-php-cve-2019-11043-rce
      - poc-go-seeyon-htmlofficeservlet-rce
      - poc-go-tomcat-cve-2020-1938
      - poc-go-tomcat-put
      - poc-go-tongda-lfi-upload-rce

log:
  level: info # 支持 debug, info, warn, error, fatal

# 配置解释见 https://chaitin.github.io/xray/#/configration/mitm
mitm:
  ca_cert: ./ca.crt
  ca_key: ./ca.key
  auth:
    username: ""
    password: ""
  restriction:
    includes: # 允许扫描的域，此处无协议
    - "*.vulnweb.com"
    excludes:
    - '*google*'
    - '*github*'
    - '*.gov.cn'
    - '*.edu.cn'
    - '*chaitin*'
  allow_ip_range: []
  queue:
    max_length: 3000
  proxy_header:
    via: "" # 如果不为空，proxy 将添加类似 Via: 1.1 $some-value-$random 的 http 头
    x_forwarded: false # 是否添加 X-Forwarded-{For,Host,Proto,Url} 四个 http 头
  upstream_proxy: "" # mitm 的全部流量继续使用 proxy

# 配置解释见 https://chaitin.github.io/xray/#/configration/basic-crawler
basic_crawler:
  max_depth: 0 # 爬虫最大深度, 0 为无限制
  max_count_of_links: 0 # 本次扫描总共爬取的最大连接数， 0 为无限制
  allow_visit_parent_path: false # 是否允许访问父目录, 如果扫描目标为 example.com/a/， 如果该项为 false, 那么就不会爬取 example.com/ 这级目录的内容
  restriction: # 和 mitm 中的写法一致, 有个点需要注意的是如果当前目标为 example.com 那么会自动添加 example.com 到 includes 中。
    includes: []
    excludes:
    - '*google*'

# 配置解释见 https://chaitin.github.io/xray/#/configration/subdomain
subdomain:
  modes: # 使用哪些方式获取子域名
    - brute # 字典爆破模式
    - api # 使用各大 api 获取
    - zone_transfer # 尝试使用域传送漏洞获取
  worker_count: 100 # 决定同时允许多少个 DNS 查询
  dns_servers: # 查询使用的 DNS server
    - 1.1.1.1
    - 8.8.8.8
  allow_recursive: false # 是否允许递归扫描，开了后如果发现 a.example.com 将继续扫描 a.example.com 的子域名
  max_depth: 5 # 最大允许的子域名深度
  main_dictionary: "" # 一级子域名字典， 绝对路径
  sub_dictionary: "" # 其它层级子域名字典， 绝对路径

# 配置解释见 https://chaitin.github.io/xray/#/configration/reverse
reverse:
  db_file_path: ""
  token: ""
  http:
    enabled: true
    listen_ip: 127.0.0.1
    listen_port: ""
  dns:
    enabled: false
    listen_ip: 0.0.0.0
    domain: ""
    is_domain_name_server: false
    # 静态解析规则
    resolve:
    - type: A # A, AAAA, TXT 三种
      record: localhost
      value: 0.0.0.0
      ttl: 60
  rmi:
    enabled: true
    listen_ip: 0.0.0.0
    listen_port: ""
  client:
    http_base_url: ""
    dns_server_ip: ""
    rmi_server_addr: ""
    remote_server: false

# 配置解释见 https://chaitin.github.io/xray/#/configration/http
http:
  proxy: "" # 漏洞扫描时使用的代理
  dial_timeout: 5 # 建立 tcp 连接的超时时间
  read_timeout: 10 # 读取 http 响应的超时时间，不可太小，否则会影响到 sql 时间盲注的判断
  fail_retries: 1 # 请求失败的重试次数，0 则不重试
  max_redirect: 5 # 单个请求最大允许的跳转数
  max_qps: 500 # 每秒最大请求数
  max_conns_per_host: 50 # 同一 host 最大允许的连接数，可以根据目标主机性能适当增大。
  max_resp_body_size: 8388608 # 8M，单个请求最大允许的响应体大小，超过该值 body 就会被截断
  headers: # 每个请求预置的 http 头
    User-Agent:
      - Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169
  cookies: # 每个请求预置的 cookie 值，效果上相当于添加了一个 Header: Cookie: key=value
    key: value
  allow_methods: # 允许使用 http 方法
    - HEAD
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
    - CONNECT
    - PROPFIND
    - MOVE
  tls_skip_verify: true # 是否验证目标网站的 https 证书。
  enable_http2: false # 是否启用 http2
`

func getDefaultXrayConfig() (xrayConf XrayConfig, err error) {
	if err = yaml.Unmarshal([]byte(xrayConfigTemplate), &xrayConf); err != nil {
		logger.Errorln(err.Error())
		return xrayConf, err
	}
	return xrayConf, nil
}

func loadXrayConfig(filename string) (xrayConf XrayConfig, err error) {
	var (
		yamlFile []byte
	)
	_, err = os.Stat(filename)
	if err != nil && os.IsExist(err) {
		if data := templateConfig(); data != nil {
			if err = ioutil.WriteFile(filename, data, 0666); err != nil {
				return xrayConf, err
			}
		}
	}
	if yamlFile, err = ioutil.ReadFile(filename); err != nil {
		logger.Errorln(err.Error())
		return xrayConf, err
	}
	if err = yaml.Unmarshal(yamlFile, &xrayConf); err != nil {
		logger.Errorln(err.Error())
		return xrayConf, err
	}
	return xrayConf, nil
}
