package main

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

func main() {
	if err := invoke(); err != nil {
		log.Fatalln(err)
	}
}

type Pipeline struct {
	name                        string
	status                      string
	lastModifiedTime            time.Time
	lastPipelineExecutionStatus string
	lastPipelineExecutionTime   time.Time
}

func GetPipelines() (*[]Pipeline, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sagemaker.New(sess)

	listPipelinesOutput, err := svc.ListPipelines(nil)
	if err != nil {
		return nil, err
	}

	pipelines := []Pipeline{}

	for _, pipline := range listPipelinesOutput.PipelineSummaries {
		name := pipline.PipelineName
		describePipelineInput := sagemaker.DescribePipelineInput{
			PipelineName: name,
		}
		output, err := svc.DescribePipeline(&describePipelineInput)
		if err != nil {
			return nil, err
		}

		pipelineName := *output.PipelineName
		pipelineStatus := *output.PipelineStatus

		pipeline := Pipeline{
			name:             pipelineName,
			status:           pipelineStatus,
			lastModifiedTime: output.LastModifiedTime.Local(),
		}

		sortBy := "CreationTime"
		sortOrder := "Descending"
		maxResults := int64(1)
		listPipelineExecutionsInput := &sagemaker.ListPipelineExecutionsInput{
			PipelineName: &pipelineName,
			SortBy:       &sortBy,
			SortOrder:    &sortOrder,
			MaxResults:   &maxResults,
		}

		listPipelineExecutionsOutput, err := svc.ListPipelineExecutions(listPipelineExecutionsInput)

		if err != nil {
			return nil, err
		}

		pipelineExecutionSummaries := listPipelineExecutionsOutput.PipelineExecutionSummaries

		if len(pipelineExecutionSummaries) > 0 {
			pipelineExecutionSummary := listPipelineExecutionsOutput.PipelineExecutionSummaries[0]
			pipeline.lastPipelineExecutionStatus = *pipelineExecutionSummary.PipelineExecutionStatus
			pipeline.lastPipelineExecutionTime = pipelineExecutionSummary.StartTime.Local()
		}

		pipelines = append(pipelines, pipeline)
	}

	return &pipelines, nil
}

func invoke() error {
	pipelines, err := GetPipelines()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Name", "PipelineStatus", "LastModifiedTime", "lastPipelineExecutionTime", "ExecutionStatus"})

	rows := []table.Row{}

	for _, pipeline := range *pipelines {
		var executionStatus string
		text.BgBlue.Sprint(pipeline.status)
		if pipeline.lastPipelineExecutionStatus == "Succeeded" {
			executionStatus = text.BgGreen.Sprint(pipeline.lastPipelineExecutionStatus)
		} else if pipeline.lastPipelineExecutionStatus == "Executing" {
			executionStatus = text.BgBlue.Sprint(pipeline.lastPipelineExecutionStatus)
		} else if pipeline.lastPipelineExecutionStatus == "Failed" {
			executionStatus = text.BgRed.Sprint(pipeline.lastPipelineExecutionStatus)
		} else {
			executionStatus = pipeline.lastPipelineExecutionStatus
		}
		rows = append(rows, table.Row{pipeline.name, pipeline.status, pipeline.lastModifiedTime, pipeline.lastPipelineExecutionTime, executionStatus})
	}

	t.AppendRows(rows)
	t.Render()

	return nil
}
