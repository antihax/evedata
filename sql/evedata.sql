CREATE DATABASE IF NOT EXISTS evedata;

USE evedata;

CREATE TABLE `alliances` (
  `allianceID` int(10) unsigned NOT NULL DEFAULT '0',
  `name` varchar(255) NOT NULL DEFAULT '',
  `shortName` varchar(45) NOT NULL DEFAULT '',
  `executorCorpID` int(10) unsigned NOT NULL DEFAULT '0',
  `corporationsCount` int(10) unsigned NOT NULL DEFAULT '0',
  `startDate` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `updated` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `deleted` tinyint(4) NOT NULL DEFAULT '0',
  `description` text NOT NULL,
  `creatorCorpID` int(10) NOT NULL,
  `creatorCharacter` int(10) NOT NULL,
  `cacheUntil` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `memberCount` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`allianceID`),
  KEY `name` (`name`),
  KEY `shortName` (`shortName`),
  KEY `executorCorpID` (`executorCorpID`),
  KEY `cacheUntil` (`cacheUntil`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `assets` (
  `locationID` bigint(20) unsigned NOT NULL,
  `typeID` smallint(5) unsigned NOT NULL DEFAULT '0',
  `quantity` int(10) unsigned NOT NULL DEFAULT '0',
  `characterID` int(10) unsigned NOT NULL DEFAULT '0',
  `locationFlag` varchar(40) NOT NULL,
  `itemID` bigint(20) unsigned NOT NULL,
  `locationType` varchar(40) NOT NULL,
  `isSingleton` tinyint(1) unsigned NOT NULL,
  PRIMARY KEY (`itemID`),
  KEY `locationID` (`locationID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `characters` (
  `characterID` int(11) unsigned NOT NULL DEFAULT '0',
  `bloodlineID` tinyint(3) unsigned NOT NULL,
  `ancestryID` tinyint(3) unsigned NOT NULL,
  `corporationID` int(11) unsigned NOT NULL,
  `allianceID` int(11) unsigned NOT NULL,
  `race` char(8) CHARACTER SET latin1 NOT NULL,
  `securityStatus` decimal(4,2) NOT NULL,
  `updated` datetime NOT NULL,
  `cacheUntil` datetime NOT NULL,
  `name` varchar(50) NOT NULL,
  `gender` varchar(50) NOT NULL,
  PRIMARY KEY (`characterID`),
  KEY `cacheUntil` (`cacheUntil`),
  KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `contactSyncs` (
  `characterID` int(11) NOT NULL,
  `source` int(11) NOT NULL,
  `destination` int(11) NOT NULL,
  `lastError` varchar(100) CHARACTER SET latin1 DEFAULT NULL,
  `nextSync` datetime NOT NULL DEFAULT '2016-10-21 01:57:40',
  PRIMARY KEY (`characterID`,`destination`,`source`),
  UNIQUE KEY `destination_UNIQUE` (`destination`),
  KEY `source` (`source`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `corporations` (
  `corporationID` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `ticker` varchar(45) CHARACTER SET latin1 NOT NULL,
  `ceoID` int(11) NOT NULL,
  `stationID` int(11) NOT NULL,
  `description` text NOT NULL,
  `allianceID` int(11) NOT NULL,
  `factionID` int(11) NOT NULL,
  `url` text CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `memberCount` int(11) NOT NULL,
  `shares` bigint(20) NOT NULL,
  `updated` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `cacheUntil` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  PRIMARY KEY (`corporationID`),
  KEY `allianceID` (`allianceID`),
  KEY `factionID` (`factionID`),
  KEY `cacheUntil` (`cacheUntil`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `crestID` (
  `id` int(11) NOT NULL,
  `crestRef` varchar(60) CHARACTER SET latin1 NOT NULL,
  `type` enum('alliance','corporation','character') CHARACTER SET latin1 NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  UNIQUE KEY `crestRef_UNIQUE` (`crestRef`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `crestTokens` (
  `characterID` int(11) NOT NULL,
  `tokenCharacterID` int(11) NOT NULL,
  `accessToken` text CHARACTER SET latin1 NOT NULL,
  `refreshToken` text CHARACTER SET latin1 NOT NULL,
  `expiry` datetime NOT NULL,
  `tokenType` varchar(100) CHARACTER SET latin1 NOT NULL,
  `lastCode` int(11) NOT NULL DEFAULT '0',
  `lastStatus` text CHARACTER SET latin1 NOT NULL,
  `characterName` varchar(100) NOT NULL,
  `request` text CHARACTER SET latin1,
  `response` text CHARACTER SET latin1,
  `assetCacheUntil` datetime NOT NULL DEFAULT '2015-01-01 00:00:00',
  `walletCacheUntil` datetime NOT NULL DEFAULT '2015-01-01 00:00:00',
  `scopes` text NOT NULL,
  PRIMARY KEY (`characterID`,`tokenCharacterID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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
  KEY `alliance` (`allianceID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `entityKillStats` (
  `id` int(10) unsigned NOT NULL,
  `kills` int(10) unsigned NOT NULL DEFAULT '0',
  `losses` int(10) unsigned NOT NULL DEFAULT '0',
  `efficiency` float NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `httpErrors` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(255) DEFAULT NULL,
  `status` smallint(6) DEFAULT NULL,
  `request` text,
  `response` text,
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1972 DEFAULT CHARSET=utf8;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `jitaPrice` (
  `itemID` int(10) unsigned NOT NULL DEFAULT '0',
  `buy` decimal(15,2) unsigned DEFAULT NULL,
  `sell` decimal(15,2) unsigned DEFAULT NULL,
  `high` decimal(15,2) DEFAULT NULL,
  `low` decimal(15,2) DEFAULT NULL,
  `mean` decimal(15,2) DEFAULT NULL,
  `quantity` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`itemID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `jumps` (
  `toSolarSystemID` int(10) unsigned NOT NULL DEFAULT '0',
  `fromSolarSystemID` int(10) unsigned NOT NULL DEFAULT '0',
  `jumps` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `secureJumps` smallint(5) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`toSolarSystemID`,`fromSolarSystemID`),
  KEY `fromSolarSystemID` (`fromSolarSystemID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `killmailAttackers` (
  `id` int(10) unsigned NOT NULL,
  `characterID` int(10) unsigned NOT NULL DEFAULT '0',
  `corporationID` int(10) unsigned NOT NULL DEFAULT '0',
  `allianceID` int(10) unsigned NOT NULL DEFAULT '0',
  `shipType` smallint(5) unsigned NOT NULL DEFAULT '0',
  `finalBlow` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `damageDone` int(8) unsigned NOT NULL DEFAULT '0',
  `weaponType` smallint(5) unsigned NOT NULL DEFAULT '0',
  `securityStatus` decimal(4,2) NOT NULL DEFAULT '0.00',
  PRIMARY KEY (`id`,`characterID`),
  KEY `allianceID` (`allianceID`),
  KEY `corporationID` (`corporationID`),
  KEY `characterID` (`characterID`),
  KEY `charcorp` (`characterID`,`corporationID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `killmailItems` (
  `id` int(10) unsigned NOT NULL,
  `itemType` smallint(5) unsigned NOT NULL DEFAULT '0',
  `flag` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `quantityDestroyed` int(8) unsigned NOT NULL DEFAULT '0',
  `quantityDropped` int(8) unsigned NOT NULL DEFAULT '0',
  `singleton` tinyint(3) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`,`itemType`),
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `killmails` (
  `id` int(9) unsigned NOT NULL,
  `solarSystemID` int(8) unsigned NOT NULL,
  `killTime` datetime NOT NULL,
  `victimCharacterID` int(10) unsigned NOT NULL DEFAULT '0',
  `victimCorporationID` int(10) unsigned NOT NULL DEFAULT '0',
  `victimAllianceID` int(10) unsigned NOT NULL DEFAULT '0',
  `hash` char(40) CHARACTER SET latin1 NOT NULL,
  `attackerCount` smallint(3) unsigned NOT NULL DEFAULT '0',
  `damageTaken` int(9) unsigned NOT NULL DEFAULT '0',
  `x` float NOT NULL,
  `y` float NOT NULL,
  `z` float NOT NULL,
  `shipType` mediumint(8) unsigned NOT NULL DEFAULT '0',
  `warID` mediumint(8) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `victimAllianceID` (`victimAllianceID`),
  KEY `victimCorporationID` (`victimCorporationID`),
  KEY `killTime` (`killTime`),
  KEY `war` (`warID`),
  KEY `victimCharacterID` (`victimCharacterID`),
  KEY `time_victim` (`killTime`,`victimAllianceID`,`victimCorporationID`,`warID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `lpOfferRequirements` (
  `offerID` int(11) NOT NULL,
  `typeID` int(11) NOT NULL,
  `quantity` int(11) NOT NULL,
  PRIMARY KEY (`offerID`,`typeID`),
  KEY `lpReqs_offerID` (`offerID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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
) ENGINE=InnoDB AUTO_INCREMENT=16275 DEFAULT CHARSET=utf8;

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
  KEY `stationID` (`stationID`),
  KEY `stationID_bid_reported` (`stationID`,`bid`,`reported`),
  KEY `duration_issued` (`duration`,`issued`),
  KEY `reported` (`reported`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `marketStations` (
  `stationName` varchar(255) DEFAULT NULL,
  `stationID` bigint(20) unsigned NOT NULL DEFAULT '0',
  `Count` bigint(21) NOT NULL DEFAULT '0',
  PRIMARY KEY (`stationID`),
  KEY `count` (`Count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `market_history` (
  `date` date NOT NULL,
  `low` decimal(20,2) unsigned NOT NULL,
  `high` decimal(20,2) unsigned NOT NULL,
  `mean` decimal(20,2) unsigned NOT NULL,
  `quantity` int(10) unsigned NOT NULL,
  `orders` smallint(4) unsigned NOT NULL,
  `itemID` smallint(5) unsigned NOT NULL,
  `regionID` mediumint(8) unsigned NOT NULL,
  PRIMARY KEY (`date`,`regionID`,`itemID`),
  KEY `regionIDDate` (`regionID`,`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `market_vol` (
  `number` bigint(21) NOT NULL DEFAULT '0',
  `quantity` decimal(36,4) DEFAULT NULL,
  `regionID` int(11) NOT NULL,
  `itemID` bigint(20) NOT NULL,
  PRIMARY KEY (`regionID`,`itemID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `states` (
  `state` varchar(45) NOT NULL,
  `value` int(11) NOT NULL,
  `nextCheck` datetime NOT NULL,
  PRIMARY KEY (`state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `structures` (
  `stationID` bigint(20) NOT NULL,
  `solarSystemID` int(11) DEFAULT NULL,
  `stationName` varchar(255) CHARACTER SET utf8 DEFAULT NULL,
  `x` float DEFAULT NULL,
  `y` float DEFAULT NULL,
  `z` float DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `marketCacheUntil` datetime DEFAULT '2016-01-01 00:00:00',
  PRIMARY KEY (`stationID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

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
  PRIMARY KEY (`refID`,`argName1`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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
  `clientTypeID` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`transactionID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `warAllies` (
  `id` mediumint(8) unsigned NOT NULL,
  `allyID` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`,`allyID`),
  KEY `allyID` (`allyID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `wars` (
  `id` mediumint(8) unsigned NOT NULL,
  `timeFinished` datetime DEFAULT NULL,
  `timeStarted` datetime DEFAULT NULL,
  `timeDeclared` datetime DEFAULT NULL,
  `openForAllies` tinyint(4) unsigned DEFAULT NULL,
  `cacheUntil` datetime DEFAULT NULL,
  `aggressorID` int(11) unsigned DEFAULT NULL,
  `defenderID` int(11) unsigned DEFAULT NULL,
  `mutual` tinyint(4) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `aggressorID` (`aggressorID`),
  KEY `defenderID` (`defenderID`),
  KEY `timeFinished_cacheUntil` (`timeFinished`,`cacheUntil`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


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
				FROM eve.mapSolarSystems M
				INNER JOIN evedata.structures S ON S.solarSystemID = M.solarSystemID
				WHERE stationID = structure
				LIMIT 1;
			
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
		