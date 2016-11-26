# \MailApi

All URIs are relative to *https://esi.tech.ccp.is/latest*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteCharactersCharacterIdMailMailId**](MailApi.md#DeleteCharactersCharacterIdMailMailId) | **Delete** /characters/{character_id}/mail/{mail_id}/ | Delete a mail
[**GetCharactersCharacterIdMail**](MailApi.md#GetCharactersCharacterIdMail) | **Get** /characters/{character_id}/mail/ | Return mail headers
[**GetCharactersCharacterIdMailLabels**](MailApi.md#GetCharactersCharacterIdMailLabels) | **Get** /characters/{character_id}/mail/labels/ | Get mail labels and unread counts
[**GetCharactersCharacterIdMailLists**](MailApi.md#GetCharactersCharacterIdMailLists) | **Get** /characters/{character_id}/mail/lists/ | Return mailing list subscriptions
[**GetCharactersCharacterIdMailMailId**](MailApi.md#GetCharactersCharacterIdMailMailId) | **Get** /characters/{character_id}/mail/{mail_id}/ | Return a mail
[**PostCharactersCharacterIdMail**](MailApi.md#PostCharactersCharacterIdMail) | **Post** /characters/{character_id}/mail/ | Send a new mail
[**PostCharactersCharacterIdMailLabels**](MailApi.md#PostCharactersCharacterIdMailLabels) | **Post** /characters/{character_id}/mail/labels/ | Create a mail label
[**PutCharactersCharacterIdMailMailId**](MailApi.md#PutCharactersCharacterIdMailMailId) | **Put** /characters/{character_id}/mail/{mail_id}/ | Update metadata about a mail


# **DeleteCharactersCharacterIdMailMailId**
> DeleteCharactersCharacterIdMailMailId($characterId, $mailId, $datasource)

Delete a mail

Delete a mail  ---  Alternate route: `/v1/characters/{character_id}/mail/{mail_id}/`  Alternate route: `/legacy/characters/{character_id}/mail/{mail_id}/`  Alternate route: `/dev/characters/{character_id}/mail/{mail_id}/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **mailId** | **int32**| An EVE mail ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCharactersCharacterIdMail**
> []GetCharactersCharacterIdMail200Ok GetCharactersCharacterIdMail($characterId, $labels, $lastMailId, $datasource)

Return mail headers

Return the 50 most recent mail headers belonging to the character that match the query criteria. Queries can be filtered by label, and last_mail_id can be used to paginate backwards.  ---  Alternate route: `/v1/characters/{character_id}/mail/`  Alternate route: `/legacy/characters/{character_id}/mail/`  Alternate route: `/dev/characters/{character_id}/mail/`   ---  This route is cached for up to 30 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **labels** | [**[]int64**](int64.md)| Fetch only mails that match one or more of the given labels | [optional] 
 **lastMailId** | **int32**| List only mail with an ID lower than the given ID, if present | [optional] 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetCharactersCharacterIdMail200Ok**](get_characters_character_id_mail_200_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCharactersCharacterIdMailLabels**
> GetCharactersCharacterIdMailLabelsOk GetCharactersCharacterIdMailLabels($characterId, $datasource)

Get mail labels and unread counts

Return a list of the users mail labels, unread counts for each label and a total unread count.  ---  Alternate route: `/v3/characters/{character_id}/mail/labels/`  Alternate route: `/dev/characters/{character_id}/mail/labels/`   ---  This route is cached for up to 30 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**GetCharactersCharacterIdMailLabelsOk**](get_characters_character_id_mail_labels_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCharactersCharacterIdMailLists**
> []GetCharactersCharacterIdMailLists200Ok GetCharactersCharacterIdMailLists($characterId, $datasource)

Return mailing list subscriptions

Return all mailing lists that the character is subscribed to   ---  Alternate route: `/v1/characters/{character_id}/mail/lists/`  Alternate route: `/legacy/characters/{character_id}/mail/lists/`  Alternate route: `/dev/characters/{character_id}/mail/lists/`   ---  This route is cached for up to 120 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**[]GetCharactersCharacterIdMailLists200Ok**](get_characters_character_id_mail_lists_200_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCharactersCharacterIdMailMailId**
> GetCharactersCharacterIdMailMailIdOk GetCharactersCharacterIdMailMailId($characterId, $mailId, $datasource)

Return a mail

Return the contents of an EVE mail  ---  Alternate route: `/v1/characters/{character_id}/mail/{mail_id}/`  Alternate route: `/legacy/characters/{character_id}/mail/{mail_id}/`  Alternate route: `/dev/characters/{character_id}/mail/{mail_id}/`   ---  This route is cached for up to 30 seconds


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **mailId** | **int32**| An EVE mail ID | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

[**GetCharactersCharacterIdMailMailIdOk**](get_characters_character_id_mail_mail_id_ok.md)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostCharactersCharacterIdMail**
> int32 PostCharactersCharacterIdMail($characterId, $mail, $datasource)

Send a new mail

Create and send a new mail  ---  Alternate route: `/v1/characters/{character_id}/mail/`  Alternate route: `/legacy/characters/{character_id}/mail/`  Alternate route: `/dev/characters/{character_id}/mail/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| The sender&#39;s character ID | 
 **mail** | [**PostCharactersCharacterIdMailMail**](PostCharactersCharacterIdMailMail.md)| The mail to send | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

**int32**

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostCharactersCharacterIdMailLabels**
> int64 PostCharactersCharacterIdMailLabels($characterId, $label, $datasource)

Create a mail label

Create a mail label  ---  Alternate route: `/v2/characters/{character_id}/mail/labels/`  Alternate route: `/legacy/characters/{character_id}/mail/labels/`  Alternate route: `/dev/characters/{character_id}/mail/labels/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **label** | [**PostCharactersCharacterIdMailLabelsLabel**](PostCharactersCharacterIdMailLabelsLabel.md)| Label to create | [optional] 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

**int64**

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PutCharactersCharacterIdMailMailId**
> PutCharactersCharacterIdMailMailId($characterId, $mailId, $contents, $datasource)

Update metadata about a mail

Update metadata about a mail  ---  Alternate route: `/v1/characters/{character_id}/mail/{mail_id}/`  Alternate route: `/legacy/characters/{character_id}/mail/{mail_id}/`  Alternate route: `/dev/characters/{character_id}/mail/{mail_id}/` 


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **characterId** | **int32**| An EVE character ID | 
 **mailId** | **int32**| An EVE mail ID | 
 **contents** | [**PutCharactersCharacterIdMailMailIdContents**](PutCharactersCharacterIdMailMailIdContents.md)| Data used to update the mail | 
 **datasource** | **string**| The server name you would like data from | [optional] [default to tranquility]

### Return type

void (empty response body)

### Authorization

[evesso](../README.md#evesso)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

