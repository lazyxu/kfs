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

export function parseDateTime(exif) {
    return moment.parseZone(exif.DateTime + " " + exif.OffsetTime, "YYYY:MM:DD HH:mm:ss ZZ");
}
