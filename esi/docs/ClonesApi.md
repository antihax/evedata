# \ClonesApi

All URIs are relative to *https://esi.tech.ccp.is/latest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCharactersCharacterIdClones**](ClonesApi.md#GetCharactersCharacterIdClones) | **Get** /characters/{character_id}/clones/ | Get clones


# **GetCharactersCharacterIdClones**
> GetCharactersCharacterIdClonesOk GetCharactersCharacterIdClones($characterId, $datasource)

Get clones

A list of the character's clones  ---  Alternate route: `/v2/characters/{character_id}/clones/`  Alternate route: `/dev/characters/{character_id}/clones/`   ---  This route is cached for up to 120 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**GetCharactersCharacterIdClonesOk**](get_characters_character_id_clones_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

