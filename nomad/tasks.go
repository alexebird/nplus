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
	//"github.com/davecgh/go-spew/spew"
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
		"ALLOC", "TASK", "JSTATUS", "ASTATUS", "TSTATE", "TFAILED", "JTYPE", "ADDR", "ALLOCID", "NODEID", "EVALID", "CREATED",
	}

	records := make([][]string, 0)

	for _, alloc := range allocs {
		allocID := alloc.ID
		allocName := alloc.Name
		evalID := alloc.EvalID
		nodeID := alloc.NodeID
		jobType := alloc.Job.Type
		taskStates := alloc.TaskStates
		jobStatus := alloc.Job.Status
		clientStatus := alloc.ClientStatus
		created := time.Unix(extractSeconds(alloc.CreateTime), extractNanoseconds(alloc.CreateTime))
		allocTaskGroup := alloc.TaskGroup

		for _, taskGroup := range alloc.Job.TaskGroups {
			// the job lists all task groups, but we only want the taskgroup associated with this alloc.
			if *taskGroup.Name != allocTaskGroup {
				continue
			}

			for _, task := range taskGroup.Tasks {
				taskState := taskStates[task.Name]
				addr := taskAddr(alloc, task.Name)

				rec := []string{
					tableme.StringifyString(allocName),
					tableme.StringifyString(task.Name),
					tableme.StringifyStringPtr(jobStatus),
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
		"ALLOC", "TASK", "JSTATUS", "ASTATUS", "TSTATE", "TFAILED", "JTYPE", "ADDR", "ALLOCID", "NODEID",
	}

	records := make([][]string, 0)

	for _, alloc := range allocs {
		allocID := alloc.ID
		allocName := alloc.Name
		nodeID := alloc.NodeID
		jobType := alloc.Job.Type
		taskStates := alloc.TaskStates
		clientStatus := alloc.ClientStatus
		jobStatus := alloc.Job.Status
		allocTaskGroup := alloc.TaskGroup

		for _, taskGroup := range alloc.Job.TaskGroups {
			// the job lists all task groups, but we only want the taskgroup associated with this alloc.
			if *taskGroup.Name != allocTaskGroup {
				continue
			}

			for _, task := range taskGroup.Tasks {
				taskState := taskStates[task.Name]
				failedState := taskState.Failed
				stateState := taskState.State
				taskFailed := colorizeTaskFailed(strconv.FormatBool(failedState))
				addr := taskAddr(alloc, task.Name)

				rec := []string{
					tableme.StringifyString(allocName),
					tableme.StringifyString(task.Name),
					tableme.StringifyStringPtr(jobStatus),
					tableme.StringifyString(clientStatus),
					tableme.StringifyString(stateState),
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
