# \UniverseApi

All URIs are relative to *https://esi.tech.ccp.is/latest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetUniversePlanetsPlanetId**](UniverseApi.md#GetUniversePlanetsPlanetId) | **Get** /universe/planets/{planet_id}/ | Get planet information
[**GetUniverseStationsStationId**](UniverseApi.md#GetUniverseStationsStationId) | **Get** /universe/stations/{station_id}/ | Get station information
[**GetUniverseStructures**](UniverseApi.md#GetUniverseStructures) | **Get** /universe/structures/ | List all public structures
[**GetUniverseStructuresStructureId**](UniverseApi.md#GetUniverseStructuresStructureId) | **Get** /universe/structures/{structure_id}/ | Get structure information
[**GetUniverseSystemsSystemId**](UniverseApi.md#GetUniverseSystemsSystemId) | **Get** /universe/systems/{system_id}/ | Get solar system information
[**GetUniverseTypesTypeId**](UniverseApi.md#GetUniverseTypesTypeId) | **Get** /universe/types/{type_id}/ | Get type information
[**PostUniverseNames**](UniverseApi.md#PostUniverseNames) | **Post** /universe/names/ | Get names and categories for a set of ID&#39;s


# **GetUniversePlanetsPlanetId**
> GetUniversePlanetsPlanetIdOk GetUniversePlanetsPlanetId($planetId, $datasource)

Get planet information

Information on a planet  ---  Alternate route: `/v1/universe/planets/{planet_id}/`  Alternate route: `/legacy/universe/planets/{planet_id}/`  Alternate route: `/dev/universe/planets/{planet_id}/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **planetId** | **int32**| An Eve planet ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**GetUniversePlanetsPlanetIdOk**](get_universe_planets_planet_id_ok.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetUniverseStationsStationId**
> GetUniverseStationsStationIdOk GetUniverseStationsStationId($stationId, $datasource)

Get station information

Public information on stations  ---  Alternate route: `/v1/universe/stations/{station_id}/`  Alternate route: `/legacy/universe/stations/{station_id}/`  Alternate route: `/dev/universe/stations/{station_id}/`   ---  This route is cached for up to 3600 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **stationId** | **int32**| An Eve station ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**GetUniverseStationsStationIdOk**](get_universe_stations_station_id_ok.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetUniverseStructures**
> []int64 GetUniverseStructures($datasource)

List all public structures

List all public structures  ---  Alternate route: `/v1/universe/structures/`  Alternate route: `/legacy/universe/structures/`  Alternate route: `/dev/universe/structures/`   ---  This route is cached for up to 3600 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

**[]int64**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetUniverseStructuresStructureId**
> GetUniverseStructuresStructureIdOk GetUniverseStructuresStructureId($structureId, $datasource)

Get structure information

Returns information on requested structure, if you are on the ACL. Otherwise, returns \"Forbidden\" for all inputs.  ---  Alternate route: `/v1/universe/structures/{structure_id}/`  Alternate route: `/legacy/universe/structures/{structure_id}/`  Alternate route: `/dev/universe/structures/{structure_id}/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **structureId** | **int64**| An Eve structure ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**GetUniverseStructuresStructureIdOk**](get_universe_structures_structure_id_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetUniverseSystemsSystemId**
> GetUniverseSystemsSystemIdOk GetUniverseSystemsSystemId($systemId, $datasource)

Get solar system information

Information on solar systems  ---  Alternate route: `/v1/universe/systems/{system_id}/`  Alternate route: `/legacy/universe/systems/{system_id}/`  Alternate route: `/dev/universe/systems/{system_id}/`   ---  This route is cached for up to 3600 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **systemId** | **int32**| An Eve solar system ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**GetUniverseSystemsSystemIdOk**](get_universe_systems_system_id_ok.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetUniverseTypesTypeId**
> GetUniverseTypesTypeIdOk GetUniverseTypesTypeId($typeId, $datasource)

Get type information

Get information on a type  ---  Alternate route: `/v1/universe/types/{type_id}/`  Alternate route: `/legacy/universe/types/{type_id}/`  Alternate route: `/dev/universe/types/{type_id}/`   ---  This route is cached for up to 3600 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **typeId** | **int32**| An Eve item type ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**GetUniverseTypesTypeIdOk**](get_universe_types_type_id_ok.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostUniverseNames**
> []PostUniverseNames200Ok PostUniverseNames($ids, $datasource)

Get names and categories for a set of ID's

Resolve a set of IDs to names and categories. Supported ID's for resolving are: Characters, Corporations, Alliances, Stations, Solar Systems, Constellations, Regions, Types.  ---  Alternate route: `/v1/universe/names/`  Alternate route: `/legacy/universe/names/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ids** | [**PostUniverseNamesIds**](PostUniverseNamesIds.md)| The ids to resolve | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]PostUniverseNames200Ok**](post_universe_names_200_ok.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

