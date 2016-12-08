# \FleetsApi

All URIs are relative to *https://esi.tech.ccp.is/latest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteFleetsFleetIdMembersMemberId**](FleetsApi.md#DeleteFleetsFleetIdMembersMemberId) | **Delete** /fleets/{fleet_id}/members/{member_id}/ | Kick fleet member
[**DeleteFleetsFleetIdSquadsSquadId**](FleetsApi.md#DeleteFleetsFleetIdSquadsSquadId) | **Delete** /fleets/{fleet_id}/squads/{squad_id}/ | Delete fleet squad
[**DeleteFleetsFleetIdWingsWingId**](FleetsApi.md#DeleteFleetsFleetIdWingsWingId) | **Delete** /fleets/{fleet_id}/wings/{wing_id}/ | Delete fleet wing
[**GetFleetsFleetId**](FleetsApi.md#GetFleetsFleetId) | **Get** /fleets/{fleet_id}/ | Get fleet information
[**GetFleetsFleetIdMembers**](FleetsApi.md#GetFleetsFleetIdMembers) | **Get** /fleets/{fleet_id}/members/ | Get fleet members
[**GetFleetsFleetIdWings**](FleetsApi.md#GetFleetsFleetIdWings) | **Get** /fleets/{fleet_id}/wings/ | Get fleet wings
[**PostFleetsFleetIdMembers**](FleetsApi.md#PostFleetsFleetIdMembers) | **Post** /fleets/{fleet_id}/members/ | Create fleet invitation
[**PostFleetsFleetIdWings**](FleetsApi.md#PostFleetsFleetIdWings) | **Post** /fleets/{fleet_id}/wings/ | Create fleet wing
[**PostFleetsFleetIdWingsWingIdSquads**](FleetsApi.md#PostFleetsFleetIdWingsWingIdSquads) | **Post** /fleets/{fleet_id}/wings/{wing_id}/squads/ | Create fleet squad
[**PutFleetsFleetId**](FleetsApi.md#PutFleetsFleetId) | **Put** /fleets/{fleet_id}/ | Update fleet
[**PutFleetsFleetIdMembersMemberId**](FleetsApi.md#PutFleetsFleetIdMembersMemberId) | **Put** /fleets/{fleet_id}/members/{member_id}/ | Move fleet member
[**PutFleetsFleetIdSquadsSquadId**](FleetsApi.md#PutFleetsFleetIdSquadsSquadId) | **Put** /fleets/{fleet_id}/squads/{squad_id}/ | Rename fleet squad
[**PutFleetsFleetIdWingsWingId**](FleetsApi.md#PutFleetsFleetIdWingsWingId) | **Put** /fleets/{fleet_id}/wings/{wing_id}/ | Rename fleet wing


# **DeleteFleetsFleetIdMembersMemberId**
> DeleteFleetsFleetIdMembersMemberId($fleetId, $memberId, $datasource)

Kick fleet member

Kick a fleet member  ---  Alternate route: `/v1/fleets/{fleet_id}/members/{member_id}/`  Alternate route: `/legacy/fleets/{fleet_id}/members/{member_id}/`  Alternate route: `/dev/fleets/{fleet_id}/members/{member_id}/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **memberId** | **int32**| The character ID of a member in this fleet | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteFleetsFleetIdSquadsSquadId**
> DeleteFleetsFleetIdSquadsSquadId($fleetId, $squadId, $datasource)

Delete fleet squad

Delete a fleet squad, only empty squads can be deleted  ---  Alternate route: `/v1/fleets/{fleet_id}/squads/{squad_id}/`  Alternate route: `/legacy/fleets/{fleet_id}/squads/{squad_id}/`  Alternate route: `/dev/fleets/{fleet_id}/squads/{squad_id}/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **squadId** | **int64**| The squad to delete | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteFleetsFleetIdWingsWingId**
> DeleteFleetsFleetIdWingsWingId($fleetId, $wingId, $datasource)

Delete fleet wing

Delete a fleet wing, only empty wings can be deleted. The wing may contain squads, but the squads must be empty  ---  Alternate route: `/v1/fleets/{fleet_id}/wings/{wing_id}/`  Alternate route: `/legacy/fleets/{fleet_id}/wings/{wing_id}/`  Alternate route: `/dev/fleets/{fleet_id}/wings/{wing_id}/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **wingId** | **int64**| The wing to delete | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetFleetsFleetId**
> GetFleetsFleetIdOk GetFleetsFleetId($fleetId, $datasource)

Get fleet information

Return details about a fleet  ---  Alternate route: `/v1/fleets/{fleet_id}/`  Alternate route: `/legacy/fleets/{fleet_id}/`  Alternate route: `/dev/fleets/{fleet_id}/`   ---  This route is cached for up to 5 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**GetFleetsFleetIdOk**](get_fleets_fleet_id_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetFleetsFleetIdMembers**
> []GetFleetsFleetIdMembers200Ok GetFleetsFleetIdMembers($fleetId, $acceptLanguage, $datasource)

Get fleet members

Return information about fleet members  ---  Alternate route: `/v1/fleets/{fleet_id}/members/`  Alternate route: `/legacy/fleets/{fleet_id}/members/`  Alternate route: `/dev/fleets/{fleet_id}/members/`   ---  This route is cached for up to 5 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **acceptLanguage** | **string**| Language to use in the response | [optional] [default to en]
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetFleetsFleetIdMembers200Ok**](get_fleets_fleet_id_members_200_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetFleetsFleetIdWings**
> []GetFleetsFleetIdWings200Ok GetFleetsFleetIdWings($fleetId, $acceptLanguage, $datasource)

Get fleet wings

Return information about wings in a fleet  ---  Alternate route: `/v1/fleets/{fleet_id}/wings/`  Alternate route: `/legacy/fleets/{fleet_id}/wings/`  Alternate route: `/dev/fleets/{fleet_id}/wings/`   ---  This route is cached for up to 5 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **acceptLanguage** | **string**| Language to use in the response | [optional] [default to en]
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetFleetsFleetIdWings200Ok**](get_fleets_fleet_id_wings_200_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostFleetsFleetIdMembers**
> PostFleetsFleetIdMembers($fleetId, $invitation, $datasource)

Create fleet invitation

Invite a character into the fleet, if a character has a CSPA charge set, it is not possible to invite them to the fleet using ESI  ---  Alternate route: `/v1/fleets/{fleet_id}/members/`  Alternate route: `/legacy/fleets/{fleet_id}/members/`  Alternate route: `/dev/fleets/{fleet_id}/members/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **invitation** | [**PostFleetsFleetIdMembersInvitation**](PostFleetsFleetIdMembersInvitation.md)| Details of the invitation | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostFleetsFleetIdWings**
> PostFleetsFleetIdWingsCreated PostFleetsFleetIdWings($fleetId, $datasource)

Create fleet wing

Create a new wing in a fleet  ---  Alternate route: `/v1/fleets/{fleet_id}/wings/`  Alternate route: `/legacy/fleets/{fleet_id}/wings/`  Alternate route: `/dev/fleets/{fleet_id}/wings/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**PostFleetsFleetIdWingsCreated**](post_fleets_fleet_id_wings_created.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostFleetsFleetIdWingsWingIdSquads**
> PostFleetsFleetIdWingsWingIdSquadsCreated PostFleetsFleetIdWingsWingIdSquads($fleetId, $wingId, $datasource)

Create fleet squad

Create a new squad in a fleet  ---  Alternate route: `/v1/fleets/{fleet_id}/wings/{wing_id}/squads/`  Alternate route: `/legacy/fleets/{fleet_id}/wings/{wing_id}/squads/`  Alternate route: `/dev/fleets/{fleet_id}/wings/{wing_id}/squads/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **wingId** | **int64**| The wing_id to create squad in | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**PostFleetsFleetIdWingsWingIdSquadsCreated**](post_fleets_fleet_id_wings_wing_id_squads_created.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutFleetsFleetId**
> PutFleetsFleetId($fleetId, $newSettings, $datasource)

Update fleet

Update settings about a fleet  ---  Alternate route: `/v1/fleets/{fleet_id}/`  Alternate route: `/legacy/fleets/{fleet_id}/`  Alternate route: `/dev/fleets/{fleet_id}/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **newSettings** | [**PutFleetsFleetIdNewSettings**](PutFleetsFleetIdNewSettings.md)| What to update for this fleet | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutFleetsFleetIdMembersMemberId**
> PutFleetsFleetIdMembersMemberId($fleetId, $memberId, $movement, $datasource)

Move fleet member

Move a fleet member around  ---  Alternate route: `/v1/fleets/{fleet_id}/members/{member_id}/`  Alternate route: `/legacy/fleets/{fleet_id}/members/{member_id}/`  Alternate route: `/dev/fleets/{fleet_id}/members/{member_id}/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **memberId** | **int32**| The character ID of a member in this fleet | 
 **movement** | [**PutFleetsFleetIdMembersMemberIdMovement**](PutFleetsFleetIdMembersMemberIdMovement.md)| Details of the invitation | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutFleetsFleetIdSquadsSquadId**
> PutFleetsFleetIdSquadsSquadId($fleetId, $squadId, $naming, $datasource)

Rename fleet squad

Rename a fleet squad  ---  Alternate route: `/v1/fleets/{fleet_id}/squads/{squad_id}/`  Alternate route: `/legacy/fleets/{fleet_id}/squads/{squad_id}/`  Alternate route: `/dev/fleets/{fleet_id}/squads/{squad_id}/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **squadId** | **int64**| The squad to rename | 
 **naming** | [**PutFleetsFleetIdSquadsSquadIdNaming**](PutFleetsFleetIdSquadsSquadIdNaming.md)| New name of the squad | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutFleetsFleetIdWingsWingId**
> PutFleetsFleetIdWingsWingId($fleetId, $wingId, $naming, $datasource)

Rename fleet wing

Rename a fleet wing  ---  Alternate route: `/v1/fleets/{fleet_id}/wings/{wing_id}/`  Alternate route: `/legacy/fleets/{fleet_id}/wings/{wing_id}/`  Alternate route: `/dev/fleets/{fleet_id}/wings/{wing_id}/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fleetId** | **int64**| ID for a fleet | 
 **wingId** | **int64**| The wing to rename | 
 **naming** | [**PutFleetsFleetIdWingsWingIdNaming**](PutFleetsFleetIdWingsWingIdNaming.md)| New name of the wing | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

