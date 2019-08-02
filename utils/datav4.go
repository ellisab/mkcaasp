package utils

type MKCaaSPCfg struct {
	Platform    string  `json: "platform"`
	Deploy      string  `json: "deploy"`
	Vmware      *VMWare `json: "vmware"`
	Skubaroot   string  `json: "skubaroot"`
	ClusterName string  `json: "clustername"`
}

type VMWare struct {
	GOVC_URL                     string
	GOVC_USERNAME                string
	GOVC_PASSWORD                string `json: "GOVC_PASSWORD"`
	GOVC_INSECURE                int
	VSPHERE_SERVER               string `json: "VSPHERE_SERVER"`
	VSPHERE_USER                 string
	VSPHERE_PASSWORD             string
	VSPHERE_ALLOW_UNVERIFIED_SSL bool
}

type TFOutput struct {
	IP_Load_Balancer *TFTag `json: ip_load_balancer`
	IP_Masters       *TFTag `json: ip_masters`
	IP_Workers       *TFTag `json: ip_workers`
}

type TFOutput_vmware struct {
	IP_Load_Balancer *TFTag `json: ip_load_balancer`
	IP_Masters       *TFTag `json: ip_masters`
	IP_Workers       *TFTag `json: ip_workers`
}

type TFOutput_openstack struct {
	IP_Load_Balancer *TFTagLoadBalancer `json: ip_load_balancer`
	IP_Masters       *TFTag             `json: ip_masters`
	IP_Workers       *TFTag             `json: ip_workers`
}

type TFTagLoadBalancer struct {
	Sensitive bool   `json: sensitive`
	Type      string `json: type`
	Value     string `json: value`
}

type TFTag struct {
	Sensitive bool     `json: sensitive`
	Type      string   `json: type`
	Value     []string `json: value`
}

type TFTag_variable struct {
	Sensitive bool          `json: sensitive`
	Type      string        `json: type`
	Value     StringOrSlice `json: value`
}

type StringOrSlice interface{}

//type ClusterCheck map[string]Node

type SkubaCluster struct {
	ClusterName string
	Diagnosis   map[string]Node
	TF_ostack   *TFOutput_openstack
	TF_vmware   *TFOutput_vmware
	Testdir     string
}

type Node struct {
	IP         string
	NodeName   string
	K8sName    string
	Role       string
	Username   string
	Network    bool
	Port22     bool
	SSH        bool
	ContHealth bool
	PackHealth bool
	RepoHealth bool
	Services   bool
	Systemd    Systemd
	K8sHealth  *K8s
}

type Systemd struct {
	CriticalChain []CriticalChain
	AnalyzeBlame  string
	AllFine       bool
}

type CriticalChain struct {
	Unit      string
	TimeAt    string
	TimeDelay string
}

type K8s struct {
}
