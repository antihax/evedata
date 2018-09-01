var package,
    urlVars = getUrlVars();

$.ajax({
    url: "https://static.evedata.org/file/evedata-killmails/" + urlVars["id"] + ".json.gz",
    dataType: 'native',
    xhrFields: {
        responseType: 'arraybuffer'
    },
    success: function (d) {
        try {
            package = $.parseJSON(pako.inflate(d, { to: 'string' }));
            console.log(package);
            $(document).ready(function () {
                populateModules(package)
                getShip(package);
                getAttackers(package);
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

function populateModules(package) {
    typeNames = package.attributes.types;
    $.each(package.attributes.modules, function (k, i) {
        switch (i.location) {
            case 27: // hiSlots
                setModuleSlot("high", "1", i);
                break;
            case 28: // hiSlots
                setModuleSlot("high", "2", i);
                break;
            case 29: // hiSlots
                setModuleSlot("high", "3", i);
                break;
            case 30: // hiSlots
                setModuleSlot("high", "4", i);
                break;
            case 31: // hiSlots
                setModuleSlot("high", "5", i);
                break;
            case 32: // hiSlots
                setModuleSlot("high", "6", i);
                break;
            case 33: // hiSlots
                setModuleSlot("high", "7", i);
                break;
            case 34: // hiSlots
                setModuleSlot("high", "8", i);
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

function getPortrait(a) {
    if (a.character_id != undefined) {
        return "character/" + a.character_id + "_64.jpg";
    } else if (a.corporation_id != undefined) {
        return "corporation/" + a.corporation_id + "_64.png";
    } else {
        return "corporation/" + a.faction_id + "_64.png";
    }
}

function getShipImage(a) {
    return a.ship_type_id + "_32.png";
}

function getWeaponImage(a) {
    return a.weapon_type_id == undefined ?
        a.ship_type_id + "_32.png" :
        a.weapon_type_id + "_32.png";
}

function getThatStuff(a) {
    var theStuff = "";

    if (a.character_id != undefined) theStuff += package.nameMap[a.character_id] + "<br>";
    if (a.corporation_id != undefined) theStuff += package.nameMap[a.corporation_id] + "<br>"
    if (a.alliance_id != undefined) theStuff += package.nameMap[a.alliance_id] + "<br>"
    if (a.faction_id != undefined) theStuff += package.nameMap[a.faction_id] + "<br>"

    return theStuff;
}

function getAttackers(package) {
    $("#numInvolved").text(simpleVal(package.killmail.attackers.length) + " Involved");
    var stripe = false;
    $.each(package.killmail.attackers, function (k, a) {
        var row = `
            <div class="row killmail" style="background-color: ${stripe ? "#06100a;" : "#16201a;"} padding: 0px;">
                <div class="col-xs-2 killmail" style="width: 64px">
                    <img src="//imageserver.eveonline.com/${getPortrait(a)}" style="width:64px; height: 64px">
                </div>
                <div class="col-xs-1 killmail" style="width: 32px">
                    <img src="//imageserver.eveonline.com/type/${getShipImage(a)}" style="width:32px; height: 32px">
                    <img src="//imageserver.eveonline.com/type/${getWeaponImage(a)}" style="width:32px; height: 32px">
                </div>
                <div class="col-xs-9 killmail" style="width: 264px;">
                    <div class="row" style="height: 64px; padding: 5px;">
                        <div class="col-xs-9">
                           ${getThatStuff(a)}
                        </div>
                        <div class="col-xs-3" style="height: 64px; text-align: right">
                            ${simpleVal(a.damage_done)}
                        </div>
                    </div>
                </div>
            </div>`;

        $("#attackers").append(row);
        stripe = !stripe;
    });
}

function getShipwebGL(package) {
    $("#shipImage").attr("src", "//imageserver.eveonline.com/Render/" + package.attributes.typeID + "_256.png")
    try {
        var mat4 = ccpwgl_int.math.mat4,
            rotation = 0.0,
            direction = 0.001,
            canvas = document.getElementById('shipCanvas');

        ccpwgl.initialize(canvas, {});

        camera = ccpwgl.createCamera(canvas, {}, true);

        scene = ccpwgl.loadScene('res:/dx9/scene/universe/m10_cube.red');
        var ship = scene.loadShip(package.dna);
        scene.loadSun('res:/fisfx/lensflare/purple_sun.red');


        var sizes = ['c', 'd', 'h', 'l', 'm', 's', 't'];
        var races = ['amarr', 'angel', 'blooodraider', 'caldari', 'concord', 'gallente', 'generic', 'jove', 'minmatar', 'ore', 'rogue', 'sansha', 'sepentis', 'sleeper', 'soct', 'soe', 'talocan'];
        var radius = 50;

        var explosions = [];
        var currentTime = 0;

        function getRandomExplosion(explodionData) {
            var size = sizes[Math.floor(Math.random() * sizes.length)];
            var race = races[Math.floor(Math.random() * races.length)];
            var explosion = scene.loadObject('res:/fisfx/deathexplosion/death_' + size + '_' + race + '.red', function () {
                this.wrappedObjects[0].Start();
                explodionData[1] = currentTime + this.wrappedObjects[0].duration;
            });
            explosion.setTransform([
                1, 0, 0, 0,
                0, 1, 0, 0,
                0, 0, 1, 0,
                Math.random() * 2 * radius - radius, Math.random() * 2 * radius - radius, Math.random() * 2 * radius - radius, 1
            ]);
            return explosion;
        }

        function spawnExplosion() {
            for (var i = 0; i < explosions.length;) {
                if (currentTime > explosions[i][1]) {
                    scene.removeObject(scene.indexOf(explosions[i][0]));
                    explosions.splice(i, 1);
                }
                else {
                    ++i;
                }
            }
            var explosion = [null, 0];
            explosion[0] = getRandomExplosion(explosion);
            explosions.push(explosion);
            window.setTimeout(spawnExplosion, 2000 + Math.random() * 2000);
        }

        //spawnExplosion();

        ccpwgl.onPreRender = function (dt) {
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
                $("#totalDamage").html(v.toFixed(0) + " DPS")
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

