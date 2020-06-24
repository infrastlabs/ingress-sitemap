package main

import (
    "fmt"
    "time"
    "flag"

    "sort"
    "encoding/json"
    "sync"

    "io/ioutil"
    "net/http"
	"bytes"
	// "os"
    
    "k8s.io/client-go/rest"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/apimachinery/pkg/apis/meta/v1" 

    log "github.com/sirupsen/logrus"
)

var (
    //read kubeconfig
    kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	
	//attach params: SERVER_URL, SYNC_TIME, MATCH_LABEL, KUBECONFIG,
	jumpserverURL      = flag.String("jumpurl", "http://jumpserver", "jumpserver url, default http://jumpserver") 
	jumpserverPushTime = flag.Int("jumptime", 3, "jump sync-push time, default 3")
	jumpserverLabel    = flag.String("jumplabel", "regist-jumpserver/enabled", "matched label to push jumpserver, default regist-jumpserver/enabled")
	
)

type (
	Row  []string
	Rows []Row
	Pods []Pod
)

func (r Rows) Len() int      { return len(r) }
func (r Rows) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Rows) Less(i, j int) bool {
	return fmt.Sprintf("%s", r[i]) < fmt.Sprintf("%s", r[j])
}

type Pod struct {
	Ns      string    `json:"ns"`
    Name    string    `json:"name"`
	Status  string    `json:"status"`
	Ip      string    `json:"ip"`
	Age     string    `json:"age"`
}

func main2() {
	flag.Parse()

    var err error
    var config *rest.Config

    if *kubeconfig != "" {
        config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
    } else {
        //get config when running inside Kubernetes
        config, err = rest.InClusterConfig()
    }

    if err != nil {
        log.Errorln(err.Error())
        return
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Errorln(err.Error())
        return
    }

    var rows Rows
    var ch chan Rows
    for {
        rows = make(Rows, 0)
        ch = make(chan Rows)

        go func() {
            for r := range ch {
                rows = append(rows, r...)
            }

            sort.Sort(rows)
            hostpush(rows)
        }()

        var wg sync.WaitGroup
        wg.Add(1)
        go func() { defer wg.Done(); getPods(ch, clientset) }()
        wg.Wait()
        close(ch)
		
        time.Sleep(time.Duration(*jumpserverPushTime) * time.Second) //(500 * time.Millisecond)
	}
}

func hostpush(rows Rows) {
	var pods Pods
	pods = make(Pods, 0)
	for _, row := range rows {
		var pod Pod
		pod.Ns = string(row[1]) //"/jumpserver/1.json"
		pod.Name = string(row[2])
		pod.Status = string(row[3])
		pod.Ip = string(row[4])
		pod.Age = string(row[5])
		pods = append(pods, pod)
	}
	if bs, err := json.Marshal(pods); err == nil {
		req := bytes.NewBuffer([]byte(bs))

		body_type := "application/json;charset=utf-8"
		resp, _ := http.Post(*jumpserverURL+"/hostpush/batch/", body_type, req)
		body, _ := ioutil.ReadAll(resp.Body)
        log.Info("bodyReturns: ", string(body))
	} else {
        log.Info("err: ", err)
		fmt.Println(err)
	}	
}

func getPods(ch chan Rows, clientset *kubernetes.Clientset) {
	pods, err := clientset.Core().Pods("").List(v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	var rows Rows
	for _, pod := range pods.Items {

		//label
		lb := pod.Labels
		if val, found2 := lb[*jumpserverLabel]; found2 {
			if val == "enabled" {

				var statuses []string
				statuses = append(statuses, string(pod.Status.Phase))
				for _, c := range pod.Status.Conditions {
					if c.Status != "True" {
						continue
					}
					statuses = append(statuses, string(c.Type))
				}
				rows = append(rows, Row{
					("[pod]"),//colorPod
					(pod.ObjectMeta.Namespace),
					pod.ObjectMeta.Name, //(fmt.Sprintf("%v", truncate(pod.ObjectMeta.Name))),
					statuses[0],// (strings.Join(statuses, " ")),
					(pod.Status.PodIP), //pod.Status.HostIP, pod.ObjectMeta.Labels),
					(shortHumanDuration(time.Since(pod.CreationTimestamp.Time))),
				})				
			}
		}

	}
	ch <- rows
}

func shortHumanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < -1 {
		return fmt.Sprintf("<invalid>")
	} else if seconds < 0 {
		return fmt.Sprintf("0s")
	} else if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	} else if minutes := int(d.Minutes()); minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	} else if hours := int(d.Hours()); hours < 24 {
		return fmt.Sprintf("%dh", hours)
	} else if hours < 24*364 {
		return fmt.Sprintf("%dd", hours/24)
	}
	return fmt.Sprintf("%dy", int(d.Hours()/24/365))
}