# \DummyApi

All URIs are relative to *https://esi.tech.ccp.is/latest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCharactersCharacterIdWalletsJournal**](DummyApi.md#GetCharactersCharacterIdWalletsJournal) | **Get** /characters/{character_id}/wallets/journal/ | Get character wallet journal
[**GetCharactersCharacterIdWalletsTransactions**](DummyApi.md#GetCharactersCharacterIdWalletsTransactions) | **Get** /characters/{character_id}/wallets/transactions/ | Get wallet transactions
[**GetCorporationsCorporationIdAssets**](DummyApi.md#GetCorporationsCorporationIdAssets) | **Get** /corporations/{corporation_id}/assets/ | Dummy Endpoint, Please Ignore
[**GetCorporationsCorporationIdAssetsAssetIdLogs**](DummyApi.md#GetCorporationsCorporationIdAssetsAssetIdLogs) | **Get** /corporations/{corporation_id}/assets/{asset_id}/logs/ | Dummy Endpoint, Please Ignore
[**GetCorporationsCorporationIdBookmarks**](DummyApi.md#GetCorporationsCorporationIdBookmarks) | **Get** /corporations/{corporation_id}/bookmarks/ | Dummy Endpoint, Please Ignore
[**GetCorporationsCorporationIdBookmarksFolders**](DummyApi.md#GetCorporationsCorporationIdBookmarksFolders) | **Get** /corporations/{corporation_id}/bookmarks/folders/ | Dummy Endpoint, Please Ignore
[**GetCorporationsCorporationIdWallets**](DummyApi.md#GetCorporationsCorporationIdWallets) | **Get** /corporations/{corporation_id}/wallets/ | Dummy Endpoint, Please Ignore
[**GetCorporationsCorporationIdWalletsWalletIdJournal**](DummyApi.md#GetCorporationsCorporationIdWalletsWalletIdJournal) | **Get** /corporations/{corporation_id}/wallets/{wallet_id}/journal/ | Dummy Endpoint, Please Ignore
[**GetCorporationsCorporationIdWalletsWalletIdTransactions**](DummyApi.md#GetCorporationsCorporationIdWalletsWalletIdTransactions) | **Get** /corporations/{corporation_id}/wallets/{wallet_id}/transactions/ | Dummy Endpoint, Please Ignore
[**GetUniversePlanetsPlanetId**](DummyApi.md#GetUniversePlanetsPlanetId) | **Get** /universe/planets/{planet_id}/ | Get planet information


# **GetCharactersCharacterIdWalletsJournal**
> []GetCharactersCharacterIdWalletsJournal200Ok GetCharactersCharacterIdWalletsJournal($characterId, $lastSeenId, $datasource)

Get character wallet journal

Returns the most recent 50 entries for the characters wallet journal. Optionally, takes an argument with a reference ID, and returns the prior 50 entries from the journal.  ---  Alternate route: `/v1/characters/{character_id}/wallets/journal/`  Alternate route: `/legacy/characters/{character_id}/wallets/journal/`  Alternate route: `/dev/characters/{character_id}/wallets/journal/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **lastSeenId** | **int64**| A journal reference ID to paginate from | [optional] 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetCharactersCharacterIdWalletsJournal200Ok**](get_characters_character_id_wallets_journal_200_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCharactersCharacterIdWalletsTransactions**
> []GetCharactersCharacterIdWalletsTransactions200Ok GetCharactersCharacterIdWalletsTransactions($characterId, $datasource)

Get wallet transactions

Gets the 50 most recent transactions in a characters wallet. Optionally, takes an argument with a transaction ID, and returns the prior 50 transactions  ---  Alternate route: `/v1/characters/{character_id}/wallets/transactions/`  Alternate route: `/legacy/characters/{character_id}/wallets/transactions/`  Alternate route: `/dev/characters/{character_id}/wallets/transactions/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetCharactersCharacterIdWalletsTransactions200Ok**](get_characters_character_id_wallets_transactions_200_ok.md)

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

# **GetCorporationsCorporationIdBookmarks**
> GetCorporationsCorporationIdBookmarks($corporationId, $datasource)

Dummy Endpoint, Please Ignore

Dummy  ---  Alternate route: `/v1/corporations/{corporation_id}/bookmarks/`  Alternate route: `/legacy/corporations/{corporation_id}/bookmarks/`  Alternate route: `/dev/corporations/{corporation_id}/bookmarks/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **corporationId** | **int32**| An EVE corporation ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCorporationsCorporationIdBookmarksFolders**
> GetCorporationsCorporationIdBookmarksFolders($corporationId, $datasource)

Dummy Endpoint, Please Ignore

Dummy  ---  Alternate route: `/v1/corporations/{corporation_id}/bookmarks/folders/`  Alternate route: `/legacy/corporations/{corporation_id}/bookmarks/folders/`  Alternate route: `/dev/corporations/{corporation_id}/bookmarks/folders/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **corporationId** | **int32**| An EVE corporation ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCorporationsCorporationIdWallets**
> GetCorporationsCorporationIdWallets($corporationId, $datasource)

Dummy Endpoint, Please Ignore

Dummy  ---  Alternate route: `/v1/corporations/{corporation_id}/wallets/`  Alternate route: `/legacy/corporations/{corporation_id}/wallets/`  Alternate route: `/dev/corporations/{corporation_id}/wallets/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **corporationId** | **int32**| An EVE corporation ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCorporationsCorporationIdWalletsWalletIdJournal**
> GetCorporationsCorporationIdWalletsWalletIdJournal($corporationId, $walletId, $datasource)

Dummy Endpoint, Please Ignore

Dummy  ---  Alternate route: `/v1/corporations/{corporation_id}/wallets/{wallet_id}/journal/`  Alternate route: `/legacy/corporations/{corporation_id}/wallets/{wallet_id}/journal/`  Alternate route: `/dev/corporations/{corporation_id}/wallets/{wallet_id}/journal/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **corporationId** | **int32**| An EVE corporation ID | 
 **walletId** | **int32**| Wallet ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCorporationsCorporationIdWalletsWalletIdTransactions**
> GetCorporationsCorporationIdWalletsWalletIdTransactions($corporationId, $walletId, $datasource)

Dummy Endpoint, Please Ignore

Dummy  ---  Alternate route: `/v1/corporations/{corporation_id}/wallets/{wallet_id}/transactions/`  Alternate route: `/legacy/corporations/{corporation_id}/wallets/{wallet_id}/transactions/`  Alternate route: `/dev/corporations/{corporation_id}/wallets/{wallet_id}/transactions/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **corporationId** | **int32**| An EVE corporation ID | 
 **walletId** | **int32**| Wallet ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

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

