package dm

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// Stat is a nested struct in dm response
type Stat struct {
	RcptUniqueClickRate  string `json:"RcptUniqueClickRate" xml:"RcptUniqueClickRate"`
	RequestCount         string `json:"requestCount" xml:"requestCount"`
	TotalNumber          string `json:"TotalNumber" xml:"TotalNumber"`
	UnavailablePercent   string `json:"unavailablePercent" xml:"unavailablePercent"`
	SucceededPercent     string `json:"succeededPercent" xml:"succeededPercent"`
	RcptClickCount       string `json:"RcptClickCount" xml:"RcptClickCount"`
	CreateTime           string `json:"CreateTime" xml:"CreateTime"`
	RcptOpenRate         string `json:"RcptOpenRate" xml:"RcptOpenRate"`
	RcptUniqueClickCount string `json:"RcptUniqueClickCount" xml:"RcptUniqueClickCount"`
	UnavailableCount     string `json:"unavailableCount" xml:"unavailableCount"`
	SuccessCount         string `json:"successCount" xml:"successCount"`
	RcptClickRate        string `json:"RcptClickRate" xml:"RcptClickRate"`
	RcptOpenCount        string `json:"RcptOpenCount" xml:"RcptOpenCount"`
	RcptUniqueOpenCount  string `json:"RcptUniqueOpenCount" xml:"RcptUniqueOpenCount"`
	FaildCount           string `json:"faildCount" xml:"faildCount"`
	RcptUniqueOpenRate   string `json:"RcptUniqueOpenRate" xml:"RcptUniqueOpenRate"`
}
