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
	"fmt"
	"strings"
)

type DummyApi struct {
	Configuration Configuration
}

func NewDummyApi() *DummyApi {
	configuration := NewConfiguration()
	return &DummyApi{
		Configuration: *configuration,
	}
}

func NewDummyApiWithBasePath(basePath string) *DummyApi {
	configuration := NewConfiguration()
	configuration.BasePath = basePath

	return &DummyApi{
		Configuration: *configuration,
	}
}

/**
 * Get character wallet journal
 * Returns the most recent 50 entries for the characters wallet journal. Optionally, takes an argument with a reference ID, and returns the prior 50 entries from the journal.  ---  Alternate route: &#x60;/v1/characters/{character_id}/wallets/journal/&#x60;  Alternate route: &#x60;/legacy/characters/{character_id}/wallets/journal/&#x60;  Alternate route: &#x60;/dev/characters/{character_id}/wallets/journal/&#x60; 
 *
 * @param characterId An EVE character ID
 * @param lastSeenId A journal reference ID to paginate from
 * @param datasource The server name you would like data from
 * @return []GetCharactersCharacterIdWalletsJournal200Ok
 */
func (a DummyApi) GetCharactersCharacterIdWalletsJournal(characterId int32, lastSeenId int64, datasource string) ([]GetCharactersCharacterIdWalletsJournal200Ok, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/characters/{character_id}/wallets/journal/"
	path = strings.Replace(path, "{"+"character_id"+"}", fmt.Sprintf("%v", characterId), -1)


	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte
	// authentication '(evesso)' required
	// oauth required
	if a.Configuration.AccessToken != ""{
		headerParams["Authorization"] =  "Bearer " + a.Configuration.AccessToken
	}
	// add default headers if any
	for key := range a.Configuration.DefaultHeader {
		headerParams[key] = a.Configuration.DefaultHeader[key]
	}
		queryParams.Add("last_seen_id", a.Configuration.APIClient.ParameterToString(lastSeenId, ""))
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
	var successPayload = new([]GetCharactersCharacterIdWalletsJournal200Ok)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Get wallet transactions
 * Gets the 50 most recent transactions in a characters wallet. Optionally, takes an argument with a transaction ID, and returns the prior 50 transactions  ---  Alternate route: &#x60;/v1/characters/{character_id}/wallets/transactions/&#x60;  Alternate route: &#x60;/legacy/characters/{character_id}/wallets/transactions/&#x60;  Alternate route: &#x60;/dev/characters/{character_id}/wallets/transactions/&#x60; 
 *
 * @param characterId An EVE character ID
 * @param datasource The server name you would like data from
 * @return []GetCharactersCharacterIdWalletsTransactions200Ok
 */
func (a DummyApi) GetCharactersCharacterIdWalletsTransactions(characterId int32, datasource string) ([]GetCharactersCharacterIdWalletsTransactions200Ok, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/characters/{character_id}/wallets/transactions/"
	path = strings.Replace(path, "{"+"character_id"+"}", fmt.Sprintf("%v", characterId), -1)


	headerParams := make(map[string]string)
	queryParams := url.Values{}
	formParams := make(map[string]string)
	var postBody interface{}
	var fileName string
	var fileBytes []byte
	// authentication '(evesso)' required
	// oauth required
	if a.Configuration.AccessToken != ""{
		headerParams["Authorization"] =  "Bearer " + a.Configuration.AccessToken
	}
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
	var successPayload = new([]GetCharactersCharacterIdWalletsTransactions200Ok)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return *successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Dummy Endpoint, Please Ignore
 * Dummy  ---  Alternate route: &#x60;/v1/corporations/{corporation_id}/assets/&#x60;  Alternate route: &#x60;/legacy/corporations/{corporation_id}/assets/&#x60;  Alternate route: &#x60;/dev/corporations/{corporation_id}/assets/&#x60; 
 *
 * @param corporationId Corporation id of the target corporation
 * @param datasource The server name you would like data from
 * @return void
 */
func (a DummyApi) GetCorporationsCorporationIdAssets(corporationId int32, datasource string) (*APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/corporations/{corporation_id}/assets/"
	path = strings.Replace(path, "{"+"corporation_id"+"}", fmt.Sprintf("%v", corporationId), -1)


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

	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return NewAPIResponse(httpResponse.RawResponse), err
	}

	return NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Dummy Endpoint, Please Ignore
 * Dummy  ---  Alternate route: &#x60;/v1/corporations/{corporation_id}/assets/{asset_id}/logs/&#x60;  Alternate route: &#x60;/legacy/corporations/{corporation_id}/assets/{asset_id}/logs/&#x60;  Alternate route: &#x60;/dev/corporations/{corporation_id}/assets/{asset_id}/logs/&#x60; 
 *
 * @param corporationId Corporation id of the target corporation
 * @param assetId Asset id
 * @param datasource The server name you would like data from
 * @return void
 */
func (a DummyApi) GetCorporationsCorporationIdAssetsAssetIdLogs(corporationId int32, assetId int32, datasource string) (*APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/corporations/{corporation_id}/assets/{asset_id}/logs/"
	path = strings.Replace(path, "{"+"corporation_id"+"}", fmt.Sprintf("%v", corporationId), -1)
	path = strings.Replace(path, "{"+"asset_id"+"}", fmt.Sprintf("%v", assetId), -1)


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

	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return NewAPIResponse(httpResponse.RawResponse), err
	}

	return NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Dummy Endpoint, Please Ignore
 * Dummy  ---  Alternate route: &#x60;/v1/corporations/{corporation_id}/bookmarks/&#x60;  Alternate route: &#x60;/legacy/corporations/{corporation_id}/bookmarks/&#x60;  Alternate route: &#x60;/dev/corporations/{corporation_id}/bookmarks/&#x60; 
 *
 * @param corporationId An EVE corporation ID
 * @param datasource The server name you would like data from
 * @return void
 */
func (a DummyApi) GetCorporationsCorporationIdBookmarks(corporationId int32, datasource string) (*APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/corporations/{corporation_id}/bookmarks/"
	path = strings.Replace(path, "{"+"corporation_id"+"}", fmt.Sprintf("%v", corporationId), -1)


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

	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return NewAPIResponse(httpResponse.RawResponse), err
	}

	return NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Dummy Endpoint, Please Ignore
 * Dummy  ---  Alternate route: &#x60;/v1/corporations/{corporation_id}/bookmarks/folders/&#x60;  Alternate route: &#x60;/legacy/corporations/{corporation_id}/bookmarks/folders/&#x60;  Alternate route: &#x60;/dev/corporations/{corporation_id}/bookmarks/folders/&#x60; 
 *
 * @param corporationId An EVE corporation ID
 * @param datasource The server name you would like data from
 * @return void
 */
func (a DummyApi) GetCorporationsCorporationIdBookmarksFolders(corporationId int32, datasource string) (*APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/corporations/{corporation_id}/bookmarks/folders/"
	path = strings.Replace(path, "{"+"corporation_id"+"}", fmt.Sprintf("%v", corporationId), -1)


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

	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return NewAPIResponse(httpResponse.RawResponse), err
	}

	return NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Dummy Endpoint, Please Ignore
 * Dummy  ---  Alternate route: &#x60;/v1/corporations/{corporation_id}/wallets/&#x60;  Alternate route: &#x60;/legacy/corporations/{corporation_id}/wallets/&#x60;  Alternate route: &#x60;/dev/corporations/{corporation_id}/wallets/&#x60; 
 *
 * @param corporationId An EVE corporation ID
 * @param datasource The server name you would like data from
 * @return void
 */
func (a DummyApi) GetCorporationsCorporationIdWallets(corporationId int32, datasource string) (*APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/corporations/{corporation_id}/wallets/"
	path = strings.Replace(path, "{"+"corporation_id"+"}", fmt.Sprintf("%v", corporationId), -1)


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

	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return NewAPIResponse(httpResponse.RawResponse), err
	}

	return NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Dummy Endpoint, Please Ignore
 * Dummy  ---  Alternate route: &#x60;/v1/corporations/{corporation_id}/wallets/{wallet_id}/journal/&#x60;  Alternate route: &#x60;/legacy/corporations/{corporation_id}/wallets/{wallet_id}/journal/&#x60;  Alternate route: &#x60;/dev/corporations/{corporation_id}/wallets/{wallet_id}/journal/&#x60; 
 *
 * @param corporationId An EVE corporation ID
 * @param walletId Wallet ID
 * @param datasource The server name you would like data from
 * @return void
 */
func (a DummyApi) GetCorporationsCorporationIdWalletsWalletIdJournal(corporationId int32, walletId int32, datasource string) (*APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/corporations/{corporation_id}/wallets/{wallet_id}/journal/"
	path = strings.Replace(path, "{"+"corporation_id"+"}", fmt.Sprintf("%v", corporationId), -1)
	path = strings.Replace(path, "{"+"wallet_id"+"}", fmt.Sprintf("%v", walletId), -1)


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

	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return NewAPIResponse(httpResponse.RawResponse), err
	}

	return NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Dummy Endpoint, Please Ignore
 * Dummy  ---  Alternate route: &#x60;/v1/corporations/{corporation_id}/wallets/{wallet_id}/transactions/&#x60;  Alternate route: &#x60;/legacy/corporations/{corporation_id}/wallets/{wallet_id}/transactions/&#x60;  Alternate route: &#x60;/dev/corporations/{corporation_id}/wallets/{wallet_id}/transactions/&#x60; 
 *
 * @param corporationId An EVE corporation ID
 * @param walletId Wallet ID
 * @param datasource The server name you would like data from
 * @return void
 */
func (a DummyApi) GetCorporationsCorporationIdWalletsWalletIdTransactions(corporationId int32, walletId int32, datasource string) (*APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/corporations/{corporation_id}/wallets/{wallet_id}/transactions/"
	path = strings.Replace(path, "{"+"corporation_id"+"}", fmt.Sprintf("%v", corporationId), -1)
	path = strings.Replace(path, "{"+"wallet_id"+"}", fmt.Sprintf("%v", walletId), -1)


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

	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return NewAPIResponse(httpResponse.RawResponse), err
	}

	return NewAPIResponse(httpResponse.RawResponse), err
}

/**
 * Get planet information
 * Information on a planet  ---  Alternate route: &#x60;/v1/universe/planets/{planet_id}/&#x60;  Alternate route: &#x60;/legacy/universe/planets/{planet_id}/&#x60;  Alternate route: &#x60;/dev/universe/planets/{planet_id}/&#x60; 
 *
 * @param planetId An Eve planet ID
 * @param datasource The server name you would like data from
 * @return *GetUniversePlanetsPlanetIdOk
 */
func (a DummyApi) GetUniversePlanetsPlanetId(planetId int32, datasource string) (*GetUniversePlanetsPlanetIdOk, *APIResponse, error) {

	var httpMethod = "Get"
	// create path and map variables
	path := a.Configuration.BasePath + "/universe/planets/{planet_id}/"
	path = strings.Replace(path, "{"+"planet_id"+"}", fmt.Sprintf("%v", planetId), -1)


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
	var successPayload = new(GetUniversePlanetsPlanetIdOk)
	httpResponse, err := a.Configuration.APIClient.CallAPI(path, httpMethod, postBody, headerParams, queryParams, formParams, fileName, fileBytes)
	if err != nil {
		return successPayload, NewAPIResponse(httpResponse.RawResponse), err
	}
	err = json.Unmarshal(httpResponse.Body(), &successPayload)
	return successPayload, NewAPIResponse(httpResponse.RawResponse), err
}

