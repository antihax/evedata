	function capitalizeFirstLetter(string) {
	    return string.charAt(0).toUpperCase() + string.slice(1);
	}
	
	function characterImage(row) {
	    return '//imageserver.eveonline.com/character/' + row.characterID + '_32.jpg';
	}

	function characterImageByID(characterID, size) {
	    return '//imageserver.eveonline.com/character/' + characterID + '_'+ size + '.jpg';
	}

	function entityImage(row) {
		if (row.type == "character") {
			return '//imageserver.eveonline.com/'+ capitalizeFirstLetter(row.type) + '/' + row.id + '_32.jpg';
		} else {
	    	return '//imageserver.eveonline.com/'+ capitalizeFirstLetter(row.type) + '/' + row.id + '_32.png';
		}
	}

	function entityTypeImage(row) {
		if (row.type == "character") {
			return '//imageserver.eveonline.com/'+ capitalizeFirstLetter(row.entityType) + '/' + row.entityID + '_32.jpg';
		} else {
	    	return '//imageserver.eveonline.com/'+ capitalizeFirstLetter(row.entityType) + '/' + row.entityID + '_32.png';
		}
	}

	function typeImage(row) {
	    return '//imageserver.eveonline.com/Type/' + row.typeID + '_32.png';
	}

	function allianceFormatter(value, row) {
		if (row.allianceID == null) {
			return '';
		}
    	var entityURL = '/alliance?id=' + row.allianceID;
    	return '<a href="' + entityURL + '">' + value + '</a>';
	}

	function corporationFormatter(value, row) {
    	var entityURL = '/corporation?id=' + row.corporationID;
    	return '<a href="' + entityURL + '">' + value + '</a>';
	}

    function zkillFormatter(value, row) {
    	var killURL = 'https://zkillboard.com/kill/' + row.killID +  '/'
    	return '<a href="' + killURL + '" target="zkill">View on ZKillBoard</a>';
    }

    function aggressorFormatter(value, row) {
    	var warURL = 'https://zkillboard.com/war/' + row.warID +  '/'
    	var entityURL = '/'+ row.aggressorType + '?id=' + row.aggressorID 
        return '<a href="' + warURL + '" target="zkill"><i class="glyphicon glyphicon-star"></i></a> <a href="' + entityURL + '">' + value + '</a>';
    }

    function defenderFormatter(value, row) {
    	var warURL = 'https://zkillboard.com/war/' + row.warID +  '/'
    	var entityURL = '/'+ row.defenderType + '?id=' + row.defenderID 
        return '<a href="' + warURL + '" target="zkill"><i class="glyphicon glyphicon-star"></i></a> <a href="' + entityURL + '">' + value + '</a>';
    }

    function capabilityFormatter(value, row) {
        return Number((row.efficiency * 100).toFixed(0))  + '% Efficiency; Kills: ' + row.kills + '; Losses: ' + row.losses;
    }

    function warsFormatter(value, row) {
        return 'Aggressor: ' + row.warAggressor + "  Defender: " + row.warDefender;
    }

    function entityFormatter(value, row) {
    	var entityURL = '/'+ row.type + '?id=' + row.id 
        return '<a href="' + entityURL + '"><img src="' + entityImage(row) + '" height=32 width=32> ' + value + '</a>';
	}

	function entityTypeFormatter(value, row) {
    	var entityURL = '/'+ row.type + '?id=' + row.entityID 
        return '<a href="' + entityURL + '"><img src="' + entityTypeImage(row) + '" height=32 width=32> ' + value + '</a>';
	}
	
	function tokenCharacterFormatter(value, row) {
        return '<img src="' + characterImageByID(row.tokenCharacterID, 32) + '" height=32 width=32> ' + value + '</a>';
    }

    function characterFormatter(value, row) {
        return '<img src="' + characterImage(row) + '" height=32 width=32 alt="'+row.characterName+'">';
	}
	
	function characterFormatterName(value, row) {
        return '<a href="/character?id=' + row.characterID +'"><img src="' + characterImage(row) + '" height=32 width=32 alt="'+row.characterName+'">'+row.characterName+'</a>';
    }

    function owner1Formatter(value, row) {
        return '<img src="' + characterImageByID(row.ownerID1, 32) + '" height=32 width=32 alt="'+row.ownerName1+'"> &nbsp;'+row.ownerName1;
    }

    function owner2Formatter(value, row) {
        return '<img src="' + characterImageByID(row.ownerID2, 32) + '" height=32 width=32 alt="'+row.ownerName2+'"> &nbsp;'+row.ownerName2;
    }

    function stationFormatter(value, row) {
    	return '<a data-toggle="tooltip" title="Set Destination" href="javascript:setEVEDestination(' + row.stationID + ')"><span class="glyphicon glyphicon-circle-arrow-right"></span></a>'
		+ '&nbsp;<a data-toggle="tooltip" title="Add Destination" href="javascript:addEVEDestination(' + row.stationID + ')"><span class="glyphicon glyphicon-plus-sign"></span></a>'
		+ '&nbsp;&nbsp;' + value;
    }

    function typeFormatter(value, row) {
    	var typeURL = '/item?id=' + row.typeID
    	return '<a data-toggle="tooltip" title="Open market in-game" href="javascript:openMarketWindow(' + row.typeID + ')"><span class="glyphicon glyphicon-circle-arrow-right"></span></a>'
		+ '&nbsp;&nbsp;<a href="' + typeURL + '"><img src="' + typeImage(row) + '" height=32 width=32> ' + value + '</a>';
    }

    function currencyFormatter(value, row) {
		return numberCommafy(value.toFixed(2));
    }

    function numberFormatter(value, row) {
		return numberCommafy(value);
    }

	function sumFormatter(data) {
		field = this.field;
		return numberCommafy(data.reduce(function(sum, row) { 
			return sum + (+row[field]);
		}, 0).toFixed(2));
	}

	function totalTextFormatter(data) {
		return 'Total';
	}

	function numberCommafy(nStr) {
		nStr += '';
		var x = nStr.split('.');
		var x1 = x[0];
		var x2 = x.length > 1 ? '.' + x[1] : '';
		var rgx = /(\d+)(\d{3})/;
		while (rgx.test(x1)) {
			x1 = x1.replace(rgx, '$1' + ',' + '$2');
		}
		return x1 + x2;
	}

	function showAlert(message, type) {
		$.growl(message,{
			// settings
			type: type,
			delay: 4000,
		});
	}

	function openMarketWindow(id) {
		if (accountInfo.cursor && accountInfo.cursor.cursorCharacterID > 0) {
			$.ajax({
				url: "/X/openMarketWindow?typeID=" + id,
				type: 'POST',
			});
		} else {
			showAlert('No characters available with UI Control. Please add characters on the account page with at least one with UI control.', 'danger');
		}
	}

	function setEVEDestination(id) {
		if (accountInfo.cursor && accountInfo.cursor.cursorCharacterID > 0) {
			$.ajax({
				url: "/X/setDestination?destinationID=" + id,
				type: 'POST',
			});
		} else {
			showAlert('No characters available with UI Control. Please add characters on the account page with at least one with UI control.', 'danger');
		}
	}

	function addEVEDestination(id) {
		if (accountInfo.cursor && accountInfo.cursor.cursorCharacterID > 0) {
			$.ajax({
				url: "/X/addDestination?destinationID=" + id,
				type: 'POST',
			});
		} else {
			showAlert('No characters available with UI Control. Please add characters on the account page with at least one with UI control.', 'danger');
		}
	}