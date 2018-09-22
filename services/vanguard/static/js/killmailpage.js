var package,
    urlVars = getUrlVars(),
    ship, canvas, camera;

    var killmail = new Killmail(urlVars["id"], function (k) {
        try {
            package = k.getKillmail();
            console.log(k.getEFT())
            console.log(package)
            $(document).ready(function () {
                getShip(package);
                populateModules(package)
                getAttackers(package);
                getTypes(package);
                getVictimInformation(package.killmail.victim);
                getSystemInfo(package.systemInfo);

                new ClipboardJS('.clipboardCopy', {
                    text: function(trigger) {
                        showAlert("copied to clipboard", "success");
                        return k.getEFT();
                    }
                });
            });
        } catch (e){
            showAlert("Failed to read killmail: " + e, "danger")
        }
    });

function setResonancePercentage(resonance, value) {
    if (!value) {
        value = 1;
    }
    value = 1 - value;
    if (value < 0) { value = 0 }
    value = (value * 100).toFixed(0);
    $('#' + resonance).css('width', value + '%').attr('aria-valuenow', value);
    $('#' + resonance).text(value + "%");
}

function setModuleSlot(type, slot, i) {
    typeNames = package.attributes.types;

    $("#" + type + slot).prepend('<img class="ttp" src="//imageserver.eveonline.com/Type/' + i.typeID + '_32.png" title="' + typeNames[i.typeID] + '" style="height: 32px; width: 32px; z-index:3000">')
    $("#" + type + slot + " .ttp").tooltipster({
        contentAsHTML: true,
        position: 'bottom',
        side: 'bottom',
        viewportAware: false,
        functionInit: function (ins) {
            ins.content(moduleToolTip(i, typeNames));
        },
    });
    if (i.chargeTypeID > 0) $("#" + type + slot + "l").prepend('<img src="//imageserver.eveonline.com/Type/' + i.chargeTypeID + '_32.png" style="height: 32px; width: 32px;">')
}

function moduleToolTip(i, t) {
    var pm = package.priceMap;
    var attr = "";
    if (pm[i.typeID] != undefined) {
        attr += `<div style="font-size: 10px;">Value: ${simpleVal(pm[i.typeID], 2)} ISK</div>`
    }


    $.each(i, function (k, v) {
        if (!k.includes("chargeRate") && !k.includes("chargeSize") &&
            !k.includes("HeatAbsorb") && !k.includes("Group") &&
            k != "optimalSigRadius" && k != "mass" && k != "hp") {
            icon = attributeNameIconMap[k];
            if (icon != undefined) {
                attr += `
            <div><img src="/i/icons/${icon}" style="height: 23px; width: 23px"> ${k}: ${v}</div>
            `
            }
        }
    });

    return `
        <div style="font-weight: bolder;">${t[i.typeID]}</div>
        ${attr}
    `;
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
    if (graphicsMap[type] && ship != undefined) {
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
            gl = canvas.getContext("webgl"),
            effect;

        ccpwgl.initialize(canvas, {});
        ccpwgl_int.store.RegisterPath('local', 'https://www.evedata.org/')

        camera = ccpwgl.createCamera(canvas, {}, true);
        scene = ccpwgl.loadScene(sceneList[Math.floor(Math.random() * sceneList.length)]);
        ship = scene.loadShip(package.dna/*, function () {
            ship.addArmorEffects('local:/i/shaders/fxarmorimpactv5.sm_hi',
                (eff) => {
                    effect = eff;

                },
                (err) => { console.log(err); throw err });
        }*/);
        scene.loadSun(sunList[Math.floor(Math.random() * sunList.length)]);

        ccpwgl.onPreRender = function (dt) {
            resizeCanvasToDisplaySize(canvas, window.devicePixelRatio);
            gl.viewport(0, 0, gl.canvas.width, gl.canvas.height);

            if (window.innerHeight != screen.height) {
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

