package nomad

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alexebird/tableme/tableme"
	"github.com/fatih/color"
	"github.com/hashicorp/nomad/api"
)

type ByCreateTime []*api.Allocation

func (s ByCreateTime) Len() int {
	return len(s)
}
func (s ByCreateTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByCreateTime) Less(i, j int) bool {
	ti := s[i].CreateTime
	tj := s[j].CreateTime
	// reverse sort
	return tj < ti
}

func extractSeconds(nstime int64) int64 {
	return int64(float64(nstime) / math.Pow(10.0, 9.0))
}

func extractNanoseconds(nstime int64) int64 {
	secs := extractSeconds(nstime)
	// convert to seconds and back to round the nanoseconds value to the nearest second
	nanos := int64(float64(secs) * math.Pow(10.0, 9.0))

	return nstime - nanos
}

func addrAsString(net *api.NetworkResource) string {
	ports := make([]api.Port, 0)
	ports = append(ports, net.ReservedPorts...)
	ports = append(ports, net.DynamicPorts...)
	portStrs := make([]string, 0)

	for _, port := range ports {
		netstr := strings.Join([]string{net.IP, strconv.FormatInt(int64(port.Value), 10)}, ":")
		netstr = fmt.Sprintf("%s(%s)", netstr, port.Label)
		portStrs = append(portStrs, netstr)
	}

	return strings.Join(portStrs, ",")
}

func taskAddr(alloc *api.Allocation, taskName string) string {
	taskResources := alloc.TaskResources
	task := taskResources[taskName]
	networks := task.Networks
	addrs := make([]string, 0)

	for _, net := range networks {
		addrs = append(addrs, addrAsString(net))
	}

	return strings.Join(addrs, ",")
}

func shortUUID(longuuid string) string {
	return strings.Split(longuuid, "-")[0]
}

func PrintTasksTableLong(allocs []*api.Allocation) {
	sort.Sort(ByCreateTime(allocs))

	headers := []string{
		"JOB", "GROUP", "TASK", "ASTATUS", "TSTATE", "TFAILED", "JTYPE", "ADDR", "ALLOCID", "NODEID", "EVALID", "CREATED",
	}

	records := make([][]string, 0)

	for _, alloc := range allocs {
		allocID := alloc.ID
		evalID := alloc.EvalID
		nodeID := alloc.NodeID
		jobID := alloc.Job.ID
		jobType := alloc.Job.Type
		taskStates := alloc.TaskStates
		clientStatus := alloc.ClientStatus
		created := time.Unix(extractSeconds(alloc.CreateTime), extractNanoseconds(alloc.CreateTime))

		for _, taskGroup := range alloc.Job.TaskGroups {
			for _, task := range taskGroup.Tasks {
				taskState := taskStates[task.Name]
				addr := taskAddr(alloc, task.Name)

				rec := []string{
					tableme.StringifyStringPtr(jobID),
					tableme.StringifyStringPtr(taskGroup.Name),
					tableme.StringifyStringPtr(&task.Name),
					tableme.StringifyString(clientStatus),
					tableme.StringifyString(taskState.State),
					tableme.StringifyBool(taskState.Failed),
					tableme.StringifyStringPtr(jobType),
					tableme.StringifyString(addr),
					tableme.StringifyString(allocID),
					tableme.StringifyString(nodeID),
					tableme.StringifyString(evalID),
					tableme.StringifyString(created.Format(time.RFC3339)),
				}

				records = append(records, rec)
			}
		}
	}

	err := tableme.TableMe(headers, records)
	if err != nil {
		os.Exit(1)
	}
}

func colorizeTaskFailed(str string) string {
	red := color.New(color.FgRed).SprintFunc()
	if str == "true" {
		return fmt.Sprintf("%s", red(str))
	} else {
		return str
	}
}

func PrintTasksTableShort(allocs []*api.Allocation) {
	sort.Sort(ByCreateTime(allocs))

	headers := []string{
		"JOB", "TASK", "ASTATUS", "TSTATE", "TFAILED", "JTYPE", "ADDR", "ALLOCID", "NODEID",
	}

	records := make([][]string, 0)

	for _, alloc := range allocs {
		allocID := alloc.ID
		nodeID := alloc.NodeID
		jobID := alloc.Job.ID
		jobType := alloc.Job.Type
		taskStates := alloc.TaskStates
		clientStatus := alloc.ClientStatus

		for _, taskGroup := range alloc.Job.TaskGroups {
			for _, task := range taskGroup.Tasks {
				taskState := taskStates[task.Name]
				taskFailed := colorizeTaskFailed(strconv.FormatBool(taskState.Failed))
				addr := taskAddr(alloc, task.Name)

				rec := []string{
					tableme.StringifyStringPtr(jobID),
					tableme.StringifyString(task.Name),
					tableme.StringifyString(clientStatus),
					tableme.StringifyString(taskState.State),
					tableme.StringifyString(taskFailed),
					tableme.StringifyStringPtr(jobType),
					tableme.StringifyStringPtr(&addr),
					tableme.StringifyString(shortUUID(allocID)),
					tableme.StringifyString(shortUUID(nodeID)),
				}

				records = append(records, rec)
			}
		}
	}

	err := tableme.TableMe(headers, records)
	if err != nil {
		os.Exit(1)
	}
}
