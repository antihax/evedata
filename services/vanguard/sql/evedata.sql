CREATE DATABASE IF NOT EXISTS evedata;

USE evedata;

CREATE TABLE `allianceHistory` (
  `recordID` int(11) NOT NULL,
  `startDate` datetime NOT NULL,
  `allianceID` int(11) NOT NULL,
  `corporationID` int(11) NOT NULL,
  `deleted` tinyint(4) NOT NULL,
  `endDate` datetime DEFAULT NULL,
  PRIMARY KEY (`recordID`),
  KEY `allianceID` (`allianceID`),
  KEY `corporationID` (`corporationID`),
  KEY `startDate` (`startDate`),
  KEY `corporationID_Date` (`corporationID`,`startDate`),
  KEY `endDate` (`endDate`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `alliances` (
  `allianceID` int(10) unsigned NOT NULL DEFAULT '0',
  `name` varchar(255) NOT NULL DEFAULT '',
  `shortName` varchar(45) NOT NULL DEFAULT '',
  `executorCorpID` int(10) unsigned NOT NULL DEFAULT '0',
  `corporationsCount` int(10) unsigned NOT NULL DEFAULT '0',
  `startDate` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `updated` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `cacheUntil` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `memberCount` int(11) NOT NULL DEFAULT '0',
  `dead` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`allianceID`),
  KEY `name` (`name`),
  KEY `executorCorpID` (`executorCorpID`),
  KEY `cacheUntil` (`cacheUntil`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `assets` (
  `locationID` bigint(20) unsigned NOT NULL,
  `typeID` smallint(5) unsigned NOT NULL DEFAULT '0',
  `quantity` int(10) NOT NULL DEFAULT '0',
  `characterID` int(10) unsigned NOT NULL DEFAULT '0',
  `locationFlag` varchar(40) NOT NULL,
  `itemID` bigint(20) unsigned NOT NULL,
  `locationType` varchar(20) NOT NULL,
  `isSingleton` tinyint(1) unsigned NOT NULL,
  PRIMARY KEY (`itemID`),
  KEY `locationID` (`locationID`),
  KEY `characterID` (`characterID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `characterAssociations` (
  `characterID` int(10) unsigned NOT NULL,
  `associateID` int(10) unsigned NOT NULL,
  `frequency` smallint(5) unsigned NOT NULL,
  `source` tinyint(4) unsigned DEFAULT NULL,
  `added` datetime DEFAULT NULL,
  PRIMARY KEY (`characterID`,`associateID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `characterKillmailAssociations` (
  `characterID` int(10) unsigned NOT NULL,
  `associateID` int(10) unsigned NOT NULL,
  `frequency` smallint(5) unsigned NOT NULL,
  `added` datetime DEFAULT NULL,
  PRIMARY KEY (`characterID`,`associateID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `characters` (
  `characterID` int(11) unsigned NOT NULL DEFAULT '0',
  `bloodlineID` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `ancestryID` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `corporationID` int(11) unsigned NOT NULL DEFAULT '0',
  `allianceID` int(11) unsigned NOT NULL DEFAULT '0',
  `race` char(8) CHARACTER SET latin1 NOT NULL DEFAULT '0',
  `securityStatus` decimal(4,2) NOT NULL DEFAULT '0.00',
  `updated` datetime NOT NULL DEFAULT '2001-01-01 00:00:00',
  `cacheUntil` datetime NOT NULL DEFAULT '2001-01-01 00:00:00',
  `name` varchar(50) NOT NULL DEFAULT '0',
  `gender` varchar(50) NOT NULL DEFAULT '0',
  `dead` tinyint(4) NOT NULL DEFAULT '0',
  `birthDate` datetime DEFAULT NULL,
  PRIMARY KEY (`characterID`),
  KEY `cacheUntil` (`cacheUntil`),
  KEY `name` (`name`),
  KEY `corporationIDCharacterID` (`corporationID`,`characterID`),
  KEY `ix_notdeadcharexpired` (`cacheUntil`,`characterID`,`dead`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `contactSyncs` (
  `characterID` int(11) NOT NULL,
  `source` int(11) NOT NULL,
  `destination` int(11) NOT NULL,
  `lastError` varchar(100) CHARACTER SET latin1 DEFAULT NULL,
  `nextSync` datetime NOT NULL DEFAULT '2016-10-21 01:57:40',
  PRIMARY KEY (`characterID`,`destination`,`source`),
  UNIQUE KEY `destination_UNIQUE` (`destination`),
  KEY `source` (`source`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `contractBids` (
  `contractID` bigint(20) unsigned NOT NULL,
  `bidID` int(10) unsigned NOT NULL,
  `dateBid` datetime NOT NULL,
  `amount` decimal(22,2) DEFAULT NULL,
  PRIMARY KEY (`contractID`,`bidID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `contractItems` (
  `recordID` bigint(20) NOT NULL,
  `contractID` int(11) NOT NULL,
  `isBPC` tinyint(3) unsigned DEFAULT NULL,
  `itemID` bigint(20) DEFAULT NULL,
  `typeID` int(11) DEFAULT NULL,
  `isIncluded` tinyint(4) DEFAULT NULL,
  `ME` tinyint(4) DEFAULT NULL,
  `TE` tinyint(4) DEFAULT NULL,
  `runs` int(11) DEFAULT NULL,
  `quantity` int(11) DEFAULT NULL,
  PRIMARY KEY (`recordID`,`contractID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `contracts` (
  `contractID` int(10) unsigned NOT NULL,
  `buyout` decimal(22,2) DEFAULT NULL,
  `collateral` decimal(22,2) DEFAULT NULL,
  `dateExpired` datetime DEFAULT NULL,
  `dateIssued` datetime NOT NULL,
  `daysToComplete` tinyint(3) unsigned DEFAULT NULL,
  `endLocationId` bigint(20) unsigned DEFAULT NULL,
  `forCorporation` tinyint(3) unsigned NOT NULL,
  `issuerCorporationID` int(10) unsigned DEFAULT NULL,
  `issuerID` int(10) unsigned NOT NULL,
  `price` decimal(22,2) DEFAULT NULL,
  `reward` decimal(22,2) DEFAULT NULL,
  `locationID` bigint(20) DEFAULT NULL,
  `title` varchar(255) COLLATE utf8_bin NOT NULL,
  `type` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `volume` decimal(22,2) DEFAULT NULL,
  PRIMARY KEY (`contractID`),
  KEY `ix_location_type_exp` (`locationID`,`type`,`dateExpired`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `corporationHistory` (
  `recordID` int(11) NOT NULL,
  `startDate` datetime NOT NULL,
  `characterID` int(11) NOT NULL,
  `corporationID` int(11) DEFAULT NULL,
  `endDate` date DEFAULT NULL,
  PRIMARY KEY (`recordID`),
  KEY `characterID` (`characterID`),
  KEY `corporationID` (`corporationID`),
  KEY `startDate` (`startDate`),
  KEY `characterID_Date` (`characterID`,`endDate`),
  KEY `endDate` (`endDate`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `corporations` (
  `corporationID` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `ticker` varchar(45) NOT NULL,
  `ceoID` int(11) NOT NULL,
  `allianceID` int(11) NOT NULL DEFAULT '0',
  `factionID` int(11) NOT NULL DEFAULT '0',
  `memberCount` int(11) NOT NULL DEFAULT '0',
  `updated` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `cacheUntil` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `dead` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`corporationID`),
  KEY `allianceID` (`allianceID`),
  KEY `factionID` (`factionID`),
  KEY `cacheUntil` (`cacheUntil`),
  KEY `name` (`name`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `crestTokens` (
  `characterID` int(11) NOT NULL,
  `tokenCharacterID` int(11) NOT NULL,
  `characterOwnerHash` varchar(32) NOT NULL DEFAULT '',
  `accessToken` varchar(100) CHARACTER SET latin1 NOT NULL,
  `refreshToken` text CHARACTER SET latin1 NOT NULL,
  `expiry` datetime NOT NULL,
  `tokenType` varchar(100) CHARACTER SET latin1 NOT NULL,
  `lastCode` int(11) NOT NULL DEFAULT '0',
  `lastStatus` text CHARACTER SET latin1 NOT NULL,
  `characterName` varchar(100) NOT NULL,
  `request` text CHARACTER SET latin1,
  `response` text CHARACTER SET latin1,
  `scopes` text NOT NULL,
  `authCharacter` tinyint(1) NOT NULL DEFAULT '0',
  `mailedError` tinyint(1) NOT NULL DEFAULT '0',
  `roles` text,
  `corporationID` int(11) NOT NULL DEFAULT '0',
  `allianceID` int(11) NOT NULL DEFAULT '0',
  `factionID` int(11) NOT NULL DEFAULT '0',
  `mailPassword` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`characterID`,`tokenCharacterID`),
  KEY `tokenCharacterID` (`tokenCharacterID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `cursorCharacter` (
  `characterID` int(11) NOT NULL,
  `cursorCharacterID` int(11) NOT NULL,
  PRIMARY KEY (`characterID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `discoveredAssets` (
  `corporationID` int(10) unsigned NOT NULL,
  `allianceID` int(11) unsigned NOT NULL,
  `typeID` int(11) unsigned NOT NULL,
  `solarSystemID` int(8) unsigned NOT NULL,
  `x` float NOT NULL,
  `y` float NOT NULL,
  `z` float NOT NULL,
  `locationID` int(10) unsigned NOT NULL,
  `lastSeen` datetime NOT NULL,
  PRIMARY KEY (`corporationID`,`typeID`,`solarSystemID`,`locationID`),
  KEY `corporation` (`corporationID`),
  KEY `alliance` (`allianceID`),
  KEY `lastSeen` (`lastSeen`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `entities` (
  `id` int(10) unsigned NOT NULL,
  `type` varchar(60) COLLATE utf8_bin NOT NULL DEFAULT 'unknown',
  PRIMARY KEY (`id`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `entityContacts` (
  `entityID` int(11) NOT NULL,
  `contactID` int(11) NOT NULL,
  `standing` decimal(4,2) DEFAULT NULL,
  PRIMARY KEY (`entityID`,`contactID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `entityKillStats` (
  `id` int(10) unsigned NOT NULL,
  `kills` int(10) unsigned NOT NULL DEFAULT '0',
  `losses` int(10) unsigned NOT NULL DEFAULT '0',
  `efficiency` float NOT NULL DEFAULT '0',
  `capKills` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `httpErrors` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(255) DEFAULT NULL,
  `status` smallint(6) DEFAULT NULL,
  `request` text,
  `response` text,
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `integrationChannels` (
  `integrationID` int(10) unsigned NOT NULL,
  `channelID` varchar(255) COLLATE utf8_bin NOT NULL,
  `services` set('locator','kill','structure','application','war') COLLATE utf8_bin NOT NULL,
  `options` text COLLATE utf8_bin NOT NULL,
  `channelName` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT 'unknown',
  PRIMARY KEY (`channelID`,`integrationID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `integrationCharacters` (
  `integrationID` int(11) NOT NULL,
  `characterID` int(11) NOT NULL,
  `integrationUserID` varchar(255) COLLATE utf8_bin DEFAULT '',
  PRIMARY KEY (`integrationID`,`characterID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `integrationDelegate` (
  `integrationID` int(11) NOT NULL,
  `characterID` int(11) NOT NULL,
  `permissions` enum('administrator') COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`integrationID`,`characterID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `integrationRoles` (
  `integrationID` int(10) unsigned NOT NULL DEFAULT '0',
  `roleID` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `roleName` varchar(45) COLLATE utf8_bin NOT NULL DEFAULT 'unknown',
  PRIMARY KEY (`integrationID`,`roleID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `integrationTokens` (
  `characterID` int(11) NOT NULL,
  `integrationUserID` varchar(255) NOT NULL DEFAULT '',
  `type` varchar(255) NOT NULL DEFAULT '',
  `integrationUserName` varchar(255) NOT NULL DEFAULT '',
  `accessToken` text CHARACTER SET latin1 NOT NULL,
  `refreshToken` text CHARACTER SET latin1 NOT NULL,
  `expiry` datetime NOT NULL,
  `tokenType` varchar(100) CHARACTER SET latin1 NOT NULL,
  `lastCode` int(11) NOT NULL DEFAULT '0',
  `lastStatus` text CHARACTER SET latin1 NOT NULL,
  `scopes` text NOT NULL,
  `mailedError` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`characterID`,`integrationUserID`,`type`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `integrations` (
  `integrationID` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `entityID` int(11) NOT NULL DEFAULT '0',
  `address` varchar(255) COLLATE utf8_bin NOT NULL,
  `authentication` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `type` enum('discord','ts3','slack') COLLATE utf8_bin NOT NULL,
  `services` set('auth') COLLATE utf8_bin NOT NULL DEFAULT '',
  `options` text COLLATE utf8_bin NOT NULL,
  `characterID` int(11) NOT NULL DEFAULT '0',
  `factionID` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`integrationID`),
  UNIQUE KEY `UNIQUE` (`address`,`authentication`)
) ENGINE=TokuDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `invMarketGroups` (
  `marketGroupID` mediumint(6) unsigned NOT NULL,
  `parentGroupID` mediumint(6) unsigned NOT NULL,
  `marketGroupName` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `description` text COLLATE utf8_bin,
  PRIMARY KEY (`marketGroupID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `iskPerLp` (
  `itemName` varchar(100) NOT NULL,
  `typeID` int(11) NOT NULL DEFAULT '0',
  `typeName` varchar(100) DEFAULT NULL,
  `lpCost` int(11) NOT NULL,
  `iskCost` int(11) NOT NULL,
  `JitaPrice` decimal(14,0) NOT NULL,
  `JitaVolume` decimal(11,0) NOT NULL,
  `itemCost` decimal(47,0) DEFAULT NULL,
  `ISKperLP` decimal(48,0) NOT NULL,
  `offerID` int(11) DEFAULT NULL,
  PRIMARY KEY (`typeID`,`itemName`),
  KEY `itemName` (`itemName`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `jitaPrice` (
  `itemID` int(10) unsigned NOT NULL DEFAULT '0',
  `buy` decimal(15,2) unsigned DEFAULT NULL,
  `sell` decimal(15,2) unsigned DEFAULT NULL,
  `high` decimal(15,2) DEFAULT NULL,
  `low` decimal(15,2) DEFAULT NULL,
  `mean` decimal(15,2) DEFAULT NULL,
  `quantity` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`itemID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `jumps` (
  `toSolarSystemID` int(10) unsigned NOT NULL DEFAULT '0',
  `fromSolarSystemID` int(10) unsigned NOT NULL DEFAULT '0',
  `jumps` tinyint(3) unsigned DEFAULT '0',
  `secureJumps` tinyint(3) unsigned DEFAULT '0',
  PRIMARY KEY (`toSolarSystemID`,`fromSolarSystemID`),
  KEY `fromSolarSystemID` (`fromSolarSystemID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `killmailAttackers` (
  `id` int(10) unsigned NOT NULL,
  `characterID` int(10) unsigned NOT NULL DEFAULT '0',
  `corporationID` int(10) unsigned NOT NULL DEFAULT '0',
  `allianceID` int(10) unsigned NOT NULL DEFAULT '0',
  `shipType` smallint(5) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`,`characterID`),
  KEY `ix_allianceID` (`allianceID`),
  KEY `ix_corporationID` (`corporationID`),
  KEY `ix_characterID` (`characterID`),
  KEY `ix_shipType` (`shipType`),
  KEY `ix_id` (`id`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `killmailAttributes` (
  `id` int(11) NOT NULL,
  `eHP` double NOT NULL,
  `DPS` double NOT NULL,
  `Alpha` double NOT NULL,
  `scanResolution` double NOT NULL,
  `signatureRadiusNoMWD` double NOT NULL,
  `signatureRadius` double NOT NULL,
  `agility` double NOT NULL,
  `warpSpeed` double NOT NULL,
  `speedNoMWD` double NOT NULL,
  `speed` double NOT NULL,
  `remoteArmorRepair` double NOT NULL,
  `remoteShieldRepair` double NOT NULL,
  `remoteEnergyTransfer` double NOT NULL,
  `energyNeutralization` double NOT NULL,
  `sensorStrength` double NOT NULL,
  `RPS` double NOT NULL,
  `CPURemaining` double NOT NULL,
  `powerRemaining` double NOT NULL,
  `capacitorNoMWD` double NOT NULL,
  `capacitor` double NOT NULL,
  `capacitorTimeNoMWD` bigint(20) NOT NULL,
  `capacitorTime` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `ix_id_cpu_pg_ehp` (`id`,`CPURemaining`,`powerRemaining`,`eHP`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `killmails` (
  `id` int(9) unsigned NOT NULL,
  `solarSystemID` int(8) unsigned NOT NULL,
  `killTime` datetime NOT NULL,
  `victimCharacterID` int(10) unsigned NOT NULL DEFAULT '0',
  `victimCorporationID` int(10) unsigned NOT NULL DEFAULT '0',
  `victimAllianceID` int(10) unsigned NOT NULL DEFAULT '0',
  `shipType` smallint(5) unsigned NOT NULL DEFAULT '0',
  `warID` mediumint(8) unsigned NOT NULL DEFAULT '0',
  `factionID` mediumint(8) unsigned NOT NULL DEFAULT '0',
  `hash` varchar(100) NOT NULL DEFAULT '',
  `x` float NOT NULL DEFAULT '0',
  `y` float NOT NULL DEFAULT '0',
  `z` float NOT NULL,
  PRIMARY KEY (`id`),
  KEY `ix_victimAllianceID` (`victimAllianceID`),
  KEY `ix_victimCorporationID` (`victimCorporationID`),
  KEY `ix_killTime` (`killTime`),
  KEY `ix_war` (`warID`),
  KEY `ix_victimCharacterID` (`victimCharacterID`),
  KEY `ix_id_killtime` (`id`,`killTime`),
  KEY `ix_ship_time` (`shipType`,`killTime`),
  KEY `ix_war_killtime` (`warID`,`killTime`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `locatedCharacters` (
  `notificationID` int(11) NOT NULL,
  `characterID` int(11) NOT NULL,
  `solarSystemID` int(11) NOT NULL,
  `constellationID` int(11) NOT NULL,
  `regionID` int(11) NOT NULL,
  `stationID` int(11) NOT NULL,
  `locatedCharacterID` int(11) NOT NULL,
  `time` datetime NOT NULL,
  PRIMARY KEY (`notificationID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `locatorShareWith` (
  `characterID` int(11) NOT NULL,
  `entityID` int(11) NOT NULL,
  PRIMARY KEY (`characterID`,`entityID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `lpOfferRequirements` (
  `offerID` int(11) NOT NULL,
  `typeID` int(11) NOT NULL,
  `quantity` int(11) NOT NULL,
  PRIMARY KEY (`offerID`,`typeID`),
  KEY `lpReqs_offerID` (`offerID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `lpOffers` (
  `offerID` int(11) NOT NULL AUTO_INCREMENT,
  `corporationID` int(11) NOT NULL,
  `typeID` int(11) DEFAULT NULL,
  `quantity` int(11) DEFAULT NULL,
  `lpCost` int(11) DEFAULT NULL,
  `akCost` int(11) DEFAULT NULL,
  `iskCost` int(11) DEFAULT NULL,
  PRIMARY KEY (`offerID`,`corporationID`),
  KEY `lpOffers_typeid` (`typeID`),
  KEY `lpOffers_corporation` (`corporationID`),
  KEY `lpOffers_corp_type` (`corporationID`,`typeID`)
) ENGINE=TokuDB AUTO_INCREMENT=16340 DEFAULT CHARSET=utf8;

CREATE TABLE `market` (
  `orderID` bigint(11) unsigned NOT NULL DEFAULT '0',
  `regionID` mediumint(8) unsigned NOT NULL DEFAULT '0',
  `stationID` bigint(15) unsigned NOT NULL DEFAULT '0',
  `typeID` smallint(5) unsigned NOT NULL DEFAULT '0',
  `bid` tinyint(4) NOT NULL DEFAULT '0',
  `price` decimal(22,2) unsigned NOT NULL DEFAULT '0.00',
  `minVolume` int(10) unsigned NOT NULL DEFAULT '0',
  `remainingVolume` int(10) unsigned NOT NULL DEFAULT '0',
  `enteredVolume` int(10) unsigned NOT NULL DEFAULT '0',
  `issued` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `duration` smallint(5) unsigned NOT NULL DEFAULT '0',
  `reported` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  PRIMARY KEY (`orderID`),
  KEY `regionID_typeID` (`regionID`,`typeID`),
  KEY `typeID` (`typeID`),
  KEY `stationID_bid_reported` (`stationID`,`bid`,`reported`),
  KEY `regionID_bid` (`regionID`,`bid`)
) ENGINE=TokuDB DEFAULT CHARSET=latin1;

CREATE TABLE `marketHistoryStatistics` (
  `itemID` int(10) unsigned NOT NULL,
  `regionID` int(10) unsigned NOT NULL,
  `low` decimal(22,2) NOT NULL,
  `mean` decimal(22,2) NOT NULL,
  `high` decimal(22,2) NOT NULL,
  `quantity` bigint(20) unsigned NOT NULL,
  `orders` int(10) unsigned NOT NULL,
  PRIMARY KEY (`itemID`,`regionID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `marketOrderHistory` (
  `orderID` bigint(20) unsigned NOT NULL,
  `locationID` bigint(20) unsigned NOT NULL,
  `typeID` smallint(5) unsigned NOT NULL,
  `volumeChange` int(11) unsigned NOT NULL,
  `volumeRemain` int(11) unsigned NOT NULL,
  `price` decimal(22,2) unsigned NOT NULL,
  `duration` smallint(5) unsigned NOT NULL,
  `isBuyOrder` tinyint(4) unsigned NOT NULL,
  `changed` datetime NOT NULL,
  PRIMARY KEY (`orderID`,`changed`),
  KEY `ix_chg_buy` (`isBuyOrder`,`changed`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `marketStations` (
  `stationName` varchar(255) DEFAULT NULL,
  `stationID` bigint(20) unsigned NOT NULL DEFAULT '0',
  `Count` bigint(21) NOT NULL DEFAULT '0',
  PRIMARY KEY (`stationID`),
  KEY `count` (`Count`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `market_history` (
  `date` date NOT NULL,
  `low` decimal(20,2) unsigned NOT NULL,
  `high` decimal(20,2) unsigned NOT NULL,
  `mean` decimal(20,2) unsigned NOT NULL,
  `quantity` bigint(20) unsigned NOT NULL,
  `orders` smallint(4) unsigned NOT NULL,
  `itemID` smallint(5) unsigned NOT NULL,
  `regionID` mediumint(8) unsigned NOT NULL,
  PRIMARY KEY (`date`,`regionID`,`itemID`),
  KEY `regionIDDate` (`regionID`,`date`),
  KEY `date` (`date`),
  KEY `itemID_regionID_date` (`itemID`,`regionID`,`date`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `market_vol` (
  `number` bigint(21) NOT NULL DEFAULT '0',
  `quantity` decimal(36,4) DEFAULT NULL,
  `regionID` int(11) NOT NULL,
  `itemID` bigint(20) NOT NULL,
  PRIMARY KEY (`regionID`,`itemID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `notifications` (
  `notificationID` int(11) NOT NULL,
  `characterID` int(11) DEFAULT NULL,
  `notificationCharacterID` int(11) DEFAULT NULL,
  `senderID` int(11) DEFAULT NULL,
  `senderType` varchar(255) DEFAULT NULL,
  `timestamp` datetime DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `text` text,
  PRIMARY KEY (`notificationID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `orders` (
  `orderid` bigint(20) unsigned NOT NULL,
  `characterID` int(10) unsigned NOT NULL,
  `duration` int(10) unsigned NOT NULL,
  `isBuyOrder` tinyint(3) unsigned NOT NULL,
  `isCorporation` tinyint(3) unsigned NOT NULL,
  `escrow` decimal(14,2) NOT NULL,
  `issued` datetime NOT NULL,
  `locationID` bigint(20) unsigned NOT NULL,
  `minVolume` int(10) unsigned NOT NULL,
  `price` decimal(14,2) unsigned NOT NULL,
  `orderRange` varchar(45) COLLATE utf8_bin NOT NULL,
  `regionID` int(10) unsigned NOT NULL,
  `typeID` int(10) unsigned NOT NULL,
  `volumeRemain` int(10) unsigned NOT NULL,
  `volumeTotal` int(10) unsigned NOT NULL,
  PRIMARY KEY (`orderid`),
  KEY `characterID` (`characterID`),
  KEY `char_region_item` (`characterID`,`regionID`,`typeID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `sharing` (
  `characterID` int(11) unsigned NOT NULL,
  `tokenCharacterID` int(11) unsigned NOT NULL,
  `entityID` int(11) unsigned NOT NULL,
  `types` set('locator','kill','structure','application','war') COLLATE utf8_bin NOT NULL,
  `ignored` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`characterID`,`tokenCharacterID`,`entityID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='For sharing character information with entities.';

CREATE TABLE `states` (
  `state` varchar(45) NOT NULL,
  `value` int(11) NOT NULL,
  `nextCheck` datetime NOT NULL DEFAULT '2011-02-01 00:00:00',
  PRIMARY KEY (`state`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `structures` (
  `stationID` bigint(20) NOT NULL,
  `solarSystemID` int(11) DEFAULT NULL,
  `stationName` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `x` float DEFAULT NULL,
  `y` float DEFAULT NULL,
  `z` float DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `marketCacheUntil` datetime DEFAULT '2016-01-01 00:00:00',
  `ownerID` int(11) NOT NULL,
  `typeID` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`stationID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `walletJournal` (
  `refID` bigint(20) unsigned NOT NULL,
  `refTypeID` int(10) unsigned NOT NULL,
  `ownerID1` int(10) unsigned NOT NULL,
  `ownerID2` int(10) unsigned NOT NULL,
  `argID1` bigint(20) unsigned NOT NULL,
  `argName1` varchar(255) NOT NULL,
  `amount` decimal(22,2) NOT NULL,
  `balance` decimal(22,2) NOT NULL,
  `reason` varchar(255) NOT NULL,
  `taxReceiverID` int(11) unsigned NOT NULL,
  `taxAmount` decimal(22,2) NOT NULL,
  `date` datetime NOT NULL,
  `characterID` int(11) unsigned NOT NULL,
  PRIMARY KEY (`refID`,`ownerID1`,`ownerID2`,`characterID`),
  KEY `charID_date` (`characterID`,`date`),
  KEY `argID1` (`argID1`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `walletJournalRefType` (
  `refTypeID` int(11) NOT NULL,
  `refTypeName` varchar(255) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`refTypeID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `walletTransactions` (
  `transactionID` bigint(20) unsigned NOT NULL,
  `transactionDateTime` datetime NOT NULL,
  `quantity` int(10) unsigned NOT NULL,
  `typeID` int(10) unsigned NOT NULL,
  `price` decimal(22,2) unsigned NOT NULL,
  `clientID` int(10) unsigned NOT NULL,
  `characterID` int(10) unsigned NOT NULL,
  `stationID` bigint(20) unsigned NOT NULL,
  `transactionType` varchar(45) COLLATE utf8_bin NOT NULL,
  `transactionFor` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `journalTransactionID` bigint(20) unsigned NOT NULL,
  `clientTypeID` bigint(20) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`transactionID`,`characterID`),
  KEY `journalID` (`journalTransactionID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `warAllies` (
  `id` mediumint(8) unsigned NOT NULL,
  `allyID` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`,`allyID`),
  KEY `allyID` (`allyID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;

CREATE TABLE `wars` (
  `id` mediumint(8) unsigned NOT NULL,
  `timeFinished` datetime DEFAULT NULL,
  `timeStarted` datetime DEFAULT NULL,
  `timeDeclared` datetime DEFAULT NULL,
  `openForAllies` tinyint(4) unsigned DEFAULT NULL,
  `cacheUntil` datetime DEFAULT NULL,
  `aggressorID` int(11) unsigned NOT NULL DEFAULT '0',
  `defenderID` int(11) unsigned NOT NULL DEFAULT '0',
  `mutual` tinyint(4) unsigned NOT NULL DEFAULT '0',
  `finished` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `aggressorID` (`aggressorID`),
  KEY `defenderID` (`defenderID`),
  KEY `timeFinished_cacheUntil` (`timeFinished`,`cacheUntil`),
  KEY `idx_lookup` (`aggressorID`,`defenderID`)
) ENGINE=TokuDB DEFAULT CHARSET=utf8;


		DELIMITER $$
		CREATE PROCEDURE atWarWith(IN entity INT)
		BEGIN
			SELECT DISTINCT IF (aggressorID = entity, defenderID, aggressorID) AS id, timeStarted, timeFinished
				FROM evedata.wars W
				LEFT OUTER JOIN evedata.warAllies A ON A.id = W.id
				WHERE (aggressorID = entity OR defenderID = entity OR allyID = entity) AND
					(timeFinished > UTC_TIMESTAMP() OR
					timeFinished = "0001-01-01 00:00:00")
			UNION
				SELECT DISTINCT allyID AS id, timeStarted, timeFinished
				FROM evedata.wars W
				INNER JOIN evedata.warAllies A ON A.id = W.id
				WHERE (aggressorID = entity) AND
					(timeFinished > UTC_TIMESTAMP() OR
					timeFinished = "0001-01-01 00:00:00");
			END$$
			DELIMITER ;
		
			DELIMITER $$
			CREATE FUNCTION alliedMilita(factionID INT UNSIGNED) RETURNS int(11)
			DETERMINISTIC
			BEGIN
			IF factionID = 500001 THEN
				RETURN 500003;
			ELSEIF factionID = 500003 THEN
				RETURN 500001;
			ELSEIF factionID = 500002 THEN  
				RETURN 500004;
			ELSEIF factionID = 500004 THEN 
				RETURN 500002;
			END IF;
			RETURN 0;
			END$$
			DELIMITER ;
			
		DELIMITER $$
		CREATE FUNCTION constellationIDBySolarSystem(system INT UNSIGNED) RETURNS int(10) unsigned
			DETERMINISTIC
		BEGIN
			DECLARE constellation int(10) unsigned;
			SELECT constellationID INTO constellation
				FROM eve.mapSolarSystems
				WHERE solarSystemID = system
				LIMIT 1;
			
		RETURN constellation;
		END$$
		DELIMITER ;
		
		DELIMITER $$
		CREATE FUNCTION closestCelestial(s INT UNSIGNED, x1 FLOAT, y1 FLOAT, z1 FLOAT) RETURNS int(10) unsigned
			DETERMINISTIC
		BEGIN
			DECLARE celestialID int(10) unsigned;
			SELECT itemID INTO celestialID
				FROM eve.mapDenormalize
				WHERE orbitID IS NOT NULL AND solarSystemID = s
				ORDER BY POW(( x1 - x), 2) + POW(( y1 - y), 2) + POW(( z1 - z), 2)
				LIMIT 1;
			
		RETURN celestialID;
		END$$
		DELIMITER ;
		DELIMITER $$
		CREATE FUNCTION regionIDBySolarSystem(system INT UNSIGNED) RETURNS int(10) unsigned
			DETERMINISTIC
		BEGIN
			DECLARE region int(10) unsigned;
			SELECT regionID INTO region
				FROM eve.mapSolarSystems
				WHERE solarSystemID = system
				LIMIT 1;
			
		RETURN region;
		END$$
		DELIMITER ;
		DELIMITER $$
		CREATE FUNCTION regionIDByStructureID(structure BIGINT UNSIGNED) RETURNS int(10) unsigned
			DETERMINISTIC
		BEGIN
			DECLARE region int(10) unsigned;
			SELECT regionID INTO region
				FROM eve.staStations
				WHERE stationID = structure
				LIMIT 1;
			IF region IS NULL THEN 
				SET region = 0;
			END IF;
		RETURN region;	
		END$$
		DELIMITER ;
		DELIMITER $$
		CREATE FUNCTION raceByID(inRaceID int UNSIGNED) RETURNS VARCHAR(20) 
			DETERMINISTIC
		BEGIN
			DECLARE race VARCHAR(20) ;
			SELECT raceName INTO race
				FROM eve.chrRaces 
				WHERE raceID = inRaceID
				LIMIT 1;
			
		RETURN race;
		END$$
		DELIMITER ;
		