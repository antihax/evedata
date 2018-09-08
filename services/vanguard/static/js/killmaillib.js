
function getKillmail(id) {
    $.ajax({
        url: "https://static.evedata.org/file/evedata-killmails/" + id + ".json.gz",
        dataType: 'native',
        xhrFields: {
            responseType: 'arraybuffer'
        },
        success: function (d) {
            var km = $.parseJSON(pako.inflate(d, { to: 'string' }));
            printKillmail(km);
        },
        failure: function (d) {
            throw "Don't know this killmail.";
        }
    });
}

function getPortrait(a) {
    if (a.character_id != undefined) {
        return "character/" + a.character_id + "_64.jpg";
    } else if (a.corporation_id != undefined) {
        return "corporation/" + a.corporation_id + "_64.png";
    } else {
        return "corporation/" + a.faction_id + "_64.png";
    }
}

function resolveEntity(a) {
    if (a.character_id != undefined) {
        return a.character_id;
    } else if (a.corporation_id != undefined) {
        return a.corporation_id;
    } else {
        return a.faction_id;
    }
}

function getMailValue(k) {
    var v = k.valueMap,
        il = k.killmail.victim.items;
        
}

function printKillmail(k) {
    console.log(k)
    var h = `
    <div style="float: left; width: 32px">
        <img src="//imageserver.eveonline.com/type/${k.killmail.victim.ship_type_id}_32.png" 
        title="${k.nameMap[k.killmail.victim.ship_type_id]}" style="width:32px; height: 32px">
    </div>
    <div style="float: left; width: 32px">
        <img src="//imageserver.eveonline.com/${getPortrait(k.killmail.victim)}" 
            title="${k.nameMap[resolveEntity(k.killmail.victim)]}" style="width:32px; height: 32px">
    </div>
`;

    $("#killmail" + k.killmail.killmail_id).html(h);
}

function killmailFormatter(value, row) {
    getKillmail(row.id);
    return `<div id="killmail${row.id}" style="border: #000 1px; padding: 0px;"></div>`;
}