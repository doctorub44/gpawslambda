package main

import (
	awslambda "github.com/aws/aws-lambda-go/lambda"
	graphproc "github.com/doctorub44/graphprocessor/graphproc"
)

//MyEvent : JSON message  sent to the Lambda
type MyEvent struct {
	Param string `json:"param"`
	Data  string `json:"data"`
}

//MyResponse :  JSON message returned by the Lambda
type MyResponse struct {
	Message string `json:"message:"`
	Data    string `json:"data"`
}

var graph graphproc.Graphline

func main() {
	awslambda.Start(EventHandler)
}

func init() {
	graph = graphproc.NewGraphline()
	graph.RegisterStage(graphproc.NormalIOC)
	graph.RegisterStage(graphproc.UrlIOC)
	graph.RegisterStage(graphproc.Ipv4IOC)
	graph.RegisterStage(graphproc.Ipv6IOC)
	graph.RegisterStage(graphproc.Md5IOC)
	graph.RegisterStage(graphproc.Sha1IOC)
	graph.RegisterStage(graphproc.Sha256IOC)
	graph.RegisterStage(graphproc.FilterWhiteListIOC)
	graph.RegisterStage(graphproc.IOCDataToJson)
	graph.RegisterStage(graphproc.IOCtoData)
	graph.RegisterStage(graphproc.AWSDownloadS3Bucket)
	graph.RegisterStage(graphproc.SelectFields)
	graph.RegisterStage(graphproc.CutFields)
}

func EventHandler(event MyEvent) (MyResponse, error) {
	var err error = nil
	var response MyResponse
	var gnames []string

	gnames, err = graph.Sequence(event.Param)
	if err != nil {
		response.Message = "Error: unable to create graphline [" + err.Error() + "]"
	} else {
		payload := new(graphproc.Payload)
		payload.Raw = make([]byte, 0, 2048)
		payload.Raw = append(payload.Raw[:0], []byte(event.Data)...)
		if err = graph.Execute(gnames[0], payload); err == nil {
			response.Message = "Success"
			response.Data = string(payload.Raw)
		} else {
			response.Message = "Error: [" + err.Error() + "]"
		}
	}

	return response, err
}
