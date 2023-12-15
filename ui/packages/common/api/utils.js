import moment from 'moment';

export function modeIsDir(mode) {
    return mode >= 2147483648;
}

export function getPerm(mode) {
    return mode >= 2147483648 ? mode - 2147483648 : mode;
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

export function isImage(name) {
    name = name.toLowerCase()
    if (name.endsWith(".jpg") || name.endsWith(".jpeg") || name.endsWith(".png") || name.endsWith(".heic")) {
        return true;
    }
    return false;
}

export function isVideo(name) {
    name = name.toLowerCase()
    if (name.endsWith(".mp4") || name.endsWith(".mov")) {
        return true;
    }
    return false;
}

export function isViewable(name) {
    name = name.toLowerCase()
    if (name.endsWith(".txt") || name.endsWith(".md") || name.endsWith(".log") ||
        name.endsWith(".go") || name.endsWith(".js") || name.endsWith(".java") || name.endsWith(".c")
    ) {
        return true;
    }
    return false;
}

export function parseShotTime(metadata) {
    let { exif, videoMetadata } = metadata;
    if (exif) {
        if (exif.DateTime) {
            return moment.parseZone(exif.DateTime + " " + exif.OffsetTime, "YYYY:MM:DD HH:mm:ss ZZ");
        } else if (exif.DateTimeOriginal) {
            return moment.parseZone(exif.DateTimeOriginal + " " + exif.OffsetTimeOriginal, "YYYY:MM:DD HH:mm:ss ZZ");
        }
        return moment.parseZone(exif.DateTimeDigitized + " " + exif.OffsetTimeDigitized, "YYYY:MM:DD HH:mm:ss ZZ");
    } else if (videoMetadata) {
        if (videoMetadata.Created) {
            return moment(videoMetadata.Created / 1000 / 1000);
        }
        return moment(videoMetadata.Modified / 1000 / 1000);
    }
    return moment.invalid();
}

export function timeSortFn(a, b) {
    if (!a.shotTime.isValid()) {
        return -1;
    } else if (!b.shotTime.isValid()) {
        return 1;
    }
    return a.shotTime.isAfter(b.shotTime) ? 1 : -1;
}

export function parseShotEquipment(metadata) {
    let { exif } = metadata;
    if (exif) {
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
    return "";
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