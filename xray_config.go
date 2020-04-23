package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// XrayConfig - xray config https://xray.cool/xray/#/configration/config
type XrayConfig struct {
	BasicCrawler struct {
		AllowVisitParentPath bool  `yaml:"allow_visit_parent_path"`
		MaxCountOfLinks      int64 `yaml:"max_count_of_links"`
		MaxDepth             int64 `yaml:"max_depth"`
		Restriction          struct {
			Excludes []string `yaml:"excludes"`
			Includes []string `yaml:"includes"`
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
		AllowIPRange []string `yaml:"allow_ip_range"`
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

func getDefaultXrayConfig() (xrayConf XrayConfig, err error) {
	defaultXrayConfigName := filepath.Join(conf.Xray.Path, "config.yaml")
	return loadXrayConfig(defaultXrayConfigName)
}

func loadXrayConfig(filename string) (xrayConf XrayConfig, err error) {
	var (
		yamlFile []byte
	)
	_, err = os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return xrayConf, err
	}
	if yamlFile, err = ioutil.ReadFile(filename); err != nil {
		logger.Errorln(err)
		return xrayConf, err
	}
	if err = yaml.Unmarshal(yamlFile, &xrayConf); err != nil {
		logger.Errorln(err)
		return xrayConf, err
	}
	return xrayConf, nil
}
