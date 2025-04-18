package monitor

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	service     string
	interval    int
	logOutput   bool
	targetURL   string // Flag for URL monitoring
	processID   int    // Flag for process monitoring by PID
	processName string // Flag for process monitoring by name
	host        string // Flag for host/IP pinging
	kubeconfig  string // New flag for kubeconfig path
)

// MonitorCmd represents the monitor command
var MonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor services, applications, hosts, k8s clusters, or system resources",
	Long: `Monitor command helps you track the status and performance of web URLs, 
local processes, system resources (CPU/memory), check host reachability, 
or verify Kubernetes cluster availability in real-time.`,
	Run: func(cmd *cobra.Command, args []string) {
		if logOutput {
			log.Printf("Starting monitoring task: %s", service)
			log.Printf("Check interval: %d seconds", interval)
		} else {
			fmt.Println("Silent mode enabled. No logs will be printed.")
		}

		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			switch service {
			case "web":
				checkURL(targetURL)
			case "process":
				checkProcess(processID, processName)
			case "cpu":
				checkCPU()
			case "memory":
				checkMemory()
			case "ping":
				checkHost(host)
			case "k8s":
				checkKubernetes()
			case "all":
				log.Println("--- Monitoring All ---")
				checkURL(targetURL)                  // Requires --url
				checkProcess(processID, processName) // Requires --pid or --process-name
				checkHost(host)                      // Requires --host
				checkKubernetes()                    // Requires kubeconfig access
				checkCPU()
				checkMemory()
				log.Println("----------------------")
			default:
				log.Printf("Unknown service type: %s. Supported types: web, process, cpu, memory, ping, k8s, all", service)
				return // Exit if service type is unknown
			}
		}
	},
}

func checkURL(url string) {
	if url == "" {
		log.Println("URL monitoring requires the --url flag to be set.")
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error checking URL %s: %v", url, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("URL %s is reachable (Status: %s)", url, resp.Status)
	} else {
		log.Printf("URL %s returned status: %s", url, resp.Status)
	}
}

func checkProcess(pid int, name string) {
	var p *process.Process
	var err error

	if pid != 0 {
		p, err = process.NewProcess(int32(pid))
		if err != nil {
			log.Printf("Error finding process by PID %d: %v", pid, err)
			return
		}
	} else if name != "" {
		processes, err := process.Processes()
		if err != nil {
			log.Printf("Error listing processes: %v", err)
			return
		}
		found := false
		for _, proc := range processes {
			procName, _ := proc.Name()
			if procName == name {
				p = proc
				found = true
				break
			}
		}
		if !found {
			log.Printf("Process with name '%s' not found.", name)
			return
		}
	} else {
		log.Println("Process monitoring requires either --pid or --process-name flag.")
		return
	}

	isRunning, err := p.IsRunning()
	if err != nil {
		log.Printf("Error checking status for process %d/%s: %v", pid, name, err)
		return
	}

	if isRunning {
		procName, _ := p.Name()
		log.Printf("Process %s (PID: %d) is running.", procName, p.Pid)
	} else {
		log.Printf("Process %d/%s is not running.", pid, name)
	}
}

func checkCPU() {
	percentages, err := cpu.Percent(time.Second, false) // Get overall CPU percentage
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
		return
	}
	if len(percentages) > 0 {
		log.Printf("CPU Usage: %.2f%%", percentages[0])
	}
}

func checkMemory() {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Error getting memory usage: %v", err)
		return
	}
	log.Printf("Memory Usage: %.2f%% (Used: %v / Total: %v)", vmStat.UsedPercent, formatBytes(vmStat.Used), formatBytes(vmStat.Total))
}

func checkHost(hostAddress string) {
	if hostAddress == "" {
		log.Println("Host/IP monitoring requires the --host flag to be set.")
		return
	}
	ports := []string{"80", "443"} // Check common web ports
	reachable := false
	for _, port := range ports {
		address := net.JoinHostPort(hostAddress, port)
		conn, err := net.DialTimeout("tcp", address, 2*time.Second) // Short timeout
		if err == nil {
			conn.Close()
			log.Printf("Host %s is reachable (checked port %s)", hostAddress, port)
			reachable = true
			break // Stop checking if reachable on one port
		}
	}

	if !reachable {
		log.Printf("Host %s appears unreachable (checked ports: %v)", hostAddress, ports)
	}
}

func checkKubernetes() {
	configPath := kubeconfig
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Error getting user home directory: %v", err)
			configPath = filepath.Join(".kube", "config")
		} else {
			configPath = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Printf("Error building kubeconfig from %s: %v", configPath, err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error creating Kubernetes clientset: %v", err)
		return
	}

	_, err = clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{Limit: 1})
	if err != nil {
		log.Printf("Error connecting to Kubernetes cluster (%s): %v", config.Host, err)
		return
	}

	log.Printf("Kubernetes cluster (%s) is reachable.", config.Host)
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func init() {
	MonitorCmd.Flags().StringVarP(&service, "service", "s", "all", "Service to monitor (web, process, cpu, memory, ping, k8s, all)")
	MonitorCmd.Flags().IntVarP(&interval, "interval", "i", 30, "Check interval in seconds")
	MonitorCmd.Flags().BoolVarP(&logOutput, "log", "l", true, "Log output to console (ignored in silent mode)")

	MonitorCmd.Flags().StringVar(&targetURL, "url", "", "URL to monitor for 'web' service")
	MonitorCmd.Flags().IntVar(&processID, "pid", 0, "Process ID (PID) to monitor for 'process' service")
	MonitorCmd.Flags().StringVar(&processName, "process-name", "", "Process name to monitor for 'process' service")
	MonitorCmd.Flags().StringVar(&host, "host", "", "Hostname or IP address to check for 'ping' service")
	MonitorCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file (optional, defaults to ~/.kube/config)")

	log.SetOutput(os.Stderr)
}
