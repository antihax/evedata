# \BookmarksApi

All URIs are relative to *https://esi.tech.ccp.is/latest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCharactersCharacterIdBookmarks**](BookmarksApi.md#GetCharactersCharacterIdBookmarks) | **Get** /characters/{character_id}/bookmarks/ | List bookmarks
[**GetCharactersCharacterIdBookmarksFolders**](BookmarksApi.md#GetCharactersCharacterIdBookmarksFolders) | **Get** /characters/{character_id}/bookmarks/folders/ | List bookmark folders


# **GetCharactersCharacterIdBookmarks**
> []GetCharactersCharacterIdBookmarks200Ok GetCharactersCharacterIdBookmarks($characterId, $datasource)

List bookmarks

List your character's personal bookmarks  ---  Alternate route: `/v1/characters/{character_id}/bookmarks/`  Alternate route: `/legacy/characters/{character_id}/bookmarks/`  Alternate route: `/dev/characters/{character_id}/bookmarks/`   ---  This route is cached for up to 3600 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetCharactersCharacterIdBookmarks200Ok**](get_characters_character_id_bookmarks_200_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCharactersCharacterIdBookmarksFolders**
> []GetCharactersCharacterIdBookmarksFolders200Ok GetCharactersCharacterIdBookmarksFolders($characterId, $datasource)

List bookmark folders

List your character's personal bookmark folders  ---  Alternate route: `/v1/characters/{character_id}/bookmarks/folders/`  Alternate route: `/legacy/characters/{character_id}/bookmarks/folders/`  Alternate route: `/dev/characters/{character_id}/bookmarks/folders/`   ---  This route is cached for up to 3600 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetCharactersCharacterIdBookmarksFolders200Ok**](get_characters_character_id_bookmarks_folders_200_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

