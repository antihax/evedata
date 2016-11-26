# \WalletApi

All URIs are relative to *https://esi.tech.ccp.is/latest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCharactersCharacterIdWallets**](WalletApi.md#GetCharactersCharacterIdWallets) | **Get** /characters/{character_id}/wallets/ | List wallets and balances
[**GetCharactersCharacterIdWalletsJournal**](WalletApi.md#GetCharactersCharacterIdWalletsJournal) | **Get** /characters/{character_id}/wallets/journal/ | Get character wallet journal
[**GetCharactersCharacterIdWalletsTransactions**](WalletApi.md#GetCharactersCharacterIdWalletsTransactions) | **Get** /characters/{character_id}/wallets/transactions/ | Get wallet transactions
[**GetCorporationsCorporationIdWallets**](WalletApi.md#GetCorporationsCorporationIdWallets) | **Get** /corporations/{corporation_id}/wallets/ | Dummy Endpoint, Please Ignore
[**GetCorporationsCorporationIdWalletsWalletIdJournal**](WalletApi.md#GetCorporationsCorporationIdWalletsWalletIdJournal) | **Get** /corporations/{corporation_id}/wallets/{wallet_id}/journal/ | Dummy Endpoint, Please Ignore
[**GetCorporationsCorporationIdWalletsWalletIdTransactions**](WalletApi.md#GetCorporationsCorporationIdWalletsWalletIdTransactions) | **Get** /corporations/{corporation_id}/wallets/{wallet_id}/transactions/ | Dummy Endpoint, Please Ignore


# **GetCharactersCharacterIdWallets**
> []GetCharactersCharacterIdWallets200Ok GetCharactersCharacterIdWallets($characterId, $datasource)

List wallets and balances

List your wallets and their balances. Characters typically have only one wallet, with wallet_id 1000 being the master wallet.  ---  Alternate route: `/v1/characters/{character_id}/wallets/`  Alternate route: `/legacy/characters/{character_id}/wallets/`  Alternate route: `/dev/characters/{character_id}/wallets/`   ---  This route is cached for up to 120 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetCharactersCharacterIdWallets200Ok**](get_characters_character_id_wallets_200_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

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

