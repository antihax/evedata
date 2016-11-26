# \AssetsApi

All URIs are relative to *https://esi.tech.ccp.is/latest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCharactersCharacterIdAssets**](AssetsApi.md#GetCharactersCharacterIdAssets) | **Get** /characters/{character_id}/assets/ | Get character assets
[**GetCorporationsCorporationIdAssets**](AssetsApi.md#GetCorporationsCorporationIdAssets) | **Get** /corporations/{corporation_id}/assets/ | Dummy Endpoint, Please Ignore
[**GetCorporationsCorporationIdAssetsAssetIdLogs**](AssetsApi.md#GetCorporationsCorporationIdAssetsAssetIdLogs) | **Get** /corporations/{corporation_id}/assets/{asset_id}/logs/ | Dummy Endpoint, Please Ignore


# **GetCharactersCharacterIdAssets**
> []GetCharactersCharacterIdAssets200Ok GetCharactersCharacterIdAssets($characterId, $datasource)

Get character assets

Return a list of the characters assets  ---  Alternate route: `/v1/characters/{character_id}/assets/`  Alternate route: `/legacy/characters/{character_id}/assets/`  Alternate route: `/dev/characters/{character_id}/assets/`   ---  This route is cached for up to 3600 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| Character id of the target character | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetCharactersCharacterIdAssets200Ok**](get_characters_character_id_assets_200_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCorporationsCorporationIdAssets**
> GetCorporationsCorporationIdAssets($corporationId, $datasource)

Dummy Endpoint, Please Ignore

Dummy  ---  Alternate route: `/v1/corporations/{corporation_id}/assets/`  Alternate route: `/legacy/corporations/{corporation_id}/assets/`  Alternate route: `/dev/corporations/{corporation_id}/assets/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **corporationId** | **int32**| Corporation id of the target corporation | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCorporationsCorporationIdAssetsAssetIdLogs**
> GetCorporationsCorporationIdAssetsAssetIdLogs($corporationId, $assetId, $datasource)

Dummy Endpoint, Please Ignore

Dummy  ---  Alternate route: `/v1/corporations/{corporation_id}/assets/{asset_id}/logs/`  Alternate route: `/legacy/corporations/{corporation_id}/assets/{asset_id}/logs/`  Alternate route: `/dev/corporations/{corporation_id}/assets/{asset_id}/logs/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **corporationId** | **int32**| Corporation id of the target corporation | 
 **assetId** | **int32**| Asset id | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

