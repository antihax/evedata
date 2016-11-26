# \KillmailsApi

All URIs are relative to *https://esi.tech.ccp.is/latest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCharactersCharacterIdKillmailsRecent**](KillmailsApi.md#GetCharactersCharacterIdKillmailsRecent) | **Get** /characters/{character_id}/killmails/recent/ | List kills and losses
[**GetKillmailsKillmailIdKillmailHash**](KillmailsApi.md#GetKillmailsKillmailIdKillmailHash) | **Get** /killmails/{killmail_id}/{killmail_hash}/ | Get a single killmail


# **GetCharactersCharacterIdKillmailsRecent**
> []GetCharactersCharacterIdKillmailsRecent200Ok GetCharactersCharacterIdKillmailsRecent($characterId, $maxCount, $maxKillId, $datasource)

List kills and losses

Return a list of character's recent kills and losses  ---  Alternate route: `/v1/characters/{character_id}/killmails/recent/`  Alternate route: `/legacy/characters/{character_id}/killmails/recent/`  Alternate route: `/dev/characters/{character_id}/killmails/recent/`   ---  This route is cached for up to 120 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **maxCount** | **int32**| How many killmails to return at maximum | [optional] [default to 50]
 **maxKillId** | **int32**| Only return killmails with ID smaller than this.  | [optional] 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetCharactersCharacterIdKillmailsRecent200Ok**](get_characters_character_id_killmails_recent_200_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetKillmailsKillmailIdKillmailHash**
> GetKillmailsKillmailIdKillmailHashOk GetKillmailsKillmailIdKillmailHash($killmailId, $killmailHash, $datasource)

Get a single killmail

Return a single killmail from its ID and hash  ---  Alternate route: `/v1/killmails/{killmail_id}/{killmail_hash}/`  Alternate route: `/legacy/killmails/{killmail_id}/{killmail_hash}/`  Alternate route: `/dev/killmails/{killmail_id}/{killmail_hash}/`   ---  This route is cached for up to 3600 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **killmailId** | **int32**| The killmail ID to be queried | 
 **killmailHash** | **string**| The killmail hash for verification | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**GetKillmailsKillmailIdKillmailHashOk**](get_killmails_killmail_id_killmail_hash_ok.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

