/* 
 * EVE Swagger Interface
 *
 * An OpenAPI for EVE Online
 *
 * OpenAPI spec version: 0.2.6.dev1
 * 
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package esi

import (
	"net/url"
	"encoding/json"
)

type SovereigntyApi struct {
	Configuration Configuration
}

func NewSovereigntyApi() *SovereigntyApi {
	configuration := NewConfiguration()
	return &SovereigntyApi{
		Configuration: *configuration,
	}
}

func NewSovereigntyApiWithBasePath(basePath string) *SovereigntyApi {
	configuration := NewConfiguration()
	configuration.BasePath = basePath

	return &SovereigntyApi{
		Configuration: *configuration,
	}
}

/**
 * List sovereignty campaigns
 * Shows sovereignty data for campaigns.  ---  Alternate route: &#x60;/v1/sovereignty/campaigns/&#x60;  Alternate route: &#x60;/legacy/sovereignty/campaigns/&#x60;  Alternate route: &#x60;/dev/sovereignty/campaigns/&#x60;   ---  This route is cached for up to 5 seconds
 *
 * @param datasource The server name you would like data from
 * @return []GetSovereigntyCampaigns200Ok
 */
func (a SovereigntyApi) GetSovereigntyCampaigns(datasource string) ([]GetSovereigntyCampaigns200Ok, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/sovereignty/campaigns/"


	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte
	// add default headers if any
	for key := range a.Configuration.DefaultHeader {
		headerParams[key] = a.Configuration.DefaultHeader[key]
	}
		queryParams.Add("datasource", a.Configuration.APIClient.ParameterToString(datasource, ""))
	

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{  }

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload = new([]GetSovereigntyCampaigns200Ok)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * List sovereignty structures
 * Shows sovereignty data for structures.  ---  Alternate route: &#x60;/v1/sovereignty/structures/&#x60;  Alternate route: &#x60;/legacy/sovereignty/structures/&#x60;  Alternate route: &#x60;/dev/sovereignty/structures/&#x60;   ---  This route is cached for up to 120 seconds
 *
 * @param datasource The server name you would like data from
 * @return []GetSovereigntyStructures200Ok
 */
func (a SovereigntyApi) GetSovereigntyStructures(datasource string) ([]GetSovereigntyStructures200Ok, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/sovereignty/structures/"


	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte
	// add default headers if any
	for key := range a.Configuration.DefaultHeader {
		headerParams[key] = a.Configuration.DefaultHeader[key]
	}
		queryParams.Add("datasource", a.Configuration.APIClient.ParameterToString(datasource, ""))
	

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{  }

	// set Content-Type header
	localVarHttpContentType := a.Configuration.APIClient.SelectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		headerParams["Content-Type"] = localVarHttpContentType
	}
	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
	}

	// set Accept header
	localVarHttpHeaderAccept := a.Configuration.APIClient.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		headerParams["Accept"] = localVarHttpHeaderAccept
	}
	var successPayload = new([]GetSovereigntyStructures200Ok)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

