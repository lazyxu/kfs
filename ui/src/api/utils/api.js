import moment from 'moment';

export function modeIsDir(mode) {
    return mode >= 2147483648;
}

export function getPerm(mode) {
    return mode >= 2147483648 ? mode - 2147483648 : 0;
}

export function isDCIM(name) {
    name = name.toLowerCase()
    if (name.endsWith(".jpg") || name.endsWith(".jpeg") || name.endsWith(".png") || name.endsWith(".heic") ||
        name.endsWith(".mp4") || name.endsWith(".mov")
    ) {
        return true;
    }
    return false;
}

export function parseShotTime(exif) {
    if (exif.DateTime) {
        return moment.parseZone(exif.DateTime + " " + exif.OffsetTime, "YYYY:MM:DD HH:mm:ss ZZ");
    } else if (exif.DateTimeOriginal) {
        return moment.parseZone(exif.DateTimeOriginal + " " + exif.OffsetTimeOriginal, "YYYY:MM:DD HH:mm:ss ZZ");
    }
    return moment.parseZone(exif.DateTimeDigitized + " " + exif.OffsetTimeDigitized, "YYYY:MM:DD HH:mm:ss ZZ");
}

export function timeSortFn(a, b) {
    if (!a.shotTime.isValid()) {
        return -1;
    } else if (!b.shotTime.isValid()) {
        return 1;
    }
    return a.shotTime.isAfter(b.shotTime) ? 1 : -1;
}

export function parseShotEquipment(exif) {
    if (exif.Model) {
        if (exif.Model.includes(exif.Make)) {
            return exif.Model;
        } else {
            return exif.Make + " " + exif.Model;
        }
    } else {
        return exif.HostComputer;
    }
}

export function toPrecent(n) {
    return (Math.floor(n * 100000) / 1000).toFixed(3) + "%";
}

export function getTransform(orientation) {
    switch (orientation) {
        case 2:
            return "rotateY(180deg)";
        case 3:
            return "rotate(180deg)";
        case 4:
            return "rotate(180deg)rotateY(180deg)";
        case 5:
            return "rotate(270deg)rotateY(180deg)";
        case 6:
            return "rotate(90deg)";
        case 7:
            return "rotate(90deg)rotateY(180deg)";
        case 8:
            return "rotate(270deg)";
    }
    return "";
}