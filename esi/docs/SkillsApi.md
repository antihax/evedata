# \SkillsApi

All URIs are relative to *https://esi.tech.ccp.is/latest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCharactersCharacterIdSkillqueue**](SkillsApi.md#GetCharactersCharacterIdSkillqueue) | **Get** /characters/{character_id}/skillqueue/ | Get character&#39;s skill queue
[**GetCharactersCharacterIdSkills**](SkillsApi.md#GetCharactersCharacterIdSkills) | **Get** /characters/{character_id}/skills/ | Get character skills


# **GetCharactersCharacterIdSkillqueue**
> []GetCharactersCharacterIdSkillqueue200Ok GetCharactersCharacterIdSkillqueue($characterId, $datasource)

Get character's skill queue

List the configured skill queue for the given character  ---  Alternate route: `/v2/characters/{character_id}/skillqueue/`  Alternate route: `/legacy/characters/{character_id}/skillqueue/`  Alternate route: `/dev/characters/{character_id}/skillqueue/`   ---  This route is cached for up to 120 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| Character id of the target character | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetCharactersCharacterIdSkillqueue200Ok**](get_characters_character_id_skillqueue_200_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCharactersCharacterIdSkills**
> GetCharactersCharacterIdSkillsOk GetCharactersCharacterIdSkills($characterId, $datasource)

Get character skills

List all trained skills for the given character  ---  Alternate route: `/v3/characters/{character_id}/skills/`  Alternate route: `/dev/characters/{character_id}/skills/`   ---  This route is cached for up to 120 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**GetCharactersCharacterIdSkillsOk**](get_characters_character_id_skills_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

