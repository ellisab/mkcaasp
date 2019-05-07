package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"mkcaasp/utils"
	"os"
)

const (
	command  = "terraform %s -var auth_url=$OS_AUTH_URL -var domain_name=$OS_USER_DOMAIN_NAME -var region_name=$OS_REGION_NAME -var project_name=$OS_PROJECT_NAME -var user_name=$OS_USERNAME -var password=$OS_PASSWORD -var-file=terraform.tfvars -auto-approve"
	howtouse = `
			 Make sure you have terraform installed and in $PATH
			 git clone https://github.com/kubic-project/automation.git
			 
			 cd automation

			 put openstack.json in the directories you want to use, for example in caaspDir and/or sesDir
			 
			 openstack.json (should reside in caaspDir and sesDir folders) should look like this:
			 {
				"AuthURL":"https://smtg:5000/v3",
				"RegionName":"Region",
				"ProjectName":"caasp",
				"UserDomainName":"users",
				"IdentityAPIVersion":"3",
				"Interface":"public",
				"Username":"user",
				"Password":"pass",
				"ProjectID":"00000000000000000000000000"
			 }

			 create in the $HOME/automation directory a file named key.json containing just a string "<your key for encrypting password>"
			 in order to put your hashed password in openstack.json, run the 1st time caasp -hash <password>

			 --------------------------------------->>>IMPORTANT!<<<-----------------------------------------------------------
			 run the utility: caasp -repo $HOME/automation -createcaasp -caaspuiinst -createses -action apply -auth openstack.json
			 ------------------------------------------------------------------------------------------------------------------


			 CHECK also the data template in utils/data.go, which sets up the parameters of your cluster in engcloud, 
			 and looks like this:
			 
			var CulsterTempl = 	image_name = "SUSE-CaaS-Platform-3.0-for-OpenStack-Cloud.x86_64-3.0.0-GM.qcow2"
							 	internal_net = "INGSOC-net"
							 	external_net = "floating"
								admin_size = "m1.large"
								master_size = "m1.medium"
								masters = {{.MastCount}}
								worker_size = "m1.medium"
								workers = {{.WorkCount}}
								workers_vol_enabled = 0
								workers_vol_size = 5
								dnsdomain = "testing.qa.caasp.suse.net"
								dnsentry = 0
								stack_name = "INGSOC"`
)

var (
	libvirt       = flag.String("tflibvirt", "", "switch for terraform-libvirt option")
	openstack     = flag.String("auth", "openstack.json", "name of the json file containing openstack variables")
	action        = flag.String("action", "apply", "terraform action to run, example: apply, destroy")
	caasp         = flag.Bool("createcaasp", false, "enables/disables caasp terraform openstack setup")
	ses           = flag.Bool("createses", false, "enables/disables ses terraform openstack setup")
	howto         = flag.Bool("usage", false, "prints usage information")
	caasptfoutput = flag.Bool("caasptfoutput", false, "loads in memory caasp terraform ouput json")
	sestfoutput   = flag.Bool("sestfoutput", false, "loads in memory ses terraform ouput json")
	caaspUIInst   = flag.Bool("caaspuiinst", false, "Configures caasp using Velum UI")
	ostkcmd       = flag.String("ostkcmd", "", "openstack command to run")
	nodes         = flag.String("nodes", "", "what is the cluster starting configuration. How many masters/workers? w1m1 or w3m1 or m3w5")
	append        = flag.String("addnodes", "", `how many more nodes to add, usage m2w2 -2 more masters, 2 more workers
	Argument must be not longer than 4 symbols (e.g. workers or masters with count more than 1 digit cannot be added; like w10m1)`)
	home = flag.String("repo", "automation", "kubic automation repo location")
	pass = flag.String("hash", "password", "the password for cloud to be hashed (and be exported into openstack.json)")
)

const (
	caaspDir = "caasp-openstack-terraform"
	sesDir   = "ses-openstack-terraform"
	output   = "terraform output -json"
)

var Cluster *utils.CaaSPCluster

func main() {
	flag.Parse()
	if *howto {
		fmt.Fprintf(os.Stdout, "%v\n", howtouse)
		os.Exit(0)
	}

	os.Chdir(*home)
	if *ostkcmd != "" {
		out1, out2 := utils.OpenstackCmd(caaspDir, *openstack)
		fmt.Printf("%s\n  %s\n", out1, out2)
	}
	os.Chdir(*home)
	if *pass != "password" {
		utils.Hashinator(*pass, *home, caaspDir)
		utils.Hashinator(*pass, *home, sesDir)
	}
	os.Chdir(*home)
	if *caasp {
		out, _ := utils.CmdRun(caaspDir, *openstack, output)
		a := utils.CAASPOut{}
		err := json.Unmarshal([]byte(out), &a)
		if err != nil {
			log.Fatal(err)
		}
		if *nodes == "" {
			*nodes = "m1w2"
		}
		Cluster = utils.NodesAdder(caaspDir, *nodes, &a, true)
		utils.TfInit(caaspDir)
		utils.CmdRun(caaspDir, *openstack, fmt.Sprintf(command, *action))
	}
	os.Chdir(*home)
	if *caaspUIInst {
		out, _ := utils.CmdRun(caaspDir, *openstack, output)
		a := utils.CAASPOut{}
		err := json.Unmarshal([]byte(out), &a)
		if err != nil {
			log.Fatal(err)
		}
		velumURL := fmt.Sprintf("https://%s.nip.io", a.IPAdminExt.Value)
		fmt.Fprintf(os.Stdout, "Velum warm up time: %2.2f Seconds\n", utils.CheckVelumUp(velumURL))
		utils.CreateAcc(&a)
		utils.FirstSetup(&a)
	}
	os.Chdir(*home)
	if *append != "" {
		out, _ := utils.CmdRun(caaspDir, *openstack, output)
		a := utils.CAASPOut{}
		err := json.Unmarshal([]byte(out), &a)
		if err != nil {
			log.Fatal(err)
		}
		velumURL := fmt.Sprintf("https://%s.nip.io", a.IPAdminExt.Value)
		fmt.Fprintf(os.Stdout, "Velum warm up time: %2.2f Seconds\n", utils.CheckVelumUp(velumURL))
		Cluster := utils.NodesAdder(caaspDir, *append, &a, false)
		utils.CmdRun(caaspDir, *openstack, fmt.Sprintf(command, *action))
		utils.InstallUI(&a, Cluster)
	}
	os.Chdir(*home)
	if *ses {
		utils.TfInit(sesDir)
		utils.CmdRun(sesDir, *openstack, fmt.Sprintf(command, *action))
	}
	os.Chdir(*home)
	if *caasptfoutput {
		utils.CmdRun(caaspDir, *openstack, output)
	}
	os.Chdir(*home)
	if *sestfoutput {
		out, _ := utils.CmdRun(sesDir, *openstack, output)
		a := utils.SESOut{}
		err := json.Unmarshal([]byte(out), &a)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		s := a.K8SSC.Value
		fmt.Println(a.IPAdminExt.Value, a.IPAdminInt.Value, a.IPMonsExt.Value, a.IPMonsExt.Value, a.IPOsdsInt.Value, "\n", a.K8SCS.Value, "\n", fmt.Sprintf("%s", s[0]))
	}

}
