var package,
    urlVars = getUrlVars(),
    ship, canvas, camera;

$.ajax({
    url: "https://static.evedata.org/file/evedata-killmails/" + urlVars["id"] + ".json.gz",
    dataType: 'native',
    xhrFields: {
        responseType: 'arraybuffer'
    },
    success: function (d) {
        try {
            package = $.parseJSON(pako.inflate(d, { to: 'string' }));
            $(document).ready(function () {
                getShip(package);
                populateModules(package)
                getAttackers(package);
                getTypes(package);
                getVictimInformation(package.killmail.victim);
                getSystemInfo(package.systemInfo);
            });
        } catch {
            showAlert("Failed to decode killmail package", "danger")
        }
    },
    failure: function (d) {
        showAlert("Don't know this killmail yet", "danger")
    }
});

function setResonancePercentage(resonance, value) {
    if (!value) {
        value = 1;
    }
    value = 1 - value;
    value = (value * 100).toFixed(0);
    $('#' + resonance).css('width', value + '%').attr('aria-valuenow', value);
    $('#' + resonance).text(value + "%");
}

function setModuleSlot(type, slot, i) {
    typeNames = package.attributes.types;
    $("#" + type + slot).prepend('<img src="//imageserver.eveonline.com/Type/' + i.typeID + '_32.png" title="' + typeNames[i.typeID] + '" style="height: 32px; width: 32px;">')
    if (i.chargeTypeID > 0) $("#" + type + slot + "l").prepend('<img src="//imageserver.eveonline.com/Type/' + i.chargeTypeID + '_32.png" style="height: 32px; width: 32px;">')
}

function rc() {
    return (Math.random() >= 0.5 ? 1 : -1) * Math.random() * 1000000;
}

function cycleModule(slot) {
    var vec4 = ccpwgl_int.math.vec4,
        quat = ccpwgl_int.math.quat,
        mat4 = ccpwgl_int.math.mat4;
    var viewProjInv = ship.getTransform();

    var pt = quat.fromValues(rc(), rc(), 4000000, 1);

    ship.setTurretTargetPosition(slot, pt);
    ship.setTurretState(slot, ccpwgl.TurretState.FIRING);

    setTimeout(cycleModule, Math.random() * 3000, slot);
}

function setModule(slot, type) {
    if (graphicsMap[type]) {
        ship.mountTurret(slot, graphicsMap[type]);
        setTimeout(cycleModule, 100, slot);
    }
}
function populateModules(package) {
    typeNames = package.attributes.types;

    $.each(package.attributes.modules, function (k, i) {
        switch (i.location) {
            case 27: // hiSlots
                setModuleSlot("high", "1", i);
                setModule(1, i.typeID)
                break;
            case 28: // hiSlots
                setModuleSlot("high", "2", i);
                setModule(2, i.typeID)
                break;
            case 29: // hiSlots
                setModuleSlot("high", "3", i);
                setModule(3, i.typeID)
                break;
            case 30: // hiSlots
                setModuleSlot("high", "4", i);
                setModule(4, i.typeID)
                break;
            case 31: // hiSlots
                setModuleSlot("high", "5", i);
                setModule(5, i.typeID)
                break;
            case 32: // hiSlots
                setModuleSlot("high", "6", i);
                setModule(6, i.typeID)
                break;
            case 33: // hiSlots
                setModuleSlot("high", "7", i);
                setModule(7, i.typeID)
                break;
            case 34: // hiSlots
                setModuleSlot("high", "8", i);
                setModule(8, i.typeID)
                break;

            case 19:
                setModuleSlot("mid", "1", i);
                break;
            case 20:
                setModuleSlot("mid", "2", i);
                break;
            case 21:
                setModuleSlot("mid", "3", i);
                break;
            case 22:
                setModuleSlot("mid", "4", i);
                break;
            case 23:
                setModuleSlot("mid", "5", i);
                break;
            case 24:
                setModuleSlot("mid", "6", i);
                break;
            case 25:
                setModuleSlot("mid", "7", i);
                break;
            case 26:
                setModuleSlot("mid", "8", i);
                break;

            case 11:
                setModuleSlot("low", "1", i);
                break;
            case 12:
                setModuleSlot("low", "2", i);
                break;
            case 13:
                setModuleSlot("low", "3", i);
                break;
            case 14:
                setModuleSlot("low", "4", i);
                break;
            case 15:
                setModuleSlot("low", "5", i);
                break;
            case 16:
                setModuleSlot("low", "6", i);
                break;
            case 17:
                setModuleSlot("low", "7", i);
                break;
            case 18:
                setModuleSlot("low", "8", i);
                break;

            case 92:
                setModuleSlot("rig", "1", i);
                break;

            case 93:
                setModuleSlot("rig", "2", i);
                break;
            case 94:
                setModuleSlot("rig", "3", i);
                break;

            case 164:
                setModuleSlot("sub", "1", i);
                break;
            case 165:
                setModuleSlot("sub", "2", i);
                break;
            case 166:
                setModuleSlot("sub", "3", i);
                break;
            case 167:
                setModuleSlot("sub", "4", i);
                break;
            case 168:
                setModuleSlot("sub", "5", i);
                break;
        }
    });
}

function getCorporationImage(a) {
    if (a.corporation_id != undefined) {
        return "corporation/" + a.corporation_id + "_32.png";
    } else {
        return "corporation/" + a.faction_id + "_32.png";
    }
}
function getAllianceImage(a) {
    if (a.alliance_id != undefined) {
        return "alliance/" + a.alliance_id + "_64.png";
    } else if (a.faction_id != undefined) {
        return "corporation/" + a.faction_id + "_64.png";
    } else {
        return "corporation/" + a.corporation_id + "_64.png";
    }
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

function getTypeImage(typeID) {
    return "type/" + typeID + "_64.png";
}

function getShipImage(a) {
    return a.ship_type_id + "_32.png";
}

function getWeaponImage(a) {
    return a.weapon_type_id == undefined ?
        a.ship_type_id + "_32.png" :
        a.weapon_type_id + "_32.png";
}

function getSystemInfo(a) {
    if (a.security == undefined) a.security = 0;
    var theStuff = "";
    theStuff += `Location: ${a.solarSystemName} (${a.security.toFixed(1)}) ${a.regionName}<br>`;

    if (a.celestialID != undefined) theStuff += "Near: " + a.celestialName + "<br>";
    $("#sysInfo").html(`${theStuff}`);
}

function getVictimInformation(a) {
    $("#victimImage").attr('src', `//imageserver.eveonline.com/${getPortrait(a)}`)

    $("#victimCorporationImage").attr('src', `//imageserver.eveonline.com/${getCorporationImage(a)}`)
    $("#victimAllianceImage").attr('src', `//imageserver.eveonline.com/${getAllianceImage(a)}`)

    $("#victim").html(getCharacterInformation(a));
}

function getLink(id, type) {
    return `<a href="/${type}?id=${id}">${package.nameMap[id]}</a>`;
}

function getCharacterInformation(a) {
    var theStuff = "";

    if (a.character_id != undefined) theStuff += getLink(a.character_id, "character") + "<br>";
    if (a.corporation_id != undefined) theStuff += getLink(a.corporation_id, "corporation") + "<br>"
    if (a.alliance_id != undefined) theStuff += getLink(a.alliance_id, "alliance") + "<br>"
    if (a.faction_id != undefined) theStuff += "<small>" + package.nameMap[a.faction_id] + "</small>"

    return theStuff;
}

function getAttackers(package) {
    $("#numInvolved").text(simpleVal(package.killmail.attackers.length) + " Involved");
    var stripe = false,
        totalDamaged = 0;

    $.each(package.killmail.attackers, function (k, a) {
        if (a.damage_done != undefined)
            totalDamaged += a.damage_done;

        var row = `
            <div class="row killmail" style="background-color: ${stripe ? "#06100a;" : "#16201a;"} padding: 0px;">
                <div style="float: left; width: 64px">
                    <img src="//imageserver.eveonline.com/${getPortrait(a)}" style="width:64px; height: 64px">
                </div>
                <div style="float: left; width: 32px">
                    <img src="//imageserver.eveonline.com/type/${getShipImage(a)}" style="width:32px; height: 32px">
                    <img src="//imageserver.eveonline.com/type/${getWeaponImage(a)}" style="width:32px; height: 32px">
                </div>
                <div style="width: 230px; float: left; ">
                    ${getCharacterInformation(a)}
                </div>
                <div style="width: 66px; float: left; text-align: right">
                    ${simpleVal(a.damage_done)}
                </div>
            </div>`;

        $("#attackers").append(row);
        stripe = !stripe;
    });
    $("#totalDamaged").html(simpleVal(totalDamaged) + " Damage");
}

function addTypeRow(typeID, dropped, quantity, value, stripe) {
    if (value == undefined)
        value = 0;
    var row = `
            <div class="row killmail" style="background-color: ${stripe ? "#06100a;" : "#16201a;"} padding: 0px;">
                <div style="float: left; width: 32px">
                    <img src="//imageserver.eveonline.com/${getTypeImage(typeID)}" style="width: 32px; height: 32px">
                </div>
                <div style="float: left; width: 182px;">
                  ${package.nameMap[typeID]}
                </div>
                <div style="width: 64px; float: left; text-align: right">
                    ${simpleVal(quantity)}
                </div>
                <div style="width: 114px; float: left; text-align: right">
                    ${simpleVal(value)}
                </div>
            </div>`;
    $("#types").append(row);
}

function getTypes(package) {
    var stripe = true,
        pm = package.priceMap,
        droppedValue = 0,
        totalValue = 0;

    if (pm[package.killmail.victim.ship_type_id]) {
        totalValue = pm[package.killmail.victim.ship_type_id];
    }

    addTypeRow(package.killmail.victim.ship_type_id, false, 1, totalValue, !stripe);

    $.each(package.killmail.victim.items, function (k, a) {
        if (pm[a.item_type_id] != undefined) {
            package.killmail.victim.items[k].value = pm[a.item_type_id] *
                (a.quantity_destroyed != undefined ? a.quantity_destroyed : a.quantity_dropped);
        } else {
            package.killmail.victim.items[k].value = 0;
        }
    });

    $.each(package.killmail.victim.items, function (k, a) {
        if (a.quantity_destroyed) {
            addTypeRow(a.item_type_id, false, a.quantity_destroyed, a.value, stripe);
            totalValue += a.value;
        } else if (a.quantity_dropped) {
            addTypeRow(a.item_type_id, true, a.quantity_dropped, a.value, stripe);
            totalValue += a.value;
            droppedValue += a.value;
        }

        stripe = !stripe;
    });
    $("#totalValue").html(simpleVal(totalValue) + " Total");
    $("#droppedValue").html(simpleVal(droppedValue) + " Dropped");
}

function resizeCanvasToDisplaySize(canvas, mult) {
    var width = Math.round(256),
        height = Math.round(256);
    if (window.innerHeight == screen.height) {
        width = screen.width;
        height = screen.height;
    }

    if (canvas.width !== width || canvas.height !== height) {
        canvas.width = width;
        canvas.height = height;
    }
}

function getShipwebGL(package) {
    $("#shipImage").attr("src", "//imageserver.eveonline.com/Render/" + package.attributes.typeID + "_256.png")
    try {
        var mat4 = ccpwgl_int.math.mat4,
            rotation = 0.0,
            direction = 0.001,
            canvas = document.getElementById('shipCanvas'),
            gl = canvas.getContext("webgl");

        ccpwgl.initialize(canvas, {});

        camera = ccpwgl.createCamera(canvas, {}, true);
        scene = ccpwgl.loadScene(sceneList[Math.floor(Math.random() * sceneList.length)]);
        ship = scene.loadShip(package.dna);
        scene.loadSun(sunList[Math.floor(Math.random() * sunList.length)]);

        ccpwgl.onPreRender = function (dt) {
            resizeCanvasToDisplaySize(canvas, window.devicePixelRatio);
            gl.viewport(0, 0, gl.canvas.width, gl.canvas.height);

            camera.rotationX += 0.01;
            camera.rotationY += direction;
            if (camera.rotationY > 1.57 & direction > 0) {
                direction = -0.001;
            } else if (camera.rotationY < -1.57 & direction < 0) {
                direction = 0.001;
            }
            if (ship.isLoaded() == true) {
                $("#shipImage").addClass("hidden");
            }
            camera.focus(ship, 5, 1);
        }
    } catch (err) {
        getShipFallback(package.attributes.typeID);
    }
}

function fullscreen() {
    var canvas = document.getElementById('shipCanvas');
    if (canvas.webkitRequestFullScreen) {
        canvas.webkitRequestFullScreen();
    }
    else {
        canvas.mozRequestFullScreen();
    }
}

function getShipFallback(typeID) {
    var canvas = document.getElementById('shipCanvas');
    canvas.parentNode.removeChild(canvas);
}

function getShip(package) {
    var a = package.attributes.ship;
    getShipwebGL(package);

    setResonancePercentage("shieldEm", a.shieldEmDamageResonance)
    setResonancePercentage("shieldTherm", a.shieldThermalDamageResonance)
    setResonancePercentage("shieldKin", a.shieldKineticDamageResonance)
    setResonancePercentage("shieldExp", a.shieldExplosiveDamageResonance)
    setResonancePercentage("armorEm", a.armorEmDamageResonance)
    setResonancePercentage("armorTherm", a.armorThermalDamageResonance)
    setResonancePercentage("armorKin", a.armorKineticDamageResonance)
    setResonancePercentage("armorExp", a.armorExplosiveDamageResonance)
    setResonancePercentage("hullEm", a.emDamageResonance)
    setResonancePercentage("hullTherm", a.thermalDamageResonance)
    setResonancePercentage("hullKin", a.kineticDamageResonance)
    setResonancePercentage("hullExp", a.explosiveDamageResonance)

    $("#ehp").text(simpleVal(a.maxEHP) + " EHP");

    $("#hullHP").text(simpleVal(a.hp) + " HP");
    $("#armorHP").text(simpleVal(a.armorHP) + " HP");
    $("#shieldHP").text(simpleVal(a.shieldCapacity) + " HP");

    $("#rps").text(simpleVal(a.maxRPS) + " EHP/s");

    $("#warpStrength").text(simpleVal(a.totalWarpScrambleStrength));
    $("#webStrength").text(simpleVal(a.stasisWebifierStrength) + "%");
    $("#targets").text(simpleVal(a.maxLockedTargets));
    $("#scanRes").text(simpleVal(a.scanResolution));

    $("#cargo").text(simpleVal(a.capacity) + " m³");

    $("#lockRange").text(simpleVal(a.maxTargetRange / 1000) + " km");
    $("#warpSpeed").text(simpleVal(a.warpSpeedMultiplier, 1) + " AU/s");

    $("#speed").text(simpleVal(a.maxVelocity) + " km/s");
    $("#mwdSpeed").text(simpleVal(a.maxVelocityMWD) + " km/s");

    $("#sigRadius").text(simpleVal(a.signatureRadius));
    $("#mwdSigRadius").text(simpleVal(a.signatureRadiusMWD));

    $.each(a, function (k, v) {
        switch (k) {
            case "hiSlots": // hiSlots
                $("#hiSlots").attr('src', '/i/fw/' + v + 'h.png');
                break;
            case "medSlots": // medSlots
                $("#medSlots").attr('src', '/i/fw/' + v + 'm.png');
                break;
            case "lowSlots": // lowSlots
                $("#lowSlots").attr('src', '/i/fw/' + v + 'l.png');
                break;
            case "rigSlots": // rigSlots
                $("#rigSlots").attr('src', '/i/fw/' + v + 'r.png');
                break;
            case "maxSubsystems": // SubSystems
                $("#subSlots").attr('src', '/i/fw/' + v + 's.png');
                break;
            case "serviceSlots": // structure services
                $("#subSlots").attr('src', '/i/fw/' + att + 's.png');
                break;
            case "capacitorCapacity":
                if (a.capacitorStable == 1.0) {
                    $("#capacitorDuration").text((100 * a.capacitorFraction).toFixed(0) + "% Stable");
                } else {
                    $("#capacitorDuration").text(convertMS(a.capacitorDuration));
                }

                $("#capacitorDetails").text(simpleVal(v) + " GJ")
                break;
            case "totalDPS":
                $("#totalDamage").html(simpleVal(v) + " DPS")
                if (a.moduleDPS) $("#offenseModule").html(simpleVal(a.moduleDPS) + " DPS<br>" + simpleVal(a.moduleAlphaDamage) + " Alpha")
                if (a.droneDPS) $("#offenseDrone").html(simpleVal(a.droneDPS) + " DPS<br>" + simpleVal(a.droneAlphaDamage) + " Alpha")
                break;

            case "scanRadarStrength":
                $("#sensorStrength").text(simpleVal(v));
                break;
            case "scanLadarStrength":
                $("#sensorStrength").text(simpleVal(v));
                break;
            case "scanMagnetometricStrength":
                $("#sensorStrength").text(simpleVal(v));
                break;
            case "scanGravimetricStrength":
                $("#sensorStrength").text(simpleVal(v));
                break;
            case "remoteArmorDamageAmountPerSecond":
                $("#remoteArmor").text(simpleVal(v) + " hp/s");
                $(".remoteRepair").removeClass('hidden');
                break;
            case "remoteStructureDamageAmountPerSecond":
                $("#remoteHull").text(simpleVal(v) + " hp/s");
                $(".remoteRepair").removeClass('hidden');
                break;
            case "remoteShieldBonusAmountPerSecond":
                $("#remoteShield").text(simpleVal(v) + " hp/s");
                $(".remoteRepair").removeClass('hidden');
                break;
            case "remotePowerTransferAmountPerSecond":
                $("#remotePower").text(simpleVal(v) + " GJ/s");
                $(".remoteRepair").removeClass('hidden');
                break;
        }
    });
}


var sceneList = [
    'res:/dx9/scene/Universe/a01_cube.red',
    'res:/dx9/scene/Universe/a02_cube.red',
    'res:/dx9/scene/Universe/a03_cube.red',
    'res:/dx9/scene/Universe/a04_cube.red',
    'res:/dx9/scene/Universe/a05_cube.red',
    'res:/dx9/scene/Universe/a06_cube.red',
    'res:/dx9/scene/Universe/a07_cube.red',
    'res:/dx9/scene/Universe/a08_cube.red',
    'res:/dx9/scene/Universe/a09_cube.red',
    'res:/dx9/scene/Universe/a10_cube.red',
    'res:/dx9/scene/Universe/a11_cube.red',
    'res:/dx9/scene/Universe/a12_cube.red',
    'res:/dx9/scene/Universe/a13_cube.red',
    'res:/dx9/scene/Universe/a14_cube.red',
    'res:/dx9/scene/Universe/a15_cube.red',
    'res:/dx9/scene/Universe/a16_cube.red',
    'res:/dx9/scene/Universe/a17_cube.red',
    'res:/dx9/scene/Universe/a18_cube.red',
    'res:/dx9/scene/Universe/c01_cube.red',
    'res:/dx9/scene/Universe/c02_cube.red',
    'res:/dx9/scene/Universe/c03_cube.red',
    'res:/dx9/scene/Universe/c04_cube.red',
    'res:/dx9/scene/Universe/c05_cube.red',
    'res:/dx9/scene/Universe/c06_cube.red',
    'res:/dx9/scene/Universe/c07_cube.red',
    'res:/dx9/scene/Universe/c08_cube.red',
    'res:/dx9/scene/Universe/c09_cube.red',
    'res:/dx9/scene/Universe/c10_cube.red',
    'res:/dx9/scene/Universe/c11_cube.red',
    'res:/dx9/scene/Universe/c12_cube.red',
    'res:/dx9/scene/Universe/c13_cube.red',
    'res:/dx9/scene/Universe/c14_cube.red',
    'res:/dx9/scene/Universe/c15_cube.red',
    'res:/dx9/scene/Universe/c16_cube.red',
    'res:/dx9/scene/Universe/c17_cube.red',
    'res:/dx9/scene/Universe/c18_cube.red',
    'res:/dx9/scene/Universe/c19_cube.red',
    'res:/dx9/scene/Universe/g01_cube.red',
    'res:/dx9/scene/Universe/g02_cube.red',
    'res:/dx9/scene/Universe/g03_cube.red',
    'res:/dx9/scene/Universe/g04_cube.red',
    'res:/dx9/scene/Universe/g05_cube.red',
    'res:/dx9/scene/Universe/g06_cube.red',
    'res:/dx9/scene/Universe/g07_cube.red',
    'res:/dx9/scene/Universe/g08_cube.red',
    'res:/dx9/scene/Universe/g09_cube.red',
    'res:/dx9/scene/Universe/g10_cube.red',
    'res:/dx9/scene/Universe/g11_cube.red',
    'res:/dx9/scene/Universe/j01_cube.red',
    'res:/dx9/scene/Universe/j02_cube.red',
    'res:/dx9/scene/Universe/m01_cube.red',
    'res:/dx9/scene/Universe/m02_cube.red',
    'res:/dx9/scene/Universe/m03_cube.red',
    'res:/dx9/scene/Universe/m04_cube.red',
    'res:/dx9/scene/Universe/m05_cube.red',
    'res:/dx9/scene/Universe/m06_cube.red',
    'res:/dx9/scene/Universe/m07_cube.red',
    'res:/dx9/scene/Universe/m08_cube.red',
    'res:/dx9/scene/Universe/m09_cube.red',
    'res:/dx9/scene/Universe/m10_cube.red',
    'res:/dx9/scene/Universe/m11_cube.red',
    'res:/dx9/scene/Universe/m12_cube.red',
    'res:/dx9/scene/Universe/m13_cube.red',
    'res:/dx9/scene/Universe/m14_cube.red',
    'res:/dx9/scene/Universe/m15_cube.red',
    'res:/dx9/scene/Universe/m16_cube.red',
    'res:/dx9/scene/Universe/m17_cube.red',
];

var sunList = [
    'res:/fisfx/lensflare/yellow_small.red',
    'res:/fisfx/lensflare/yellow.red',
    'res:/fisfx/lensflare/white_tiny.red',
    'res:/fisfx/lensflare/white.red',
    'res:/fisfx/lensflare/red.red',
    'res:/fisfx/lensflare/purple_sun.red',
    'res:/fisfx/lensflare/pink_sun_small.red',
    'res:/fisfx/lensflare/pink_hazy.red',
    'res:/fisfx/lensflare/orange_sun.red',
    'res:/fisfx/lensflare/orange_radiating.red',
    'res:/fisfx/lensflare/orange.red',
    'res:/fisfx/lensflare/blue_sun.red',
    'res:/fisfx/lensflare/blue_small.red',
    'res:/fisfx/lensflare/blue.red',
];

var graphicsMap = {
    "6": "res:/dx9/model/celestial/sun/sun_yellow_01a.red",
    "7": "res:/dx9/model/celestial/sun/sun_orange_01a.red",
    "8": "res:/dx9/model/celestial/sun/sun_red_01a.red",
    "9": "res:/dx9/model/celestial/sun/sun_blue_01a.red",
    "10": "res:/dx9/model/celestial/sun/sun_white_01a.red",
    "11": "res:/dx9/model/worldobject/planet/Earthlike.red",
    "12": "res:/dx9/model/worldobject/planet/IcePlanet.red",
    "13": "res:/dx9/model/worldobject/planet/GasGiant.red",
    "14": "res:/dx9/model/worldobject/planet/moon.red",
    "23": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "60": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCL.red",
    "202": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "203": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "204": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "205": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "206": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "207": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "208": "res:/dx9/model/turret/launcher/heavy/heavy_impact_inferno.red",
    "209": "res:/dx9/model/turret/launcher/heavy/heavy_impact_scourge.red",
    "210": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "211": "res:/dx9/model/turret/launcher/light/light_impact_inferno.red",
    "212": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "213": "res:/dx9/model/turret/launcher/light/light_impact_nova.red",
    "265": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "266": "res:/dx9/model/turret/launcher/rocket/rocket_impact_scourge.red",
    "267": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_scourge.red",
    "269": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "450": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Gatling_T1.red",
    "451": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Dual_T1.red",
    "452": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Dual_T1.red",
    "453": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Medium_T1.red",
    "454": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Medium_T1.red",
    "455": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Quad_T1.red",
    "456": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Focused_T1.red",
    "457": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Focused_T1.red",
    "458": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Heavy_T1.red",
    "459": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Heavy_T1.red",
    "460": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "461": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "462": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "463": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "464": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "482": "res:/dx9/model/turret/mining/standard/standard_t1.red",
    "483": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "484": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_125mm_T1.red",
    "485": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_150mm_T1.red",
    "486": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_200mm_T1.red",
    "487": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_250mm_T1.red",
    "488": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_280mm_T1.red",
    "489": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_180mmDual_T1.red",
    "490": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_220mm_T1.red",
    "491": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_425mm_T1.red",
    "492": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_650mm_T1.red",
    "493": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_720mm_T1.red",
    "494": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_425mmDual_T1.red",
    "495": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_650mmDual_T1.red",
    "496": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_800mmHeavy_T1.red",
    "497": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "498": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1400mm_T1.red",
    "499": "res:/dx9/model/Turret/Launcher/Light/Light_T1.red",
    "501": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "503": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "561": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_75mm_T1.red",
    "562": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Electron_T1.red",
    "563": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Ion_T1.red",
    "564": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Neutron_T1.red",
    "565": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_150mm_T1.red",
    "566": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Electron_T1.red",
    "567": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_150mmDual_T1.red",
    "568": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Neutron_T1.red",
    "569": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Ion_T1.red",
    "570": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_250mm_T1.red",
    "571": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Electron_T1.red",
    "572": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "573": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Neutron_T1.red",
    "574": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "575": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Ion_T1.red",
    "788": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "790": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "805": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",

    "1810": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "1811": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "1814": "res:/dx9/model/turret/launcher/light/light_impact_nova.red",
    "1816": "res:/dx9/model/turret/launcher/light/light_impact_inferno.red",
    "1818": "res:/dx9/model/turret/launcher/heavy/heavy_impact_scourge.red",
    "1820": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "1822": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "1824": "res:/dx9/model/turret/launcher/heavy/heavy_impact_inferno.red",
    "1826": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "1828": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "1829": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "1830": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "1832": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "1875": "res:/dx9/model/Turret/Launcher/RapidLight/RapidLight_T1.red",
    "1877": "res:/dx9/model/turret/launcher/rapidlight/rapidlight_t1.red",
    "2014": "res:/dx9/model/worldobject/planet/Ocean.red",
    "2015": "res:/dx9/model/worldobject/planet/LavaPlanet.red",
    "2016": "res:/dx9/model/worldobject/planet/SandStormPlanet.red",
    "2017": "res:/dx9/model/worldobject/planet/ThunderStormPlanet.red",
    "2063": "res:/dx9/model/WorldObject/planet/PlasmaPlanet.red",
    "2064": "res:/dx9/model/celestial/wormhole/evegate.red",
    "2129": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2130": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2131": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2132": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2133": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2134": "res:/dx9/model/Structure/Planetary/Gas/Command/CommG_T1/CommG_T1.red",
    "2135": "res:/dx9/model/Structure/Planetary/Gas/Command/CommG_T1/CommG_T1.red",
    "2136": "res:/dx9/model/Structure/Planetary/Gas/Command/CommG_T1/CommG_T1.red",
    "2137": "res:/dx9/model/Structure/Planetary/Gas/Command/CommG_T1/CommG_T1.red",
    "2138": "res:/dx9/model/Structure/Planetary/Gas/Command/CommG_T1/CommG_T1.red",
    "2139": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2140": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2141": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2142": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2143": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2144": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2145": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2146": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2147": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2148": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2149": "res:/dx9/model/Structure/Planetary/Oceanic/Command/CommO_T1/CommO_T1.red",
    "2150": "res:/dx9/model/Structure/Planetary/Oceanic/Command/CommO_T1/CommO_T1.red",
    "2151": "res:/dx9/model/Structure/Planetary/Oceanic/Command/CommO_T1/CommO_T1.red",
    "2152": "res:/dx9/model/Structure/Planetary/Oceanic/Command/CommO_T1/CommO_T1.red",
    "2153": "res:/dx9/model/Structure/Planetary/Oceanic/Command/CommO_T1/CommO_T1.red",
    "2154": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2155": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2156": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2157": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2158": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2159": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2160": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2165": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "2178": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_nova.red",
    "2180": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_scourge.red",
    "2182": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_inferno.red",
    "2188": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_mjolnir.red",
    "2210": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_nova.red",
    "2212": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "2218": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "2254": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2256": "res:/dx9/model/Structure/Planetary/Terrestrial/Spaceport/PortT_T1/PortT_T1.red",
    "2257": "res:/dx9/model/Structure/Planetary/Hostile/Storage/StorH_T1/StorH_T1.red",
    "2404": "res:/dx9/model/turret/launcher/light/light_t1.red",
    "2407": "res:/FisFX/Celestial/SpatialRift_Rt_01a.red",
    "2409": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T2/ExtT_T2.red",
    "2410": "res:/dx9/model/turret/launcher/heavy/heavy_t1.red",
    "2411": "res:/dx9/model/turret/launcher/heavy/heavy_t1.red",
    "2412": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T2/ExtT_T2.red",
    "2413": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T3/ExtH_T3.red",
    "2414": "res:/dx9/model/Structure/Planetary/Oceanic/Extractor/ExtO_T2/ExtO_T2.red",
    "2415": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T3/ExtH_T3.red",
    "2416": "res:/dx9/model/Structure/Planetary/Gas/Extractor/ExtG_T1/ExtG_T1.red",
    "2417": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T1/ExtH_T1.red",
    "2418": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T1/ExtH_T1.red",
    "2419": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T1/ExtH_T1.red",
    "2420": "res:/dx9/model/turret/launcher/torpedo/torpedo_t1.red",
    "2421": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "2422": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T3/ExtH_T3.red",
    "2423": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T3/ExtH_T3.red",
    "2424": "res:/dx9/model/Structure/Planetary/Gas/Extractor/ExtG_T1/ExtG_T1.red",
    "2425": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T1/ExtH_T1.red",
    "2426": "res:/dx9/model/Structure/Planetary/Gas/Extractor/ExtG_T1/ExtG_T1.red",
    "2427": "res:/dx9/model/Structure/Planetary/Gas/Extractor/ExtG_T1/ExtG_T1.red",
    "2428": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "2429": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "2430": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T1/ExtT_T1.red",
    "2431": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "2432": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "2433": "res:/dx9/model/Structure/Planetary/Gas/Extractor/ExtG_T1/ExtG_T1.red",
    "2434": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "2435": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T1/ExtT_T1.red",
    "2438": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "2439": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "2440": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "2441": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "2442": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "2443": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "2448": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T1/ExtH_T1.red",
    "2449": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T3/ExtT_T3.red",
    "2450": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T3/ExtT_T3.red",
    "2451": "res:/dx9/model/Structure/Planetary/Oceanic/Extractor/ExtO_T3/ExtO_T3.red",
    "2452": "res:/dx9/model/Structure/Planetary/Oceanic/Extractor/ExtO_T3/ExtO_T3.red",
    "2453": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T4/ExtT_T4.red",
    "2458": "res:/dx9/model/Structure/Planetary/Oceanic/Extractor/ExtO_T4/ExtO_T4.red",
    "2459": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T4/ExtT_T4.red",
    "2460": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T4/ExtT_T4.red",
    "2461": "res:/dx9/model/Structure/Planetary/Oceanic/Extractor/ExtO_T4/ExtO_T4.red",
    "2462": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T4/ExtT_T4.red",
    "2469": "res:/dx9/model/Structure/Planetary/Hostile/Processing/ProcH_T1/ProcH_T1.red",
    "2470": "res:/dx9/model/Structure/Planetary/Hostile/Processing/ProcH_T2/ProcH_T2.red",
    "2471": "res:/dx9/model/Structure/Planetary/Hostile/Processing/ProcH_T1/ProcH_T1.red",
    "2472": "res:/dx9/model/Structure/Planetary/Hostile/Processing/ProcH_T2/ProcH_T2.red",
    "2473": "res:/dx9/model/Structure/Planetary/Terrestrial/Processing/ProcT_T1/ProcT_T1.red",
    "2474": "res:/dx9/model/Structure/Planetary/Terrestrial/Processing/ProcT_T2/ProcT_T2.red",
    "2475": "res:/dx9/model/Structure/Planetary/Terrestrial/Processing/ProcT_T4/ProcT_T4.red",
    "2480": "res:/dx9/model/Structure/Planetary/Terrestrial/Processing/ProcT_T2/ProcT_T2.red",
    "2481": "res:/dx9/model/Structure/Planetary/Terrestrial/Processing/ProcT_T1/ProcT_T1.red",
    "2482": "res:/dx9/model/Structure/Planetary/Terrestrial/Processing/ProcT_T4/ProcT_T4.red",
    "2483": "res:/dx9/model/Structure/Planetary/Hostile/Processing/ProcH_T1/ProcH_T1.red",
    "2484": "res:/dx9/model/Structure/Planetary/Hostile/Processing/ProcH_T2/ProcH_T2.red",
    "2485": "res:/dx9/model/Structure/Planetary/Oceanic/Processing/ProcO_T2/ProcO_T2.red",
    "2490": "res:/dx9/model/Structure/Planetary/Oceanic/Processing/ProcO_T1/ProcO_T1.red",
    "2491": "res:/dx9/model/Structure/Planetary/Hostile/Processing/ProcH_T2/ProcH_T2.red",
    "2492": "res:/dx9/model/Structure/Planetary/Gas/Processing/ProcG_T1/ProcG_T1.red",
    "2493": "res:/dx9/model/Structure/Planetary/Hostile/Processing/ProcH_T1/ProcH_T1.red",
    "2494": "res:/dx9/model/Structure/Planetary/Gas/Processing/ProcG_T2/ProcG_T2.red",
    "2506": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_mjolnir.red",
    "2507": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_mjolnir.red",
    "2508": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_nova.red",
    "2510": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_inferno.red",
    "2512": "res:/dx9/model/turret/launcher/rocket/rocket_impact_mjolnir.red",
    "2514": "res:/dx9/model/turret/launcher/rocket/rocket_impact_inferno.red",
    "2516": "res:/dx9/model/turret/launcher/rocket/rocket_impact_nova.red",
    "2524": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2525": "res:/dx9/model/Structure/Planetary/Oceanic/Command/CommO_T1/CommO_T1.red",
    "2533": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2534": "res:/dx9/model/Structure/Planetary/Gas/Command/CommG_T1/CommG_T1.red",
    "2535": "res:/dx9/model/Structure/Planetary/Oceanic/Storage/StorO_T1/StorO_T1.red",
    "2536": "res:/dx9/model/Structure/Planetary/Gas/Storage/StorG_T1/StorG_T1.red",
    "2541": "res:/dx9/model/Structure/Planetary/Terrestrial/Storage/StorT_T1/StorT_T1.red",
    "2542": "res:/dx9/model/Structure/Planetary/Oceanic/Spaceport/PortO_T1/PortO_T1.red",
    "2543": "res:/dx9/model/Structure/Planetary/Gas/Spaceport/PortG_T1/PortG_T1.red",
    "2544": "res:/dx9/model/Structure/Planetary/Terrestrial/Spaceport/PortT_T1/PortT_T1.red",
    "2549": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2550": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2551": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2552": "res:/dx9/model/Structure/Planetary/Hostile/Spaceport/PortH_T1/PortH_T1.red",
    "2555": "res:/dx9/model/Structure/Planetary/Hostile/Spaceport/PortH_T1/PortH_T1.red",
    "2556": "res:/dx9/model/Structure/Planetary/Hostile/Spaceport/PortH_T1/PortH_T1.red",
    "2557": "res:/dx9/model/Structure/Planetary/Hostile/Spaceport/PortH_T1/PortH_T1.red",
    "2558": "res:/dx9/model/Structure/Planetary/Hostile/Storage/StorH_T1/StorH_T1.red",
    "2560": "res:/dx9/model/Structure/Planetary/Hostile/Storage/StorH_T1/StorH_T1.red",
    "2561": "res:/dx9/model/Structure/Planetary/Hostile/Storage/StorH_T1/StorH_T1.red",
    "2562": "res:/dx9/model/Structure/Planetary/Terrestrial/Storage/StorT_T1/StorT_T1.red",
    "2574": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2576": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2577": "res:/dx9/model/Structure/Planetary/Hostile/Command/CommH_T1/CommH_T1.red",
    "2578": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2581": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2582": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2585": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2586": "res:/dx9/model/Structure/Planetary/Terrestrial/Command/CommT_T1/CommT_T1.red",
    "2596": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "2602": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCL.red",
    "2612": "res:/dx9/model/Celestial/rock/Hollow/Hollow.red",
    "2613": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "2618": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "2620": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "2621": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "2622": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "2628": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "2629": "res:/dx9/model/turret/launcher/heavy/heavy_impact_scourge.red",
    "2637": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "2647": "res:/dx9/model/turret/launcher/light/light_impact_inferno.red",
    "2655": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "2679": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_scourge.red",
    "2684": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCL.red",
    "2686": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCL.red",
    "2692": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCL.red",
    "2712": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCL.red",
    "2801": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_nova.red",
    "2811": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_inferno.red",
    "2817": "res:/dx9/model/turret/launcher/rocket/rocket_impact_mjolnir.red",
    "2848": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T2/ExtT_T2.red",
    "2852": "res:/FisFX/Deployable/CynoBeacon_Rt_T1a.red",
    "2865": "res:/dx9/model/turret/projectile/artil/l/artil_1200mmheavy_t1.red",
    "2873": "res:/dx9/model/turret/projectile/auto/s/auto_125mm_t1.red",
    "2881": "res:/dx9/model/turret/projectile/auto/s/auto_150mm_t1.red",
    "2889": "res:/dx9/model/turret/projectile/auto/s/auto_200mm_t1.red",
    "2894": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "2897": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_220mm_T1.red",
    "2905": "res:/dx9/model/turret/projectile/artil/s/artil_250mm_t1.red",
    "2913": "res:/dx9/model/turret/projectile/auto/m/auto_425mm_t1.red",
    "2921": "res:/dx9/model/turret/projectile/artil/m/artil_650mm_t1.red",
    "2927": "res:/dx9/model/Artifact/Universal/ConcentricRings/UCR.red",
    "2929": "res:/dx9/model/turret/projectile/auto/l/auto_800mmheavy_t1.red",
    "2937": "res:/dx9/model/turret/projectile/auto/m/auto_180mmdual_t1.red",
    "2945": "res:/dx9/model/turret/projectile/auto/l/auto_425mmdual_t1.red",
    "2953": "res:/dx9/model/turret/projectile/auto/l/auto_650mmdual_t1.red",
    "2961": "res:/dx9/model/turret/projectile/artil/l/artil_1400mm_t1.red",
    "2969": "res:/dx9/model/turret/projectile/artil/m/artil_720mm_t1.red",
    "2977": "res:/dx9/model/turret/projectile/artil/s/artil_280mm_t1.red",
    "2985": "res:/dx9/model/turret/energy/beam/l/beam_heavydual_t1.red",
    "2993": "res:/dx9/model/turret/energy/beam/s/beam_dual_t1.red",
    "3001": "res:/dx9/model/turret/energy/pulse/s/pulse_dual_t1.red",
    "3009": "res:/dx9/model/turret/energy/beam/m/beam_focused_t1.red",
    "3017": "res:/dx9/model/turret/energy/pulse/s/pulse_gatling_t1.red",
    "3025": "res:/dx9/model/turret/energy/beam/m/beam_heavy_t1.red",
    "3033": "res:/dx9/model/turret/energy/beam/s/beam_medium_t1.red",
    "3041": "res:/dx9/model/turret/energy/pulse/s/pulse_medium_t1.red",
    "3049": "res:/dx9/model/turret/energy/beam/l/beam_mega_t1.red",
    "3057": "res:/dx9/model/turret/energy/pulse/l/pulse_mega_t1.red",
    "3060": "res:/dx9/model/Structure/Planetary/Gas/Extractor/ExtG_T1/ExtG_T1.red",
    "3061": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T3/ExtH_T3.red",
    "3062": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "3063": "res:/dx9/model/Structure/Planetary/Oceanic/Extractor/ExtO_T2/ExtO_T2.red",
    "3064": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T2/ExtH_T2.red",
    "3065": "res:/dx9/model/turret/energy/beam/l/beam_tachyon_t1.red",
    "3067": "res:/dx9/model/Structure/Planetary/Hostile/Extractor/ExtH_T3/ExtH_T3.red",
    "3068": "res:/dx9/model/Structure/Planetary/Terrestrial/Extractor/ExtT_T2/ExtT_T2.red",
    "3073": "res:/dx9/model/worldobject/cloud/redcloud.red",
    "3074": "res:/dx9/model/turret/hybrid/rail/s/rail_150mm_t1.red",
    "3082": "res:/dx9/model/turret/hybrid/rail/m/rail_250mm_t1.red",
    "3090": "res:/dx9/model/turret/hybrid/rail/l/rail_425mm_t1.red",
    "3098": "res:/dx9/model/turret/hybrid/rail/s/rail_75mm_t1.red",
    "3106": "res:/dx9/model/turret/hybrid/rail/m/rail_150mmdual_t1.red",
    "3114": "res:/dx9/model/turret/hybrid/rail/l/rail_250mmdual_t1.red",
    "3122": "res:/dx9/model/turret/hybrid/blast/l/blast_electron_t1.red",
    "3130": "res:/dx9/model/turret/hybrid/blast/m/blast_electron_t1.red",
    "3138": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Ion_T1.red",
    "3146": "res:/dx9/model/turret/hybrid/blast/m/blast_neutron_t1.red",
    "3154": "res:/dx9/model/turret/hybrid/blast/l/blast_ion_t1.red",
    "3162": "res:/dx9/model/turret/hybrid/blast/s/blast_electron_t1.red",
    "3170": "res:/dx9/model/turret/hybrid/blast/s/blast_ion_t1.red",
    "3178": "res:/dx9/model/turret/hybrid/blast/s/blast_neutron_t1.red",
    "3186": "res:/dx9/model/turret/hybrid/blast/l/blast_neutron_t1.red",
    "3285": "res:/dx9/model/turret/energy/beam/m/beam_quad_t1.red",
    "3293": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCM.red",
    "3296": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCL.red",
    "3297": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "3465": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCL.red",
    "3466": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCM.red",
    "3467": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCS.red",
    "3468": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "3502": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "3512": "res:/dx9/model/turret/energy/pulse/m/pulse_focused_t1.red",
    "3520": "res:/dx9/model/turret/energy/pulse/m/pulse_heavy_t1.red",
    "3546": "res:/dx9/model/turret/hybrid/blast/xl/blast_ion_t1.red",
    "3550": "res:/dx9/model/turret/hybrid/rail/xl/rail_1000mmdual_t1.red",
    "3559": "res:/dx9/model/turret/energy/pulse/xl/pulse_gigadual_t1.red",
    "3561": "res:/dx9/model/turret/energy/beam/xl/beam_gigadual_t1.red",
    "3563": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_t1.red",
    "3564": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "3565": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_t1.red",
    "3570": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "3571": "res:/dx9/model/turret/projectile/artil/xl/artil_3500mmquad_t1.red",
    "3573": "res:/dx9/model/turret/projectile/auto/xl/auto_6x2500mm_t1.red",
    "3620": "res:/dx9/model/Deployables/Starbase/ShieldGenerator/ShieldGenerator_T1.red",
    "3634": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Gatling_T1.red",
    "3636": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_125mm_T1.red",
    "3638": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_75mm_T1.red",
    "3640": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Electron_T1.red",
    "3651": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "3796": "res:/dx9/model/celestial/sun/sun_blue_sun_01a.red",
    "3797": "res:/dx9/model/celestial/sun/sun_pink_hazy_01a.red",
    "3798": "res:/dx9/model/celestial/sun/sun_orange_sun_01a.red",
    "3799": "res:/dx9/model/celestial/sun/sun_pink_sun_small_01a.red",
    "3800": "res:/dx9/model/celestial/sun/sun_orange_radiating_01a.red",
    "3801": "res:/dx9/model/celestial/sun/sun_blue_small_01a.red",
    "3802": "res:/dx9/model/celestial/sun/sun_yellow_small_01a.red",
    "3803": "res:/dx9/model/celestial/sun/sun_white_tiny_01a.red",
    "3807": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",

    "4011": "res:/dx9/model/Structure/Drone/Defense/Wall_Bunker.red",
    "4044": "res:/dx9/model/turret/drone/sleeperdrone_small_t1.red",
    "4045": "res:/dx9/model/turret/drone/sleeperdrone_medium_t1.red",
    "4049": "res:/dx9/model/turret/drone/sleeperdrone_large_t1.red",
    "4147": "res:/dx9/model/turret/energy/pulse/l/pulse_heavydual_t1.red",
    "4248": "res:/fisfx/generic/warp_disruption/warp_disruption_field_generator.red",
    "4250": "res:/dx9/model/turret/tractor/s/tractor_s_t1.red",
    "4252": "res:/dx9/model/turret/tractor/xl/tractor_xl_t1.red",
    "4256": "res:/dx9/model/turret/launcher/bomb/bomb_t1.red",
    "4257": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "5175": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Gatling_T1.red",
    "5177": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Gatling_T1.red",
    "5179": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Gatling_T1.red",
    "5181": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Gatling_T1.red",
    "5215": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Dual_T1.red",
    "5217": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Dual_T1.red",
    "5219": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Dual_T1.red",
    "5221": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Dual_T1.red",
    "5231": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "5233": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "5235": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "5237": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "5239": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "5241": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "5243": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "5245": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "6631": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Dual_T1.red",
    "6633": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Dual_T1.red",
    "6635": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Dual_T1.red",
    "6637": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Dual_T1.red",
    "6671": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Medium_T1.red",
    "6673": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Medium_T1.red",
    "6675": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Medium_T1.red",
    "6677": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Medium_T1.red",
    "6715": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Medium_T1.red",
    "6717": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Medium_T1.red",
    "6719": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Medium_T1.red",
    "6721": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Medium_T1.red",
    "6757": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Quad_T1.red",
    "6759": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Quad_T1.red",
    "6761": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Quad_T1.red",
    "6763": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Quad_T1.red",
    "6805": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Focused_T1.red",
    "6807": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Focused_T1.red",
    "6809": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Focused_T1.red",
    "6811": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Focused_T1.red",
    "6859": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Focused_T1.red",
    "6861": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Focused_T1.red",
    "6863": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Focused_T1.red",
    "6865": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Focused_T1.red",
    "6919": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Heavy_T1.red",
    "6921": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Heavy_T1.red",
    "6923": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Heavy_T1.red",
    "6925": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Heavy_T1.red",
    "6959": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Heavy_T1.red",
    "6961": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Heavy_T1.red",
    "6963": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Heavy_T1.red",
    "6965": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Heavy_T1.red",
    "6999": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "7001": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "7003": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "7005": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "7043": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "7045": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "7047": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "7049": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "7083": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "7085": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "7087": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "7089": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "7123": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "7125": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "7127": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "7131": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "7167": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "7169": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "7171": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "7173": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "7247": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_75mm_T1.red",
    "7249": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_75mm_T1.red",
    "7251": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_75mm_T1.red",
    "7253": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_75mm_T1.red",
    "7287": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_150mm_T1.red",
    "7289": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_150mm_T1.red",
    "7291": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_150mm_T1.red",
    "7293": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_150mm_T1.red",
    "7327": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_150mmDual_T1.red",
    "7329": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_150mmDual_T1.red",
    "7331": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_150mmDual_T1.red",
    "7333": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_150mmDual_T1.red",
    "7367": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_250mm_T1.red",
    "7369": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_250mm_T1.red",
    "7371": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_250mm_T1.red",
    "7373": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_250mm_T1.red",
    "7407": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "7409": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "7411": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "7413": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "7447": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "7449": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "7451": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "7453": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "7487": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Electron_T1.red",
    "7489": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Electron_T1.red",
    "7491": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Electron_T1.red",
    "7493": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Electron_T1.red",
    "7535": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Ion_T1.red",
    "7537": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Ion_T1.red",
    "7539": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Ion_T1.red",
    "7541": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Ion_T1.red",
    "7579": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Neutron_T1.red",
    "7581": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Neutron_T1.red",
    "7583": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Neutron_T1.red",
    "7585": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Neutron_T1.red",
    "7619": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Electron_T1.red",
    "7621": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Electron_T1.red",
    "7623": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Electron_T1.red",
    "7625": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Electron_T1.red",
    "7663": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Ion_T1.red",
    "7665": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Ion_T1.red",
    "7667": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Ion_T1.red",
    "7669": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Ion_T1.red",
    "7703": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Neutron_T1.red",
    "7705": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Neutron_T1.red",
    "7707": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Neutron_T1.red",
    "7709": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Neutron_T1.red",
    "7743": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Electron_T1.red",
    "7745": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Electron_T1.red",
    "7747": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Electron_T1.red",
    "7749": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Electron_T1.red",
    "7783": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Neutron_T1.red",
    "7785": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Neutron_T1.red",
    "7787": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Neutron_T1.red",
    "7789": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Neutron_T1.red",
    "7827": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Ion_T1.red",
    "7829": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Ion_T1.red",
    "7831": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Ion_T1.red",
    "7833": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Ion_T1.red",
    "7993": "res:/dx9/model/turret/launcher/light/light_t1.red",
    "7997": "res:/dx9/model/turret/launcher/heavy/heavy_t1.red",
    "8001": "res:/dx9/model/turret/launcher/torpedo/torpedo_t1.red",
    "8007": "res:/dx9/model/turret/launcher/rapidlight/rapidlight_t1.red",
    "8023": "res:/dx9/model/Turret/Launcher/RapidLight/RapidLight_T1.red",
    "8025": "res:/dx9/model/Turret/Launcher/RapidLight/RapidLight_T1.red",
    "8027": "res:/dx9/model/Turret/Launcher/RapidLight/RapidLight_T1.red",
    "8089": "res:/dx9/model/Turret/Launcher/Light/Light_T1.red",
    "8091": "res:/dx9/model/Turret/Launcher/Light/Light_T1.red",
    "8093": "res:/dx9/model/Turret/Launcher/Light/Light_T1.red",
    "8101": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "8103": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "8105": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "8113": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "8115": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "8117": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "8759": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_125mm_T1.red",
    "8785": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_125mm_T1.red",
    "8787": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_125mm_T1.red",
    "8789": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_125mm_T1.red",
    "8815": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_150mm_T1.red",
    "8817": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_150mm_T1.red",
    "8819": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_150mm_T1.red",
    "8821": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_150mm_T1.red",
    "8863": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_200mm_T1.red",
    "8865": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_200mm_T1.red",
    "8867": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_200mm_T1.red",
    "8869": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_200mm_T1.red",
    "8903": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_250mm_T1.red",
    "8905": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_250mm_T1.red",
    "8907": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_250mm_T1.red",
    "8909": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_250mm_T1.red",
    "9071": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_180mmDual_T1.red",
    "9073": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_180mmDual_T1.red",
    "9091": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_180mmDual_T1.red",
    "9093": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_180mmDual_T1.red",
    "9127": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_220mm_T1.red",
    "9129": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_220mm_T1.red",
    "9131": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_220mm_T1.red",
    "9133": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_220mm_T1.red",
    "9135": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_425mm_T1.red",
    "9137": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_425mm_T1.red",
    "9139": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_425mm_T1.red",
    "9141": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_425mm_T1.red",
    "9207": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_650mm_T1.red",
    "9209": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_650mm_T1.red",
    "9211": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_650mm_T1.red",
    "9213": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_650mm_T1.red",
    "9247": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_425mmDual_T1.red",
    "9249": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_425mmDual_T1.red",
    "9251": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_425mmDual_T1.red",
    "9253": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_425mmDual_T1.red",
    "9287": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_650mmDual_T1.red",
    "9289": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_650mmDual_T1.red",
    "9291": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_650mmDual_T1.red",
    "9293": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_650mmDual_T1.red",
    "9327": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_800mmHeavy_T1.red",
    "9329": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_800mmHeavy_T1.red",
    "9331": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_800mmHeavy_T1.red",
    "9333": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_800mmHeavy_T1.red",
    "9367": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "9369": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "9371": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "9373": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "9377": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "9411": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_280mm_T1.red",
    "9413": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_280mm_T1.red",
    "9415": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_280mm_T1.red",
    "9417": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_280mm_T1.red",
    "9419": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_720mm_T1.red",
    "9421": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_720mm_T1.red",
    "9451": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_720mm_T1.red",
    "9453": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_720mm_T1.red",
    "9455": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_720mm_T1.red",
    "9457": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_720mm_T1.red",
    "9491": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1400mm_T1.red",
    "9493": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1400mm_T1.red",
    "9495": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1400mm_T1.red",
    "9497": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1400mm_T1.red",

    "10065": "res:/dx9/model/worldobject/cloud/darkcloud.red",
    "10066": "res:/dx9/model/worldobject/cloud/darkgreencloud.red",
    "10067": "res:/dx9/model/worldobject/cloud/dustcloud.red",
    "10068": "res:/dx9/model/worldobject/cloud/ionloud.red",
    "10069": "res:/dx9/model/worldobject/cloud/sparkcloud.red",
    "10119": "res:/dx9/model/Structure/Universal/DamageDampener/UDD1.red",
    "10128": "res:/dx9/model/worldobject/cloud/darkgreyclouds.red",
    "10129": "res:/dx9/model/worldobject/cloud/darkgreyturbulentclouds.red",
    "10130": "res:/fisfx/cloud/electricclouds_v2.red",
    "10131": "res:/dx9/model/worldobject/cloud/fireclouds.red",
    "10132": "res:/dx9/model/worldobject/cloud/plasmaclouds.red",
    "10137": "res:/dx9/model/Celestial/rock/Coral/Coral1.red",
    "10167": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "10231": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "10232": "res:/dx9/model/worldobject/cloud/debrisstorm.red",
    "10233": "res:/dx9/model/worldobject/cloud/meteorstorm.red",
    "10261": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "10262": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "10263": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "10267": "res:/dx9/model/Celestial/rock/Coral/Coral2.red",
    "10629": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "10631": "res:/dx9/model/turret/launcher/rocket/rocket_t1.red",
    "10678": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_125mm_T1.red",
    "10680": "res:/dx9/model/turret/hybrid/rail/s/rail_125mm_t1.red",
    "10688": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_125mm_T1.red",
    "10690": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_125mm_T1.red",
    "10692": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_125mm_T1.red",
    "10694": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_125mm_T1.red",
    "10753": "res:/dx9/model/worldobject/cloud/gaseoussuperlight.red",
    "10754": "res:/dx9/model/worldobject/cloud/longwispyorange.red",
    "10755": "res:/dx9/model/worldobject/cloud/yellowoval.red",
    "10756": "res:/dx9/model/worldobject/cloud/dustybrownstreak.red",
    "10757": "res:/dx9/model/worldobject/cloud/gaseouscloudlight.red",
    "10758": "res:/dx9/model/worldobject/cloud/wispyaqua.red",
    "10759": "res:/dx9/model/worldobject/cloud/cloud1.red",
    "10760": "res:/dx9/model/worldobject/cloud/greendusty.red",
    "10761": "res:/dx9/model/worldobject/cloud/cloud2.red",
    "10762": "res:/dx9/model/worldobject/cloud/aquapuff.red",
    "10763": "res:/dx9/model/worldobject/cloud/gaseuoscloud.red",
    "10764": "res:/dx9/model/worldobject/cloud/dustybrown.red",
    "10765": "res:/dx9/model/worldobject/cloud/greengasescloud.red",
    "10782": "res:/dx9/model/Celestial/rock/RockFormation/RF1/RF1.red",
    "10783": "res:/dx9/model/Celestial/rock/RockFormation/RF2/RF2.red",
    "10784": "res:/dx9/model/Celestial/rock/RockFormation/RF3/RF3.red",
    "10785": "res:/dx9/model/Celestial/rock/RockFormation/RF4/RF4.red",
    "10786": "res:/dx9/model/Celestial/rock/RockFormation/RF5/RF5.red",
    "10787": "res:/dx9/model/Celestial/rock/RockFormation/RF6/RF6.red",
    "11072": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "11369": "res:/dx9/model/Artifact/Universal/ConcentricRings/UCR.red",
    "11488": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCH.red",
    "11489": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCG.red",
    "11490": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCC.red",
    "11546": "res:/dx9/model/worldobject/cloud/greendusty.red",
    "11589": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "11590": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "12028": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "12108": "res:/dx9/model/Turret/Mining/Deep_Core/DeepCore_T1.red",
    "12344": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_200mm_T1.red",
    "12346": "res:/dx9/model/turret/hybrid/rail/m/rail_200mm_t1.red",
    "12354": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "12356": "res:/dx9/model/turret/hybrid/rail/l/rail_350mm_t1.red",
    "12554": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "12555": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "12556": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "12602": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "12847": "res:/dx9/model/Deployables/Generic/Container/Cryo/CryoCS.red",
    "12850": "res:/dx9/model/Deployables/Generic/Container/Cryo/CryoCS.red",
    "12851": "res:/dx9/model/Deployables/Generic/Container/Cryo/CryoCS.red",
    "12852": "res:/dx9/model/Deployables/Generic/Container/Cryo/CryoCS.red",
    "12853": "res:/dx9/model/Deployables/Generic/Container/Cryo/CryoCS.red",
    "12854": "res:/dx9/model/Deployables/Generic/Container/Cryo/CryoCS.red",
    "12856": "res:/dx9/model/Deployables/Generic/Container/Cryo/CryoCS.red",
    "13119": "res:/dx9/model/turret/launcher/rocket/rocket_impact_mjolnir.red",
    "13200": "res:/dx9/model/worldobject/forcefield/forcefield_eml.red",
    "13320": "res:/dx9/model/Turret/Launcher/Cruise/Cruise_T1.red",
    "13321": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "13773": "res:/dx9/model/turret/projectile/auto/s/auto_125mm_t1.red",
    "13774": "res:/dx9/model/turret/projectile/artil/l/artil_1200mmheavy_t1.red",
    "13775": "res:/dx9/model/turret/projectile/artil/l/artil_1400mm_t1.red",
    "13776": "res:/dx9/model/turret/projectile/auto/s/auto_150mm_t1.red",
    "13777": "res:/dx9/model/turret/projectile/auto/s/auto_200mm_t1.red",
    "13778": "res:/dx9/model/turret/projectile/auto/m/auto_220mm_t1.red",
    "13779": "res:/dx9/model/turret/projectile/artil/s/artil_250mm_t1.red",
    "13781": "res:/dx9/model/turret/projectile/artil/s/artil_280mm_t1.red",
    "13782": "res:/dx9/model/turret/projectile/auto/m/auto_425mm_t1.red",
    "13783": "res:/dx9/model/turret/projectile/artil/m/artil_650mm_t1.red",
    "13784": "res:/dx9/model/turret/projectile/artil/m/artil_720mm_t1.red",
    "13785": "res:/dx9/model/turret/projectile/auto/l/auto_800mmheavy_t1.red",
    "13786": "res:/dx9/model/turret/projectile/auto/m/auto_180mmdual_t1.red",
    "13787": "res:/dx9/model/turret/projectile/auto/l/auto_425mmdual_t1.red",
    "13788": "res:/dx9/model/turret/projectile/auto/l/auto_650mmdual_t1.red",
    "13791": "res:/dx9/model/turret/energy/pulse/l/pulse_heavydual_t1.red",
    "13793": "res:/dx9/model/turret/energy/beam/l/beam_heavydual_t1.red",
    "13795": "res:/dx9/model/turret/energy/beam/s/beam_dual_t1.red",
    "13797": "res:/dx9/model/turret/energy/pulse/s/pulse_dual_t1.red",
    "13799": "res:/dx9/model/turret/energy/beam/m/beam_focused_t1.red",
    "13801": "res:/dx9/model/turret/energy/pulse/m/pulse_focused_t1.red",
    "13803": "res:/dx9/model/turret/energy/pulse/s/pulse_gatling_t1.red",
    "13805": "res:/dx9/model/turret/energy/beam/m/beam_heavy_t1.red",
    "13807": "res:/dx9/model/turret/energy/pulse/m/pulse_heavy_t1.red",
    "13809": "res:/dx9/model/turret/energy/beam/s/beam_medium_t1.red",
    "13811": "res:/dx9/model/turret/energy/pulse/s/pulse_medium_t1.red",
    "13813": "res:/dx9/model/turret/energy/beam/l/beam_mega_t1.red",
    "13815": "res:/dx9/model/turret/energy/pulse/l/pulse_mega_t1.red",
    "13817": "res:/dx9/model/turret/energy/beam/l/beam_tachyon_t1.red",
    "13819": "res:/dx9/model/turret/energy/beam/m/beam_quad_t1.red",
    "13820": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "13821": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "13822": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Dual_T1.red",
    "13823": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Dual_T1.red",
    "13824": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Focused_T1.red",
    "13825": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Focused_T1.red",
    "13826": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Gatling_T1.red",
    "13827": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Heavy_T1.red",
    "13828": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Heavy_T1.red",
    "13829": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Medium_T1.red",
    "13830": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Medium_T1.red",
    "13831": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "13832": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "13833": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Quad_T1.red",
    "13834": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "13856": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_nova.red",
    "13864": "res:/dx9/model/turret/hybrid/rail/s/rail_125mm_t1.red",
    "13865": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_125mm_T1.red",
    "13866": "res:/dx9/model/turret/hybrid/rail/s/rail_150mm_t1.red",
    "13867": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_150mm_T1.red",
    "13868": "res:/dx9/model/turret/hybrid/rail/m/rail_200mm_t1.red",
    "13870": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_200mm_T1.red",
    "13872": "res:/dx9/model/turret/hybrid/rail/m/rail_250mm_t1.red",
    "13873": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_250mm_T1.red",
    "13874": "res:/dx9/model/turret/hybrid/rail/l/rail_350mm_t1.red",
    "13876": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "13878": "res:/dx9/model/turret/hybrid/rail/l/rail_425mm_t1.red",
    "13879": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "13880": "res:/dx9/model/turret/hybrid/rail/m/rail_150mmdual_t1.red",
    "13881": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_150mmDual_T1.red",
    "13882": "res:/dx9/model/turret/hybrid/rail/l/rail_250mmdual_t1.red",
    "13883": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "13884": "res:/dx9/model/turret/hybrid/blast/m/blast_electron_t1.red",
    "13885": "res:/dx9/model/turret/hybrid/blast/m/blast_ion_t1.red",
    "13886": "res:/dx9/model/turret/hybrid/blast/s/blast_electron_t1.red",
    "13887": "res:/dx9/model/turret/hybrid/blast/s/blast_ion_t1.red",
    "13888": "res:/dx9/model/turret/hybrid/blast/s/blast_neutron_t1.red",
    "13889": "res:/dx9/model/turret/hybrid/blast/l/blast_electron_t1.red",
    "13890": "res:/dx9/model/turret/hybrid/blast/l/blast_ion_t1.red",
    "13891": "res:/dx9/model/turret/hybrid/blast/l/blast_neutron_t1.red",
    "13892": "res:/dx9/model/turret/hybrid/blast/m/blast_neutron_t1.red",
    "13893": "res:/dx9/model/turret/hybrid/rail/s/rail_75mm_t1.red",
    "13894": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_75mm_T1.red",
    "13919": "res:/dx9/model/turret/launcher/rapidlight/rapidlight_t1.red",
    "13920": "res:/dx9/model/Turret/Launcher/RapidLight/RapidLight_T1.red",
    "13921": "res:/dx9/model/turret/launcher/heavy/heavy_t1.red",
    "13922": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "13923": "res:/dx9/model/turret/launcher/torpedo/torpedo_t1.red",
    "13924": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "13925": "res:/dx9/model/turret/launcher/light/light_t1.red",
    "13926": "res:/dx9/model/Turret/Launcher/Light/Light_T1.red",
    "13927": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "13929": "res:/dx9/model/Turret/Launcher/Cruise/Cruise_T1.red",
    "13931": "res:/dx9/model/turret/launcher/rocket/rocket_t1.red",
    "13933": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",

    "14272": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_200mm_T1.red",
    "14274": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_200mm_T1.red",
    "14276": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_200mm_T1.red",
    "14278": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_200mm_T1.red",
    "14280": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "14282": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "14284": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "14286": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "14375": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Electron_T1.red",
    "14377": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Electron_T1.red",
    "14379": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Ion_T1.red",
    "14381": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Ion_T1.red",
    "14383": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Neutron_T1.red",
    "14385": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Neutron_T1.red",
    "14387": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "14389": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "14391": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "14393": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "14395": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "14397": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "14399": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "14401": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "14403": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "14405": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "14407": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "14409": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "14411": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "14413": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "14415": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "14417": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "14419": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "14421": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "14423": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "14425": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "14427": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "14429": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "14431": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "14433": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "14435": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "14437": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "14439": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "14441": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "14443": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "14445": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "14447": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "14449": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "14451": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "14453": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "14455": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "14457": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_800mmHeavy_T1.red",
    "14459": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_800mmHeavy_T1.red",
    "14461": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "14463": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "14465": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1400mm_T1.red",
    "14467": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1400mm_T1.red",
    "14469": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_425mmDual_T1.red",
    "14471": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_425mmDual_T1.red",
    "14473": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_650mmDual_T1.red",
    "14475": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_650mmDual_T1.red",
    "14516": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "14518": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "14520": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "14522": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "14524": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "14525": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "14526": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "14527": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "14672": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "14674": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "14676": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "14678": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "14680": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "14681": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "14682": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "14683": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "15399": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "15401": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "15403": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "15421": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "15423": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "15427": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "15429": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "15443": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "15445": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1400mm_T1.red",
    "15449": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "15814": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "15815": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_150mmDual_T1.red",
    "15816": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_75mm_T1.red",
    "15817": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "15818": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "15820": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_250mm_T1.red",
    "15821": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_200mm_T1.red",
    "15823": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_150mm_T1.red",
    "15824": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_125mm_T1.red",
    "15825": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Neutron_T1.red",
    "15826": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Neutron_T1.red",
    "15827": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Ion_T1.red",
    "15828": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Electron_T1.red",
    "15829": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Ion_T1.red",
    "15830": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Neutron_T1.red",
    "15831": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Ion_T1.red",
    "15832": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Electron_T1.red",
    "15833": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Electron_T1.red",
    "15834": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_250mmDual_T1.red",
    "15835": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_150mmDual_T1.red",
    "15836": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_75mm_T1.red",
    "15837": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "15838": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "15840": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_250mm_T1.red",
    "15841": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_200mm_T1.red",
    "15843": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_150mm_T1.red",
    "15844": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_125mm_T1.red",
    "15845": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "15846": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Quad_T1.red",
    "15847": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "15848": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "15849": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Medium_T1.red",
    "15850": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Medium_T1.red",
    "15851": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Heavy_T1.red",
    "15852": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Heavy_T1.red",
    "15853": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Gatling_T1.red",
    "15854": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Focused_T1.red",
    "15855": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Focused_T1.red",
    "15856": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Dual_T1.red",
    "15857": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Dual_T1.red",
    "15858": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "15859": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "15860": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "15861": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Quad_T1.red",
    "15862": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "15863": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "15864": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Medium_T1.red",
    "15865": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Medium_T1.red",
    "15866": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Heavy_T1.red",
    "15867": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Heavy_T1.red",
    "15868": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Gatling_T1.red",
    "15869": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Focused_T1.red",
    "15870": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Focused_T1.red",
    "15871": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Dual_T1.red",
    "15872": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Dual_T1.red",
    "15873": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "15874": "res:/dx9/model/Turret/Energy/Beam/L/Beam_HeavyDual_T1.red",
    "16025": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "16041": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCC.red",
    "16042": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCC.red",
    "16043": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCG.red",
    "16044": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCG.red",
    "16045": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCL.red",
    "16046": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_125mm_T1.red",
    "16047": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "16048": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1400mm_T1.red",
    "16049": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_150mm_T1.red",
    "16050": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_200mm_T1.red",
    "16051": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_220mm_T1.red",
    "16052": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_250mm_T1.red",
    "16053": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_280mm_T1.red",
    "16054": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_425mm_T1.red",
    "16055": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_650mm_T1.red",
    "16056": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_720mm_T1.red",
    "16057": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_800mmHeavy_T1.red",
    "16058": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_180mmDual_T1.red",
    "16059": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_425mmDual_T1.red",
    "16060": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_650mmDual_T1.red",
    "16061": "res:/dx9/model/Turret/Launcher/RapidLight/RapidLight_T1.red",
    "16062": "res:/dx9/model/Turret/Launcher/Cruise/Cruise_T1.red",
    "16063": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "16064": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "16065": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "16067": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "16068": "res:/dx9/model/Turret/Launcher/Light/Light_T1.red",
    "16103": "res:/fisfx/generic/forcefield/posforcefield.red",
    "16128": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Medium_T1.red",
    "16129": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_HeavyDual_T1.red",
    "16131": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Heavy_T1.red",
    "16132": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_150mm_T1.red",
    "16133": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_250mm_T1.red",
    "16134": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_350mm_T1.red",
    "16136": "res:/dx9/model/Turret/Launcher/Light/Light_T1.red",
    "16137": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "16138": "res:/dx9/model/Turret/Launcher/Cruise/Cruise_T1.red",
    "16148": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "16149": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_650mm_T1.red",
    "16150": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_200mm_T1.red",
    "16223": "res:/dx9/model/Deployables/Starbase/ShieldGenerator/ShieldGenerator_T1.red",
    "16278": "res:/dx9/model/Turret/Mining/Ice/Ice_T1.red",
    "16513": "res:/dx9/model/Turret/Launcher/Cruise/Cruise_T1.red",
    "16515": "res:/dx9/model/Turret/Launcher/Cruise/Cruise_T1.red",
    "16517": "res:/dx9/model/Turret/Launcher/Cruise/Cruise_T1.red",
    "16519": "res:/dx9/model/Turret/Launcher/Cruise/Cruise_T1.red",
    "16521": "res:/dx9/model/turret/launcher/rocket/rocket_t1.red",
    "16523": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "16525": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "16527": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "16720": "res:/dx9/model/Deployables/Generic/Container/Ammo/AmmoCS.red",
    "16721": "res:/dx9/model/Deployables/Generic/Container/Armor/ArmorCS.red",
    "16722": "res:/dx9/model/Deployables/Generic/Container/Electronic/ElectronicCS.red",
    "16723": "res:/dx9/model/Deployables/Generic/Container/Mineral/MineralCS.red",
    "16724": "res:/dx9/model/Deployables/Generic/Container/Rogue/RogueCS.red",
    "16725": "res:/dx9/model/Deployables/Generic/Container/Weapon/WeaponCS.red",
    "16734": "res:/dx9/model/Structure/Universal/Construction/GB2/UC_GB2.red",
    "16849": "res:/dx9/model/Artifact/Universal/ConcentricRings/UCR.red",
    "17184": "res:/dx9/model/Deployables/Starbase/ShieldGenerator/ShieldGenerator_T1.red",
    "17185": "res:/dx9/model/Deployables/Starbase/ShieldGenerator/ShieldGenerator_T1.red",
    "17186": "res:/dx9/model/Deployables/Starbase/ShieldGenerator/ShieldGenerator_T1.red",
    "17187": "res:/dx9/model/Deployables/Starbase/ShieldGenerator/ShieldGenerator_T1.red",
    "17274": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "17276": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "17277": "res:/dx9/model/Deployables/Generic/Container/Electronic/ElectronicCS.red",
    "17363": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCS.red",
    "17364": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCM.red",
    "17365": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCL.red",
    "17366": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCH.red",
    "17367": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCG.red",
    "17368": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCC.red",
    "17482": "res:/dx9/model/Turret/Mining/Strip/Strip_T1.red",
    "17484": "res:/dx9/model/Turret/Launcher/RapidLight/RapidLight_T1.red",
    "17485": "res:/dx9/model/Turret/Launcher/Cruise/Cruise_T1.red",
    "17486": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "17487": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "17488": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "17490": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "17491": "res:/dx9/model/Turret/Launcher/Light/Light_T1.red",
    "17748": "res:/dx9/model/Structure/Universal/Construction/GB2/UC_GB2.red",
    "17763": "res:/dx9/model/worldobject/forcefield/forcefield_eml.red",
    "17774": "res:/dx9/model/worldobject/enviroments/env_icefield.red",
    "17786": "res:/dx9/model/Deployables/Starbase/ShieldGenerator/ShieldGenerator_T1.red",
    "17798": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "17857": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_mjolnir.red",
    "17858": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_mjolnir.red",
    "17859": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_scourge.red",
    "17861": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_inferno.red",
    "17863": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_nova.red",
    "17908": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "17912": "res:/dx9/model/turret/mining/strip/strip_t1.red",
    "17979": "res:/dx9/model/Deployables/Starbase/ShieldGenerator/ShieldGenerator_T1.red",
    "17996": "res:/dx9/model/Deployables/Generic/Container/Electronic/ElectronicCS.red",
    "18012": "res:/dx9/model/Structure/Drone/Defense/Wall_Bunker.red",
    "18013": "res:/dx9/model/Structure/Drone/Defense/Wall_Elevator.red",
    "18014": "res:/dx9/model/Structure/Drone/Defense/Wall_Junction.red",
    "18015": "res:/dx9/model/Structure/Drone/Defense/Wall_Lookout.red",
    "18016": "res:/dx9/model/Structure/Drone/Defense/Wall_MissileBattery.red",
    "18017": "res:/dx9/model/Structure/Drone/Defense/Wall_StraightX.red",
    "18018": "res:/dx9/model/Structure/Drone/Defense/Wall_StraightXZ.red",
    "18019": "res:/dx9/model/Structure/Drone/Defense/Wall_StraightZ.red",
    "18020": "res:/dx9/model/Structure/Drone/Defense/Wall_StraightZX.red",
    "18021": "res:/dx9/model/Structure/Drone/Defense/Wall_Sentry.red",
    "18023": "res:/dx9/model/Structure/Drone/Defense/Wall_Sentry.red",
    "18028": "res:/dx9/model/Structure/Drone/Defense/Wall_Sentry.red",
    "18031": "res:/dx9/model/Structure/Drone/Defense/Wall_MissileBattery.red",
    "18032": "res:/dx9/model/Structure/Drone/Defense/Wall_MissileBattery.red",
    "18033": "res:/dx9/model/Structure/Drone/Defense/Wall_MissileBattery.red",
    "18035": "res:/dx9/model/Structure/Drone/Defense/Wall_MissileBattery.red",
    "18068": "res:/dx9/model/turret/mining/deep_core/deepcore_t1.red",
    "18582": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "18626": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "18628": "res:/dx9/model/Deployables/Generic/Container/Rogue/RogueCS.red",
    "18629": "res:/dx9/model/Deployables/Generic/Container/Rogue/RogueCS.red",
    "18635": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "18637": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "19373": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "19425": "res:/dx9/model/Deployables/Generic/Container/Rogue/RogueCS.red",
    "19426": "res:/dx9/model/Structure/Drone/Defense/Wall_Bunker.red",
    "19591": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "19658": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "19660": "res:/dx9/model/turret/launcher/festival/festival_t1.red",
    "19703": "res:/dx9/model/Structure/Drone/Defense/Wall_MissileBattery.red",
    "19713": "res:/dx9/model/worldobject/enviroments/env_argon.red",
    "19728": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "19739": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "19740": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "19743": "res:/dx9/model/worldobject/forcefield/forcefield_eml.red",
    "19746": "res:/dx9/model/worldobject/enviroments/env_dust_devil.red",
    "19747": "res:/dx9/model/worldobject/enviroments/env_electricstorm.red",
    "19748": "res:/dx9/model/worldobject/enviroments/env_ghostworld.red",
    "19749": "res:/dx9/model/worldobject/enviroments/env_heaven.red",
    "19750": "res:/dx9/model/worldobject/enviroments/env_hell.red",
    "19751": "res:/dx9/model/worldobject/enviroments/env_icefield.red",
    "19752": "res:/dx9/model/worldobject/enviroments/env_krypton.red",
    "19753": "res:/dx9/model/worldobject/enviroments/env_poison.red",
    "19754": "res:/dx9/model/worldobject/enviroments/env_radon.red",
    "19755": "res:/dx9/model/worldobject/enviroments/env_storm.red",
    "19756": "res:/dx9/model/worldobject/enviroments/env_xenon.red",
    "19794": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "19960": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "20077": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "20138": "res:/dx9/model/Turret/Launcher/HeavyAssault/HeavyAssault_T1.red",
    "20182": "res:/dx9/model/Structure/Drone/Defense/Wall_Bunker.red",
    "20306": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_mjolnir.red",
    "20307": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_scourge.red",
    "20308": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_inferno.red",
    "20444": "res:/dx9/model/turret/energy/pulse/xl/pulse_gigadual_t1.red",
    "20446": "res:/dx9/model/turret/energy/beam/xl/beam_gigadual_t1.red",
    "20448": "res:/dx9/model/turret/hybrid/rail/xl/rail_1000mmdual_t1.red",
    "20450": "res:/dx9/model/turret/hybrid/blast/xl/blast_ion_t1.red",
    "20452": "res:/dx9/model/turret/projectile/auto/xl/auto_6x2500mm_t1.red",
    "20454": "res:/dx9/model/turret/projectile/artil/xl/artil_3500mmquad_t1.red",
    "20539": "res:/dx9/model/Turret/Launcher/CitadelTorpedo/CitadelTorpedo_T1.red",
    "20540": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "20541": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCL.red",
    "20587": "res:/dx9/model/Turret/Hybrid/Rail/S/Rail_150mm_T1.red",
    "20589": "res:/dx9/model/Turret/Hybrid/Rail/M/Rail_250mm_T1.red",
    "20591": "res:/dx9/model/Turret/Hybrid/Rail/L/Rail_425mm_T1.red",
    "20593": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "20595": "res:/dx9/model/Turret/Launcher/Light/Light_T1.red",
    "20597": "res:/dx9/model/Turret/Launcher/RapidLight/RapidLight_T1.red",
    "20599": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "20600": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "20601": "res:/dx9/model/turret/launcher/cruise/cruise_t1.red",
    "20602": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "20603": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "20604": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "21052": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "21094": "res:/FisFX/Deployable/CynoBeacon_Rt_T1a.red",
    "21099": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "21542": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "21545": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_200mm_T1.red",
    "21547": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_250mm_T1.red",
    "21549": "res:/dx9/model/Turret/Projectile/Artil/S/Artil_280mm_T1.red",
    "21551": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_425mm_T1.red",
    "21553": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_650mm_T1.red",
    "21555": "res:/dx9/model/Turret/Projectile/Artil/M/Artil_720mm_T1.red",
    "21557": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_800mmHeavy_T1.red",
    "21559": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1200mmHeavy_T1.red",
    "21561": "res:/dx9/model/Turret/Projectile/Artil/L/Artil_1400mm_T1.red",
    "21841": "res:/dx9/model/turret/mining/standard/standard_t1.red",
    "21867": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_nova.red",
    "21891": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "22168": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "22229": "res:/dx9/model/turret/mining/ice/ice_t1.red",
    "22318": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "22564": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "22565": "res:/dx9/model/Turret/Launcher/Light/Light_T1.red",
    "22566": "res:/dx9/model/Turret/Launcher/RapidLight/RapidLight_T1.red",
    "22567": "res:/dx9/model/Turret/Launcher/Heavy/Heavy_T1.red",
    "22568": "res:/dx9/model/Turret/Launcher/Cruise/Cruise_T1.red",
    "22569": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "22711": "res:/dx9/model/worldobject/forcefield/forcefield_eml.red",
    "22778": "res:/fisfx/generic/warp_disruption/warp_disruption_probe.red",
    "22849": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "22899": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Electron_T1.red",
    "22901": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Ion_T1.red",
    "22903": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Neutron_T1.red",
    "22905": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Electron_T1.red",
    "22907": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Ion_T1.red",
    "22909": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Neutron_T1.red",
    "22911": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Electron_T1.red",
    "22913": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Ion_T1.red",
    "22915": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Neutron_T1.red",
    "22921": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "22923": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "23220": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "23544": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCC.red",
    "23595": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "23666": "res:/dx9/model/Structure/Drone/Defense/Wall_Sentry.red",
    "23688": "res:/dx9/model/Structure/Drone/Defense/Wall_Bunker.red",
    "23689": "res:/dx9/model/Structure/Drone/Defense/Wall_Sentry.red",
    "23691": "res:/dx9/model/Structure/Drone/Defense/Wall_Bunker.red",
    "23733": "res:/dx9/model/worldobject/cloud/fireclouds.red",
    "23752": "res:/dx9/model/Celestial/rock/Shattered/Cloven.red",
    "23753": "res:/dx9/model/Celestial/rock/Shattered/Cloven_Red.red",
    "23754": "res:/dx9/model/Celestial/Rock/Shattered/Crystal_Blue.red",
    "23755": "res:/dx9/model/Celestial/Rock/Shattered/Crystal_Clear.red",
    "23756": "res:/dx9/model/Celestial/Rock/Shattered/Crystal_Orange.red",
    "23767": "res:/dx9/model/worldobject/cloud/fireclouds.red",
    "23828": "res:/FisFX/Celestial/SpatialRift_Rt_01a.red",
    "23834": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Dual_T1.red",
    "23836": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Medium_T1.red",
    "23838": "res:/dx9/model/Turret/Energy/Beam/S/Beam_Medium_T1.red",
    "23840": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Focused_T1.red",
    "23842": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Heavy_T1.red",
    "23844": "res:/dx9/model/Turret/Energy/Beam/M/Beam_Heavy_T1.red",
    "23846": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "23848": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Mega_T1.red",
    "23850": "res:/dx9/model/Turret/Energy/Beam/L/Beam_Tachyon_T1.red",
    "23904": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "24124": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "24279": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "24280": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "24281": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "24305": "res:/dx9/model/Turret/Mining/Modulated_Strip/Modulated_Strip_T1.red",
    "24348": "res:/dx9/model/turret/tractor/s/tractor_s_t1.red",
    "24445": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCL.red",
    "24447": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "24471": "res:/dx9/model/turret/launcher/rocket/rocket_impact_scourge.red",
    "24473": "res:/dx9/model/turret/launcher/rocket/rocket_impact_nova.red",
    "24475": "res:/dx9/model/turret/launcher/rocket/rocket_impact_inferno.red",
    "24477": "res:/dx9/model/turret/launcher/rocket/rocket_impact_scourge.red",
    "24478": "res:/dx9/model/turret/launcher/rocket/rocket_impact_nova.red",
    "24479": "res:/dx9/model/turret/launcher/rocket/rocket_impact_inferno.red",
    "24480": "res:/FisFX/Celestial/SpatialRift_Rt_01a.red",
    "24486": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_inferno.red",
    "24488": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_nova.red",
    "24490": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_mjolnir.red",
    "24492": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_scourge.red",
    "24493": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_mjolnir.red",
    "24494": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_inferno.red",
    "24495": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "24497": "res:/dx9/model/turret/launcher/light/light_impact_nova.red",
    "24499": "res:/dx9/model/turret/launcher/light/light_impact_inferno.red",
    "24501": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "24503": "res:/dx9/model/turret/launcher/light/light_impact_nova.red",
    "24505": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "24507": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "24509": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "24511": "res:/dx9/model/turret/launcher/heavy/heavy_impact_inferno.red",
    "24513": "res:/dx9/model/turret/launcher/heavy/heavy_impact_scourge.red",
    "24515": "res:/dx9/model/turret/launcher/heavy/heavy_impact_inferno.red",
    "24517": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "24519": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_nova.red",
    "24521": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_scourge.red",
    "24523": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_mjolnir.red",
    "24525": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_inferno.red",
    "24527": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_mjolnir.red",
    "24529": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_scourge.red",
    "24531": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "24532": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "24533": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "24534": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "24535": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "24536": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "24537": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "24539": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "24541": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "24578": "res:/dx9/model/worldobject/cloud/debrisstormflat.red",
    "24608": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "24644": "res:/dx9/model/turret/tractor/xl/tractor_xl_t1.red",
    "24713": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "24750": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "25236": "res:/dx9/model/worldobject/cloud/booster/cloud1.red",
    "25244": "res:/dx9/model/worldobject/cloud/booster/cloud2.red",
    "25245": "res:/dx9/model/worldobject/cloud/booster/cloud3.red",
    "25246": "res:/dx9/model/worldobject/cloud/booster/cloud4.red",
    "25247": "res:/dx9/model/worldobject/cloud/booster/cloud5.red",
    "25248": "res:/dx9/model/worldobject/cloud/booster/cloud6.red",
    "25249": "res:/dx9/model/worldobject/cloud/booster/cloud7.red",
    "25250": "res:/dx9/model/worldobject/cloud/booster/cloud8.red",
    "25266": "res:/dx9/model/Turret/Mining/Gas/Gas_T1.red",
    "25268": "res:/dx9/model/worldobject/cloud/booster/cloud1.red",
    "25273": "res:/dx9/model/worldobject/cloud/booster/cloud2.red",
    "25274": "res:/dx9/model/worldobject/cloud/booster/cloud6.red",
    "25275": "res:/dx9/model/worldobject/cloud/booster/cloud8.red",
    "25276": "res:/dx9/model/worldobject/cloud/booster/cloud4.red",
    "25277": "res:/dx9/model/worldobject/cloud/booster/cloud3.red",
    "25278": "res:/dx9/model/worldobject/cloud/booster/cloud5.red",
    "25279": "res:/dx9/model/worldobject/cloud/booster/cloud7.red",
    "25384": "res:/dx9/model/Deployables/Starbase/ShieldGenerator/ShieldGenerator_T1.red",
    "25393": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCG.red",
    "25471": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "25526": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "25540": "res:/dx9/model/Turret/Mining/Gas/Gas_T1.red",
    "25542": "res:/dx9/model/Turret/Mining/Gas/Gas_T1.red",
    "25707": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_t1.red",
    "25709": "res:/dx9/model/Turret/Launcher/HeavyAssault/HeavyAssault_T1.red",
    "25711": "res:/dx9/model/Turret/Launcher/HeavyAssault/HeavyAssault_T1.red",
    "25713": "res:/dx9/model/Turret/Launcher/HeavyAssault/HeavyAssault_T1.red",
    "25715": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_t1.red",
    "25716": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "25812": "res:/dx9/model/turret/mining/gas/gas_t1.red",
    "25816": "res:/dx9/model/worldobject/enviroments/env_argon.red",
    "25860": "res:/dx9/model/worldobject/enviroments/env_argon.red",
    "25861": "res:/dx9/model/Turret/Salvage/S/Salvage_S_T1.red",
    "25879": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCG.red",
    "25880": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "26135": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26145": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26148": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26149": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26150": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26151": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26152": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26153": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26154": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26164": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26166": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26272": "res:/FisFX/Celestial/SpatialRift_Rt_01a.red",
    "26467": "res:/dx9/model/worldobject/cloud/booster/cloud1.red",
    "26605": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26715": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26719": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "26727": "res:/dx9/model/Deployables/Starbase/ShieldGenerator/ShieldGenerator_T1.red",
    "26968": "res:/dx9/model/Structure/Drone/Defense/Wall_Bunker.red",
    "26983": "res:/dx9/model/Turret/Salvage/S/Salvage_S_T1.red",
    "27292": "res:/dx9/model/Artifact/Universal/ConcentricRings/UCR.red",
    "27313": "res:/dx9/model/turret/launcher/rocket/rocket_impact_inferno.red",
    "27315": "res:/dx9/model/turret/launcher/rocket/rocket_impact_inferno.red",
    "27317": "res:/dx9/model/turret/launcher/rocket/rocket_impact_inferno.red",
    "27319": "res:/dx9/model/turret/launcher/rocket/rocket_impact_mjolnir.red",
    "27321": "res:/dx9/model/turret/launcher/rocket/rocket_impact_mjolnir.red",
    "27323": "res:/dx9/model/turret/launcher/rocket/rocket_impact_mjolnir.red",
    "27325": "res:/dx9/model/turret/launcher/rocket/rocket_impact_nova.red",
    "27327": "res:/dx9/model/turret/launcher/rocket/rocket_impact_nova.red",
    "27329": "res:/dx9/model/turret/launcher/rocket/rocket_impact_nova.red",
    "27331": "res:/dx9/model/turret/launcher/rocket/rocket_impact_scourge.red",
    "27333": "res:/dx9/model/turret/launcher/rocket/rocket_impact_scourge.red",
    "27335": "res:/dx9/model/turret/launcher/rocket/rocket_impact_scourge.red",
    "27337": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_mjolnir.red",
    "27339": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_mjolnir.red",
    "27341": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_mjolnir.red",
    "27343": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_scourge.red",
    "27345": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_scourge.red",
    "27347": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_scourge.red",
    "27349": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_inferno.red",
    "27351": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_inferno.red",
    "27353": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "27355": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_inferno.red",
    "27357": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_nova.red",
    "27359": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_nova.red",
    "27361": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "27363": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_nova.red",
    "27365": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "27367": "res:/dx9/model/turret/launcher/light/light_impact_inferno.red",
    "27369": "res:/dx9/model/turret/launcher/light/light_impact_inferno.red",
    "27371": "res:/dx9/model/turret/launcher/light/light_impact_inferno.red",
    "27373": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "27375": "res:/dx9/model/turret/launcher/light/light_impact_nova.red",
    "27377": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "27379": "res:/dx9/model/turret/launcher/light/light_impact_nova.red",
    "27381": "res:/dx9/model/turret/launcher/light/light_impact_nova.red",
    "27383": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "27385": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "27387": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "27389": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "27391": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "27393": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_nova.red",
    "27395": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "27397": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_nova.red",
    "27399": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "27401": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_nova.red",
    "27403": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_inferno.red",
    "27405": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_inferno.red",
    "27407": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_inferno.red",
    "27409": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "27411": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_scourge.red",
    "27413": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_scourge.red",
    "27415": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_scourge.red",
    "27417": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_mjolnir.red",
    "27419": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_mjolnir.red",
    "27421": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_mjolnir.red",
    "27423": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "27425": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "27427": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "27429": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "27431": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "27433": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "27435": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "27437": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "27439": "res:/dx9/model/turret/launcher/heavy/heavy_impact_scourge.red",
    "27441": "res:/dx9/model/turret/launcher/heavy/heavy_impact_scourge.red",
    "27443": "res:/dx9/model/turret/launcher/heavy/heavy_impact_scourge.red",
    "27445": "res:/dx9/model/turret/launcher/heavy/heavy_impact_inferno.red",
    "27447": "res:/dx9/model/turret/launcher/heavy/heavy_impact_inferno.red",
    "27449": "res:/dx9/model/turret/launcher/heavy/heavy_impact_inferno.red",
    "27451": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "27453": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "27455": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "27457": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "27459": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "27461": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "27463": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "27465": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "27467": "res:/dx9/model/turret/launcher/cruise/cruise_impact_scourge.red",
    "27469": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "27471": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "27473": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "27475": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "27477": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "27479": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "27481": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "27483": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "27485": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "27487": "res:/dx9/model/turret/launcher/heavy/heavy_impact_scourge.red",
    "27489": "res:/dx9/model/turret/launcher/heavy/heavy_impact_scourge.red",
    "27491": "res:/dx9/model/turret/launcher/heavy/heavy_impact_scourge.red",
    "27493": "res:/dx9/model/turret/launcher/heavy/heavy_impact_scourge.red",
    "27495": "res:/dx9/model/turret/launcher/heavy/heavy_impact_inferno.red",
    "27497": "res:/dx9/model/turret/launcher/heavy/heavy_impact_inferno.red",
    "27499": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "27501": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "27503": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "27505": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "27507": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "27509": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "27511": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "27513": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "27515": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "27517": "res:/dx9/model/turret/launcher/light/light_impact_inferno.red",
    "27519": "res:/dx9/model/turret/launcher/light/light_impact_inferno.red",
    "27521": "res:/dx9/model/turret/launcher/light/light_impact_inferno.red",
    "27523": "res:/dx9/model/turret/launcher/light/light_impact_nova.red",
    "27525": "res:/dx9/model/turret/launcher/light/light_impact_nova.red",
    "27527": "res:/dx9/model/turret/launcher/light/light_impact_nova.red",
    "27803": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCC.red",
    "27883": "res:/dx9/model/turret/launcher/rocket/rocket_impact_mjolnir.red",
    "27884": "res:/dx9/model/turret/launcher/rocket/rocket_impact_mjolnir.red",
    "27885": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "27886": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "27887": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_mjolnir.red",
    "27888": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_impact_mjolnir.red",
    "27889": "res:/dx9/model/turret/launcher/heavy/heavy_impact_mjolnir.red",
    "27890": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "27891": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_mjolnir.red",
    "27892": "res:/dx9/model/turret/launcher/torpedo/torpedo_impact_mjolnir.red",
    "27893": "res:/dx9/model/turret/launcher/cruise/cruise_impact_nova.red",
    "27894": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "27915": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "27953": "res:/dx9/model/Structure/Drone/Defense/Wall_Lookout.red",
    "27954": "res:/dx9/model/Structure/Drone/Defense/Wall_Lookout.red",
    "27955": "res:/dx9/model/Structure/Drone/Defense/Wall_Lookout.red",
    "27956": "res:/dx9/model/Structure/Drone/Defense/Wall_Lookout.red",
    "28155": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "28156": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "28157": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "28158": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "28159": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "28166": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "28356": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "28369": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "28375": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_t1.red",
    "28377": "res:/dx9/model/Turret/Launcher/HeavyAssault/HeavyAssault_T1.red",
    "28379": "res:/dx9/model/Turret/Launcher/HeavyAssault/HeavyAssault_T1.red",
    "28381": "res:/dx9/model/Turret/Launcher/HeavyAssault/HeavyAssault_T1.red",
    "28383": "res:/dx9/model/Turret/Launcher/HeavyAssault/HeavyAssault_T1.red",
    "28508": "res:/dx9/model/Structure/Drone/Defense/Wall_MissileBattery.red",
    "28511": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "28513": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "28565": "res:/dx9/model/Turret/Launcher/HeavyAssault/HeavyAssault_T1.red",
    "28629": "res:/dx9/model/worldobject/cloud/booster/cloud1.red",
    "28630": "res:/dx9/model/worldobject/cloud/booster/cloud8.red",
    "28650": "res:/FisFX/Deployable/CynoBeaconBO_Rt_T1a.red",
    "28654": "res:/fisfx/generic/warp_disruption/warp_disruption_field_generator.red",
    "28694": "res:/dx9/model/worldobject/cloud/booster/cloud1.red",
    "28695": "res:/dx9/model/worldobject/cloud/booster/cloud7.red",
    "28696": "res:/dx9/model/worldobject/cloud/booster/cloud8.red",
    "28697": "res:/dx9/model/worldobject/cloud/booster/cloud2.red",
    "28698": "res:/dx9/model/worldobject/cloud/booster/cloud3.red",
    "28699": "res:/dx9/model/worldobject/cloud/booster/cloud4.red",
    "28700": "res:/dx9/model/worldobject/cloud/booster/cloud5.red",
    "28701": "res:/dx9/model/worldobject/cloud/booster/cloud6.red",
    "28748": "res:/dx9/model/Turret/Mining/Deep_Core/DeepCore_T1.red",
    "28750": "res:/dx9/model/Turret/Mining/Standard/Standard_T1.red",
    "28752": "res:/dx9/model/Turret/Mining/Ice/Ice_T1.red",
    "28754": "res:/dx9/model/Turret/Mining/Strip/Strip_T1.red",
    "28788": "res:/dx9/model/turret/mining/gas/gas_t1.red",
    "28841": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "28859": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "28860": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "28861": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "28862": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "28868": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCG.red",
    "28972": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCC.red",
    "29181": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "29182": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "29183": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "29184": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "29230": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "29250": "res:/dx9/model/worldobject/forcefield/forcefield_eml.red",
    "29251": "res:/dx9/model/Deployables/Starbase/ShieldGenerator/ShieldGenerator_T1.red",
    "29321": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCL.red",
    "29324": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "29445": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "29464": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCC.red",
    "29585": "res:/dx9/model/Structure/Drone/Defense/Wall_StraightXZ.red",
    "29586": "res:/dx9/model/Structure/Drone/Defense/Wall_StraightZX.red",
    "29589": "res:/dx9/model/Structure/Drone/Defense/Wall_MissileBattery.red",
    "29590": "res:/dx9/model/Structure/Drone/Defense/Wall_Elevator.red",
    "29591": "res:/dx9/model/Structure/Drone/Defense/Wall_StraightZ.red",
    "29592": "res:/dx9/model/Structure/Drone/Defense/Wall_Junction.red",
    "29593": "res:/dx9/model/Structure/Drone/Defense/Wall_Lookout.red",
    "29594": "res:/dx9/model/Structure/Drone/Defense/Wall_StraightX.red",
    "29616": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_nova.red",
    "29618": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_inferno.red",
    "29620": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_scourge.red",
    "29622": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_mjolnir.red",
    "29939": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "30221": "res:/dx9/model/turret/launcher/festival/festival_impact_snow.red",
    "30222": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "30370": "res:/dx9/model/worldobject/cloud/booster/cloud5.red",
    "30371": "res:/dx9/model/worldobject/cloud/booster/cloud5.red",
    "30372": "res:/dx9/model/worldobject/cloud/booster/cloud7.red",
    "30373": "res:/dx9/model/worldobject/cloud/booster/cloud4.red",
    "30374": "res:/dx9/model/worldobject/cloud/booster/cloud3.red",
    "30375": "res:/dx9/model/worldobject/cloud/booster/cloud3.red",
    "30376": "res:/dx9/model/worldobject/cloud/booster/cloud1.red",
    "30377": "res:/dx9/model/worldobject/cloud/booster/cloud2.red",
    "30378": "res:/dx9/model/worldobject/cloud/booster/cloud2.red",
    "30388": "res:/dx9/model/Structure/Gallente/Arena/Arena_GA_Mainstructures/Arena_GA_Mainstructure.red",
    "30412": "res:/dx9/model/Structure/Caldari/Arena/Arena_CA_Mainstructures/Arena_CA_Mainstructure.red",
    "30413": "res:/dx9/model/Structure/Gallente/Arena/Arena_GA_Mainstructures/Arena_GA_Smallstructures.red",
    "30414": "res:/dx9/model/Structure/Caldari/Arena/Arena_CA_Mainstructures/Arena_CA_Smallstructures.red",
    "30419": "res:/dx9/model/Structure/Caldari/Arena/Arena_CA_CenterPiece/Arena_CA_CenterPiece.red",
    "30426": "res:/dx9/model/turret/launcher/cruise/cruise_impact_mjolnir.red",
    "30428": "res:/dx9/model/turret/launcher/heavy/heavy_impact_nova.red",
    "30430": "res:/dx9/model/turret/launcher/light/light_impact_nova.red",
    "30434": "res:/dx9/model/Structure/Minmatar/Arena/Arena_MI_CenterPiece/Arena_MI_CenterPiece.red",
    "30435": "res:/dx9/model/Structure/Gallente/Arena/Arena_GA_CenterPiece/Arena_GA_CenterPiece.red",
    "30440": "res:/dx9/model/worldobject/cloud/redcloud.red",
    "30449": "res:/dx9/model/worldobject/cloud/booster/cloud1.red",
    "30451": "res:/dx9/model/Structure/Amarr/Arena/Arena_Amarr_Center/Arena_Amarr_Center.red",
    "30452": "res:/dx9/model/Structure/Amarr/Arena/Arena_Amarr_MainStructure/Arena_Amarr_MainStructure.red",
    "30453": "res:/dx9/model/Structure/Minmatar/Arena/Arena_MI_Mainstructures/Arena_MI_Mainstructure.red",
    "30455": "res:/dx9/model/Structure/Amarr/Arena/Arena_Amarr_MainStructure/Amarr_Small_Structure.red",
    "30456": "res:/dx9/model/Structure/Minmatar/Arena/Arena_MI_Mainstructures/Arena_MI_Smallstructures.red",
    "30512": "res:/dx9/model/Structure/Sleeper/Defense/SL_DroneBunker/SL_Db1/SL_Db1.red",
    "30513": "res:/dx9/model/Structure/Sleeper/Defense/SL_DroneBunker/SL_Db1/SL_Db1.red",
    "30574": "res:/dx9/Model/WorldObject/Sun/Magnetar.red",
    "30575": "res:/dx9/Model/WorldObject/Sun/BlackHole.red",
    "30576": "res:/dx9/Model/WorldObject/Sun/RedGiant.red",
    "30577": "res:/dx9/Model/WorldObject/Sun/Pulsar.red",
    "30579": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30583": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30584": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30642": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30643": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30644": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30645": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30646": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30647": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30648": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30649": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30657": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30658": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30659": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30660": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30661": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30662": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30663": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30664": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30665": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30666": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30667": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30668": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30669": "res:/dx9/Model/WorldObject/Sun/WolfRayet.red",
    "30670": "res:/dx9/Model/WorldObject/Sun/Cataclysmic.red",
    "30671": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30672": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30673": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30674": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30675": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30676": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30677": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30678": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30679": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30680": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30681": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30682": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30683": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30684": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30685": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30686": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30687": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30688": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30689": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30690": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30691": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30692": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30693": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30694": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30695": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30696": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30697": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30698": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30699": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30700": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30701": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30702": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30703": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30704": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30705": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30706": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30707": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30708": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30709": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30710": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30711": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30712": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30713": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30714": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30715": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30762": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "30794": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCL.red",
    "30797": "res:/dx9/model/Structure/Talocan/Defense/TA_LinkageStructure/TA_Ls1/TA_Ls1.red",
    "30798": "res:/dx9/model/Structure/Talocan/Defense/TA_LinkageStructure/TA_Ls2/TA_Ls2.red",
    "30806": "res:/dx9/model/Structure/Talocan/Defense/TA_LinkageStructure/TA_Ls1/TA_Ls1.red",
    "30807": "res:/dx9/model/Structure/Talocan/Defense/TA_LinkageStructure/TA_Ls2/TA_Ls2.red",
    "30820": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "30831": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "30836": "res:/dx9/model/turret/salvage/s/salvage_s_t1.red",
    "30889": "res:/dx9/model/WorldObject/planet/ShatteredPlanet.red",
    "30903": "res:/dx9/model/Structure/Talocan/TowerStructure/TA_Ts1/TA_Ts1.red",
    "30905": "res:/dx9/model/Structure/Talocan/Defense/TA_LinkageStructure/TA_Ls1/TA_Ls1.red",
    "30967": "res:/dx9/model/Deployables/Generic/Container/Secure/SecureCS.red",
    "31961": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "32128": "res:/dx9/model/Structure/Universal/Construction/GB2/UC_GB2.red",
    "32211": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "32353": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "32367": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "32368": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "32369": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "32378": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "32379": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "32386": "res:/dx9/model/Celestial/Wormhole/Wormhole_Violent.red",
    "32387": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "32436": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_scourge.red",
    "32438": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_nova.red",
    "32440": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_inferno.red",
    "32442": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_mjolnir.red",
    "32443": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_mjolnir.red",
    "32444": "res:/dx9/model/Turret/Launcher/CitadelCruise/CitadelCruise_T1.red",
    "32445": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "32461": "res:/dx9/model/Turret/Launcher/Light/Light_T1.red",
    "32463": "res:/dx9/model/turret/launcher/light/light_impact_scourge.red",
    "32782": "res:/dx9/model/turret/launcher/light/light_impact_mjolnir.red",
    "32894": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "32895": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "32993": "res:/dx9/model/turret/launcher/festival/festival_impact_firework1.red",
    "32994": "res:/dx9/model/turret/launcher/festival/festival_impact_firework3.red",
    "32995": "res:/dx9/model/turret/launcher/festival/festival_impact_firework2.red",
    "33003": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCL.red",
    "33005": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCL.red",
    "33007": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCL.red",
    "33009": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCL.red",
    "33011": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCL.red",
    "33014": "res:/dx9/model/worldobject/cloud/amarrprimestationcloud.red",
    "33186": "res:/dx9/model/Artifact/Colony/Battleship/ColB1/ColB1_T1_Wreck/ColB1_T1_Wreck.red",
    "33234": "res:/dx9/model/Artifact/Colony/Battleship/ColB1/ColB1_T1_Wreck/ColB1_T1_Wreck.red",
    "33237": "res:/dx9/model/Artifact/Colony/Battleship/ColB1/ColB1_T1_Wreck/ColB1_T1_Wreck.red",
    "33238": "res:/dx9/model/Artifact/Colony/Battleship/ColB1/ColB1_T1_Wreck/ColB1_T1_Wreck.red",
    "33245": "res:/dx9/model/Artifact/Colony/Battleship/ColB1/ColB1_T1_Wreck/ColB1_T1_Wreck.red",
    "33246": "res:/dx9/model/Artifact/Colony/Battleship/ColB1/ColB1_T1_Wreck/ColB1_T1_Wreck.red",
    "33254": "res:/dx9/model/Artifact/Colony/Battleship/ColB1/ColB1_T1_Wreck/ColB1_T1_Wreck.red",
    "33255": "res:/dx9/model/Artifact/Colony/Battleship/ColB1/ColB1_T1_Wreck/ColB1_T1_Wreck.red",
    "33262": "res:/dx9/model/Artifact/Colony/Battleship/ColB1/ColB1_T1_Wreck/ColB1_T1_Wreck.red",
    "33263": "res:/dx9/model/Artifact/Colony/Battleship/ColB1/ColB1_T1_Wreck/ColB1_T1_Wreck.red",
    "33367": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "33368": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "33440": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33441": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33442": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33446": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33447": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "33448": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33449": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "33450": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33451": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "33452": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33453": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33454": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33455": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33456": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33457": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33458": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33459": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33460": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "33461": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33462": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33463": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33464": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33465": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33466": "res:/dx9/model/turret/launcher/rapidheavy/rapidheavy_t1.red",
    "33569": "res:/dx9/model/turret/launcher/festival/festival_impact_snow.red",
    "33571": "res:/dx9/model/turret/launcher/festival/festival_impact_firework1.red",
    "33572": "res:/dx9/model/turret/launcher/festival/festival_impact_firework3.red",
    "33573": "res:/dx9/model/turret/launcher/festival/festival_impact_firework2.red",
    "33589": "res:/dx9/model/Deployables/Mobile/dscandisruptor/MDD_T1.red",
    "34134": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34135": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34136": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34137": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34138": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34139": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34140": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34260": "res:/fisfx/generic/warp_disruption/warp_disruption_probe.red",
    "34272": "res:/dx9/model/Turret/Energy/Pulse/S/Pulse_Medium_T1.red",
    "34273": "res:/dx9/model/turret/energy/pulse/s/pulse_medium_t1.red",
    "34274": "res:/dx9/model/Turret/Energy/Pulse/M/Pulse_Heavy_T1.red",
    "34275": "res:/dx9/model/turret/energy/pulse/m/pulse_heavy_t1.red",
    "34276": "res:/dx9/model/Turret/Energy/Pulse/L/Pulse_Mega_T1.red",
    "34277": "res:/dx9/model/turret/energy/pulse/l/pulse_mega_t1.red",
    "34278": "res:/dx9/model/Turret/Hybrid/Blast/S/Blast_Neutron_T1.red",
    "34279": "res:/dx9/model/turret/hybrid/blast/s/blast_neutron_t1.red",
    "34280": "res:/dx9/model/Turret/Hybrid/Blast/M/Blast_Neutron_T1.red",
    "34281": "res:/dx9/model/turret/hybrid/blast/m/blast_neutron_t1.red",
    "34282": "res:/dx9/model/Turret/Hybrid/Blast/L/Blast_Neutron_T1.red",
    "34283": "res:/dx9/model/turret/hybrid/blast/l/blast_neutron_t1.red",
    "34284": "res:/dx9/model/Turret/Projectile/Auto/S/Auto_200mm_T1.red",
    "34285": "res:/dx9/model/turret/projectile/auto/s/auto_200mm_t1.red",
    "34286": "res:/dx9/model/Turret/Projectile/Auto/M/Auto_425mm_T1.red",
    "34287": "res:/dx9/model/turret/projectile/auto/m/auto_425mm_t1.red",
    "34288": "res:/dx9/model/Turret/Projectile/Auto/L/Auto_800mmHeavy_T1.red",
    "34289": "res:/dx9/model/turret/projectile/auto/l/auto_800mmheavy_t1.red",
    "34290": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "34291": "res:/dx9/model/turret/launcher/rocket/rocket_t1.red",
    "34292": "res:/dx9/model/Turret/Launcher/HeavyAssault/HeavyAssault_T1.red",
    "34293": "res:/dx9/model/turret/launcher/heavyassault/heavyassault_t1.red",
    "34294": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "34295": "res:/dx9/model/turret/launcher/torpedo/torpedo_t1.red",
    "34301": "res:/dx9/model/Structure/Universal/DamageDampener/UDD1.red",
    "34305": "res:/dx9/model/Structure/Universal/DamageDampener/UDD1.red",
    "34310": "res:/dx9/model/outposts/jove/jo1/jo1_t1.red",
    "34315": "res:/dx9/model/structure/talocan/defense/ta_linkagestructure/ta_ls1/cloaked/ta_ls1_cloaked.red",
    "34316": "res:/dx9/model/structure/sleeper/engineeringstation/sl_es1/cloaked/sl_es1_cloaked.red",
    "34330": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "34331": "res:/dx9/model/celestial/sun/sun_purple_sun_01a.red",
    "34333": "res:/dx9/model/worldobject/cloud/UnknownCloud.red",
    "34338": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34366": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34367": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34368": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34369": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34370": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34371": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34372": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34375": "res:/dx9/model/Structure/Talocan/Defense/TA_LinkageStructure/TA_Ls1/TA_Ls1.red",
    "34387": "res:/FisFX/Celestial/SpatialRift_Rt_01a.red",
    "34404": "res:/dx9/model/Artifact/Universal/ConcentricRings/UCR.red",
    "34405": "res:/dx9/model/Artifact/Universal/ConcentricRings/UCR.red",
    "34406": "res:/dx9/model/Structure/Universal/DamageDampener/UDD1.red",
    "34439": "res:/dx9/model/celestial/wormhole/wormhole.red",
    "34494": "res:/dx9/model/celestial/wormhole/wormhole_sleeper.red",
    "34580": "res:/dx9/model/turret/jove/l/jove_l_t1.red",
    "34593": "res:/fisfx/module/entosislink_st_t1a.red",
    "34595": "res:/fisfx/module/entosislink_st_t1a.red",
    "34826": "res:/fisfx/module/entosislink_st_t1a.red",
    "35645": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "35650": "res:/dx9/model/celestial/wormhole/wormhole_sleeper.red",
    "35651": "res:/dx9/model/celestial/wormhole/wormhole_sleeper.red",
    "35652": "res:/dx9/model/celestial/wormhole/wormhole_sleeper.red",
    "35653": "res:/dx9/model/celestial/wormhole/wormhole_sleeper.red",
    "35654": "res:/dx9/model/celestial/wormhole/wormhole_sleeper.red",
    "35921": "res:/dx9/model/turret/structure/antiship/antiship_t1.red",
    "35923": "res:/dx9/model/turret/structure/launcher/launcher_t1.red",
    "35924": "res:/dx9/model/turret/Structure/UtilityA/UtilityA_Neutralizer_t1.red",
    "35925": "res:/dx9/model/turret/Structure/UtilityB/UtilityB_Neutralizer_t1.red",
    "35928": "res:/dx9/model/turret/structure/doomsday/doomsday_t1.red",
    "35939": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_Bumping_T1.red",
    "35940": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_ECM_T1.red",
    "35941": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_SensorDampener_T1.red",
    "35943": "res:/dx9/model/turret/structure/utilityc/utilityc_stasisweb_t1.red",
    "35945": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_TrackingDisruptor_T1.red",
    "35947": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_TargetPainter_T1.red",
    "35949": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_WarpScrambler_T1.red",
    "36464": "res:/fisfx/environment/cloudvortex_01a.red",
    "36523": "res:/dx9/model/Deployables/Mobile/dscandisruptor/MDD_T1.red",
    "37288": "res:/dx9/model/turret/launcher/rapidtorpedo/rapidtorpedo_t1.red",
    "37289": "res:/dx9/model/turret/projectile/auto/xl/auto_4x800mm_t1.red",
    "37290": "res:/dx9/model/turret/energy/pulse/xl/pulse_quadmega_t1.red",
    "37291": "res:/dx9/model/turret/hybrid/blast/xl/blast_quadneutron_t1.red",
    "37292": "res:/dx9/model/turret/launcher/rapidtorpedo/rapidtorpedo_t1.red",
    "37293": "res:/dx9/model/turret/launcher/rapidtorpedo/rapidtorpedo_t1.red",
    "37294": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_t1.red",
    "37295": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_t1.red",
    "37296": "res:/dx9/model/turret/energy/pulse/xl/pulse_quadmega_t1.red",
    "37297": "res:/dx9/model/turret/energy/pulse/xl/pulse_quadmega_t1.red",
    "37298": "res:/dx9/model/turret/energy/pulse/xl/pulse_gigadual_t1.red",
    "37299": "res:/dx9/model/turret/energy/beam/xl/beam_gigadual_t1.red",
    "37300": "res:/dx9/model/turret/hybrid/blast/xl/blast_quadneutron_t1.red",
    "37301": "res:/dx9/model/turret/hybrid/blast/xl/blast_quadneutron_t1.red",
    "37302": "res:/dx9/model/turret/hybrid/blast/xl/blast_ion_t1.red",
    "37303": "res:/dx9/model/turret/hybrid/rail/xl/rail_1000mmdual_t1.red",
    "37304": "res:/dx9/model/turret/projectile/auto/xl/auto_4x800mm_t1.red",
    "37305": "res:/dx9/model/turret/projectile/auto/xl/auto_4x800mm_t1.red",
    "37306": "res:/dx9/model/turret/projectile/auto/xl/auto_6x2500mm_t1.red",
    "37307": "res:/dx9/model/turret/projectile/artil/xl/artil_3500mmquad_t1.red",
    "37450": "res:/dx9/model/turret/mining/standardice/standardice_T1.red",
    "37451": "res:/dx9/model/turret/mining/standardice/standardice_t1.red",
    "37452": "res:/dx9/model/turret/mining/standardice/standardice_T1.red",
    "37608": "res:/fisfx/generic/warp_disruption/warp_disruption_field_generator.red",
    "37610": "res:/fisfx/generic/warp_disruption/warp_disruption_field_generator.red",
    "37611": "res:/fisfx/generic/warp_disruption/warp_disruption_field_generator.red",
    "37612": "res:/fisfx/generic/warp_disruption/warp_disruption_field_generator.red",
    "37613": "res:/fisfx/generic/warp_disruption/warp_disruption_field_generator.red",
    "37843": "res:/dx9/model/turret/structure/antiship/antiship_impact_longrange.red",
    "37844": "res:/dx9/model/turret/structure/antiship/antiship_impact_mediumrange.red",
    "37845": "res:/dx9/model/turret/structure/antiship/antiship_impact_shortrange.red",
    "37849": "res:/dx9/model/turret/structure/launcher/launcher_impact_longrange.red",
    "37850": "res:/dx9/model/turret/structure/launcher/launcher_impact_shortrange.red",
    "37851": "res:/dx9/model/turret/structure/launcher/launcher_impact_edrainmediumrange.red",
    "37855": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "37856": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "37857": "res:/dx9/model/turret/launcher/cruise/cruise_impact_inferno.red",
    "40307": "res:/fisfx/module/entosislink_st_t1a.red",
    "40308": "res:/fisfx/module/entosislink_st_t1a.red",
    "40309": "res:/fisfx/module/entosislink_st_t1a.red",
    "40310": "res:/fisfx/module/entosislink_st_t1a.red",
    "41063": "res:/dx9/model/turret/energy/pulse/xl/pulse_quadmega_t1.red",
    "41064": "res:/dx9/model/turret/energy/pulse/xl/pulse_quadmega_t1.red",
    "41065": "res:/dx9/model/turret/energy/pulse/xl/pulse_quadmega_t1.red",
    "41066": "res:/dx9/model/turret/energy/pulse/xl/pulse_quadmega_t1.red",
    "41067": "res:/dx9/model/turret/energy/pulse/xl/pulse_quadmega_t1.red",
    "41076": "res:/dx9/model/turret/hybrid/blast/xl/blast_quadneutron_t1.red",
    "41077": "res:/dx9/model/turret/hybrid/blast/xl/blast_quadneutron_t1.red",
    "41078": "res:/dx9/model/turret/hybrid/blast/xl/blast_quadneutron_t1.red",
    "41079": "res:/dx9/model/turret/hybrid/blast/xl/blast_quadneutron_t1.red",
    "41080": "res:/dx9/model/turret/projectile/auto/xl/auto_4x800mm_t1.red",
    "41081": "res:/dx9/model/turret/projectile/auto/xl/auto_4x800mm_t1.red",
    "41082": "res:/dx9/model/turret/projectile/auto/xl/auto_4x800mm_t1.red",
    "41083": "res:/dx9/model/turret/projectile/auto/xl/auto_4x800mm_t1.red",
    "41099": "res:/dx9/model/turret/energy/pulse/xl/pulse_gigadual_t1.red",
    "41100": "res:/dx9/model/turret/energy/pulse/xl/pulse_gigadual_t1.red",
    "41101": "res:/dx9/model/turret/energy/pulse/xl/pulse_gigadual_t1.red",
    "41102": "res:/dx9/model/turret/energy/pulse/xl/pulse_gigadual_t1.red",
    "41103": "res:/dx9/model/turret/energy/pulse/xl/pulse_gigadual_t1.red",
    "41104": "res:/dx9/model/turret/energy/pulse/xl/pulse_gigadual_t1.red",
    "41114": "res:/dx9/model/turret/energy/beam/xl/beam_gigadual_t1.red",
    "41115": "res:/dx9/model/turret/energy/beam/xl/beam_gigadual_t1.red",
    "41116": "res:/dx9/model/turret/energy/beam/xl/beam_gigadual_t1.red",
    "41117": "res:/dx9/model/turret/energy/beam/xl/beam_gigadual_t1.red",
    "41118": "res:/dx9/model/turret/energy/beam/xl/beam_gigadual_t1.red",
    "41119": "res:/dx9/model/turret/energy/beam/xl/beam_gigadual_t1.red",
    "41126": "res:/dx9/model/turret/hybrid/blast/xl/blast_ion_t1.red",
    "41127": "res:/dx9/model/turret/hybrid/blast/xl/blast_ion_t1.red",
    "41128": "res:/dx9/model/turret/hybrid/blast/xl/blast_ion_t1.red",
    "41129": "res:/dx9/model/turret/hybrid/blast/xl/blast_ion_t1.red",
    "41130": "res:/dx9/model/turret/hybrid/blast/xl/blast_ion_t1.red",
    "41138": "res:/dx9/model/turret/hybrid/rail/xl/rail_1000mmdual_t1.red",
    "41139": "res:/dx9/model/turret/hybrid/rail/xl/rail_1000mmdual_t1.red",
    "41140": "res:/dx9/model/turret/hybrid/rail/xl/rail_1000mmdual_t1.red",
    "41141": "res:/dx9/model/turret/hybrid/rail/xl/rail_1000mmdual_t1.red",
    "41142": "res:/dx9/model/turret/hybrid/rail/xl/rail_1000mmdual_t1.red",
    "41150": "res:/dx9/model/turret/projectile/auto/xl/auto_6x2500mm_t1.red",
    "41151": "res:/dx9/model/turret/projectile/auto/xl/auto_6x2500mm_t1.red",
    "41152": "res:/dx9/model/turret/projectile/auto/xl/auto_6x2500mm_t1.red",
    "41153": "res:/dx9/model/turret/projectile/auto/xl/auto_6x2500mm_t1.red",
    "41154": "res:/dx9/model/turret/projectile/auto/xl/auto_6x2500mm_t1.red",
    "41156": "res:/dx9/model/turret/projectile/artil/xl/artil_3500mmquad_t1.red",
    "41157": "res:/dx9/model/turret/projectile/artil/xl/artil_3500mmquad_t1.red",
    "41158": "res:/dx9/model/turret/projectile/artil/xl/artil_3500mmquad_t1.red",
    "41159": "res:/dx9/model/turret/projectile/artil/xl/artil_3500mmquad_t1.red",
    "41160": "res:/dx9/model/turret/projectile/artil/xl/artil_3500mmquad_t1.red",
    "41174": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_t1.red",
    "41175": "res:/dx9/model/Turret/Launcher/CitadelCruise/CitadelCruise_T1.red",
    "41176": "res:/dx9/model/Turret/Launcher/CitadelCruise/CitadelCruise_T1.red",
    "41177": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41178": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41179": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41180": "res:/dx9/model/Turret/Launcher/CitadelTorpedo/CitadelTorpedo_T1.red",
    "41181": "res:/dx9/model/Turret/Launcher/CitadelTorpedo/CitadelTorpedo_T1.red",
    "41182": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_t1.red",
    "41183": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41184": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41185": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41187": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41190": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41223": "res:/dx9/model/turret/launcher/rapidtorpedo/rapidtorpedo_t1.red",
    "41224": "res:/dx9/model/turret/launcher/rapidtorpedo/rapidtorpedo_t1.red",
    "41225": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41226": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41227": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41228": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41229": "res:/dx9/model/turret/Launcher/Torpedo/Torpedo_T1.red",
    "41233": "res:/fisfx/module/modular_anomalybeacon_rt_t1a.red",
    "41270": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_scourge.red",
    "41271": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_scourge.red",
    "41272": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_nova.red",
    "41273": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_nova.red",
    "41274": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_mjolnir.red",
    "41275": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_mjolnir.red",
    "41276": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_inferno.red",
    "41277": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_inferno.red",
    "41282": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_mjolnir.red",
    "41283": "res:/dx9/model/turret/launcher/citadeltorpedo/citadeltorpedo_impact_mjolnir.red",
    "41286": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_inferno.red",
    "41287": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_inferno.red",
    "41290": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_scourge.red",
    "41291": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_scourge.red",
    "41294": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_nova.red",
    "41295": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_nova.red",
    "41298": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_mjolnir.red",
    "41299": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_mjolnir.red",
    "41300": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_mjolnir.red",
    "41301": "res:/dx9/model/turret/launcher/citadelcruise/citadelcruise_impact_mjolnir.red",
    "41342": "res:/fisfx/module/modular_anomalybeacon_rt_t1a.red",
    "41343": "res:/fisfx/module/modular_anomalybeacon_rt_t1a.red",
    "41344": "res:/fisfx/module/modular_anomalybeacon_rt_t1a.red",
    "41345": "res:/fisfx/module/modular_anomalybeacon_rt_t1a.red",
    "41346": "res:/fisfx/module/modular_anomalybeacon_rt_t1a.red",
    "41347": "res:/fisfx/module/modular_anomalybeacon_rt_t1a.red",
    "41348": "res:/fisfx/module/modular_anomalybeacon_rt_t1a.red",
    "41349": "res:/fisfx/module/modular_anomalybeacon_rt_t1a.red",
    "41402": "res:/fisfx/generic/warp_disruption/warp_disruption_probe.red",
    "41540": "res:/fisfx/module/modular_anomalybeacon_rt_t1a.red",
    "41548": "res:/dx9/model/turret/launcher/bomb/bomb_impact_scourge.red",
    "41549": "res:/dx9/model/turret/launcher/bomb/bomb_impact_mjolnir.red",
    "41550": "res:/dx9/model/turret/launcher/bomb/bomb_impact_inferno.red",
    "41551": "res:/dx9/model/turret/launcher/bomb/bomb_impact_nova.red",
    "41567": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "42532": "res:/FisFX/Environment/cloudvortex_01b.red",
    "42536": "res:/FisFX/Generic/lightning/npe_customlightning_01a.red",
    "42697": "res:/FisFX/Generic/bubble/drifterbubble_01a.red",
    "42814": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "42901": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "42902": "res:/FisFX/Celestial/SpatialRift_Rt_01a.red",
    "43549": "res:/FisFX/Environment/cloudvortex_01b.red",
    "43917": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "43918": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "44066": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "44102": "res:/dx9/model/Turret/Launcher/Rocket/Rocket_T1.red",
    "44105": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "44109": "res:/dx9/model/turret/launcher/holiday/holiday_t1.red",
    "44110": "res:/dx9/model/turret/launcher/holiday/holiday_impact_winter1.red",
    "44260": "res:/dx9/model/turret/launcher/holiday/holiday_impact_amarr1.red",
    "44261": "res:/dx9/model/turret/launcher/holiday/holiday_impact_angels1.red",
    "44262": "res:/dx9/model/turret/launcher/holiday/holiday_impact_bloodraiders1.red",
    "44263": "res:/dx9/model/turret/launcher/holiday/holiday_impact_caldari1.red",
    "44264": "res:/dx9/model/turret/launcher/holiday/holiday_impact_celebration1.red",
    "44265": "res:/dx9/model/turret/launcher/holiday/holiday_impact_celebration2.red",
    "44266": "res:/dx9/model/turret/launcher/holiday/holiday_impact_celebration3.red",
    "44267": "res:/dx9/model/turret/launcher/holiday/holiday_impact_crimsonharvest1.red",
    "44268": "res:/dx9/model/turret/launcher/holiday/holiday_impact_gallente1.red",
    "44269": "res:/dx9/model/turret/launcher/holiday/holiday_impact_guristas1.red",
    "44270": "res:/dx9/model/turret/launcher/holiday/holiday_impact_minmatar1.red",
    "44271": "res:/dx9/model/turret/launcher/holiday/holiday_impact_sansha1.red",
    "44272": "res:/dx9/model/turret/launcher/holiday/holiday_impact_serpentis1.red",
    "45009": "res:/dx9/model/turret/structure/moonminer/MoonMiner_T1.red",
    "45030": "res:/dx9/model/celestial/sun/sun_yellow_01b.red",
    "45031": "res:/dx9/model/celestial/sun/sun_orange_01b.red",
    "45032": "res:/dx9/model/celestial/sun/sun_orange_01c.red",
    "45033": "res:/dx9/model/celestial/sun/sun_red_01b.red",
    "45034": "res:/dx9/model/celestial/sun/sun_blue_01b.red",
    "45035": "res:/dx9/model/celestial/sun/sun_white_01b.red",
    "45036": "res:/dx9/model/celestial/sun/sun_pink_hazy_01b.red",
    "45037": "res:/dx9/model/celestial/sun/sun_orange_sun_01b.red",
    "45038": "res:/dx9/model/celestial/sun/sun_pink_sun_small_01b.red",
    "45039": "res:/dx9/model/celestial/sun/sun_orange_radiating_01b.red",
    "45040": "res:/dx9/model/celestial/sun/sun_orange_radiating_01c.red",
    "45041": "res:/dx9/model/celestial/sun/sun_yellow_small_01b.red",
    "45042": "res:/dx9/model/celestial/sun/sun_white_tiny_01b.red",
    "45046": "res:/dx9/model/celestial/sun/sun_blue_01c.red",
    "45047": "res:/dx9/model/celestial/sun/sun_yellow_01c.red",
    "46335": "res:/dx9/model/Deployables/Generic/Container/Standard/StandardCS.red",
    "46565": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "46566": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "46567": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "46568": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "46569": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "46570": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "46571": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "46572": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "46573": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "46574": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "46575": "res:/dx9/model/turret/Structure/UtilityB/UtilityB_Neutralizer_t1.red",
    "46577": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_ECM_T1.red",
    "47066": "res:/dx9/model/testcases/dungeonexploration/clouds/cloudplanes/cloud_01.red",
    "47067": "res:/dx9/model/testcases/dungeonexploration/clouds/cloudplanes/cloud_02.red",
    "47152": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "47203": "res:/dx9/model/testcases/dungeonexploration/clouds/storm_01/storm_01.red",
    "47264": "res:/dx9/model/turret/launcher/holiday/holiday_impact_valentine1.red",
    "47272": "res:/dx9/model/turret/atomic/s/atomic_s01_t1.red",
    "47273": "res:/dx9/model/turret/atomic/m/atomic_m01_t1.red",
    "47274": "res:/dx9/model/turret/atomic/l/atomic_l01_t1.red",
    "47300": "res:/dx9/model/turret/launcher/holiday/holiday_impact_capsuleerday1.red",
    "47301": "res:/dx9/model/turret/launcher/holiday/holiday_impact_crimsonharvest2.red",
    "47302": "res:/dx9/model/turret/launcher/holiday/holiday_impact_winter2.red",
    "47323": "res:/dx9/model/turret/structure/antiship/antiship_t1.red",
    "47325": "res:/dx9/model/turret/structure/launcher/launcher_t1.red",
    "47329": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "47330": "res:/dx9/model/turret/Structure/UtilityA/UtilityA_Neutralizer_t1.red",
    "47332": "res:/dx9/model/turret/Structure/UtilityB/UtilityB_Neutralizer_t1.red",
    "47334": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_WarpScrambler_T1.red",
    "47338": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_ECM_T1.red",
    "47351": "res:/dx9/model/turret/structure/utilityc/utilityc_stasisweb_t1.red",
    "47364": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_TrackingDisruptor_T1.red",
    "47366": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_TargetPainter_T1.red",
    "47368": "res:/dx9/model/turret/Structure/UtilityC/UtilityC_SensorDampener_T1.red",
    "47380": "res:/fisfx/environment/cloudcover/cc_abyssal_darkness_s3_01a.red",
    "47383": "res:/fisfx/environment/cloudcover/cc_abyssal_electricstorm_s3_01a.red",
    "47386": "res:/fisfx/environment/cloudcover/cc_abyssal_caustictoxin_s3_01a.red",
    "47389": "res:/fisfx/environment/cloudcover/cc_abyssal_xenongas_s3_01a.red",
    "47392": "res:/fisfx/environment/cloudcover/cc_abyssal_infernal_s3_01a.red",
    "47398": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "47399": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "47400": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "47401": "res:/dx9/model/Structure/Drone/Defense/Wall_Bunker.red",
    "47430": "res:/fisfx/lensflare/abyssal_dark.red",
    "47431": "res:/fisfx/lensflare/abyssal_blue.red",
    "47432": "res:/fisfx/lensflare/abyssal_red.red",
    "47433": "res:/fisfx/lensflare/abyssal_green.red",
    "47434": "res:/fisfx/lensflare/abyssal_purple.red",
    "47435": "res:/fisfx/lensflare/abyssal_white.red",
    "47436": "res:/fisfx/environment/aoe/aoe_causticwave_10k_rt_01a.red",
    "47439": "res:/fisfx/environment/aoe/aoe_bioluminescence_10k_rt_01a.red",
    "47440": "res:/fisfx/environment/aoe/aoe_bioluminescence_30k_rt_01a.red",
    "47441": "res:/fisfx/environment/aoe/aoe_bioluminescence_80k_rt_01a.red",
    "47446": "res:/dx9/model/turret/jove/m/jove_m_t1.red",
    "47451": "res:/dx9/model/celestial/environment/planet/Moon/moon.red",
    "47452": "res:/dx9/model/celestial/environment/planet/sandstorm/sandstorm.red",
    "47453": "res:/dx9/model/celestial/environment/planet/Lava/lava.red",
    "47454": "res:/dx9/model/celestial/environment/planet/Plasma/plasma.red",
    "47455": "res:/dx9/model/celestial/environment/planet/Gas/gas.red",
    "47456": "res:/dx9/model/celestial/environment/planet/Terrestrial/terrestrial.red",
    "47458": "res:/fisfx/travel/abyssal/Abyssal_SpaceTear_Trace_01a.red",
    "47460": "res:/FisFX/Deployable/CynoBeaconBO_Rt_T1a.red",
    "47461": "res:/FisFX/Deployable/CynoBeaconBO_Rt_T1a.red",
    "47462": "res:/FisFX/Deployable/CynoBeaconBO_Rt_T1a.red",
    "47463": "res:/FisFX/Deployable/CynoBeacon_Rt_T1a.red",
    "47464": "res:/FisFX/Deployable/CynoBeacon_Rt_T1a.red",
    "47465": "res:/fisfx/environment/boundary/boundary_bubble.red",
    "47467": "res:/fisfx/environment/aoe/aoe_causticwave_30k_rt_01a.red",
    "47468": "res:/fisfx/environment/aoe/aoe_causticwave_80k_rt_01a.red",
    "47472": "res:/fisfx/environment/aoe/aoe_filament_30k_rt_01a.red",
    "47473": "res:/fisfx/environment/aoe/aoe_filament_80k_rt_01a.red",
    "47488": "res:/dx9/model/Structure/Drone/Defense/Wall_MissileBattery.red",
    "47620": "res:/fisfx/environment/aoe/aoe_filament_10k_rt_01a.red",
    "47656": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/5810/clouds_5810.red",
    "47657": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6047/clouds_6047.red",
    "47658": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6200/clouds_6200.red",
    "47659": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6201/clouds_6201.red",
    "47660": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6202/clouds_6202.red",
    "47661": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6048/clouds_6048.red",
    "47662": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6049/clouds_6049.red",
    "47663": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6050/clouds_6050.red",
    "47664": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6051/clouds_6051.red",
    "47665": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6052/clouds_6052.red",
    "47666": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6203/clouds_6203.red",
    "47667": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6204/clouds_6204.red",
    "47668": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6205/clouds_6205.red",
    "47669": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6206/clouds_6206.red",
    "47670": "res:/dx9/model/celestial/environment/cloud/volumetric/dungeons/6207/clouds_6207.red",
    "47671": "res:/dx9/model/celestial/environment/rock/asteroidset_01/debris/as1_debris_01.red",
    "47672": "res:/dx9/model/celestial/environment/rock/asteroidset_01/debris/as1_debris_02.red",
    "47673": "res:/dx9/model/celestial/environment/rock/asteroidset_01/debris/as1_debris_03.red",
    "47674": "res:/dx9/model/celestial/environment/rock/pillarset_01/debris/ps1_debris_01.red",
    "47675": "res:/dx9/model/celestial/environment/rock/pillarset_01/debris/ps1_debris_02.red",
    "47676": "res:/dx9/model/celestial/environment/rock/pillarset_01/debris/ps1_debris_03.red",
    "47677": "res:/dx9/model/celestial/environment/crystal/crystalset_01/debris/cs1_debris_01.red",
    "47678": "res:/dx9/model/celestial/environment/crystal/crystalset_01/debris/cs1_debris_02.red",
    "47679": "res:/dx9/model/celestial/environment/crystal/crystalset_01/debris/cs1_debris_03.red",
    "47760": "res:/fisfx/lensflare/abyssal_yellow.red",
    "47833": "res:/fisfx/travel/abyssal/Abyssal_SpaceTear_Trace_01b.red",
    "47834": "res:/fisfx/travel/abyssal/Abyssal_SpaceTear_Trace_01a.red",
    "47861": "res:/dx9/model/WorldObject/Beacon/TrigBeacon01.red",
    "47912": "res:/dx9/model/turret/atomic/s/atomic_s01_t1.red",
    "47913": "res:/dx9/model/turret/atomic/s/atomic_s01_t1.red",
    "47914": "res:/dx9/model/turret/atomic/s/atomic_s01_t1.red",
    "47915": "res:/dx9/model/turret/atomic/s/atomic_s01_t1.red",
    "47916": "res:/dx9/model/turret/atomic/m/atomic_m01_t1.red",
    "47917": "res:/dx9/model/turret/atomic/m/atomic_m01_t1.red",
    "47918": "res:/dx9/model/turret/atomic/m/atomic_m01_t1.red",
    "47919": "res:/dx9/model/turret/atomic/m/atomic_m01_t1.red",
    "47920": "res:/dx9/model/turret/atomic/l/atomic_l01_t1.red",
    "47921": "res:/dx9/model/turret/atomic/l/atomic_l01_t1.red",
    "47922": "res:/dx9/model/turret/atomic/l/atomic_l01_t1.red",
    "47923": "res:/dx9/model/turret/atomic/l/atomic_l01_t1.red",
    "47949": "res:/dx9/model/celestial/environment/planet/plasma_blood/plasma_blood.red",
    "47950": "res:/dx9/model/celestial/environment/planet/sandstorm_yellow/sandstorm_yellow.red",
    "48079": "res:/fisfx/travel/abyssal/Abyssal_SpaceTear_Trace_01a.red",
    "48084": "res:/fisfx/travel/abyssal/Abyssal_SpaceTear_Trace_01a.red",
    "48085": "res:/fisfx/travel/abyssal/Abyssal_SpaceTear_Trace_01a.red",
    "48093": "res:/FisFX/Deployable/CynoBeaconBO_Rt_T1a.red",
    "48094": "res:/FisFX/Deployable/CynoBeaconBO_Rt_T1a.red",
    "48594": "res:/dx9/model/Structure/Drone/Defense/Wall_MissileBattery.red"
};