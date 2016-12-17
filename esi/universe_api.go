/* 
 * EVE Swagger Interface
 *
 * An OpenAPI for EVE Online
 *
 * OpenAPI spec version: 0.3.2.dev3
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
	"strings"
	"time"
	"errors"
	"golang.org/x/net/context"
	"encoding/json"
	"fmt"
)

var _ context.Context

type UniverseApiService service


/**
 * Get station information
 * Public information on stations  ---  Alternate route: &#x60;/v1/universe/stations/{station_id}/&#x60;  Alternate route: &#x60;/legacy/universe/stations/{station_id}/&#x60;  Alternate route: &#x60;/dev/universe/stations/{station_id}/&#x60;   ---  This route is cached for up to 3600 seconds
 *

 * @param stationId An Eve station ID
 * @param optional (nil or map[string]interface{}) with one or more of:
 *     @param "datasource" (string) The server name you would like data from
 * @return GetUniverseStationsStationIdOk
 */
func (a UniverseApiService) GetUniverseStationsStationId(stationId int32, localVarOptionals map[string]interface{}) (GetUniverseStationsStationIdOk,  time.Time, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	 	successPayload  GetUniverseStationsStationIdOk
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/universe/stations/{station_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"station_id"+"}", fmt.Sprintf("%v", stationId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarTempParam, localVarOk := localVarOptionals["datasource"].(string); localVarOptionals != nil && localVarOk {
		localVarQueryParams.Add("datasource", a.client.parameterToString(localVarTempParam, ""))
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
		}

	// set Accept header
	localVarHttpHeaderAccept := a.client.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}


	 r, err := a.client.prepareRequest(nil, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return successPayload, time.Now(), err
	 }

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return successPayload, time.Now(), err
	 }
	 defer localVarHttpResponse.Body.Close()
	 if localVarHttpResponse.StatusCode >= 300 {
		return successPayload, time.Now(), errors.New(localVarHttpResponse.Status)
	 }
	
	if err = json.NewDecoder(localVarHttpResponse.Body).Decode(&successPayload); err != nil {
	 	return successPayload, time.Now(), err
	}

	expires := cacheExpires(localVarHttpResponse)
	return successPayload, expires, err
}

/**
 * List all public structures
 * List all public structures  ---  Alternate route: &#x60;/v1/universe/structures/&#x60;  Alternate route: &#x60;/legacy/universe/structures/&#x60;  Alternate route: &#x60;/dev/universe/structures/&#x60;   ---  This route is cached for up to 3600 seconds
 *

 * @param optional (nil or map[string]interface{}) with one or more of:
 *     @param "datasource" (string) The server name you would like data from
 * @return []int64
 */
func (a UniverseApiService) GetUniverseStructures(localVarOptionals map[string]interface{}) ([]int64,  time.Time, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	 	successPayload  []int64
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/universe/structures/"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarTempParam, localVarOk := localVarOptionals["datasource"].(string); localVarOptionals != nil && localVarOk {
		localVarQueryParams.Add("datasource", a.client.parameterToString(localVarTempParam, ""))
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
		}

	// set Accept header
	localVarHttpHeaderAccept := a.client.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}


	 r, err := a.client.prepareRequest(nil, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return successPayload, time.Now(), err
	 }

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return successPayload, time.Now(), err
	 }
	 defer localVarHttpResponse.Body.Close()
	 if localVarHttpResponse.StatusCode >= 300 {
		return successPayload, time.Now(), errors.New(localVarHttpResponse.Status)
	 }
	
	if err = json.NewDecoder(localVarHttpResponse.Body).Decode(&successPayload); err != nil {
	 	return successPayload, time.Now(), err
	}

	expires := cacheExpires(localVarHttpResponse)
	return successPayload, expires, err
}

/**
 * Get structure information
 * Returns information on requested structure, if you are on the ACL. Otherwise, returns \&quot;Forbidden\&quot; for all inputs.  ---  Alternate route: &#x60;/v1/universe/structures/{structure_id}/&#x60;  Alternate route: &#x60;/legacy/universe/structures/{structure_id}/&#x60;  Alternate route: &#x60;/dev/universe/structures/{structure_id}/&#x60; 
 *
 * @param ctx context.Context Authentication Context 
 * @param structureId An Eve structure ID
 * @param optional (nil or map[string]interface{}) with one or more of:
 *     @param "datasource" (string) The server name you would like data from
 * @return GetUniverseStructuresStructureIdOk
 */
func (a UniverseApiService) GetUniverseStructuresStructureId(ctx context.Context, structureId int64, localVarOptionals map[string]interface{}) (GetUniverseStructuresStructureIdOk,  time.Time, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	 	successPayload  GetUniverseStructuresStructureIdOk
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/universe/structures/{structure_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"structure_id"+"}", fmt.Sprintf("%v", structureId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarTempParam, localVarOk := localVarOptionals["datasource"].(string); localVarOptionals != nil && localVarOk {
		localVarQueryParams.Add("datasource", a.client.parameterToString(localVarTempParam, ""))
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
		}

	// set Accept header
	localVarHttpHeaderAccept := a.client.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}


	 r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return successPayload, time.Now(), err
	 }

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return successPayload, time.Now(), err
	 }
	 defer localVarHttpResponse.Body.Close()
	 if localVarHttpResponse.StatusCode >= 300 {
		return successPayload, time.Now(), errors.New(localVarHttpResponse.Status)
	 }
	
	if err = json.NewDecoder(localVarHttpResponse.Body).Decode(&successPayload); err != nil {
	 	return successPayload, time.Now(), err
	}

	expires := cacheExpires(localVarHttpResponse)
	return successPayload, expires, err
}

/**
 * Get solar system information
 * Information on solar systems  ---  Alternate route: &#x60;/v1/universe/systems/{system_id}/&#x60;  Alternate route: &#x60;/legacy/universe/systems/{system_id}/&#x60;  Alternate route: &#x60;/dev/universe/systems/{system_id}/&#x60;   ---  This route is cached for up to 3600 seconds
 *

 * @param systemId An Eve solar system ID
 * @param optional (nil or map[string]interface{}) with one or more of:
 *     @param "datasource" (string) The server name you would like data from
 * @return GetUniverseSystemsSystemIdOk
 */
func (a UniverseApiService) GetUniverseSystemsSystemId(systemId int32, localVarOptionals map[string]interface{}) (GetUniverseSystemsSystemIdOk,  time.Time, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	 	successPayload  GetUniverseSystemsSystemIdOk
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/universe/systems/{system_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"system_id"+"}", fmt.Sprintf("%v", systemId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarTempParam, localVarOk := localVarOptionals["datasource"].(string); localVarOptionals != nil && localVarOk {
		localVarQueryParams.Add("datasource", a.client.parameterToString(localVarTempParam, ""))
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
		}

	// set Accept header
	localVarHttpHeaderAccept := a.client.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}


	 r, err := a.client.prepareRequest(nil, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return successPayload, time.Now(), err
	 }

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return successPayload, time.Now(), err
	 }
	 defer localVarHttpResponse.Body.Close()
	 if localVarHttpResponse.StatusCode >= 300 {
		return successPayload, time.Now(), errors.New(localVarHttpResponse.Status)
	 }
	
	if err = json.NewDecoder(localVarHttpResponse.Body).Decode(&successPayload); err != nil {
	 	return successPayload, time.Now(), err
	}

	expires := cacheExpires(localVarHttpResponse)
	return successPayload, expires, err
}

/**
 * Get type information
 * Get information on a type  ---  Alternate route: &#x60;/v1/universe/types/{type_id}/&#x60;  Alternate route: &#x60;/legacy/universe/types/{type_id}/&#x60;  Alternate route: &#x60;/dev/universe/types/{type_id}/&#x60;   ---  This route is cached for up to 3600 seconds
 *

 * @param typeId An Eve item type ID
 * @param optional (nil or map[string]interface{}) with one or more of:
 *     @param "datasource" (string) The server name you would like data from
 * @return GetUniverseTypesTypeIdOk
 */
func (a UniverseApiService) GetUniverseTypesTypeId(typeId int32, localVarOptionals map[string]interface{}) (GetUniverseTypesTypeIdOk,  time.Time, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	 	successPayload  GetUniverseTypesTypeIdOk
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/universe/types/{type_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"type_id"+"}", fmt.Sprintf("%v", typeId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarTempParam, localVarOk := localVarOptionals["datasource"].(string); localVarOptionals != nil && localVarOk {
		localVarQueryParams.Add("datasource", a.client.parameterToString(localVarTempParam, ""))
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
		}

	// set Accept header
	localVarHttpHeaderAccept := a.client.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}


	 r, err := a.client.prepareRequest(nil, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return successPayload, time.Now(), err
	 }

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return successPayload, time.Now(), err
	 }
	 defer localVarHttpResponse.Body.Close()
	 if localVarHttpResponse.StatusCode >= 300 {
		return successPayload, time.Now(), errors.New(localVarHttpResponse.Status)
	 }
	
	if err = json.NewDecoder(localVarHttpResponse.Body).Decode(&successPayload); err != nil {
	 	return successPayload, time.Now(), err
	}

	expires := cacheExpires(localVarHttpResponse)
	return successPayload, expires, err
}

/**
 * Get names and categories for a set of ID&#39;s
 * Resolve a set of IDs to names and categories. Supported ID&#39;s for resolving are: Characters, Corporations, Alliances, Stations, Solar Systems, Constellations, Regions, Types.  ---  Alternate route: &#x60;/v1/universe/names/&#x60;  Alternate route: &#x60;/legacy/universe/names/&#x60; 
 *

 * @param ids The ids to resolve
 * @param optional (nil or map[string]interface{}) with one or more of:
 *     @param "datasource" (string) The server name you would like data from
 * @return []PostUniverseNames200Ok
 */
func (a UniverseApiService) PostUniverseNames(ids PostUniverseNamesIds, localVarOptionals map[string]interface{}) ([]PostUniverseNames200Ok,  time.Time, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	 	successPayload  []PostUniverseNames200Ok
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/universe/names/"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarTempParam, localVarOk := localVarOptionals["datasource"].(string); localVarOptionals != nil && localVarOk {
		localVarQueryParams.Add("datasource", a.client.parameterToString(localVarTempParam, ""))
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{
		"application/json",
		}

	// set Accept header
	localVarHttpHeaderAccept := a.client.SelectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	// body params
	 localVarPostBody = &ids


	 r, err := a.client.prepareRequest(nil, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return successPayload, time.Now(), err
	 }

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return successPayload, time.Now(), err
	 }
	 defer localVarHttpResponse.Body.Close()
	 if localVarHttpResponse.StatusCode >= 300 {
		return successPayload, time.Now(), errors.New(localVarHttpResponse.Status)
	 }
	
	if err = json.NewDecoder(localVarHttpResponse.Body).Decode(&successPayload); err != nil {
	 	return successPayload, time.Now(), err
	}

	expires := cacheExpires(localVarHttpResponse)
	return successPayload, expires, err
}

