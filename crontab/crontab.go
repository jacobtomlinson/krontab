package crontab

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"

	batchv1beta1 "k8s.io/api/batch/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	batchv1beta1Types "k8s.io/client-go/kubernetes/typed/batch/v1beta1"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/jacobtomlinson/krontab/input"
	"github.com/jacobtomlinson/krontab/template"
)

var namespace string
var kubeconfig *string
var clientset *kubernetes.Clientset
var cronjobsClient batchv1beta1Types.CronJobInterface

// KronJob represents a job
type KronJob struct {
	Template string
	Name     string
	Timing   string
	Command  string
}

// Create a new cronjob on the cluster
func (k KronJob) Create() error {
	templateYaml, err := template.GetTemplate(k.Template)
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(templateYaml), nil, nil)
	if err != nil {
		panic(err.Error())
	}
	cronjob := obj.(*batchv1beta1.CronJob)
	cronjob.Spec.JobTemplate.Name = k.Name
	cronjob.Name = k.Name
	cronjob.Spec.Schedule = k.Timing
	cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Name = k.Name
	cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Command = strings.Split(k.Command, " ")
	if cronjob.Annotations == nil {
		cronjob.Annotations = make(map[string]string)
	}
	cronjob.Annotations["krontabTemplate"] = k.Template
	_, err = cronjobsClient.Create(cronjob)
	if err != nil {
		panic(err.Error())
	}
	return err
}

// Exists checks whether a job exists on the cluster
func (k KronJob) Exists() bool {
	exists := false
	for _, existingJob := range ListKronJobs() {
		if k.Name == existingJob.Name {
			exists = true
		}
	}
	return exists
}

// Delete deletes a job from the cluster
func (k KronJob) Delete() error {
	deletePolicy := metav1.DeletePropagationForeground
	err := cronjobsClient.Delete(k.Name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	return err
}

// Update updates a job on the cluster
func (k KronJob) Update() error {
	cronjob, err := cronjobsClient.Get(k.Name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	if cronjob.Annotations["krontabTemplate"] != k.Template {
		// TODO Check new template exists before deleting
		err = k.Delete()
		if err != nil {
			panic(err)
		}
		err = k.Create()
		if err != nil {
			panic(err)
		}
		return nil
	}

	cronjob.Spec.JobTemplate.Name = k.Name
	cronjob.Name = k.Name
	cronjob.Spec.Schedule = k.Timing
	cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Name = k.Name
	cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Command = strings.Split(k.Command, " ")
	_, err = cronjobsClient.Update(cronjob)
	return err
}

// EditCrontab generates the current crontab, allows the user to edit and then applies the changes
func EditCrontab() {
	crontabHeader := `# Welcome to krontab, a crontab like editor for Kubernetes cron jobs

`
	// TODO Finish crontab edit blurb
	rawKrontab := input.UserInput(crontabHeader + BuildCrontab())
	jobs, err := ParseCrontab(rawKrontab)
	for _, job := range jobs {
		if job.Exists() {
			err = job.Update()
			if err != nil {
				panic(err)
			}
		} else {
			err = job.Create()
			if err != nil {
				panic(err)
			}
		}
	}
	existingJobs := ListKronJobs()
	for _, existingJob := range existingJobs {
		found := false
		for _, job := range jobs {
			if existingJob.Name == job.Name {
				found = true
			}
		}
		if !found {
			existingJob.Delete()
		}
	}
}

// ListCrontab generates the current crontab and shows it to the user
func ListCrontab() {
	fmt.Printf(`# Welcome to krontab, a crontab like editor for Kubernetes cron jobs

`) // TODO Finish list blurb
	fmt.Println(BuildCrontab())
}

// ListCronJobs gets a list of Kubernetes CronJob resources
func ListCronJobs() []batchv1beta1.CronJob {
	cronjobs, err := clientset.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return cronjobs.Items
}

// ListKronJobs gets a list of Kubernetes CronJob resources in KronJob format
func ListKronJobs() []KronJob {
	var jobs []KronJob
	cronjobs := ListCronJobs()

	for _, job := range cronjobs {
		jobs = append(jobs, KronJob{
			job.Annotations["krontabTemplate"],
			job.Spec.JobTemplate.Name,
			job.Spec.Schedule,
			strings.TrimSpace(strings.Join(job.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Command, " ")),
		})
	}

	return jobs
}

// BuildCrontab constructs a string representation of Kubernetes CronJob resources in a crontab format
func BuildCrontab() string {
	var output []string
	cronjobs := ListKronJobs()

	var jobTemplateGroups map[string][]KronJob
	jobTemplateGroups = make(map[string][]KronJob)

	for _, job := range cronjobs {
		jobTemplateGroups[job.Template] = append(jobTemplateGroups[job.Template], job)
	}

	for template, jobs := range jobTemplateGroups {
		output = append(output, fmt.Sprintf("\n# template: %s", template))
		for _, job := range jobs {
			output = append(output, fmt.Sprintf("%s %s  # name: %s", job.Timing, job.Command, job.Name))
		}
	}
	return strings.Join(output, "\n")
}

// ParseCrontab reads the crontab and parses it into jobs
func ParseCrontab(crontab string) ([]KronJob, error) {
	scanner := bufio.NewScanner(strings.NewReader(crontab))
	var line string
	var jobs []KronJob
	template := "default"
	for scanner.Scan() {
		line = scanner.Text()
		if isBlankLine(line) {
			continue
		}
		if isComment(line) {
			parsedTemplate, err := parseTemplateYaml(uncomment(line))
			if err == nil {
				template = parsedTemplate
			}
			// TODO Validate that the template exists
			continue
		}
		jobs = append(jobs, parseKronJob(line, template))
	}
	return jobs, nil
}

func isBlankLine(line string) bool {
	line = strings.TrimSpace(line)
	return len(line) <= 0
}

func isComment(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, "#")
}

func parseTemplateYaml(line string) (string, error) {
	type Template struct {
		Template string
	}
	var template Template
	err := yaml.Unmarshal([]byte(line), &template)
	return template.Template, err
}

func parseNameYaml(line string) (string, error) {
	type Name struct {
		Name string
	}
	var name Name
	err := yaml.Unmarshal([]byte(line), &name)
	return name.Name, err
}

func uncomment(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimLeft(line, "#")
	line = strings.TrimSpace(line)
	return line
}

func parseKronJob(line string, template string) KronJob {
	slices := strings.Split(line, "#")
	cronjob := slices[0]
	name := uuid.NewV4().String()
	if len(slices) > 1 {
		config, err := parseNameYaml(slices[1])
		if err == nil {
			name = config
		}
	}
	slices = strings.Split(cronjob, " ")
	timing := strings.Join(slices[:5], " ")
	command := strings.Join(slices[5:], " ")
	return KronJob{
		template,
		name,
		timing,
		command,
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func init() {
	namespaceEnv, envSet := os.LookupEnv("KRONTAB_NAMESPACE")
	if envSet {
		namespace = namespaceEnv
	} else {
		namespace = apiv1.NamespaceDefault
	}

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	newClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	clientset = newClientset
	cronjobsClient = clientset.BatchV1beta1().CronJobs(namespace)
}
