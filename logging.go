package easygo

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//NewLogClient returns a client for writing logs and errors to the mongo instance at mongoURL, under the database Logs, under the collection serviceName
func NewLogClient(mongoURL string, serviceName string) (*LogClient, error) {
	mongoCli, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	err = mongoCli.Connect(context.TODO())
	if err != nil {
		return nil, err
	}
	logCollection := mongoCli.Database("Logs").Collection(serviceName)

	return &LogClient{MongoCli: mongoCli, LogCollection: logCollection}, nil
}

//NewLogClientConnect returns a client for writing logs and errors to the mongo connection all ready established.
func NewLogClientConnect(client *mongo.Client, serviceName string) (*LogClient, error) {
	logCollection := client.Database("Logs").Collection(serviceName)

	return &LogClient{MongoCli: client, LogCollection: logCollection}, nil
}

//WriteLog writes a log to the mongo specified during ConnectLogging()
func (c *LogClient) WriteLog(msg string, metaData map[string]interface{}) error {
	if metaData == nil {
		metaData = map[string]interface{}{}
	}
	curTime := time.Now()
	pc, _, line, _ := runtime.Caller(1)
	locStr := runtime.FuncForPC(pc).Name() + " line " + strconv.Itoa(line)
	newLog := log{Message: msg, Location: locStr, Date: curTime, MetaData: metaData}

	_, err := c.LogCollection.InsertOne(context.Background(), newLog)
	if err != nil {
		return err
	}

	fmt.Println(curTime.Format("2006-01-02 15:04:05"), "| Log |", locStr, "|", msg, "|", fmt.Sprintf("%+v", metaData))

	return nil
}

//WriteErr writes an error log to the mongo specified during ConnectLogging()
func (c *LogClient) WriteErr(err error, Metadata interface{}) error {
	if Metadata == nil {
		Metadata = map[string]interface{}{}
	}
	errStr := err.Error()
	curTime := time.Now()
	pc, _, line, _ := runtime.Caller(1)
	locStr := runtime.FuncForPC(pc).Name() + " line " + strconv.Itoa(line)
	newErrLog := errorLog{Error: errStr, Date: curTime, Location: locStr, MetaData: Metadata}

	_, err = c.LogCollection.InsertOne(context.Background(), newErrLog)
	if err != nil {
		return err
	}

	fmt.Println(curTime.Format("2006-01-02 15:04:05"), "| Error |", locStr, "|", errStr, "|", fmt.Sprintf("%+v", Metadata))

	return nil
}

//MakeMetadata Takes the request and generates a map of its contents to send in the errLog
func MakeMetadata(r *http.Request) map[string]interface{} {
	newMap := map[string]interface{}{}
	rMap := map[string]interface{}{}
	rMap["requestMethod"] = r.Method
	rMap["requestHeader"] = r.Header
	rMap["requestURL"] = r.URL
	rMap["requestURI"] = r.RequestURI
	rMap["requestRemoteAddr"] = r.RemoteAddr
	rMap["requestHost"] = r.Host
	newMap["request"] = rMap
	return newMap
}

type errorLog struct {
	Error    string      `json:"error"  bson:"error"`
	Location string      `json:"location"  bson:"location"`
	Date     time.Time   `json:"date"  bson:"date"`
	MetaData interface{} `json:"metaData"  bson:"metaData"`
}

type log struct {
	Message  string
	Location string      `json:"location"  bson:"location"`
	Date     time.Time   `json:"date"  bson:"date"`
	MetaData interface{} `json:"metaData"  bson:"metaData"`
}

//LogClient stores the pointers to the mongo connections to keep them alive
type LogClient struct {
	MongoCli      *mongo.Client
	LogCollection *mongo.Collection
}
