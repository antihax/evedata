   function corporationFormatter(value, row) {
    	var entityURL = '/corporation?id=' + row.corporationID;
        return '<a href="' + entityURL + '">' + value + '</a>';
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
        return '<a href="' + entityURL + '">' + value + '</a>';
    }