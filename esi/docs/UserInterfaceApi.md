# \UserInterfaceApi

All URIs are relative to *https://esi.tech.ccp.is/latest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PostUiAutopilotWaypoint**](UserInterfaceApi.md#PostUiAutopilotWaypoint) | **Post** /ui/autopilot/waypoint/ | Set Autopilot Waypoint
[**PostUiOpenwindowContract**](UserInterfaceApi.md#PostUiOpenwindowContract) | **Post** /ui/openwindow/contract/ | Open Contract Window
[**PostUiOpenwindowInformation**](UserInterfaceApi.md#PostUiOpenwindowInformation) | **Post** /ui/openwindow/information/ | Open Information Window
[**PostUiOpenwindowMarketdetails**](UserInterfaceApi.md#PostUiOpenwindowMarketdetails) | **Post** /ui/openwindow/marketdetails/ | Open Market Details
[**PostUiOpenwindowNewmail**](UserInterfaceApi.md#PostUiOpenwindowNewmail) | **Post** /ui/openwindow/newmail/ | Open New Mail Window


# **PostUiAutopilotWaypoint**
> PostUiAutopilotWaypoint(ctx, destinationId, clearOtherWaypoints, addToBeginning, optional)
Set Autopilot Waypoint

Set a solar system as autopilot waypoint

---

Alternate route: `/v2/ui/autopilot/waypoint/`

Alternate route: `/dev/ui/autopilot/waypoint/`


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context containing the authentication | nil if no authentication
  **destinationId** | **int64**| The destination to travel to, can be solar system, station or structure&#39;s id | 
  **clearOtherWaypoints** | **bool**| Whether clean other waypoints beforing adding this one | [default to false]
  **addToBeginning** | **bool**| Whether this solar system should be added to the beginning of all waypoints | [default to false]
 **optional** | **map[string]interface{}** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a map[string]interface{}.

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **destinationId** | **int64**| The destination to travel to, can be solar system, station or structure&#39;s id | 
 **clearOtherWaypoints** | **bool**| Whether clean other waypoints beforing adding this one | [default to false]
 **addToBeginning** | **bool**| Whether this solar system should be added to the beginning of all waypoints | [default to false]
 **datasource** | **string**| The server name you would like data from | [default to tranquility]

### Return type

 (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostUiOpenwindowContract**
> PostUiOpenwindowContract(ctx, contractId, optional)
Open Contract Window

Open the contract window inside the client

---

Alternate route: `/v1/ui/openwindow/contract/`

Alternate route: `/legacy/ui/openwindow/contract/`

Alternate route: `/dev/ui/openwindow/contract/`


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context containing the authentication | nil if no authentication
  **contractId** | **int32**| The contract to open | 
 **optional** | **map[string]interface{}** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a map[string]interface{}.

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **contractId** | **int32**| The contract to open | 
 **datasource** | **string**| The server name you would like data from | [default to tranquility]

### Return type

 (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostUiOpenwindowInformation**
> PostUiOpenwindowInformation(ctx, targetId, optional)
Open Information Window

Open the information window for a character, corporation or alliance inside the client

---

Alternate route: `/v1/ui/openwindow/information/`

Alternate route: `/legacy/ui/openwindow/information/`

Alternate route: `/dev/ui/openwindow/information/`


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context containing the authentication | nil if no authentication
  **targetId** | **int32**| The target to open | 
 **optional** | **map[string]interface{}** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a map[string]interface{}.

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **targetId** | **int32**| The target to open | 
 **datasource** | **string**| The server name you would like data from | [default to tranquility]

### Return type

 (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostUiOpenwindowMarketdetails**
> PostUiOpenwindowMarketdetails(ctx, typeId, optional)
Open Market Details

Open the market details window for a specific typeID inside the client

---

Alternate route: `/v1/ui/openwindow/marketdetails/`

Alternate route: `/legacy/ui/openwindow/marketdetails/`

Alternate route: `/dev/ui/openwindow/marketdetails/`


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context containing the authentication | nil if no authentication
  **typeId** | **int32**| The item type to open in market window | 
 **optional** | **map[string]interface{}** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a map[string]interface{}.

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **typeId** | **int32**| The item type to open in market window | 
 **datasource** | **string**| The server name you would like data from | [default to tranquility]

### Return type

 (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostUiOpenwindowNewmail**
> PostUiOpenwindowNewmail(ctx, newMail, optional)
Open New Mail Window

Open the New Mail window, according to settings from the request if applicable

---

Alternate route: `/v1/ui/openwindow/newmail/`

Alternate route: `/legacy/ui/openwindow/newmail/`

Alternate route: `/dev/ui/openwindow/newmail/`


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context containing the authentication | nil if no authentication
  **newMail** | [**PostUiOpenwindowNewmailNewMail**](PostUiOpenwindowNewmailNewMail.md)| The details of mail to create | 
 **optional** | **map[string]interface{}** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a map[string]interface{}.

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **newMail** | [**PostUiOpenwindowNewmailNewMail**](PostUiOpenwindowNewmailNewMail.md)| The details of mail to create | 
 **datasource** | **string**| The server name you would like data from | [default to tranquility]

### Return type

 (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

