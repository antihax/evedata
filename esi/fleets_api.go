/* 
 * EVE Swagger Interface
 *
 * An OpenAPI for EVE Online
 *
 * OpenAPI spec version: 0.3.1
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

	"encoding/json"
	"fmt"
)

type FleetsApiService service


/**
 * Kick fleet member
 * Kick a fleet member  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/members/{member_id}/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/members/{member_id}/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/members/{member_id}/&#x60; 
 *
 * @param fleetId ID for a fleet 
 * @param memberId The character ID of a member in this fleet 
 * @param datasource(string) The server name you would like data from 
 * @return nil
 */
func (a FleetsApiService) DeleteFleetsFleetIdMembersMemberId(ts TokenSource, fleetId int64, memberId int32, datasource interface{}) ( error) {
	var (
		localVarHttpMethod = strings.ToUpper("Delete")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/members/{member_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"member_id"+"}", fmt.Sprintf("%v", memberId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return err
	 }
	return err
}

/**
 * Delete fleet squad
 * Delete a fleet squad, only empty squads can be deleted  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/squads/{squad_id}/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/squads/{squad_id}/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/squads/{squad_id}/&#x60; 
 *
 * @param fleetId ID for a fleet 
 * @param squadId The squad to delete 
 * @param datasource(string) The server name you would like data from 
 * @return nil
 */
func (a FleetsApiService) DeleteFleetsFleetIdSquadsSquadId(ts TokenSource, fleetId int64, squadId int64, datasource interface{}) ( error) {
	var (
		localVarHttpMethod = strings.ToUpper("Delete")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/squads/{squad_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"squad_id"+"}", fmt.Sprintf("%v", squadId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return err
	 }
	return err
}

/**
 * Delete fleet wing
 * Delete a fleet wing, only empty wings can be deleted. The wing may contain squads, but the squads must be empty  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/wings/{wing_id}/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/wings/{wing_id}/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/wings/{wing_id}/&#x60; 
 *
 * @param fleetId ID for a fleet 
 * @param wingId The wing to delete 
 * @param datasource(string) The server name you would like data from 
 * @return nil
 */
func (a FleetsApiService) DeleteFleetsFleetIdWingsWingId(ts TokenSource, fleetId int64, wingId int64, datasource interface{}) ( error) {
	var (
		localVarHttpMethod = strings.ToUpper("Delete")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/wings/{wing_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"wing_id"+"}", fmt.Sprintf("%v", wingId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return err
	 }
	return err
}

/**
 * Get fleet information
 * Return details about a fleet  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/&#x60;   ---  This route is cached for up to 5 seconds
 *
 * @param fleetId ID for a fleet 
 * @param datasource(string) The server name you would like data from 
 * @return *GetFleetsFleetIdOk
 */
func (a FleetsApiService) GetFleetsFleetId(ts TokenSource, fleetId int64, datasource interface{}) (*GetFleetsFleetIdOk,  error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return nil, err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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
	 var successPayload = new(GetFleetsFleetIdOk)

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return successPayload, err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return successPayload, err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return successPayload, err
	 }

	 defer localVarHttpResponse.Body.Close()
	 if err = json.NewDecoder(localVarHttpResponse.Body).Decode(&successPayload); err != nil {
	 	return nil, err
     }

	return successPayload, err
}

/**
 * Get fleet members
 * Return information about fleet members  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/members/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/members/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/members/&#x60;   ---  This route is cached for up to 5 seconds
 *
 * @param fleetId ID for a fleet 
 * @param acceptLanguage(string) Language to use in the response 
 * @param datasource(string) The server name you would like data from 
 * @return []GetFleetsFleetIdMembers200Ok
 */
func (a FleetsApiService) GetFleetsFleetIdMembers(ts TokenSource, fleetId int64, acceptLanguage interface{}, datasource interface{}) ([]GetFleetsFleetIdMembers200Ok,  error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/members/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(acceptLanguage, "string", "acceptLanguage"); err != nil {
		return nil, err
	}
	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return nil, err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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
	// header params "Accept-Language"
	localVarHeaderParams["Accept-Language"] = a.client.parameterToString(acceptLanguage, "")
	 var successPayload = new([]GetFleetsFleetIdMembers200Ok)

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return *successPayload, err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return *successPayload, err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return *successPayload, err
	 }

	 defer localVarHttpResponse.Body.Close()
	 if err = json.NewDecoder(localVarHttpResponse.Body).Decode(&successPayload); err != nil {
	 	return nil, err
     }

	return *successPayload, err
}

/**
 * Get fleet wings
 * Return information about wings in a fleet  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/wings/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/wings/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/wings/&#x60;   ---  This route is cached for up to 5 seconds
 *
 * @param fleetId ID for a fleet 
 * @param acceptLanguage(string) Language to use in the response 
 * @param datasource(string) The server name you would like data from 
 * @return []GetFleetsFleetIdWings200Ok
 */
func (a FleetsApiService) GetFleetsFleetIdWings(ts TokenSource, fleetId int64, acceptLanguage interface{}, datasource interface{}) ([]GetFleetsFleetIdWings200Ok,  error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/wings/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(acceptLanguage, "string", "acceptLanguage"); err != nil {
		return nil, err
	}
	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return nil, err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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
	// header params "Accept-Language"
	localVarHeaderParams["Accept-Language"] = a.client.parameterToString(acceptLanguage, "")
	 var successPayload = new([]GetFleetsFleetIdWings200Ok)

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return *successPayload, err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return *successPayload, err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return *successPayload, err
	 }

	 defer localVarHttpResponse.Body.Close()
	 if err = json.NewDecoder(localVarHttpResponse.Body).Decode(&successPayload); err != nil {
	 	return nil, err
     }

	return *successPayload, err
}

/**
 * Create fleet invitation
 * Invite a character into the fleet, if a character has a CSPA charge set, it is not possible to invite them to the fleet using ESI  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/members/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/members/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/members/&#x60; 
 *
 * @param fleetId ID for a fleet 
 * @param invitation Details of the invitation 
 * @param datasource(string) The server name you would like data from 
 * @return nil
 */
func (a FleetsApiService) PostFleetsFleetIdMembers(ts TokenSource, fleetId int64, invitation PostFleetsFleetIdMembersInvitation, datasource interface{}) ( error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/members/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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
	 localVarPostBody = &invitation

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return err
	 }
	return err
}

/**
 * Create fleet wing
 * Create a new wing in a fleet  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/wings/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/wings/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/wings/&#x60; 
 *
 * @param fleetId ID for a fleet 
 * @param datasource(string) The server name you would like data from 
 * @return *PostFleetsFleetIdWingsCreated
 */
func (a FleetsApiService) PostFleetsFleetIdWings(ts TokenSource, fleetId int64, datasource interface{}) (*PostFleetsFleetIdWingsCreated,  error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/wings/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return nil, err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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
	 var successPayload = new(PostFleetsFleetIdWingsCreated)

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return successPayload, err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return successPayload, err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return successPayload, err
	 }

	 defer localVarHttpResponse.Body.Close()
	 if err = json.NewDecoder(localVarHttpResponse.Body).Decode(&successPayload); err != nil {
	 	return nil, err
     }

	return successPayload, err
}

/**
 * Create fleet squad
 * Create a new squad in a fleet  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/wings/{wing_id}/squads/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/wings/{wing_id}/squads/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/wings/{wing_id}/squads/&#x60; 
 *
 * @param fleetId ID for a fleet 
 * @param wingId The wing_id to create squad in 
 * @param datasource(string) The server name you would like data from 
 * @return *PostFleetsFleetIdWingsWingIdSquadsCreated
 */
func (a FleetsApiService) PostFleetsFleetIdWingsWingIdSquads(ts TokenSource, fleetId int64, wingId int64, datasource interface{}) (*PostFleetsFleetIdWingsWingIdSquadsCreated,  error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/wings/{wing_id}/squads/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"wing_id"+"}", fmt.Sprintf("%v", wingId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return nil, err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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
	 var successPayload = new(PostFleetsFleetIdWingsWingIdSquadsCreated)

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return successPayload, err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return successPayload, err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return successPayload, err
	 }

	 defer localVarHttpResponse.Body.Close()
	 if err = json.NewDecoder(localVarHttpResponse.Body).Decode(&successPayload); err != nil {
	 	return nil, err
     }

	return successPayload, err
}

/**
 * Update fleet
 * Update settings about a fleet  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/&#x60; 
 *
 * @param fleetId ID for a fleet 
 * @param newSettings What to update for this fleet 
 * @param datasource(string) The server name you would like data from 
 * @return nil
 */
func (a FleetsApiService) PutFleetsFleetId(ts TokenSource, fleetId int64, newSettings PutFleetsFleetIdNewSettings, datasource interface{}) ( error) {
	var (
		localVarHttpMethod = strings.ToUpper("Put")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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
	 localVarPostBody = &newSettings

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return err
	 }
	return err
}

/**
 * Move fleet member
 * Move a fleet member around  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/members/{member_id}/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/members/{member_id}/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/members/{member_id}/&#x60; 
 *
 * @param fleetId ID for a fleet 
 * @param memberId The character ID of a member in this fleet 
 * @param movement Details of the invitation 
 * @param datasource(string) The server name you would like data from 
 * @return nil
 */
func (a FleetsApiService) PutFleetsFleetIdMembersMemberId(ts TokenSource, fleetId int64, memberId int32, movement PutFleetsFleetIdMembersMemberIdMovement, datasource interface{}) ( error) {
	var (
		localVarHttpMethod = strings.ToUpper("Put")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/members/{member_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"member_id"+"}", fmt.Sprintf("%v", memberId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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
	 localVarPostBody = &movement

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return err
	 }
	return err
}

/**
 * Rename fleet squad
 * Rename a fleet squad  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/squads/{squad_id}/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/squads/{squad_id}/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/squads/{squad_id}/&#x60; 
 *
 * @param fleetId ID for a fleet 
 * @param squadId The squad to rename 
 * @param naming New name of the squad 
 * @param datasource(string) The server name you would like data from 
 * @return nil
 */
func (a FleetsApiService) PutFleetsFleetIdSquadsSquadId(ts TokenSource, fleetId int64, squadId int64, naming PutFleetsFleetIdSquadsSquadIdNaming, datasource interface{}) ( error) {
	var (
		localVarHttpMethod = strings.ToUpper("Put")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/squads/{squad_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"squad_id"+"}", fmt.Sprintf("%v", squadId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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
	 localVarPostBody = &naming

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return err
	 }
	return err
}

/**
 * Rename fleet wing
 * Rename a fleet wing  ---  Alternate route: &#x60;/v1/fleets/{fleet_id}/wings/{wing_id}/&#x60;  Alternate route: &#x60;/legacy/fleets/{fleet_id}/wings/{wing_id}/&#x60;  Alternate route: &#x60;/dev/fleets/{fleet_id}/wings/{wing_id}/&#x60; 
 *
 * @param fleetId ID for a fleet 
 * @param wingId The wing to rename 
 * @param naming New name of the wing 
 * @param datasource(string) The server name you would like data from 
 * @return nil
 */
func (a FleetsApiService) PutFleetsFleetIdWingsWingId(ts TokenSource, fleetId int64, wingId int64, naming PutFleetsFleetIdWingsWingIdNaming, datasource interface{}) ( error) {
	var (
		localVarHttpMethod = strings.ToUpper("Put")
		localVarPostBody interface{}
		localVarFileName string
		localVarFileBytes []byte
	)

	// create path and map variables
	localVarPath := "https://esi.tech.ccp.is/latest/fleets/{fleet_id}/wings/{wing_id}/"
	localVarPath = strings.Replace(localVarPath, "{"+"fleet_id"+"}", fmt.Sprintf("%v", fleetId), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"wing_id"+"}", fmt.Sprintf("%v", wingId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if err := a.client.typeCheckParameter(datasource, "string", "datasource"); err != nil {
		return err
	}
	if datasource != nil {
		localVarQueryParams.Add("datasource", a.client.parameterToString(datasource, ""))
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
	 localVarPostBody = &naming

	 r, err := a.client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes, "application/json")
	 if err != nil {
		  return err
	 }

	if ts != nil {
		if t, err := ts.Token(); err != nil {
			return err
		} else if t != nil {
			t.SetAuthHeader(r)
		}
	}

	 localVarHttpResponse, err := a.client.callAPI(r)
	 if err != nil || localVarHttpResponse == nil {
		  return err
	 }
	return err
}

