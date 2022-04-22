package elastic

import (
	"bytes"
	"encoding/json"
	"errors"
	"eurus-backend/foundation"
	"eurus-backend/foundation/log"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ElasticSearchData interface {
	GetPath() string
	GetIndex() string
}

type ElasticSearchDataBase struct { //implements ElasticSearchData
	Path  string `json:"path"`
	Index string `json:"index"`
}

func (me *ElasticSearchDataBase) GetPath() string {
	return me.Path
}

func (me *ElasticSearchDataBase) GetIndex() string {
	return me.Index
}

type ElasticSearch interface {
	IsValid() bool
	GetError() error
	InsertLog(data ElasticSearchData) error
}
type elasticSearchByApi struct {
	url string
}

func NewElasticSearchByApi(hostUrl string) ElasticSearch {
	elasticSearch := new(elasticSearchByApi)
	elasticSearch.url = hostUrl
	return elasticSearch
}

func (me *elasticSearchByApi) InsertLog(data ElasticSearchData) error {

	if data.GetPath() == "" {
		return foundation.NewErrorWithMessage(foundation.InvalidArgument, "Path is empty")
	}
	url := me.url + "/" + data.GetPath()

	if data.GetIndex() != "" {
		url = url + "/" + data.GetIndex()
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		resData, _ := ioutil.ReadAll(res.Body)
		content := string(resData)
		errStr := fmt.Sprintf("HTTP status: %d, content: %s", res.StatusCode, content)
		return errors.New(errStr)
	}
	return nil
}

func (me *elasticSearchByApi) IsValid() bool {
	return true
}

func (me *elasticSearchByApi) GetError() error {
	return nil
}

type elasticSearchByLogFile struct {
	Error      error
	filePath   string
	LoggerName string
	logger     *logrus.Logger
}

type RawFormatter struct {
}

func (me *RawFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

func NewElasticSearchByLogFile(filePath string) ElasticSearch {
	elasticSearch := new(elasticSearchByLogFile)
	elasticSearch.filePath = filePath

	dir, fileName := path.Split(filePath)

	index := strings.Index(fileName, ".")
	var extension string
	if index >= 0 {
		extension = fileName[index:]
		fileName = fileName[:index]
	}

	elasticSearch.LoggerName = fileName
	curTime := time.Now()

	logger, err := log.NewLoggerWithFileNamePattern(fileName, dir, fileName+"_%Y-%m-%dT"+curTime.Format("15_02_01")+extension, filePath, logrus.InfoLevel)

	elasticSearch.Error = err
	elasticSearch.logger = logger
	elasticSearch.logger.SetReportCaller(false)
	if elasticSearch.logger != nil {
		elasticSearch.logger.SetFormatter(&RawFormatter{})

	}
	return elasticSearch
}

func (me *elasticSearchByLogFile) IsValid() bool {
	return me.Error == nil
}

type ElasticSearchFileLog struct {
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
	Path      string      `json:"path"`
	Index     string      `json:"index"`
}

func (me *elasticSearchByLogFile) InsertLog(data ElasticSearchData) error {

	logData := new(ElasticSearchFileLog)
	logData.Timestamp = time.Now().Format(time.RFC3339)
	logData.Data = data
	logData.Index = data.GetIndex()
	logData.Path = data.GetPath()

	logMsg, _ := json.Marshal(logData)
	me.logger.Println(string(logMsg))

	return nil
}

func (me *elasticSearchByLogFile) GetError() error {
	return me.Error
}
