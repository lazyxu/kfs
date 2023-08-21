import moment from 'moment';

export function modeIsDir(mode) {
    return mode >= 2147483648;
}

export function getPerm(mode) {
    return mode >= 2147483648 ? mode - 2147483648 : 0;
}

export function isDCIM(name) {
    name = name.toLowerCase()
    if (name.endsWith(".jpg") || name.endsWith(".jpeg") || name.endsWith(".png")) {
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
