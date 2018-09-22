function capitalizeFirstLetter(string) {
	if (string) {
		return string.charAt(0).toUpperCase() + string.slice(1);
	}
	return string;
}

Array.prototype.inArray = function (comparer) {
	for (var i = 0; i < this.length; i++) {
		if (comparer(this[i])) return true;
	}
	return false;
};

function dateFormatter(value, row) {
	let date = new Date(value);
	if (date.getTime() == 0) { 
		return ""
	}
	return date.toUTCString();
}

function convertMS(ms) {
	var d, h, m, s;
	s = Math.floor(ms / 1000);
	m = Math.floor(s / 60);
	s = s % 60;
	h = Math.floor(m / 60);
	m = m % 60;
	d = Math.floor(h / 24);
	h = h % 24;
	h += d * 24;
	m += h * 60;
	return m + 'm ' + s + "s";
}

Array.prototype.pushIfNotExist = function (element, comparer) {
	if (!this.inArray(comparer)) {
		this.push(element);
	}
};

function characterImage(row) {
	return '//imageserver.eveonline.com/character/' + row.characterID + '_32.jpg';
}

function characterImageByID(characterID, size) {
	return '//imageserver.eveonline.com/character/' + characterID + '_' + size + '.jpg';
}

function entityImage(row) {
	if (row.type == "character") {
		return '//imageserver.eveonline.com/' + capitalizeFirstLetter(row.type) + '/' + row.id + '_32.jpg';
	} else {
		return '//imageserver.eveonline.com/' + capitalizeFirstLetter(row.type) + '/' + row.id + '_32.png';
	}
}

function entityTypeImage(row) {
	if (row.type == "character") {
		return '//imageserver.eveonline.com/' + capitalizeFirstLetter(row.entityType) + '/' + row.entityID + '_32.jpg';
	} else {
		return '//imageserver.eveonline.com/' + capitalizeFirstLetter(row.entityType) + '/' + row.entityID + '_32.png';
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
	var killURL = 'https://zkillboard.com/kill/' + row.killID + '/'
	return '<a href="' + killURL + '" target="zkill">View on ZKillBoard</a>';
}

function aggressorFormatter(value, row) {
	var warURL = 'https://zkillboard.com/war/' + row.warID + '/'
	var entityURL = '/' + row.aggressorType + '?id=' + row.aggressorID
	return '<a href="' + warURL + '" target="zkill"><i class="glyphicon glyphicon-star"></i></a> <a href="' + entityURL + '">' + value + '</a>';
}

function defenderFormatter(value, row) {
	var warURL = 'https://zkillboard.com/war/' + row.warID + '/'
	var entityURL = '/' + row.defenderType + '?id=' + row.defenderID
	return '<a href="' + warURL + '" target="zkill"><i class="glyphicon glyphicon-star"></i></a> <a href="' + entityURL + '">' + value + '</a>';
}

function capabilityFormatter(value, row) {
	return 'Kills: ' + row.kills + '; Losses: ' + row.losses + '<br>' + Number((row.efficiency * 100).toFixed(0)) + '% Efficiency; ' + Number((row.capProbability * 100).toFixed(0)) + '% Hot Drop';
}

function warsFormatter(value, row) {
	if (row.warAggressor + row.warDefender == 0) {
		return 'None';
	}
	return 'Agg: ' + row.warAggressor + "<br>Def: " + row.warDefender;
}

function entityFormatter(value, row) {
	var entityURL = '/' + row.type + '?id=' + row.id
	return '<a href="' + entityURL + '"><img src="' + entityImage(row) + '" height=32 width=32> ' + value + '</a>';
}

function entityTypeFormatter(value, row) {
	var entityURL = '/' + row.entityType + '?id=' + row.entityID
	return '<a href="' + entityURL + '"><img src="' + entityTypeImage(row) + '" height=32 width=32> ' + value + '</a>';
}

function tokenCharacterFormatter(value, row) {
	return '<img class="rounded-8" src="' + characterImageByID(row.tokenCharacterID, 32) + '" height=32 width=32> ' + value + '</a>';
}

function characterFormatter(value, row) {
	return '<img class="rounded-8" src="' + characterImage(row) + '" height=32 width=32 alt="' + row.characterName + '">';
}

function characterFormatterName(value, row) {
	return '<a href="/character?id=' + row.characterID + '"><img class="rounded-8" src="' + characterImage(row) + '" height=32 width=32 alt="' + row.characterName + '">' + row.characterName + '</a>';
}

function owner1Formatter(value, row) {
	return '<img class="rounded-8" src="' + characterImageByID(row.ownerID1, 32) + '" height=32 width=32 alt="' + row.ownerName1 + '"> &nbsp;' + row.ownerName1;
}

function owner2Formatter(value, row) {
	return '<img class="rounded-8" src="' + characterImageByID(row.ownerID2, 32) + '" height=32 width=32 alt="' + row.ownerName2 + '"> &nbsp;' + row.ownerName2;
}

function stationFormatter(value, row) {
	return '<a data-toggle="tooltip" title="Set Destination" href="javascript:setEVEDestination(' + row.stationID + ')"><span class="glyphicon glyphicon-circle-arrow-right"></span></a>' +
		'&nbsp;<a data-toggle="tooltip" title="Add Destination" href="javascript:addEVEDestination(' + row.stationID + ')"><span class="glyphicon glyphicon-plus-sign"></span></a>' +
		'&nbsp;&nbsp;' + value;
}

function typeFormatter(value, row) {
	var typeURL = '/item?id=' + row.typeID
	return '<a data-toggle="tooltip" title="Open market in-game" href="javascript:openMarketWindow(' + row.typeID + ')"><span class="glyphicon glyphicon-circle-arrow-right"></span></a>' +
		'&nbsp;&nbsp;<a href="' + typeURL + '"><img class="rounded-8" src="' + typeImage(row) + '" height=25 width=25></a> &nbsp;<a href="' + typeURL + '">' + value + '</a>';
}


function currencyFormatter(value, row) {
	return numberCommafy(value.toFixed(2));
}

function numberRound0Formatter(value, row) {
	return numberCommafy(value.toFixed(0));
}

function numberFormatter(value, row) {
	return numberCommafy(value);
}

function escapeFormatter(value, row) {
	return escapeHtml(value);
}

function sumFormatter(data) {
	field = this.field;
	return numberCommafy(data.reduce(function (sum, row) {
		return sum + (+row[field]);
	}, 0).toFixed(2));
}

function totalTextFormatter(data) {
	return 'Total';
}

function simpleVal(nStr, decimals = 0) {
	if (nStr == undefined) {
		return 0;
	}
	return numberCommafy(nStr.toFixed(decimals))
}

function numberCommafy(nStr) {

	if (nStr == undefined) {
		return 0;
	}
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
	$.growl(message, {
		// settings
		type: type,
		delay: 4000,
	});
}

function getUrlVars() {
	var vars = [],
		hash;
	var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
	for (var i = 0; i < hashes.length; i++) {
		hash = hashes[i].split('=');
		vars.push(hash[0]);
		vars[hash[0]] = hash[1];
	}
	return vars;
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

var entityMap = {
	'&': '&amp;',
	'<': '&lt;',
	'>': '&gt;',
	'"': '&quot;',
	"'": '&#39;',
	'/': '&#x2F;',
	'`': '&#x60;',
	'=': '&#x3D;'
};

function escapeHtml(string) {
	return String(string).replace(/[&<>"'`=\/]/g, function (s) {
		return entityMap[s];
	});
}

function nsToTime(nanoseconds) {
	var duration = nanoseconds / 1000;
	var seconds = parseInt((duration / 1000) % 60)
		, minutes = parseInt((duration / (1000 * 60)) % 60)
		, hours = parseInt((duration / (1000 * 60 * 60)) % 24);
	minutes = (minutes < 10) ? "0" + minutes : minutes;
	seconds = (seconds < 10) ? "0" + seconds : seconds;
	return hours + "h " + minutes + "m " + seconds + "s ";
}

function killmailTypeFormatter(value, row) {
	var typeURL = '/killmail?id=' + row.id
	return '<a data-toggle="tooltip" title="Open zkillboard" href="https://zkillboard.com/kill/' + row.id + '" target="zkill"><img src="https://zkillboard.com/img/wreck.png" height=16 width=16></a>' +
		'&nbsp;&nbsp;<a href="' + typeURL + '" target="zkill"><img src="' + typeImage(row) + '" height=25 width=25></a> &nbsp;<a href="' + typeURL + '" target="zkill">' + value + '</a>';
}

function kmCapacitorFormatter(value, row) {
	if (row.capacitorNoMWD > 0) {
		return (row.capacitorNoMWD * 100).toFixed(0) + "%";
	} else {
		return nsToTime(row.capacitorTimeNoMWD);
	}
}


