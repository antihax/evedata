var Killmail = function (id, completeFunc) {
    var kmURL = "https://static.evedata.org/file/evedata-killmails/",
        imgURL = "//imageserver.eveonline.com/",
        killmail;

    $.ajax({
        url: kmURL + id + ".json.gz",
        dataType: 'native',
        xhrFields: {
            responseType: 'arraybuffer'
        },
        success: function (d) {
            killmail = $.parseJSON(pako.inflate(d, { to: 'string' }));
            completeFunc(returnPackage());
        },
        failure: function (d) {
            throw "Don't know this killmail.";
        }
    });

    function returnPackage() {
        return {
            getKillmail: function () { return killmail },
            getVictim: function () { return killmail.killmail.victim },
            getAttackerCount: function () { return killmail.killmail.attackers.length },
            getName: function (id) { return killmail.nameMap[id] },
            getPrice: function (id) { return killmail.priceMap[id] != undefined ? killmail.priceMap[id] : 0 },
            getShipAttributes: function () { return killmail.attributes.ship },
            getPortait: function (a) {
                if (a.character_id != undefined) {
                    return "character/" + a.character_id + "_64.jpg";
                } else if (a.corporation_id != undefined) {
                    return "corporation/" + a.corporation_id + "_64.png";
                } else {
                    return "corporation/" + a.faction_id + "_64.png";
                }
            },
            convertMS: function (ms) {
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
            },
            resolveEntity: function (a) {
                if (a.character_id != undefined) {
                    return a.character_id;
                } else if (a.corporation_id != undefined) {
                    return a.corporation_id;
                } else {
                    return a.faction_id;
                }
            },

            formatValue: function (v) {
                switch (true) {
                    case v > 1000000000000:
                        return {
                            value: v / 1000000000000,
                            indicator: "t",
                        };
                    case v > 1000000000:
                        return {
                            value: v / 1000000000,
                            indicator: "b",
                        };
                    case v > 1000000:
                        return {
                            value: v / 1000000,
                            indicator: "m",
                        };
                    case v > 1000:
                        return {
                            value: v / 1000,
                            indicator: "k",
                        };
                    default:
                        return {
                            value: v,
                            indicator: "",
                        };
                }
                return v
            },

            getMailValue: function () {
                var v = killmail.priceMap,
                    il = killmail.killmail.victim.items,
                    ship = killmail.killmail.victim.ship_type_id,
                    value = 0;

                $.each(il, function (i, a) {
                    $.each(a.items, function (i, a) {
                        if (v[a.item_type_id] != undefined) {
                            value += v[a.item_type_id] *
                                (a.quantity_destroyed != undefined ? a.quantity_destroyed : a.quantity_dropped);
                        }
                    });
                    if (v[a.item_type_id] != undefined) {
                        value += v[a.item_type_id] *
                            (a.quantity_destroyed != undefined ? a.quantity_destroyed : a.quantity_dropped);
                    }
                });
                if (v[ship] != undefined) {
                    value += v[ship];
                }
                return value;
            },
        }
    }
}

